// Copyright 2016 Circonus, Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

var (
	testAlert = Alert{
		CID:                "/alert/1234",
		AcknowledgementCID: "/acknowledgement/1234",
		AlertURL:           "https://example.circonus.com/fault-detection?alert_id=1234",
		BrokerCID:          "/broker/1234",
		CheckCID:           "/check/1234",
		CheckName:          "foo bar",
		ClearedOn:          &[]uint{1483033602}[0],
		ClearedValue:       &[]string{"1234"}[0],
		Maintenance:        []string{},
		MetricLinkURL:      &[]string{"http://example.com/docs/what_to_do_when/foo_bar_failure.html"}[0],
		MetricName:         "baz",
		MetricNotes:        &[]string{"blah blah blah"}[0],
		OccurredOn:         1483033102,
		RuleSetCID:         "/rule_set/1234_baz",
		Severity:           2,
		Tags:               []string{"cat:tag"},
		Value:              "5678",
	}
)

func testAlertServer() *httptest.Server {
	f := func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if path == "/alert/1234" {
			switch r.Method {
			case "GET":
				ret, err := json.Marshal(testAlert)
				if err != nil {
					panic(err)
				}
				w.WriteHeader(200)
				w.Header().Set("Content-Type", "application/json")
				fmt.Fprintln(w, string(ret))
			default:
				w.WriteHeader(404)
				fmt.Fprintln(w, fmt.Sprintf("not found: %s %s", r.Method, path))
			}
		} else if path == "/alert" {
			switch r.Method {
			case "GET":
				c := []Alert{testAlert}
				ret, err := json.Marshal(c)
				if err != nil {
					panic(err)
				}
				w.WriteHeader(200)
				w.Header().Set("Content-Type", "application/json")
				fmt.Fprintln(w, string(ret))
			default:
				w.WriteHeader(404)
				fmt.Fprintln(w, fmt.Sprintf("not found: %s %s", r.Method, path))
			}
		} else {
			w.WriteHeader(404)
			fmt.Fprintln(w, fmt.Sprintf("not found: %s %s", r.Method, path))
		}
	}

	return httptest.NewServer(http.HandlerFunc(f))
}

func TestFetchAlert(t *testing.T) {
	server := testAlertServer()
	defer server.Close()

	ac := &Config{
		TokenKey: "abc123",
		TokenApp: "test",
		URL:      server.URL,
	}
	apih, err := NewAPI(ac)
	if err != nil {
		t.Errorf("Expected no error, got '%v'", err)
	}

	t.Log("without CID")
	{
		cid := ""
		expectedError := errors.New("Invalid alert CID [none]")
		_, err := apih.FetchAlert(CIDType(&cid))
		if err == nil {
			t.Fatalf("Expected error")
		}
		if err.Error() != expectedError.Error() {
			t.Fatalf("Expected %+v got '%+v'", expectedError, err)
		}
	}

	t.Log("with valid CID")
	{
		cid := "/alert/1234"
		alert, err := apih.FetchAlert(CIDType(&cid))
		if err != nil {
			t.Fatalf("Expected no error, got '%v'", err)
		}

		actualType := reflect.TypeOf(alert)
		expectedType := "*api.Alert"
		if actualType.String() != expectedType {
			t.Fatalf("Expected %s, got %s", expectedType, actualType.String())
		}

		if alert.CID != testAlert.CID {
			t.Fatalf("CIDs do not match: %+v != %+v\n", alert, testAlert)
		}
	}

	t.Log("with invalid CID")
	{
		cid := "/invalid"
		expectedError := errors.New("Invalid alert CID [/invalid]")
		_, err := apih.FetchAlert(CIDType(&cid))
		if err == nil {
			t.Fatalf("Expected error")
		}
		if err.Error() != expectedError.Error() {
			t.Fatalf("Expected %+v got '%+v'", expectedError, err)
		}
	}
}

func TestFetchAlerts(t *testing.T) {
	server := testAlertServer()
	defer server.Close()

	ac := &Config{
		TokenKey: "abc123",
		TokenApp: "test",
		URL:      server.URL,
	}
	apih, err := NewAPI(ac)
	if err != nil {
		t.Errorf("Expected no error, got '%v'", err)
	}

	alerts, err := apih.FetchAlerts()
	if err != nil {
		t.Fatalf("Expected no error, got '%v'", err)
	}

	actualType := reflect.TypeOf(alerts)
	expectedType := "*[]api.Alert"
	if actualType.String() != expectedType {
		t.Fatalf("Expected %s, got %s", expectedType, actualType.String())
	}

}
