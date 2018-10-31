// Copyright 2016 Circonus, Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package circonusgometrics

import (
	"testing"
)

func TestSetText(t *testing.T) {
	t.Log("Testing gauge.SetText")

	cm := &CirconusMetrics{text: make(map[string]string)}

	cm.SetText("foo", "bar")

	val, ok := cm.text["foo"]
	if !ok {
		t.Errorf("Expected to find foo")
	}

	if val != "bar" {
		t.Errorf("Expected 'bar', found '%s'", val)
	}
}

func TestSetTextWithTags(t *testing.T) {
	t.Log("Testing gauge.SetTextWithTags")

	cm := &CirconusMetrics{text: make(map[string]string)}

	metricName := "foo"
	tags := Tags{{"foo", "bar"}, {"baz", "qux"}}
	streamTagMetricName := MetricNameWithStreamTags("foo", tags)

	cm.SetTextWithTags(metricName, tags, "bar")

	val, ok := cm.text[streamTagMetricName]
	if !ok {
		t.Fatalf("%s with %v tags not found (%s) (%#v)", metricName, tags, streamTagMetricName, cm.text)
	}

	if val != "bar" {
		t.Fatalf("Expected 'bar', found '%s'", val)
	}
}

func TestSetTextValue(t *testing.T) {
	t.Log("Testing gauge.SetTextValue")

	cm := &CirconusMetrics{}
	cm.text = make(map[string]string)
	cm.SetTextValue("foo", "bar")

	val, ok := cm.text["foo"]
	if !ok {
		t.Errorf("Expected to find foo")
	}

	if val != "bar" {
		t.Errorf("Expected 'bar', found '%s'", val)
	}
}

func TestSetTextValueWithTags(t *testing.T) {
	t.Log("Testing gauge.SetTextValueWithTags")

	cm := &CirconusMetrics{text: make(map[string]string)}

	metricName := "foo"
	tags := Tags{{"foo", "bar"}, {"baz", "qux"}}
	streamTagMetricName := MetricNameWithStreamTags("foo", tags)

	cm.SetTextValueWithTags(metricName, tags, "bar")

	val, ok := cm.text[streamTagMetricName]
	if !ok {
		t.Fatalf("%s with %v tags not found (%s) (%#v)", metricName, tags, streamTagMetricName, cm.text)
	}

	if val != "bar" {
		t.Fatalf("Expected 'bar', found '%s'", val)
	}
}

func TestRemoveText(t *testing.T) {
	t.Log("Testing text.RemoveText")

	cm := &CirconusMetrics{text: make(map[string]string)}

	cm.SetText("foo", "bar")

	val, ok := cm.text["foo"]
	if !ok {
		t.Errorf("Expected to find foo")
	}

	if val != "bar" {
		t.Errorf("Expected 'bar', found '%s'", val)
	}

	cm.RemoveText("foo")

	val, ok = cm.text["foo"]
	if ok {
		t.Errorf("Expected NOT to find foo")
	}

	if val != "" {
		t.Errorf("Expected '', found '%s'", val)
	}
}

func TestRemoveTextWithTags(t *testing.T) {
	t.Log("Testing gauge.RemoveTextWithTags")

	cm := &CirconusMetrics{text: make(map[string]string)}

	metricName := "foo"
	tags := Tags{{"foo", "bar"}, {"baz", "qux"}}
	streamTagMetricName := MetricNameWithStreamTags("foo", tags)

	cm.SetTextWithTags(metricName, tags, "bar")

	val, ok := cm.text[streamTagMetricName]
	if !ok {
		t.Fatalf("%s with %v tags not found (%s) (%#v)", metricName, tags, streamTagMetricName, cm.text)
	}

	if val != "bar" {
		t.Fatalf("Expected 'bar', found '%s'", val)
	}

	cm.RemoveTextWithTags(metricName, tags)

	val, ok = cm.text[streamTagMetricName]
	if ok {
		t.Fatalf("expected NOT to find %s", streamTagMetricName)
	}
	if val != "" {
		t.Fatalf("expected '' got (%s)", val)
	}
}

func TestSetTextFunc(t *testing.T) {
	t.Log("Testing text.SetTextFunc")

	tf := func() string {
		return "bar"
	}

	cm := &CirconusMetrics{textFuncs: make(map[string]func() string)}

	cm.SetTextFunc("foo", tf)

	val, ok := cm.textFuncs["foo"]
	if !ok {
		t.Errorf("Expected to find foo")
	}

	if val() != "bar" {
		t.Errorf("Expected 'bar', found '%s'", val())
	}
}

func TestSetTextFuncWithTags(t *testing.T) {
	t.Log("Testing text.SetTextFuncWithTags")

	tf := func() string {
		return "bar"
	}

	cm := &CirconusMetrics{textFuncs: make(map[string]func() string)}

	metricName := "foo"
	tags := Tags{{"foo", "bar"}, {"baz", "qux"}}
	streamTagMetricName := MetricNameWithStreamTags("foo", tags)

	cm.SetTextFuncWithTags(metricName, tags, tf)

	val, ok := cm.textFuncs[streamTagMetricName]
	if !ok {
		t.Fatalf("%s with %v tags not found (%s) (%#v)", metricName, tags, streamTagMetricName, cm.text)
	}

	if val() != "bar" {
		t.Fatalf("expected 'bar', got (%s)", val())
	}
}

func TestRemoveTextFunc(t *testing.T) {
	t.Log("Testing text.RemoveTextFunc")

	tf := func() string {
		return "bar"
	}

	cm := &CirconusMetrics{textFuncs: make(map[string]func() string)}

	cm.SetTextFunc("foo", tf)

	val, ok := cm.textFuncs["foo"]
	if !ok {
		t.Errorf("Expected to find foo")
	}

	if val() != "bar" {
		t.Errorf("Expected 'bar', found '%s'", val())
	}

	cm.RemoveTextFunc("foo")

	val, ok = cm.textFuncs["foo"]
	if ok {
		t.Errorf("Expected NOT to find foo")
	}

	if val != nil {
		t.Errorf("Expected nil, found %s", val())
	}

}

func TestRemoveTextFuncWithTags(t *testing.T) {
	t.Log("Testing text.RemoveTextFuncWithTags")

	tf := func() string {
		return "bar"
	}

	cm := &CirconusMetrics{textFuncs: make(map[string]func() string)}

	metricName := "foo"
	tags := Tags{{"foo", "bar"}, {"baz", "qux"}}
	streamTagMetricName := MetricNameWithStreamTags("foo", tags)

	cm.SetTextFuncWithTags(metricName, tags, tf)

	val, ok := cm.textFuncs[streamTagMetricName]
	if !ok {
		t.Fatalf("%s with %v tags not found (%s) (%#v)", metricName, tags, streamTagMetricName, cm.text)
	}

	if val() != "bar" {
		t.Fatalf("expected 'bar', got (%s)", val())
	}

	cm.RemoveTextFuncWithTags(metricName, tags)

	val, ok = cm.textFuncs[streamTagMetricName]
	if ok {
		t.Fatalf("expected NOT to find %s", streamTagMetricName)
	}
	if val != nil {
		t.Fatalf("expected nil got (%v)", val())
	}
}
