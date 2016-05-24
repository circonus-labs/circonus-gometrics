package circonusgometrics

import (
	"github.com/circonus-labs/circonusllhist"
)

// Reset removes all existing counters and gauges.
func Reset() {
	cm.Lock()
	defer cm.Unlock()

	gm.Lock()
	defer gm.Unlock()

	hm.Lock()
	defer hm.Unlock()

	counters = make(map[string]uint64)
	counterFuncs = make(map[string]func() uint64)
	gauges = make(map[string]func() int64)
	histograms = make(map[string]*Histogram)
}

// snapshot returns a copy of the values of all registered counters and gauges.
func snapshot() (c map[string]uint64, g map[string]int64, h map[string]*circonusllhist.Histogram) {
	cm.Lock()
	defer cm.Unlock()

	gm.Lock()
	defer gm.Unlock()

	hm.Lock()
	defer hm.Unlock()

	c = make(map[string]uint64, len(counters)+len(counterFuncs))
	for n, v := range counters {
		c[n] = v
	}

	for n, f := range counterFuncs {
		c[n] = f()
	}

	g = make(map[string]int64, len(gauges))
	for n, f := range gauges {
		g[n] = f()
	}

	h = make(map[string]*circonusllhist.Histogram, len(histograms))
	for n, hist := range histograms {
		h[n] = hist.hist.CopyAndReset()
	}

	return
}
