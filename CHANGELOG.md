# v3.0.0-alpha.4

* add: missing SetHistogramDurationWithTags
* upd: go-apiclient v0.5.1
* fix: remove cgm v2 dependency from DEPRECATED api package
* upd: retryablehttp v0.5.0

# v3.0.0-alpha.3

* add: RecordDuration, RecordDurationWithTags, SetHistogramDuration

# v3.0.0-alpha.2

* upd: circllhist v0.1.2

# v3.0.0-alpha.1

* fix: enable check management for add tags test
* fix: api.circonus.com hostname (accidentally changed during switch to apiclient)

# v3.0.0-alpha

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
