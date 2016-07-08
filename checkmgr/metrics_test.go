package checkmgr

import (
	"testing"
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
