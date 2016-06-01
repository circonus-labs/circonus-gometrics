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

/*
type Counter string

// Add increments the counter by one.
func (c Counter) Add() {
	c.AddN(1)
}

// AddN increments the counter by N.
func (c Counter) AddN(delta uint64) {
	cm.Lock()
	defer cm.Unlock()
	counters[string(c)] += delta
}

// SetFunc sets the counter's value to the lazily-called return value of the
// given function.
func (c Counter) SetFunc(f func() uint64) {
	cm.Lock()
	defer cm.Unlock()
	counterFuncs[string(c)] = f
}

// Remove removes the given counter.
func (c Counter) Remove() {
	cm.Lock()
	defer cm.Unlock()
	delete(counters, string(c))
	delete(counterFuncs, string(c))
}
*/
