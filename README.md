# Circonus metrics tracking for Go applications

This library supports named counters, gauges and histograms.
It also provides convenience wrappers for registering latency
instrumented functions with Go's builtin http server.

Initializing only requires setting an ApiToken.

## Example

**rough and simple**

```go
package main

import (
	"log"
    "time"
	"math/rand"

	"github.com/circonus-labs/circonus-gometrics"
)

func main() {

	metrics := circonusgometrics.NewCirconusMetrics()

    // from circonus UI tokens page
	metrics.ApiToken = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxx"

    // app name associated with token
    //metrics.ApiApp = ""                           // default: 'circonus-gometrics'

    // fqdn of the circonus api server
    //metrics.ApiHost = ""                          // default: 'api.circonus.com'

    // interval at which metrics should be sent to circonus
    //metrics.Interval = 60 * time.seconds          // default: 10 seconds

    // submission url for a previously created httptrap check
    //metrics.SubmissionUrl = "https://..."         // precedence 1

    // a specific **check** id (not check bundle id)
    //metrics.CheckId = 0                           // precedence 2

    // if neither a submission url nor check id are provided, an attempt will be made to find an existing
    // httptrap check by using the circonus api to search for a check matching the following criteria:
    //      an active check,
    //      of type httptrap,
    //      where the target/host is equal to metrics.InstanceId - see below
    //      and the check has a tag equal to metrics.SearchTag - see below

    // an identifier for the 'group of metrics emitted by this process or service'
	//metrics.InstanceId = "centos7.gmtest"          // default: 'hostname':'program name'

    // a specific tag which, when coupled with the instanceid serves to identify the
    // origin and/or grouping of the metrics
    //metrics.SearchTag = "service:gmtest"          // default: service:'program name'

    // if an applicable check is NOT specified or found, an attempt will be made to automatically create one

    // "GROUP ID" for a specific broker from the Brokers page in circonus ui
    // metrics.BrokerGroupId = 58938                // default: random enterprise broker

    // used to select a broker with the same tag (e.g. can be used to dictate that a broker
    // serving a specific location should be used. "dc:sfo", "location:new_york", "zone:us-west")
    // if more than one broker has the tag, one will be selected randomly from the resulting list
    // metrics.BrokerSelectTag = ""                 // default: not used

    // if no BrokerGroupId or BrokerSelectTag is specified a broker will be selected randomly
    // from the list of brokers available to the api token. enterprise brokers take precedence
    // viable brokers are "active" and have the "httptrap" module enabled.

    // additional tags to add to an automatically created check (array of strings)
    // metrics.Tags = []string{"category:tag", "category:tag"} // default: none

    // specific check secret
    // metrics.CheckSecret = "a!secret"             // default: randomly generated

    // custom logger
    //metrics.Log =                                 // default: discards messages (unless debug is true)

    // emit debugging messages
	metrics.Debug = true                           // default: false

	src := rand.NewSource(time.Now().UnixNano())
	rnd := rand.New(src)

	metrics.Start()

	for i := 1; i < 60; i++ {
		metrics.Timing("ding", rnd.Float64() * 10)
		metrics.Increment("dong")
		metrics.Gauge("dang", 10)
		time.Sleep(1000 * time.Millisecond)
	}

    // ensure last bit are flushed (or if it has run for less than interval)
	metrics.Flush()

}
```

# untested

### HTTP Handler wrapping

```
http.HandleFunc("/", metrics.TrackHTTPLatency("/", handler_func))
```

### HTTP latency example

```
package main

import (
    "fmt"
    "net/http"
    metrics "github.com/circonus-labs/circonus-gometrics"
)

func main() {
    metrics.WithAuthToken("9fdd5432-5308-4691-acd1-6bf1f7a20f73")
    metrics.WithCheckId(115010)
    metrics.Start()

    http.HandleFunc("/", metrics.TrackHTTPLatency("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
    }))
    http.ListenAndServe(":8080", http.DefaultServeMux)
}

```
