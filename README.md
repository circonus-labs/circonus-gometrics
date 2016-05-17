# Circonus metrics tracking for Go applications

This library supports named counters, gauges and histograms.
It also provides convenience wrappers for registering latency
instrumented functions with Go's builtin http server.

Initializing only requires you set the AuthToken and CheckId and "Start" the metrics reporter.

### Counters

```
metrics.Counter("widgets").Add()
metrics.Counter("widgets").AddN(1)
metrics.Counter("lazy_widgets").SetFunc(func () int64 {
    return total_widgets_sold
})
```

### Gauges

```
metrics.Gauge("temperature").Set(78.2)
metrics.Gauge("lazy_gauge").SetFunc(func () int64 {
    return some_value
})
```

### Histograms

```
hist := metrics.NewHistogram("read_size")
hist.RecordValue(513)
hist.RecordValue(562)
```

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
