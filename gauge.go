package circonusgometrics

// A Gauge is an instantaneous measurement of a value.
//
// Use a gauge to track metrics which increase and decrease (e.g., amount of
// free memory).

// Set a gauge to a value
func (m *CirconusMetrics) Gauge(metric string, val int64) {
	m.SetGauge(metric, val)
}

// Set a gauge to a value
func (m *CirconusMetrics) SetGauge(metric string, val int64) {
	m.gm.Lock()
	defer m.gm.Unlock()
	m.gauges[metric] = val
}

// Remove a gauge
func (m *CirconusMetrics) RemoveGauge(metric string) {
	m.gm.Lock()
	defer m.gm.Unlock()
	delete(m.gauges, metric)
}

// Set a gauge to a function [called at flush interval]
func (m *CirconusMetrics) SetGaugeFunc(metric string, fn func() int64) {
	m.gfm.Lock()
	defer m.gfm.Unlock()
	m.gaugeFuncs[metric] = fn
}

// Remove a gauge function
func (m *CirconusMetrics) RemoveGaugeFunc(metric string) {
	m.gfm.Lock()
	defer m.gfm.Unlock()
	delete(m.gaugeFuncs, metric)
}
