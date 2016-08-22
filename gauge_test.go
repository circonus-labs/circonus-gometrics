// Copyright 2016 Circonus, Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package circonusgometrics

import (
	"testing"
)

func TestGauge(t *testing.T) {
	t.Log("Testing gauge.Gauge")

	cm := &CirconusMetrics{}
	cm.gauges = make(map[string]string)

	// int
	cm.Gauge("foo", 1)
	val, ok := cm.gauges["foo"]
	if !ok {
		t.Errorf("Expected to find foo")
	}

	if val != "1" {
		t.Errorf("Expected 1, found %s", val)
	}

	// uint
	cm.Gauge("foo", uint64(10))
	val, ok = cm.gauges["foo"]
	if !ok {
		t.Errorf("Expected to find foo")
	}

	if val != "10" {
		t.Errorf("Expected 10, found %s", val)
	}

	// float
	cm.Gauge("foo", 3.12)
	val, ok = cm.gauges["foo"]
	if !ok {
		t.Errorf("Expected to find foo")
	}

	if val[0:4] != "3.12" {
		t.Errorf("Expected 3.12, found %s", val)
	}

}

func TestSetGauge(t *testing.T) {
	t.Log("Testing gauge.SetGauge")

	cm := &CirconusMetrics{}
	cm.gauges = make(map[string]string)
	cm.SetGauge("foo", 10)

	val, ok := cm.gauges["foo"]
	if !ok {
		t.Errorf("Expected to find foo")
	}

	if val != "10" {
		t.Errorf("Expected 10, found %s", val)
	}
}

func TestRemoveGauge(t *testing.T) {
	t.Log("Testing gauge.RemoveGauge")

	cm := &CirconusMetrics{}
	cm.gauges = make(map[string]string)
	cm.Gauge("foo", 5)

	val, ok := cm.gauges["foo"]
	if !ok {
		t.Errorf("Expected to find foo")
	}

	if val != "5" {
		t.Errorf("Expected 5, found %s", val)
	}

	cm.RemoveGauge("foo")

	val, ok = cm.gauges["foo"]
	if ok {
		t.Errorf("Expected NOT to find foo")
	}

	if val != "" {
		t.Errorf("Expected '', found '%s'", val)
	}
}

func TestSetGaugeFunc(t *testing.T) {
	t.Log("Testing gauge.SetGaugeFunc")

	gf := func() int64 {
		return 1
	}
	cm := &CirconusMetrics{}
	cm.gaugeFuncs = make(map[string]func() int64)
	cm.SetGaugeFunc("foo", gf)

	val, ok := cm.gaugeFuncs["foo"]
	if !ok {
		t.Errorf("Expected to find foo")
	}

	if val() != 1 {
		t.Errorf("Expected 1, found %d", val())
	}
}

func TestRemoveGaugeFunc(t *testing.T) {
	t.Log("Testing gauge.RemoveGaugeFunc")

	gf := func() int64 {
		return 1
	}
	cm := &CirconusMetrics{}
	cm.gaugeFuncs = make(map[string]func() int64)
	cm.SetGaugeFunc("foo", gf)

	val, ok := cm.gaugeFuncs["foo"]
	if !ok {
		t.Errorf("Expected to find foo")
	}

	if val() != 1 {
		t.Errorf("Expected 1, found %d", val())
	}

	cm.RemoveGaugeFunc("foo")

	val, ok = cm.gaugeFuncs["foo"]
	if ok {
		t.Errorf("Expected NOT to find foo")
	}

	if val != nil {
		t.Errorf("Expected nil, found %v", val())
	}

}
