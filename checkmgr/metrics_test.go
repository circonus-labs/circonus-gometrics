package checkmgr

import (
	"testing"

	"github.com/circonus-labs/circonus-gometrics/api"
)

func TestIsMetricActive(t *testing.T) {
	t.Log("Testing correct return from IsMetricActive")

	cm := &CheckManager{}

	cm.activeMetrics = map[string]bool{
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
	t.Log("Testing correct return from IsMetricActive")

	cm := &CheckManager{}
	cm.checkBundle = &api.CheckBundle{}
	cm.checkBundle.Metrics = []api.CheckBundleMetric{
		api.CheckBundleMetric{
			Name:   "foo",
			Type:   "text",
			Status: "active",
		},
	}

	cm.activeMetrics = make(map[string]bool)
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
