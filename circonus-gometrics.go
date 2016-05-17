// Package circonusgometrics provides instrumentation for your applications in the form
// of counters, gauges and histograms and allows you to publish them to
// Circonus
//
// Counters
//
// A counter is a monotonically-increasing, unsigned, 64-bit integer used to
// represent the number of times an event has occurred. By tracking the deltas
// between measurements of a counter over intervals of time, an aggregation
// layer can derive rates, acceleration, etc.
//
// Gauges
//
// A gauge returns instantaneous measurements of something using signed, 64-bit
// integers. This value does not need to be monotonic.
//
// Histograms
//
// A histogram tracks the distribution of a stream of values (e.g. the number of
// seconds it takes to handle requests).  Circonus can calculate complex
// analytics on these.
//
// Reporting
//
// A period push to a Circonus httptrap is confgurable.

package circonusgometrics

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/circonus-labs/circonusllhist"
)

// A Counter is a monotonically increasing unsigned integer.
//
// Use a counter to derive rates (e.g., record total number of requests, derive
// requests per second).
type Counter string

// Add increments the counter by one.
func (c Counter) Add() {
	c.AddN(1)
}

// AddN increments the counter by N.
func (c Counter) AddN(delta uint64) {
	cm.Lock()
	defer cm.Unlock()
	counters[string(c)] += delta
}

// SetFunc sets the counter's value to the lazily-called return value of the
// given function.
func (c Counter) SetFunc(f func() uint64) {
	cm.Lock()
	defer cm.Unlock()
	counterFuncs[string(c)] = f
}

// Remove removes the given counter.
func (c Counter) Remove() {
	cm.Lock()
	defer cm.Unlock()
	delete(counters, string(c))
	delete(counterFuncs, string(c))
}

// A Gauge is an instantaneous measurement of a value.
//
// Use a gauge to track metrics which increase and decrease (e.g., amount of
// free memory).
type Gauge string

// Set the gauge's value to the given value.
func (g Gauge) Set(value int64) {
	gm.Lock()
	defer gm.Unlock()

	gauges[string(g)] = func() int64 {
		return value
	}
}

// SetFunc sets the gauge's value to the lazily-called return value of the given
// function.
func (g Gauge) SetFunc(f func() int64) {
	gm.Lock()
	defer gm.Unlock()

	gauges[string(g)] = f
}

// Remove removes the given gauge.
func (g Gauge) Remove() {
	gm.Lock()
	defer gm.Unlock()

	delete(gauges, string(g))
}

// Reset removes all existing counters and gauges.
func Reset() {
	cm.Lock()
	defer cm.Unlock()

	gm.Lock()
	defer gm.Unlock()

	hm.Lock()
	defer hm.Unlock()

	counters = make(map[string]uint64)
	counterFuncs = make(map[string]func() uint64)
	gauges = make(map[string]func() int64)
	histograms = make(map[string]*Histogram)
}

// snapshot returns a copy of the values of all registered counters and gauges.
func snapshot() (c map[string]uint64, g map[string]int64, h map[string]*circonusllhist.Histogram) {
	cm.Lock()
	defer cm.Unlock()

	gm.Lock()
	defer gm.Unlock()

	hm.Lock()
	defer hm.Unlock()

	c = make(map[string]uint64, len(counters)+len(counterFuncs))
	for n, v := range counters {
		c[n] = v
	}

	for n, f := range counterFuncs {
		c[n] = f()
	}

	g = make(map[string]int64, len(gauges))
	for n, f := range gauges {
		g[n] = f()
	}

	h = make(map[string]*circonusllhist.Histogram, len(histograms))
	for n, hist := range histograms {
		h[n] = hist.hist.CopyAndReset()
	}

	return
}

// NewHistogram returns a new Circonus histogram that accumulates until reported on.
func NewHistogram(name string) *Histogram {
	hm.Lock()
	defer hm.Unlock()

	if hist, ok := histograms[name]; ok {
		return hist
	}

	hist := &Histogram{
		name: name,
		hist: circonusllhist.New(),
	}
	histograms[name] = hist
	return hist
}

// Remove removes the given histogram.
func (h *Histogram) Remove() {
	hm.Lock()
	defer hm.Unlock()
	delete(histograms, h.name)
}

type hname string // unexported to prevent collisions

// A Histogram measures the distribution of a stream of values.
type Histogram struct {
	name string
	hist *circonusllhist.Histogram
	rw   sync.RWMutex
}

// Name returns the name of the histogram
func (h *Histogram) Name() string {
	return h.name
}

// RecordValue records the given value
func (h *Histogram) RecordValue(v float64) {
	h.rw.Lock()
	defer h.rw.Unlock()

	h.hist.RecordValue(v)
}

var (
	counters      = make(map[string]uint64)
	counterFuncs  = make(map[string]func() uint64)
	gauges        = make(map[string]func() int64)
	histograms    = make(map[string]*Histogram)
	rootCA        = x509.NewCertPool()
	apiUrl        = "api.circonus.com"
	authtoken     = ""
	checkId       = 0
	submissionUrl = ""
	interval      = 10 * time.Second

	cm, gm, hm sync.Mutex
)

func submit(output map[string]interface{}) {
	str, err := json.Marshal(output)
	if err == nil {
		trapCall(str)
	}
}

// WithInterval sets the reporting interval to Circonus (Default 10s)
func WithInterval(_interval time.Duration) {
	interval = _interval
}

// WithAuthToken sets the Auth Token for Circonus API services
func WithAuthToken(_authtoken string) {
	authtoken = _authtoken
}

// WithCheckId sets the Circonus check id from which to determine the endpoint
func WithCheckId(_id int) {
	checkId = _id
}

// WithSubmissionUrl sets the endpoint explicitly (not needed if WithCheckId is set)
func WithSubmissionUrl(url string) {
	submissionUrl = url
}

func apiCall(url string) map[string]interface{} {
	client := &http.Client{}
	req, err := http.NewRequest("GET", strings.Join([]string{"https://", apiUrl, url}, ""), nil)
	if err != nil {
		return nil
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("X-Circonus-Auth-Token", authtoken)
	req.Header.Add("X-Circonus-App-Name", "circonus-cip")
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error fetching %s: %s\n", url, err)
		return nil
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	var response map[string]interface{}
	json.Unmarshal(body, &response)
	if resp.StatusCode != 200 {
		log.Printf("response: %v\n", response)
		return nil
	}
	return response
}

func parsePEMCertificates(pemData []byte) ([]*x509.Certificate, error) {
	var certs []*x509.Certificate
	for {
		var der *pem.Block
		der, pemData = pem.Decode(pemData)
		if der == nil {
			break
		}
		if der.Type == "CERTIFICATE" {
			dcerts, err := x509.ParseCertificates(der.Bytes)
			if err != nil {
				return nil, err
			}
			certs = append(certs, dcerts...)
		}
	}
	return certs, nil
}
func getCertChain() {
	caDetails := apiCall("/v2/pki/ca.crt")
	val, ok := caDetails["contents"]
	if !ok {
		log.Printf("Error fetching ca.crt\n")
		setRootCA(circonusCA)
		return
	}

	setRootCA([]byte(val.(string)))
	log.Print("Circonusgometrics fetched CA.")
}

func setRootCA(val []byte) {
	certs, err := parsePEMCertificates(val)
	if err != nil {
		return
	}
	for _, cert := range certs {
		rootCA.AddCert(cert)
	}
}

func getCheck(id int) {
	url := strings.Join([]string{"/v2/check/", strconv.Itoa(id)}, "")
	checkDetails := apiCall(url)
	details, ok := checkDetails["_details"]
	if !ok {
		log.Printf("Cannot find submission URL at %s\n", url)
		return
	}
	dmap := details.(map[string]interface{})
	val := dmap["submission_url"]
	submissionUrl = val.(string)
}

func trapCall(payload []byte) (int, error) {
	tr := &http.Transport{
		TLSClientConfig:    &tls.Config{RootCAs: rootCA},
		DisableCompression: true,
	}
	client := &http.Client{Transport: tr}
	req, err := http.NewRequest("POST", submissionUrl, bytes.NewBuffer(payload))
	if err != nil {
		return 0, err
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("X-Circonus-Auth-Token", authtoken)
	req.Header.Add("X-Circonus-App-Name", "circonus-cip")
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	var response map[string]interface{}
	json.Unmarshal(body, &response)
	if resp.StatusCode != 200 {
		return 0, errors.New("bad response code: " + strconv.Itoa(resp.StatusCode))
	}
	switch v := response["stats"].(type) {
	case float64:
		return int(v), nil
	case int:
		return v, nil
	default:
	}
	return 0, errors.New("bad response type")
}

// Start starts a perdiodic submission process of all metrics collected
func Start() {
	go func() {
		getCertChain()
		if checkId != 0 {
			getCheck(checkId)
		}
		for _ = range time.NewTicker(interval).C {
			counters, gauges, histograms := snapshot()
			output := make(map[string]interface{})
			for name, value := range counters {
				output[name] = map[string]interface{}{
					"_type":  "n",
					"_value": value,
				}
			}
			for name, value := range gauges {
				output[name] = map[string]interface{}{
					"_type":  "n",
					"_value": value,
				}
			}
			for name, value := range histograms {
				output[name] = map[string]interface{}{
					"_type":  "n",
					"_value": value.DecStrings(),
				}
			}
			submit(output)
		}
	}()
}

// TrackHTTPLatency wraps Handler functions registered with an http.ServerMux tracking latencies.
// Metrics are of the for go`HTTP`<method>`<name>`latency and are tracked in a histogram in units
// of seconds (as a float64) providing nanosecond ganularity.
func TrackHTTPLatency(name string, handler func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		start := time.Now().UnixNano()
		handler(rw, req)
		elapsed := time.Now().UnixNano() - start
		hist := NewHistogram("go`HTTP`" + req.Method + "`" + name + "`latency")
		hist.RecordValue(float64(elapsed) / float64(time.Second))
	}
}
