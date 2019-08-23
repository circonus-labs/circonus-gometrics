// Copyright 2016 Circonus, Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package circonusgometrics

import (
	"strings"
	"testing"
)

func TestFlush(t *testing.T) {
	server := testServer()
	defer server.Close()

	submissionURL := server.URL + "/metrics_endpoint"

	t.Log("Already flushing")
	{
		cfg := &Config{}
		cfg.CheckManager.Check.SubmissionURL = submissionURL
		cm, err := NewCirconusMetrics(cfg)
		if err != nil {
			t.Errorf("Expected no error, got '%v'", err)
		}

		cm.flushing = true
		cm.Flush()
	}

	t.Log("No metrics")
	{
		cfg := &Config{}
		cfg.CheckManager.Check.SubmissionURL = submissionURL
		cm, err := NewCirconusMetrics(cfg)
		if err != nil {
			t.Errorf("Expected no error, got '%v'", err)
		}

		cm.Flush()
	}

	t.Log("counter")
	{
		cfg := &Config{}
		cfg.CheckManager.Check.SubmissionURL = submissionURL
		cm, err := NewCirconusMetrics(cfg)
		if err != nil {
			t.Errorf("Expected no error, got '%v'", err)
		}

		cm.Set("foo", 30)

		cm.Flush()
	}

	t.Log("gauge")
	{
		cfg := &Config{}
		cfg.CheckManager.Check.SubmissionURL = submissionURL
		cm, err := NewCirconusMetrics(cfg)
		if err != nil {
			t.Errorf("Expected no error, got '%v'", err)
		}

		cm.SetGauge("foo", 30)

		cm.Flush()
	}

	t.Log("histogram")
	{
		cfg := &Config{}
		cfg.CheckManager.Check.SubmissionURL = submissionURL
		cm, err := NewCirconusMetrics(cfg)
		if err != nil {
			t.Errorf("Expected no error, got '%v'", err)
		}

		cm.Timing("foo", 30.28)

		cm.Flush()
	}

	t.Log("text")
	{
		cfg := &Config{}
		cfg.CheckManager.Check.SubmissionURL = submissionURL
		cm, err := NewCirconusMetrics(cfg)
		if err != nil {
			t.Errorf("Expected no error, got '%v'", err)
		}

		cm.SetText("foo", "bar")

		cm.Flush()
	}
}

func TestFlushMetrics(t *testing.T) {
	cfg := &Config{}
	cfg.CheckManager.Check.SubmissionURL = "none"
	cfg.Interval = "0"

	t.Log("Already flushing")
	{
		cm, err := NewCirconusMetrics(cfg)
		if err != nil {
			t.Errorf("Expected no error, got '%v'", err)
		}

		cm.flushing = true
		metrics := cm.FlushMetrics()
		if len(*metrics) != 0 {
			t.Fatal("expected 0 metrics")
		}
	}

	t.Log("No metrics")
	{
		cm, err := NewCirconusMetrics(cfg)
		if err != nil {
			t.Errorf("Expected no error, got '%v'", err)
		}

		metrics := cm.FlushMetrics()
		if len(*metrics) != 0 {
			t.Fatal("expected 0 metrics")
		}
	}

	t.Log("counter")
	{
		cm, err := NewCirconusMetrics(cfg)
		if err != nil {
			t.Errorf("Expected no error, got '%v'", err)
		}

		cm.Set("foo", 30)

		metrics := cm.FlushMetrics()
		if len(*metrics) == 0 {
			t.Fatal("expected 1 metric")
		}

		m, mok := (*metrics)["foo"]
		if !mok {
			t.Fatalf("'foo' not found in %v", metrics)
		}
		if m.Type != "L" {
			t.Fatalf("'Type' not correct %v", m)
		}
		if m.Value.(uint64) != 30 {
			t.Fatalf("'Value' not correct %v", m)
		}

	}

	t.Log("gauge")
	{
		cm, err := NewCirconusMetrics(cfg)
		if err != nil {
			t.Errorf("Expected no error, got '%v'", err)
		}

		v := int64(30)
		cm.SetGauge("foo", v)

		metrics := cm.FlushMetrics()
		if len(*metrics) == 0 {
			t.Fatal("expected 1 metric")
		}

		m, mok := (*metrics)["foo"]
		if !mok {
			t.Fatalf("'foo' not found in %v", metrics)
		}
		if m.Type != "l" {
			t.Fatalf("'Type' not correct %v", m)
		}
		if m.Value.(int64) != v {
			t.Fatalf("'Value' not correct, expected %v got %v", v, m.Value)
		}
	}

	t.Log("histogram")
	{
		cm, err := NewCirconusMetrics(cfg)
		if err != nil {
			t.Errorf("Expected no error, got '%v'", err)
		}

		cm.Timing("foo", 30.28)

		metrics := cm.FlushMetrics()
		if len(*metrics) == 0 {
			t.Fatal("expected 1 metric")
		}

		m, mok := (*metrics)["foo"]
		if !mok {
			t.Fatalf("'foo' not found in %v", metrics)
		}
		if m.Type != "h" {
			t.Fatalf("'Type' not correct %v", m)
		}
		if len(m.Value.([]string)) != 1 {
			t.Fatal("expected 1 value")
		}
		if m.Value.([]string)[0] != "H[3.0e+01]=1" {
			t.Fatalf("'Value' not correct %v", m)
		}
	}

	t.Log("text")
	{
		cm, err := NewCirconusMetrics(cfg)
		if err != nil {
			t.Errorf("Expected no error, got '%v'", err)
		}

		cm.SetText("foo", "bar")

		metrics := cm.FlushMetrics()
		if len(*metrics) == 0 {
			t.Fatal("expected 1 metric")
		}

		m, mok := (*metrics)["foo"]
		if !mok {
			t.Fatalf("'foo' not found in %v", metrics)
		}
		if m.Type != "s" {
			t.Fatalf("'Type' not correct %v", m)
		}
		if m.Value != "bar" {
			t.Fatalf("'Value' not correct %v", m)
		}
	}
}

func TestPromOutput(t *testing.T) {
	cfg := &Config{}
	cfg.CheckManager.Check.SubmissionURL = "none"
	cfg.Interval = "0"

	t.Log("No metrics")
	{
		cm, err := NewCirconusMetrics(cfg)
		if err != nil {
			t.Errorf("Expected no error, got '%v'", err)
		}

		b, err := cm.PromOutput()
		if err == nil {
			t.Fatal("expected error")
		}
		if b != nil {
			t.Fatalf("expected nil, got (%v)", b.String())
		}
	}

	t.Log("counter")
	{
		cm, err := NewCirconusMetrics(cfg)
		if err != nil {
			t.Errorf("Expected no error, got '%v'", err)
		}

		cm.Set("foo", 30)

		metrics := cm.FlushMetrics()
		if len(*metrics) == 0 {
			t.Fatal("expected 1 metric")
		}

		m, mok := (*metrics)["foo"]
		if !mok {
			t.Fatalf("'foo' not found in %v", metrics)
		}
		if m.Type != "L" {
			t.Fatalf("'Type' not correct %v", m)
		}
		if m.Value.(uint64) != 30 {
			t.Fatalf("'Value' not correct %v", m)
		}

		b, err := cm.PromOutput()
		if err != nil {
			t.Fatalf("expected no error, got (%s)", err)
		}
		if b == nil {
			t.Fatal("expected not nil")
		}
		expect := "foo 30"
		if !strings.HasPrefix(b.String(), expect) {
			t.Fatalf("expected prefix (%s) got (%s)", expect, b.String())
		}
	}

	t.Log("gauge")
	{
		cm, err := NewCirconusMetrics(cfg)
		if err != nil {
			t.Errorf("Expected no error, got '%v'", err)
		}

		v := int(30)
		cm.SetGauge("foo", v)

		metrics := cm.FlushMetrics()
		if len(*metrics) == 0 {
			t.Fatal("expected 1 metric")
		}

		m, mok := (*metrics)["foo"]
		if !mok {
			t.Fatalf("'foo' not found in %v", metrics)
		}
		if m.Type != "i" {
			t.Fatalf("'Type' not correct %v", m)
		}
		if m.Value.(int) != v {
			t.Fatalf("'Value' not correct, expected %v got %v", v, m.Value)
		}

		b, err := cm.PromOutput()
		if err != nil {
			t.Fatalf("expected no error, got (%s)", err)
		}
		if b == nil {
			t.Fatal("expected not nil")
		}
		expect := "foo 30"
		if !strings.HasPrefix(b.String(), expect) {
			t.Fatalf("expected prefix (%s) got (%s)", expect, b.String())
		}
	}

	t.Log("histogram")
	{
		cm, err := NewCirconusMetrics(cfg)
		if err != nil {
			t.Errorf("Expected no error, got '%v'", err)
		}

		cm.Timing("foo", 30.28)

		metrics := cm.FlushMetrics()
		if len(*metrics) == 0 {
			t.Fatal("expected 1 metric")
		}

		m, mok := (*metrics)["foo"]
		if !mok {
			t.Fatalf("'foo' not found in %v", metrics)
		}
		if m.Type != "h" {
			t.Fatalf("'Type' not correct %v", m)
		}
		if len(m.Value.([]string)) != 1 {
			t.Fatal("expected 1 value")
		}
		if m.Value.([]string)[0] != "H[3.0e+01]=1" {
			t.Fatalf("'Value' not correct %v", m)
		}

		b, err := cm.PromOutput()
		if err != nil {
			t.Fatalf("expected no error, got (%s)", err)
		}
		if b == nil {
			t.Fatal("expected not nil")
		}
		expect := ""
		if b.String() != expect {
			t.Fatalf("expected prefix (%s) got (%s)", expect, b.String())
		}
	}

	t.Log("text")
	{
		cm, err := NewCirconusMetrics(cfg)
		if err != nil {
			t.Errorf("Expected no error, got '%v'", err)
		}

		cm.SetText("foo", "bar")

		metrics := cm.FlushMetrics()
		if len(*metrics) == 0 {
			t.Fatal("expected 1 metric")
		}

		m, mok := (*metrics)["foo"]
		if !mok {
			t.Fatalf("'foo' not found in %v", metrics)
		}
		if m.Type != "s" {
			t.Fatalf("'Type' not correct %v", m)
		}
		if m.Value != "bar" {
			t.Fatalf("'Value' not correct %v", m)
		}

		b, err := cm.PromOutput()
		if err != nil {
			t.Fatalf("expected no error, got (%s)", err)
		}
		if b == nil {
			t.Fatal("expected not nil")
		}
		expect := ""
		if b.String() != expect {
			t.Fatalf("expected prefix (%s) got (%s)", expect, b.String())
		}
	}
}

func TestPackageMetrics(t *testing.T) {
	t.Log("Testing packageMetrics")

	cfg := &Config{}
	cfg.CheckManager.Check.SubmissionURL = "none"
	cfg.Interval = "0"

	cm, err := New(cfg)
	if err != nil {
		t.Errorf("Expected no error, got '%v'", err)
	}

	cm.flushing = false
	newMetrics, output := cm.packageMetrics()
	if len(newMetrics) != 0 && len(output) != 0 {
		t.Fatal("expected 0 metrics")
	}
}

func TestReset(t *testing.T) {
	t.Log("Testing util.Reset")

	cm := &CirconusMetrics{}

	cm.counters = make(map[string]uint64)
	cm.counterFuncs = make(map[string]func() uint64)
	cm.Increment("foo")

	// cm.gauges = make(map[string]string)
	cm.gauges = make(map[string]interface{})
	cm.gaugeFuncs = make(map[string]func() int64)
	cm.Gauge("foo", 1)

	cm.histograms = make(map[string]*Histogram)
	cm.Timing("foo", 1)

	cm.text = make(map[string]string)
	cm.textFuncs = make(map[string]func() string)
	cm.SetText("foo", "bar")

	if len(cm.counters) != 1 {
		t.Errorf("Expected 1, found %d", len(cm.counters))
	}

	if len(cm.gauges) != 1 {
		t.Errorf("Expected 1, found %d", len(cm.gauges))
	}

	if len(cm.histograms) != 1 {
		t.Errorf("Expected 1, found %d", len(cm.histograms))
	}

	if len(cm.text) != 1 {
		t.Errorf("Expected 1, found %d", len(cm.text))
	}

	cm.Reset()

	if len(cm.counters) != 0 {
		t.Errorf("Expected 0, found %d", len(cm.counters))
	}

	if len(cm.gauges) != 0 {
		t.Errorf("Expected 0, found %d", len(cm.gauges))
	}

	if len(cm.histograms) != 0 {
		t.Errorf("Expected 0, found %d", len(cm.histograms))
	}

	if len(cm.text) != 0 {
		t.Errorf("Expected 0, found %d", len(cm.text))
	}
}

func TestSnapshot(t *testing.T) {
	t.Log("Testing util.snapshot")

	cm := &CirconusMetrics{}

	cm.resetCounters = true
	cm.counters = make(map[string]uint64)
	cm.counterFuncs = make(map[string]func() uint64)
	cm.Increment("foo")

	cm.resetGauges = true
	// cm.gauges = make(map[string]string)
	cm.gauges = make(map[string]interface{})
	cm.gaugeFuncs = make(map[string]func() int64)
	cm.Gauge("foo", 1)

	cm.resetHistograms = true
	cm.histograms = make(map[string]*Histogram)
	cm.Timing("foo", 1)

	cm.resetText = true
	cm.text = make(map[string]string)
	cm.textFuncs = make(map[string]func() string)
	cm.SetText("foo", "bar")

	if len(cm.counters) != 1 {
		t.Errorf("Expected 1, found %d", len(cm.counters))
	}

	if len(cm.gauges) != 1 {
		t.Errorf("Expected 1, found %d", len(cm.gauges))
	}

	if len(cm.histograms) != 1 {
		t.Errorf("Expected 1, found %d", len(cm.histograms))
	}

	if len(cm.text) != 1 {
		t.Errorf("Expected 1, found %d", len(cm.text))
	}

	counters, gauges, histograms, text := cm.snapshot()

	if len(counters) != 1 {
		t.Errorf("Expected 1, found %d", len(counters))
	}

	if len(gauges) != 1 {
		t.Errorf("Expected 1, found %d", len(gauges))
	}

	if len(histograms) != 1 {
		t.Errorf("Expected 1, found %d", len(histograms))
	}

	if len(text) != 1 {
		t.Errorf("Expected 1, found %d", len(text))
	}
}

func TestSnapshotWithoutHistogramReset(t *testing.T) {
	t.Log("Testing util.snapshot with no histogram reset")

	cm := &CirconusMetrics{}

	cm.resetHistograms = false
	cm.histograms = make(map[string]*Histogram)
	cm.Timing("foo", 1)

	if len(cm.histograms) != 1 {
		t.Errorf("Expected 1, found %d", len(cm.histograms))
	}

	_, _, histograms, _ := cm.snapshot()

	if len(histograms) != 1 {
		t.Errorf("Expected 1, found %d", len(histograms))
	}

	if len(cm.histograms) != 1 {
		t.Errorf("Expected 1, found %d", len(cm.histograms))
	}

	if len(cm.histograms["foo"].hist.DecStrings()) != 1 {
		t.Errorf("Expected 1, found %d", len(cm.histograms))
	}
}
