package circonusgometrics

// A Counter is a monotonically increasing unsigned integer.
//
// Use a counter to derive rates (e.g., record total number of requests, derive
// requests per second).

func (m *CirconusMetrics) Increment(metric string) {
	m.Add(metric, 1)
}

func (m *CirconusMetrics) IncrementByValue(metric string, val uint64) {
	m.Add(metric, val)
}

func (m *CirconusMetrics) Add(metric string, val uint64) {
	m.cm.Lock()
	defer m.cm.Unlock()
	m.counters[metric] += val
}

func (m *CirconusMetrics) RemoveCounter(metric string) {
	m.cm.Lock()
	defer m.cm.Unlock()
	delete(m.counters, metric)
}

func (m *CirconusMetrics) SetCounterFunc(metric string, fn func() uint64) {
	m.cfm.Lock()
	defer m.cfm.Unlock()
	m.counterFuncs[metric] = fn
}

func (m *CirconusMetrics) RemoveCounterFunc(metric string) {
	m.cfm.Lock()
	defer m.cfm.Unlock()
	delete(m.counterFuncs, metric)
}
