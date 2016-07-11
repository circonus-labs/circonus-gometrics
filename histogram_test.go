package circonusgometrics

import (
	"reflect"
	"testing"
)

func TestTiming(t *testing.T) {
	t.Log("Testing histogram.Timing")

	cm := &CirconusMetrics{}
	cm.histograms = make(map[string]*Histogram)
	cm.Timing("foo", 1)

	hist, ok := cm.histograms["foo"]
	if !ok {
		t.Errorf("Expected to find foo")
	}

	if hist == nil {
		t.Errorf("Expected *Histogram, found %d", hist)
	}

	val := hist.hist.DecStrings()
	if len(val) != 1 {
		t.Errorf("Expected 1, found '0'", val)
	}

	expectedVal := "H[1.0e+00]=1"
	if val[0] != expectedVal {
		t.Errorf("Expected '%s', found '%s'", expectedVal, val[0])
	}
}

func TestRecordValue(t *testing.T) {
	t.Log("Testing histogram.RecordValue")

	cm := &CirconusMetrics{}
	cm.histograms = make(map[string]*Histogram)
	cm.RecordValue("foo", 1)

	hist, ok := cm.histograms["foo"]
	if !ok {
		t.Errorf("Expected to find foo")
	}

	if hist == nil {
		t.Errorf("Expected *Histogram, found %d", hist)
	}

	val := hist.hist.DecStrings()
	if len(val) != 1 {
		t.Errorf("Expected 1, found '0'", val)
	}

	expectedVal := "H[1.0e+00]=1"
	if val[0] != expectedVal {
		t.Errorf("Expected '%s', found '%s'", expectedVal, val[0])
	}
}

func TestSetHistogramValue(t *testing.T) {
	t.Log("Testing histogram.SetHistogramValue")

	cm := &CirconusMetrics{}
	cm.histograms = make(map[string]*Histogram)
	cm.SetHistogramValue("foo", 1)

	hist, ok := cm.histograms["foo"]
	if !ok {
		t.Errorf("Expected to find foo")
	}

	if hist == nil {
		t.Errorf("Expected *Histogram, found %d", hist)
	}

	val := hist.hist.DecStrings()
	if len(val) != 1 {
		t.Errorf("Expected 1, found '0'", val)
	}

	expectedVal := "H[1.0e+00]=1"
	if val[0] != expectedVal {
		t.Errorf("Expected '%s', found '%s'", expectedVal, val[0])
	}
}

func TestNewHistogram(t *testing.T) {
	t.Log("Testing histogram.NewHistogram")

	cm := &CirconusMetrics{}
	cm.histograms = make(map[string]*Histogram)
	hist1 := cm.NewHistogram("foo")

	actualType := reflect.TypeOf(hist1)
	expectedType := "*circonusgometrics.Histogram"
	if actualType.String() != expectedType {
		t.Errorf("Expected %s, got %s", expectedType, actualType.String())
	}
}

func TestRemoveHistogram(t *testing.T) {
	t.Log("Testing histogram.RemoveHistogram")

	cm := &CirconusMetrics{}
	cm.histograms = make(map[string]*Histogram)
	cm.SetHistogramValue("foo", 1)

	hist, ok := cm.histograms["foo"]
	if !ok {
		t.Errorf("Expected to find foo")
	}

	if hist == nil {
		t.Errorf("Expected *Histogram, found %d", hist)
	}

	val := hist.hist.DecStrings()
	if len(val) != 1 {
		t.Errorf("Expected 1, found '0'", val)
	}

	expectedVal := "H[1.0e+00]=1"
	if val[0] != expectedVal {
		t.Errorf("Expected '%s', found '%s'", expectedVal, val[0])
	}

	cm.RemoveHistogram("foo")

	hist, ok = cm.histograms["foo"]
	if ok {
		t.Errorf("Expected NOT to find foo")
	}

	if hist != nil {
		t.Errorf("Expected nil, found %d", hist)
	}
}
