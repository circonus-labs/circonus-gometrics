// Copyright 2016 Circonus, Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package checkmgr

import (
	"testing"

	"github.com/circonus-labs/circonus-gometrics/api"
)

func TestIsMetricActive(t *testing.T) {

	cm := &CheckManager{}

	cm.availableMetrics = map[string]bool{
		"foo": true,
	}

	t.Log("'foo' in active metric list")
	{
		if !cm.IsMetricActive("foo") {
			t.Error("Expected true")
		}
	}

	t.Log("'bar' not in active metric list")
	{
		if cm.IsMetricActive("bar") {
			t.Error("Expected false")
		}
	}
}

func TestActivateMetric(t *testing.T) {
	cm := &CheckManager{}
	cm.checkBundle = &api.CheckBundle{}
	cm.checkBundle.Metrics = []api.CheckBundleMetric{
		api.CheckBundleMetric{
			Name:   "foo",
			Type:   "numeric",
			Status: "active",
		},
		api.CheckBundleMetric{
			Name:   "bar",
			Type:   "numeric",
			Status: "available",
		},
	}

	cm.availableMetrics = make(map[string]bool)
	cm.forceMetricActivation = false

	cm.inventoryMetrics()

	t.Log("'foo' already active")
	{
		if cm.ActivateMetric("foo") {
			t.Error("Expected false")
		}
	}

	t.Log("'bar' in list but not active [force=false]")
	{
		if cm.ActivateMetric("bar") {
			t.Error("Expected false")
		}
	}

	t.Log("'baz' not in list")
	{
		if !cm.ActivateMetric("baz") {
			t.Error("Expected true")
		}
	}

	cm.forceMetricActivation = true

	t.Log("'bar' in list but not active [force=true]")
	{
		if !cm.ActivateMetric("bar") {
			t.Error("Expected true")
		}
	}
}

func TestInventoryMetrics(t *testing.T) {
	cm := &CheckManager{}
	cm.checkBundle = &api.CheckBundle{}
	cm.checkBundle.Metrics = []api.CheckBundleMetric{
		api.CheckBundleMetric{
			Name:   "foo",
			Type:   "numeric",
			Status: "active",
		},
		api.CheckBundleMetric{
			Name:   "bar",
			Type:   "numeric",
			Status: "available",
		},
	}

	cm.availableMetrics = make(map[string]bool)
	cm.inventoryMetrics()

	expectedMetrics := make(map[string]bool)
	expectedMetrics["foo"] = true
	expectedMetrics["bar"] = false

	t.Log("'foo', in inventory and active")
	{
		active, exists := cm.availableMetrics["foo"]
		if !active {
			t.Fatalf("Expected active")
		}
		if !exists {
			t.Fatalf("Expected exists")
		}
	}

	t.Log("'bar', in inventory and not active")
	{
		active, exists := cm.availableMetrics["bar"]
		if active {
			t.Fatalf("Expected not active")
		}
		if !exists {
			t.Fatalf("Expected exists")
		}
	}

	t.Log("'baz', not in inventory and not active")
	{
		active, exists := cm.availableMetrics["baz"]
		if active {
			t.Fatalf("Expected not active")
		}
		if exists {
			t.Fatalf("Expected not exists")
		}
	}
}

func TestAddMetricTags(t *testing.T) {
	cm := &CheckManager{}
	cm.checkBundle = &api.CheckBundle{}
	cm.metricTags = make(map[string][]string)

	t.Log("no tags")
	{
		if cm.AddMetricTags("foo", []string{}, false) {
			t.Fatalf("Expected false")
		}
	}

	t.Log("no metric named 'foo'")
	{
		if !cm.AddMetricTags("foo", []string{"cat:tag"}, false) {
			t.Fatalf("Expected true")
		}
	}

	cm.checkBundle.Metrics = []api.CheckBundleMetric{
		api.CheckBundleMetric{
			Name:   "bar",
			Type:   "numeric",
			Status: "active",
		},
		api.CheckBundleMetric{
			Name:   "baz",
			Type:   "numeric",
			Status: "active",
			Tags:   []string{"cat1:tag1"},
		},
	}

	t.Log("metric named 'bar', add tag")
	{
		// append, zero current
		if !cm.AddMetricTags("bar", []string{"cat:tag"}, true) {
			t.Fatalf("Expected true")
		}
		// replace any existing
		if !cm.AddMetricTags("bar", []string{"cat:tag"}, false) {
			t.Fatalf("Expected true")
		}
	}

	t.Log("metric named 'baz', add tag")
	{
		// append, one current
		if !cm.AddMetricTags("baz", []string{"cat:tag"}, true) {
			t.Fatalf("Expected true")
		}
		// append, tag already exists
		if cm.AddMetricTags("baz", []string{"cat:tag"}, true) {
			t.Fatalf("Expected false")
		}
	}
}

func TestAddNewMetrics(t *testing.T) {
	cm := &CheckManager{}

	newMetrics := make(map[string]*api.CheckBundleMetric)

	newMetrics["foo"] = &api.CheckBundleMetric{
		Name:   "foo",
		Type:   "numeric",
		Status: "active",
	}

	t.Log("no check bundle")
	{
		if cm.addNewMetrics(newMetrics) {
			t.Fatalf("Expected false")
		}
	}

	cm.checkBundle = &api.CheckBundle{}
	t.Log("no check bundle metrics")
	{
		if !cm.addNewMetrics(newMetrics) {
			t.Fatalf("Expected true")
		}
		if !cm.forceCheckUpdate {
			t.Fatal("Expected forceCheckUpdate to be true")
		}
	}
}
