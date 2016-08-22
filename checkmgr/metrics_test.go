// Copyright 2016 Circonus, Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package checkmgr

import (
	"testing"

	"github.com/circonus-labs/circonus-gometrics/api"
)

func TestIsMetricActive(t *testing.T) {
	t.Log("Testing correct return from IsMetricActive")

	cm := &CheckManager{}

	cm.availableMetrics = map[string]bool{
		"foo": true,
	}

	t.Log("Testing for 'foo', foo in list")
	if !cm.IsMetricActive("foo") {
		t.Error("Expected true")
	}

	t.Log("Testing for 'bar', bar not in list")
	if cm.IsMetricActive("bar") {
		t.Error("Expected false")
	}
}

func TestInventoryMetrics(t *testing.T) {
	t.Log("Testing correct return from InventoryMetrics")

	cm := &CheckManager{}
	cm.checkBundle = &api.CheckBundle{}
	cm.checkBundle.Metrics = []api.CheckBundleMetric{
		api.CheckBundleMetric{
			Name:   "foo",
			Type:   "numeric",
			Status: "active",
		},
	}

	cm.availableMetrics = make(map[string]bool)

	t.Log("Testing for 'foo', foo not in list")
	if cm.IsMetricActive("foo") {
		t.Error("Expected false")
	}

	t.Log("Inventory metrics in check bundle")
	cm.inventoryMetrics()

	t.Log("Testing for 'foo', foo in list")
	if !cm.IsMetricActive("foo") {
		t.Error("Expected true")
	}

	t.Log("Testing for 'bar', bar not in list")
	if cm.IsMetricActive("bar") {
		t.Error("Expected false")
	}
}

func TestActivateMetric(t *testing.T) {
	t.Log("Testing correct return from ActivateMetric")

	cm := &CheckManager{}
	cm.checkBundle = &api.CheckBundle{}
	cm.checkBundle.Metrics = []api.CheckBundleMetric{
		api.CheckBundleMetric{
			Name:   "foo",
			Type:   "numeric",
			Status: "active",
		},
	}

	cm.availableMetrics = make(map[string]bool)
	cm.forceMetricActivation = false

	t.Log("Testing for 'foo', foo not in list")
	if !cm.ActivateMetric("foo") {
		t.Error("Expected true")
	}

	t.Log("Inventory metrics in check bundle")
	cm.inventoryMetrics()

	t.Log("Testing for 'foo', foo in list")
	if cm.ActivateMetric("foo") {
		t.Error("Expected false")
	}

	cm.checkBundle.Metrics = []api.CheckBundleMetric{
		api.CheckBundleMetric{
			Name:   "bar",
			Type:   "numeric",
			Status: "available",
		},
	}

	t.Log("Testing for 'bar', bar not in list")
	if !cm.ActivateMetric("bar") {
		t.Error("Expected true")
	}

	t.Log("Inventory metrics in check bundle")
	cm.inventoryMetrics()

	t.Log("Testing for 'bar', bar in list[false]")
	if cm.ActivateMetric("bar") {
		t.Error("Expected false")
	}

	t.Log("Change forceMetricActivation to true")
	cm.forceMetricActivation = true

	t.Log("Testing for 'bar', bar in list[false]")
	if !cm.ActivateMetric("bar") {
		t.Error("Expected true")
	}
}
