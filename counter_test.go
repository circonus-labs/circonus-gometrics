package circonusgometrics

import (
	"testing"
)

func TestIncrement(t *testing.T) {
	t.Log("Testing counter.Increment")

	cm := &CirconusMetrics{}
	cm.counters = make(map[string]uint64)
	cm.Increment("foo")

	val, ok := cm.counters["foo"]
	if !ok {
		t.Errorf("Expected to find foo")
	}

	if val != 1 {
		t.Errorf("Expected 1, found %d", val)
	}
}

func TestIncrementByValue(t *testing.T) {
	t.Log("Testing counter.IncrementByValue")

	cm := &CirconusMetrics{}
	cm.counters = make(map[string]uint64)
	cm.IncrementByValue("foo", 10)

	val, ok := cm.counters["foo"]
	if !ok {
		t.Errorf("Expected to find foo")
	}

	if val != 10 {
		t.Errorf("Expected 1, found %d", val)
	}
}

func TestAdd(t *testing.T) {
	t.Log("Testing counter.Add")

	cm := &CirconusMetrics{}
	cm.counters = make(map[string]uint64)
	cm.Add("foo", 5)

	val, ok := cm.counters["foo"]
	if !ok {
		t.Errorf("Expected to find foo")
	}

	if val != 5 {
		t.Errorf("Expected 1, found %d", val)
	}
}

func TestRemoveCounter(t *testing.T) {
	t.Log("Testing counter.RemoveCounter")

	cm := &CirconusMetrics{}
	cm.counters = make(map[string]uint64)
	cm.Increment("foo")

	val, ok := cm.counters["foo"]
	if !ok {
		t.Errorf("Expected to find foo")
	}

	if val != 1 {
		t.Errorf("Expected 1, found %d", val)
	}

	cm.RemoveCounter("foo")

	val, ok = cm.counters["foo"]
	if ok {
		t.Errorf("Expected NOT to find foo")
	}

	if val != 0 {
		t.Errorf("Expected 0, found %d", val)
	}
}

func TestSetCounterFunc(t *testing.T) {
	t.Log("Testing counter.SetCounterFunc")

	cf := func() uint64 {
		return 1
	}
	cm := &CirconusMetrics{}
	cm.counterFuncs = make(map[string]func() uint64)
	cm.SetCounterFunc("foo", cf)

	val, ok := cm.counterFuncs["foo"]
	if !ok {
		t.Errorf("Expected to find foo")
	}

	if val() != 1 {
		t.Errorf("Expected 1, found %d", val())
	}
}

func TestRemoveCounterFunc(t *testing.T) {
	t.Log("Testing counter.RemoveCounterFunc")

	cf := func() uint64 {
		return 1
	}
	cm := &CirconusMetrics{}
	cm.counterFuncs = make(map[string]func() uint64)
	cm.SetCounterFunc("foo", cf)

	val, ok := cm.counterFuncs["foo"]
	if !ok {
		t.Errorf("Expected to find foo")
	}

	if val() != 1 {
		t.Errorf("Expected 1, found %d", val())
	}

	cm.RemoveCounterFunc("foo")

	val, ok = cm.counterFuncs["foo"]
	if ok {
		t.Errorf("Expected NOT to find foo")
	}

	if val != nil {
		t.Errorf("Expected nil, found %v", val())
	}

}
