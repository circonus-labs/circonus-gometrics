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
	"sync"
	"time"
)

var (
	counters      = make(map[string]uint64)
	counterFuncs  = make(map[string]func() uint64)
	gauges        = make(map[string]func() int64)
	histograms    = make(map[string]*Histogram)
	rootCA        = x509.NewCertPool()
	apiUrl        = "api.circonus.com"
	authtoken     = ""
	authapp       = "circonus-gometrics"
	checkId       = 0
	submissionUrl = ""
	interval      = 10 * time.Second

	cm, gm, hm sync.Mutex
)

// WithInterval sets the reporting interval to Circonus (Default 10s)
func WithInterval(_interval time.Duration) {
	if _interval > 0 {
		interval = _interval
	}
}

// WithAuthToken sets the Auth Token for Circonus API services
func WithAuthToken(_authtoken string) {
	authtoken = _authtoken
}

// WithAuthApp sets the Auth App Name for Circonus API services
func WithAuthApp(_authapp string) {
	if _authapp != "" {
		authapp = _authapp
	}
}

// WithApiHost sets the API host for Circonus API services
func WithApiHost(_apihost string) {
	if _apihost != "" {
		apiUrl = _apihost
	}
}

// WithCheckId sets the Circonus check id from which to determine the endpoint
func WithCheckId(_id int) {
	checkId = _id
}

// WithSubmissionUrl sets the endpoint explicitly (not needed if WithCheckId is set)
func WithSubmissionUrl(url string) {
	if url != "" {
		submissionUrl = url
	}
}

// Start starts a perdiodic submission process of all metrics collected
func Start() {
	go func() {
		loadCACert()
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
