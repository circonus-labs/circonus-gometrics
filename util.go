package circonusgometrics

import (
	"github.com/circonus-labs/circonusllhist"
)

// Reset removes all existing counters and gauges.
func (m *CirconusMetrics) Reset() {
	m.cm.Lock()
	defer m.cm.Unlock()

	m.cfm.Lock()
	defer m.cfm.Unlock()

	m.gm.Lock()
	defer m.gm.Unlock()

	m.hm.Lock()
	defer m.hm.Unlock()

	m.counters = make(map[string]uint64)
	m.counterFuncs = make(map[string]func() uint64)
	//	m.gauges = make(map[string]func() int64)
	m.gauges = make(map[string]int64)
	m.histograms = make(map[string]*Histogram)
}

// snapshot returns a copy of the values of all registered counters and gauges.
func (m *CirconusMetrics) snapshot() (c map[string]uint64, g map[string]int64, h map[string]*circonusllhist.Histogram) {
	m.cm.Lock()
	defer m.cm.Unlock()

	m.cfm.Lock()
	defer m.cfm.Unlock()

	m.gm.Lock()
	defer m.gm.Unlock()

	m.hm.Lock()
	defer m.hm.Unlock()

	c = make(map[string]uint64, len(m.counters)+len(m.counterFuncs))
	for n, v := range m.counters {
		c[n] = v
	}

	for n, f := range m.counterFuncs {
		c[n] = f()
	}

	// g = make(map[string]int64, len(m.gauges))
	// for n, f := range m.gauges {
	// 	g[n] = f()
	// }

	g = make(map[string]int64, len(m.gauges))
	for n, v := range m.gauges {
		g[n] = v
	}

	h = make(map[string]*circonusllhist.Histogram, len(m.histograms))
	for n, hist := range m.histograms {
		h[n] = hist.hist.CopyAndReset()
	}

	return
}
