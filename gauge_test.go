package circonusgometrics

import (
	"testing"
)

func TestGauge(t *testing.T) {
	t.Log("Testing gauge.Gauge")

	cm := &CirconusMetrics{}
	cm.gauges = make(map[string]int64)
	cm.Gauge("foo", 1)

	val, ok := cm.gauges["foo"]
	if !ok {
		t.Errorf("Expected to find foo")
	}

	if val != 1 {
		t.Errorf("Expected 1, found %d", val)
	}
}

func TestSetGauge(t *testing.T) {
	t.Log("Testing gauge.SetGauge")

	cm := &CirconusMetrics{}
	cm.gauges = make(map[string]int64)
	cm.SetGauge("foo", 10)

	val, ok := cm.gauges["foo"]
	if !ok {
		t.Errorf("Expected to find foo")
	}

	if val != 10 {
		t.Errorf("Expected 10, found %d", val)
	}
}

func TestRemoveGauge(t *testing.T) {
	t.Log("Testing gauge.RemoveGauge")

	cm := &CirconusMetrics{}
	cm.gauges = make(map[string]int64)
	cm.Gauge("foo", 5)

	val, ok := cm.gauges["foo"]
	if !ok {
		t.Errorf("Expected to find foo")
	}

	if val != 5 {
		t.Errorf("Expected 5, found %d", val)
	}

	cm.RemoveGauge("foo")

	val, ok = cm.gauges["foo"]
	if ok {
		t.Errorf("Expected NOT to find foo")
	}

	if val != 0 {
		t.Errorf("Expected 0, found %d", val)
	}
}
