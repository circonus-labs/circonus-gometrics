package circonusgometrics

// A Gauge is an instantaneous measurement of a value.
//
// Use a gauge to track metrics which increase and decrease (e.g., amount of
// free memory).
type Gauge string

// Set the gauge's value to the given value.
func (g Gauge) Set(value int64) {
	gm.Lock()
	defer gm.Unlock()

	gauges[string(g)] = func() int64 {
		return value
	}
}

// SetFunc sets the gauge's value to the lazily-called return value of the given
// function.
func (g Gauge) SetFunc(f func() int64) {
	gm.Lock()
	defer gm.Unlock()

	gauges[string(g)] = f
}

// Remove removes the given gauge.
func (g Gauge) Remove() {
	gm.Lock()
	defer gm.Unlock()

	delete(gauges, string(g))
}
