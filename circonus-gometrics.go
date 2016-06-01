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
	"io/ioutil"
	"log"
	"os"
	"path"
	"sync"
	"time"
)

const (
	// a few sensible defaults
	defaultApiHost  = "api.circonus.com"
	defaultApiApp   = "circonus-gometrics"
	defaultInterval = 10 * time.Second
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
	ApiToken      string
	SubmissionUrl string
	CheckId       int
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
	ready   bool
	trapUrl string
	trapCN  string
	trapmu  sync.Mutex

	certPool      *x509.CertPool
	cert          []byte
	checkBundle   *CheckBundle
	activeMetrics map[string]bool
	checkType     string

	counters map[string]uint64
	cm       sync.Mutex

	counterFuncs map[string]func() uint64
	cfm          sync.Mutex

	gauges map[string]int64
	gm     sync.Mutex

	gaugeFuncs map[string]func() int64
	gfm        sync.Mutex

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
		InstanceId:    fmt.Sprintf("%s:%s", hn, an),
		SearchTag:     fmt.Sprintf("service:%s", an),
		ApiHost:       defaultApiHost,
		ApiApp:        defaultApiApp,
		Interval:      defaultInterval,
		Log:           log.New(ioutil.Discard, "", log.LstdFlags),
		Debug:         false,
		ready:         false,
		trapUrl:       "",
		activeMetrics: make(map[string]bool),
		counters:      make(map[string]uint64),
		counterFuncs:  make(map[string]func() uint64),
		gauges:        make(map[string]int64),
		gaugeFuncs:    make(map[string]func() int64),
		histograms:    make(map[string]*Histogram),
		certPool:      x509.NewCertPool(),
		checkType:     "httptrap",
	}

}

// Start starts a perdiodic submission process of all metrics collected
func (m *CirconusMetrics) Start() {
	go func() {
		if m.Debug {
			m.Log = log.New(os.Stderr, "", log.LstdFlags)
		}
		m.loadCACert()
		if !m.ready {
			m.initializeTrap()
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
	if !m.ready {
		if err := m.initializeTrap(); err != nil {
			m.Log.Println("Unable to initialize check, NOT flushing metrics.")
			return
		}
	}

	// check for new metrics and enable them automatically
	newMetrics := make(map[string]*CheckBundleMetric)

	counters, gauges, histograms := m.snapshot()
	output := make(map[string]interface{})
	for name, value := range counters {
		output[name] = map[string]interface{}{
			"_type":  "n",
			"_value": value,
		}
		if _, ok := m.activeMetrics[name]; !ok {
			newMetrics[name] = &CheckBundleMetric{
				Name:   name,
				Type:   "numeric",
				Status: "active",
			}
		}
	}

	for name, value := range gauges {
		output[name] = map[string]interface{}{
			"_type":  "n",
			"_value": value,
		}
		if _, ok := m.activeMetrics[name]; !ok {
			newMetrics[name] = &CheckBundleMetric{
				Name:   name,
				Type:   "numeric",
				Status: "active",
			}
		}
	}

	for name, value := range histograms {
		output[name] = map[string]interface{}{
			"_type":  "n",
			"_value": value.DecStrings(),
		}
		if _, ok := m.activeMetrics[name]; !ok {
			newMetrics[name] = &CheckBundleMetric{
				Name:   name,
				Type:   "histogram",
				Status: "active",
			}
		}
	}

	m.submit(output, newMetrics)
}
