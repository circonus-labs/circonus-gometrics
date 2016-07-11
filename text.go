package circonusgometrics

// A Text metric is an arbitrary string
//

// Set a text metric
func (m *CirconusMetrics) SetText(metric string, val string) {
	m.SetTextValue(metric, val)
}

// Set a text metric
func (m *CirconusMetrics) SetTextValue(metric string, val string) {
	m.tm.Lock()
	defer m.tm.Unlock()
	m.text[metric] = val
}

// Remove a text metric
func (m *CirconusMetrics) RemoveText(metric string) {
	m.tm.Lock()
	defer m.tm.Unlock()
	delete(m.text, metric)
}

// Set a text metric to a function [called at flush interval]
func (m *CirconusMetrics) SetTextFunc(metric string, fn func() string) {
	m.tfm.Lock()
	defer m.tfm.Unlock()
	m.textFuncs[metric] = fn
}

// Remove a text metric function
func (m *CirconusMetrics) RemoveTextFunc(metric string) {
	m.tfm.Lock()
	defer m.tfm.Unlock()
	delete(m.textFuncs, metric)
}
