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
	"crypto/x509"
	"fmt"
	"log"
	"os"
	"path"
	"sync"
	"time"
)

var (
	// a few sensible defaults
	defaultApiHost  = "api.circonus.com"
	defaultApiApp   = "circonus-gometrics"
	defaultInterval = 10 * time.Second
	checkType       = "httptrap"
)

// a few words about: "BrokerGroupId"
//
// calling it this because the instructions for how to get into the UI and FIND this value are more straight-forward:
//
// log into ui
// navigate to brokers page
// identify which broker you need to use
// click the little down arrow in the circle on the right-hand side of the line for the broker you'd like to use
// use the value from the "GROUP ID:" field under "Broker Details" in the drop-down afetr clicking the down arrow
//
// ... or ...
//
// log into ui
// navigate to brokers page
// identify which broker you need to use
// click the hamburger menu icon (three lines to the left of the broker name)
// click "view API object" from the drop-down menu
// look for "_cid" field, use integer value after "/broker/" e.g. "/broker/35" would be 35
//

type CirconusMetrics struct {
	SubmissionUrl string
	ApiToken      string
	ApiApp        string
	ApiHost       string
	InstanceId    string
	SearchTag     string
	BrokerGroupId int
	Tags          []string
	CheckSecret   string

	Interval time.Duration
	Log      *log.Logger
	Debug    bool

	// internals
	trapUrl string

	certPool *x509.CertPool
	cert     []byte
	check    Check

	counters map[string]uint64
	cm       sync.Mutex

	counterFuncs map[string]func() uint64
	cfm          sync.Mutex

	//gauges map[string]func() int64
	gauges map[string]int64
	gm     sync.Mutex

	histograms map[string]*Histogram
	hm         sync.Mutex
}

func NewCirconusMetrics() *CirconusMetrics {
	_, an := path.Split(os.Args[0])
	hn, err := os.Hostname()
	if err != nil {
		hn = "unknown"
	}

	return &CirconusMetrics{
		InstanceId:   fmt.Sprintf("%s:%s", hn, an),
		SearchTag:    fmt.Sprintf("service:%s", an),
		ApiHost:      defaultApiHost,
		ApiApp:       defaultApiApp,
		Interval:     defaultInterval,
		Log:          log.New(os.Stderr, "", log.LstdFlags),
		Debug:        false,
		counterFuncs: make(map[string]func() uint64),
		counters:     make(map[string]uint64),
		//		gauges:       make(map[string]func() int64),
		gauges:     make(map[string]int64),
		histograms: make(map[string]*Histogram),
		certPool:   x509.NewCertPool(),
	}

}

// Start starts a perdiodic submission process of all metrics collected
func (m *CirconusMetrics) Start() {
	go func() {
		m.loadCACert()
		if m.trapUrl == "" {
			url, err := m.getTrapUrl()
			if err != nil {
				m.Log.Printf("%+v\n", err)
			}
			m.trapUrl = url
		}
	}()

	go func() {
		for _ = range time.NewTicker(m.Interval).C {
			m.Flush()
		}
	}()
}

func (m *CirconusMetrics) Flush() {
	m.Log.Println("Flushing")
	if m.trapUrl == "" {
		url, err := m.getTrapUrl()
		if err != nil {
			m.Log.Printf("%+v\n", err)
		}
		m.trapUrl = url
	}
	counters, gauges, histograms := m.snapshot()
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
	m.submit(output)
}
