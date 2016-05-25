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
	"log"
	"os"
	"sync"
	"time"
)

var (
	// a few sensible defaults
	defaultApiHost  = "api.circonus.com"
	defaultApiApp   = "circonus-gometrics"
	defaultInterval = 10 * time.Second
	/*
	   if no Broker set by user
	   1. call api for list of brokers (available to account for which token is valid)
	   2. iterate list, eliminate any brokers which do not support "httptrap"
	   3. prioritize, enterprise then others
	   4. selecting
	       first, ... first one from the resulting list
	       random, a random one from the resulting list
	       time, iterate the list and time connections to each, pick fastest

	   these methods are fundamentally flawed w/re to WHERE the origin is and WHERE the broker is from a networking perspective...
	   .. consider adding "fastest" and doing conn timings for the brokers in the list selecting the fastest (still flawed from ops/load
	      perspective but, potentially less-so from a durability of metric collection perspective).
	*/
	defaultBrokerSelectionMethod = "time" // time, random or first

	// internals

	counters     = make(map[string]uint64)
	counterFuncs = make(map[string]func() uint64)
	gauges       = make(map[string]func() int64)
	histograms   = make(map[string]*Histogram)
	cm, gm, hm   sync.Mutex

	checkType         = "httptrap"
	circonusSearchTag = "service:consul"
	rootCA            = x509.NewCertPool()
	// note, public broker, for inside customers pelase specify a broker group id
	circonusTrapBroker = 35
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
	ApiApp        string
	ApiHost       string
	BrokerGroupId int
	Tags          []string
	InstanceId    string
	Interval      time.Duration
	CheckSecret   string
	TrapUrl       string
	Log           *log.Logger
	Debug         bool
}

func NewCirconusMetrics() *CirconusMetrics {

	return &CirconusMetrics{
		ApiHost:  defaultApiHost,
		ApiApp:   defaultApiApp,
		Interval: defaultInterval,
		Log:      log.New(os.Stderr, "", log.LstdFlags),
		Debug:    false,
	}

}

func (m *CirconusMetrics) Test() {
	m.loadCACert()
	if m.TrapUrl == "" {
		url, err := m.getTrapUrl()
		if err != nil {
			m.Log.Printf("%+v\n", err)
		}
		m.TrapUrl = url
	}
}

// Start starts a perdiodic submission process of all metrics collected
func (m *CirconusMetrics) Start() {
	go func() {
		m.loadCACert()
		if m.TrapUrl == "" {
			url, err := m.getTrapUrl()
			if err != nil {
				m.Log.Printf("%+v\n", err)
			}
			m.TrapUrl = url
		}
		for _ = range time.NewTicker(m.Interval).C {
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
	}()
}
