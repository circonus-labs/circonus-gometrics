package circonusgometrics

import (
	"testing"
)

func TestSetText(t *testing.T) {
	t.Log("Testing gauge.SetText")

	cm := &CirconusMetrics{}
	cm.text = make(map[string]string)
	cm.SetText("foo", "bar")

	val, ok := cm.text["foo"]
	if !ok {
		t.Errorf("Expected to find foo")
	}

	if val != "bar" {
		t.Errorf("Expected 'bar', found '%d'", val)
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
		t.Errorf("Expected 'bar', found '%d'", val)
	}
}

func TestRemoveText(t *testing.T) {
	t.Log("Testing text.RemoveText")

	cm := &CirconusMetrics{}
	cm.text = make(map[string]string)
	cm.SetText("foo", "bar")

	val, ok := cm.text["foo"]
	if !ok {
		t.Errorf("Expected to find foo")
	}

	if val != "bar" {
		t.Errorf("Expected 'bar', found '%d'", val)
	}

	cm.RemoveText("foo")

	val, ok = cm.text["foo"]
	if ok {
		t.Errorf("Expected NOT to find foo")
	}

	if val != "" {
		t.Errorf("Expected '', found '%d'", val)
	}
}
