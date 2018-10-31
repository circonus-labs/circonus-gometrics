# v3.0.0

* add: helper functions for metrics `*WithTags` e.g. `TimingWithTags(metricName,tagList,val)`
* upd: default new checks to use metric_filters
* add: metric_filters support
* upd: dependencies (circonusllhist v0.1.0)
* upd: change histograms from type 'n' to type 'h' in submissions
* upd: DEPRECATED github.com/circonus-labs/circonus-gometrics/api
* upd: switch to using github.com/circonus-labs/go-apiclient
* upd: merge other metric tag functions into tags
* add: helper methods for handling tags (for new stream tags syntax and old check_bundle.metrics.metric.tags)
* upd: merge other metric output functions into metric_output
* upd: merge util into metric_output (methods in util are specifically for working with metric outputs)
