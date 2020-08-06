// Copyright 2016 Circonus, Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package circonusgometrics

import (
	"testing"
)

func TestSet(t *testing.T) {
	t.Log("Testing counter.Set")

	cm := &CirconusMetrics{counters: make(map[string]uint64)}

	cm.Set("foo", 30)

	val, ok := cm.counters["foo"]
	if !ok {
		t.Errorf("Expected to find foo")
	}

	if val != 30 {
		t.Errorf("Expected 30, found %d", val)
	}

	cm.Set("foo", 10)

	val, ok = cm.counters["foo"]
	if !ok {
		t.Errorf("Expected to find foo")
	}

	if val != 10 {
		t.Errorf("Expected 10, found %d", val)
	}
}

func TestSetWithTags(t *testing.T) {
	t.Log("Testing counter.SetWithTags")

	cm := &CirconusMetrics{counters: make(map[string]uint64)}

	metricName := "foo"
	tags := Tags{{"foo", "bar"}, {"baz", "qux"}}
	streamTagMetricName := cm.MetricNameWithStreamTags("foo", tags)

	cm.SetWithTags(metricName, tags, 30)

	val, ok := cm.counters[streamTagMetricName]
	if !ok {
		t.Fatalf("%s with %v tags not found (%s) (%#v)", metricName, tags, streamTagMetricName, cm.counters)
	}

	if val != 30 {
		t.Fatalf("expected 30 got (%d)", val)
	}

	cm.SetWithTags(metricName, tags, 10)

	val, ok = cm.counters[streamTagMetricName]
	if !ok {
		t.Fatalf("%s with %v tags not found (%s) (%#v)", metricName, tags, streamTagMetricName, cm.counters)
	}

	if val != 10 {
		t.Fatalf("expected 10 got (%d)", val)
	}
}

func TestIncrement(t *testing.T) {
	t.Log("Testing counter.Increment")

	cm := &CirconusMetrics{counters: make(map[string]uint64)}

	cm.Increment("foo")

	val, ok := cm.counters["foo"]
	if !ok {
		t.Errorf("Expected to find foo")
	}

	if val != 1 {
		t.Errorf("Expected 1, found %d", val)
	}
}

func TestIncrementWithTags(t *testing.T) {
	t.Log("Testing counter.IncrementWithTags")

	cm := &CirconusMetrics{counters: make(map[string]uint64)}

	metricName := "foo"
	tags := Tags{{"foo", "bar"}, {"baz", "qux"}}
	streamTagMetricName := cm.MetricNameWithStreamTags("foo", tags)

	cm.IncrementWithTags(metricName, tags)

	val, ok := cm.counters[streamTagMetricName]
	if !ok {
		t.Fatalf("%s with %v tags not found (%s) (%#v)", metricName, tags, streamTagMetricName, cm.counters)
	}

	if val != 1 {
		t.Fatalf("expected 1 got (%d)", val)
	}
}

func TestIncrementByValue(t *testing.T) {
	t.Log("Testing counter.IncrementByValue")

	cm := &CirconusMetrics{counters: make(map[string]uint64)}

	cm.IncrementByValue("foo", 10)

	val, ok := cm.counters["foo"]
	if !ok {
		t.Errorf("Expected to find foo")
	}

	if val != 10 {
		t.Errorf("Expected 1, found %d", val)
	}
}

func TestIncrementByValueWithTags(t *testing.T) {
	t.Log("Testing counter.IncrementByValueWithTags")

	cm := &CirconusMetrics{counters: make(map[string]uint64)}

	metricName := "foo"
	tags := Tags{{"foo", "bar"}, {"baz", "qux"}}
	streamTagMetricName := cm.MetricNameWithStreamTags("foo", tags)

	cm.IncrementByValueWithTags(metricName, tags, 10)

	val, ok := cm.counters[streamTagMetricName]
	if !ok {
		t.Fatalf("%s with %v tags not found (%s) (%#v)", metricName, tags, streamTagMetricName, cm.counters)
	}

	if val != 10 {
		t.Fatalf("expected 10 got (%d)", val)
	}
}

func TestAdd(t *testing.T) {
	t.Log("Testing counter.Add")

	cm := &CirconusMetrics{counters: make(map[string]uint64)}

	cm.Set("foo", 2)
	cm.Add("foo", 3)

	val, ok := cm.counters["foo"]
	if !ok {
		t.Fatal("Expected to find foo")
	}

	if val != 5 {
		t.Fatalf("Expected 1, found %d", val)
	}
}

func TestAddWithTags(t *testing.T) {
	t.Log("Testing counter.AddWithTags")

	cm := &CirconusMetrics{counters: make(map[string]uint64)}

	metricName := "foo"
	tags := Tags{{"foo", "bar"}, {"baz", "qux"}}
	streamTagMetricName := cm.MetricNameWithStreamTags("foo", tags)

	cm.SetWithTags(metricName, tags, 30)

	val, ok := cm.counters[streamTagMetricName]
	if !ok {
		t.Fatalf("%s with %v tags not found (%s) (%#v)", metricName, tags, streamTagMetricName, cm.counters)
	}

	if val != 30 {
		t.Fatalf("expected 30, got %d", val)
	}

	cm.AddWithTags(metricName, tags, 1)

	val, ok = cm.counters[streamTagMetricName]
	if !ok {
		t.Fatalf("%s with %v tags not found (%s) (%#v)", metricName, tags, streamTagMetricName, cm.counters)
	}

	if val != 31 {
		t.Fatalf("expected 31, got %d", val)
	}
}

func TestRemoveCounter(t *testing.T) {
	t.Log("Testing counter.RemoveCounter")

	cm := &CirconusMetrics{counters: make(map[string]uint64)}

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

func TestRemoveCounterWithTags(t *testing.T) {
	t.Log("Testing counter.RemoveCounterWithTags")

	cm := &CirconusMetrics{counters: make(map[string]uint64)}

	metricName := "foo"
	tags := Tags{{"foo", "bar"}, {"baz", "qux"}}
	streamTagMetricName := cm.MetricNameWithStreamTags("foo", tags)

	cm.IncrementWithTags(metricName, tags)

	val, ok := cm.counters[streamTagMetricName]
	if !ok {
		t.Fatalf("%s with %v tags not found (%s) (%#v)", metricName, tags, streamTagMetricName, cm.counters)
	}

	if val != 1 {
		t.Fatalf("expected 1 got (%d)", val)
	}

	cm.RemoveCounterWithTags(metricName, tags)

	val, ok = cm.counters[streamTagMetricName]
	if ok {
		t.Fatalf("expected NOT to find %s", streamTagMetricName)
	}

	if val != 0 {
		t.Fatalf("expected 0 got (%d)", val)
	}
}

func TestSetCounterFunc(t *testing.T) {
	t.Log("Testing counter.SetCounterFunc")

	cf := func() uint64 {
		return 1
	}

	cm := &CirconusMetrics{counterFuncs: make(map[string]func() uint64)}

	cm.SetCounterFunc("foo", cf)

	val, ok := cm.counterFuncs["foo"]
	if !ok {
		t.Errorf("Expected to find foo")
	}

	if val() != 1 {
		t.Errorf("Expected 1, found %d", val())
	}
}

func TestSetCounterFuncWithTags(t *testing.T) {
	t.Log("Testing counter.SetCounterFuncWithTags")

	cf := func() uint64 {
		return 1
	}

	cm := &CirconusMetrics{counterFuncs: make(map[string]func() uint64)}

	metricName := "foo"
	tags := Tags{{"foo", "bar"}, {"baz", "qux"}}
	streamTagMetricName := cm.MetricNameWithStreamTags("foo", tags)

	cm.SetCounterFuncWithTags(metricName, tags, cf)

	val, ok := cm.counterFuncs[streamTagMetricName]
	if !ok {
		t.Fatalf("%s with %v tags not found (%s) (%#v)", metricName, tags, streamTagMetricName, cm.counterFuncs)
	}

	if val() != 1 {
		t.Fatalf("expected 1 got (%d)", val())
	}
}

func TestRemoveCounterFunc(t *testing.T) {
	t.Log("Testing counter.RemoveCounterFunc")

	cf := func() uint64 {
		return 1
	}

	cm := &CirconusMetrics{counterFuncs: make(map[string]func() uint64)}

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

func TestRemoveCounterFuncWithTags(t *testing.T) {
	t.Log("Testing counter.RemoveCounterFuncWithTags")

	cf := func() uint64 {
		return 1
	}

	cm := &CirconusMetrics{counterFuncs: make(map[string]func() uint64)}

	metricName := "foo"
	tags := Tags{{"foo", "bar"}, {"baz", "qux"}}
	streamTagMetricName := cm.MetricNameWithStreamTags("foo", tags)

	cm.SetCounterFuncWithTags(metricName, tags, cf)

	val, ok := cm.counterFuncs[streamTagMetricName]
	if !ok {
		t.Fatalf("%s with %v tags not found (%s) (%#v)", metricName, tags, streamTagMetricName, cm.counterFuncs)
	}

	if val() != 1 {
		t.Fatalf("expected 1 got (%d)", val())
	}

	cm.RemoveCounterFuncWithTags(metricName, tags)

	val, ok = cm.counterFuncs[streamTagMetricName]
	if ok {
		t.Fatalf("expected NOT to find (%s)", streamTagMetricName)
	}

	if val != nil {
		t.Fatalf("expected nil got (%v)", val())
	}

}

func TestGetCounterTest(t *testing.T) {
	t.Log("Testing counter.GetCounterTest")

	cm := &CirconusMetrics{counters: make(map[string]uint64)}

	cm.Set("foo", 10)

	val, err := cm.GetCounterTest("foo")
	if err != nil {
		t.Errorf("Expected no error %v", err)
	}
	if val != 10 {
		t.Errorf("Expected 10 got %v", val)
	}

	_, err = cm.GetCounterTest("bar")
	if err == nil {
		t.Error("Expected error")
	}

}
