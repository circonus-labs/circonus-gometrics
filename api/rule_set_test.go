// Copyright 2016 Circonus, Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

var (
	testRuleset = Ruleset{
		CID:      "/rule_set/1234_tt_firstbyte",
		CheckCID: "/check/1234",
		ContactGroups: map[uint8][]string{
			1: []string{"/contact_group/1234", "/contact_group/5678"},
			2: []string{"/contact_group/1234"},
			3: []string{"/contact_group/1234"},
			4: []string{},
			5: []string{},
		},
		Derive:     "",
		Link:       "http://example.com/how2fix/webserver_down/",
		MetricName: "tt_firstbyte",
		MetricType: "numeric",
		Notes:      "Determine if the HTTP request is taking too long to start (or is down.)  Don't fire if ping is already alerting",
		Parent:     "1233_ping",
		Rules: []RulesetRule{
			RulesetRule{
				Criteria:          "on absence",
				Severity:          1,
				Value:             "300",
				Wait:              5,
				WindowingDuration: 300,
				WindowingFunction: "",
			},
			RulesetRule{
				Criteria: "max value",
				Severity: 2,
				Value:    "1000",
				Wait:     5,
			},
		},
	}
)

func testRulesetServer() *httptest.Server {
	f := func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if path == "/rule_set/1234_tt_firstbyte" {
			switch r.Method {
			case "GET":
				ret, err := json.Marshal(testRuleset)
				if err != nil {
					panic(err)
				}
				w.WriteHeader(200)
				w.Header().Set("Content-Type", "application/json")
				fmt.Fprintln(w, string(ret))
			case "PUT":
				defer r.Body.Close()
				b, err := ioutil.ReadAll(r.Body)
				if err != nil {
					panic(err)
				}
				w.WriteHeader(200)
				w.Header().Set("Content-Type", "application/json")
				fmt.Fprintln(w, string(b))
			case "DELETE":
				w.WriteHeader(200)
				w.Header().Set("Content-Type", "application/json")
			default:
				w.WriteHeader(404)
				fmt.Fprintln(w, fmt.Sprintf("not found: %s %s", r.Method, path))
			}
		} else if path == "/rule_set" {
			switch r.Method {
			case "GET":
				reqURL := r.URL.String()
				var c []Ruleset
				if reqURL == "/rule_set?search=request%60latency_ms" {
					c = []Ruleset{testRuleset}
				} else if reqURL == "/rule_set?f_tags_has=service%3Aweb" {
					c = []Ruleset{testRuleset}
				} else if reqURL == "/rule_set?f_tags_has=service%3Aweb&search=request%60latency_ms" {
					c = []Ruleset{testRuleset}
				} else if reqURL == "/rule_set" {
					c = []Ruleset{testRuleset}
				} else {
					c = []Ruleset{}
				}
				if len(c) > 0 {
					ret, err := json.Marshal(c)
					if err != nil {
						panic(err)
					}
					w.WriteHeader(200)
					w.Header().Set("Content-Type", "application/json")
					fmt.Fprintln(w, string(ret))
				} else {
					w.WriteHeader(404)
					fmt.Fprintln(w, fmt.Sprintf("not found: %s %s", r.Method, reqURL))
				}
			case "POST":
				defer r.Body.Close()
				_, err := ioutil.ReadAll(r.Body)
				if err != nil {
					panic(err)
				}
				ret, err := json.Marshal(testRuleset)
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

func TestNewRuleset(t *testing.T) {
	bundle := NewRuleset()
	actualType := reflect.TypeOf(bundle)
	expectedType := "*api.Ruleset"
	if actualType.String() != expectedType {
		t.Fatalf("Expected %s, got %s", expectedType, actualType.String())
	}
}

func TestFetchRuleset(t *testing.T) {
	server := testRulesetServer()
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
		expectedError := errors.New("Invalid rule set CID [none]")
		_, err := apih.FetchRuleset(CIDType(&cid))
		if err == nil {
			t.Fatalf("Expected error")
		}
		if err.Error() != expectedError.Error() {
			t.Fatalf("Expected %+v got '%+v'", expectedError, err)
		}
	}

	t.Log("with valid CID")
	{
		cid := "/rule_set/1234_tt_firstbyte"
		ruleset, err := apih.FetchRuleset(CIDType(&cid))
		if err != nil {
			t.Fatalf("Expected no error, got '%v'", err)
		}

		actualType := reflect.TypeOf(ruleset)
		expectedType := "*api.Ruleset"
		if actualType.String() != expectedType {
			t.Fatalf("Expected %s, got %s", expectedType, actualType.String())
		}

		if ruleset.CID != testRuleset.CID {
			t.Fatalf("CIDs do not match: %+v != %+v\n", ruleset, testRuleset)
		}
	}

	t.Log("with invalid CID")
	{
		cid := "/invalid"
		expectedError := errors.New("Invalid rule set CID [/invalid]")
		_, err := apih.FetchRuleset(CIDType(&cid))
		if err == nil {
			t.Fatalf("Expected error")
		}
		if err.Error() != expectedError.Error() {
			t.Fatalf("Expected %+v got '%+v'", expectedError, err)
		}
	}
}

func TestFetchRulesets(t *testing.T) {
	server := testRulesetServer()
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

	rulesets, err := apih.FetchRulesets()
	if err != nil {
		t.Fatalf("Expected no error, got '%v'", err)
	}

	actualType := reflect.TypeOf(rulesets)
	expectedType := "*[]api.Ruleset"
	if actualType.String() != expectedType {
		t.Fatalf("Expected %s, got %s", expectedType, actualType.String())
	}

}

func TestCreateRuleset(t *testing.T) {
	server := testRulesetServer()
	defer server.Close()

	var apih *API

	ac := &Config{
		TokenKey: "abc123",
		TokenApp: "test",
		URL:      server.URL,
	}
	apih, err := NewAPI(ac)
	if err != nil {
		t.Errorf("Expected no error, got '%v'", err)
	}

	ruleset, err := apih.CreateRuleset(&testRuleset)
	if err != nil {
		t.Fatalf("Expected no error, got '%v'", err)
	}

	actualType := reflect.TypeOf(ruleset)
	expectedType := "*api.Ruleset"
	if actualType.String() != expectedType {
		t.Fatalf("Expected %s, got %s", expectedType, actualType.String())
	}
}

func TestUpdateRuleset(t *testing.T) {
	server := testRulesetServer()
	defer server.Close()

	var apih *API

	ac := &Config{
		TokenKey: "abc123",
		TokenApp: "test",
		URL:      server.URL,
	}
	apih, err := NewAPI(ac)
	if err != nil {
		t.Errorf("Expected no error, got '%v'", err)
	}

	t.Log("valid Ruleset")
	{
		ruleset, err := apih.UpdateRuleset(&testRuleset)
		if err != nil {
			t.Fatalf("Expected no error, got '%v'", err)
		}

		actualType := reflect.TypeOf(ruleset)
		expectedType := "*api.Ruleset"
		if actualType.String() != expectedType {
			t.Fatalf("Expected %s, got %s", expectedType, actualType.String())
		}
	}

	t.Log("Test with invalid CID")
	{
		expectedError := errors.New("Invalid rule set CID [/invalid]")
		x := &Ruleset{CID: "/invalid"}
		_, err := apih.UpdateRuleset(x)
		if err == nil {
			t.Fatal("Expected an error")
		}
		if err.Error() != expectedError.Error() {
			t.Fatalf("Expected %+v got '%+v'", expectedError, err)
		}
	}
}

func TestDeleteRuleset(t *testing.T) {
	server := testRulesetServer()
	defer server.Close()

	var apih *API

	ac := &Config{
		TokenKey: "abc123",
		TokenApp: "test",
		URL:      server.URL,
	}
	apih, err := NewAPI(ac)
	if err != nil {
		t.Errorf("Expected no error, got '%v'", err)
	}

	t.Log("valid Ruleset")
	{
		_, err := apih.DeleteRuleset(&testRuleset)
		if err != nil {
			t.Fatalf("Expected no error, got '%v'", err)
		}
	}

	t.Log("Test with invalid CID")
	{
		expectedError := errors.New("Invalid rule set CID [/invalid]")
		x := &Ruleset{CID: "/invalid"}
		_, err := apih.UpdateRuleset(x)
		if err == nil {
			t.Fatal("Expected an error")
		}
		if err.Error() != expectedError.Error() {
			t.Fatalf("Expected %+v got '%+v'", expectedError, err)
		}
	}
}

func TestSearchRulesets(t *testing.T) {
	server := testRulesetServer()
	defer server.Close()

	var apih *API

	ac := &Config{
		TokenKey: "abc123",
		TokenApp: "test",
		URL:      server.URL,
	}
	apih, err := NewAPI(ac)
	if err != nil {
		t.Errorf("Expected no error, got '%v'", err)
	}

	search := SearchQueryType("request`latency_ms")
	filter := SearchFilterType(map[string][]string{"f_tags_has": []string{"service:web"}})

	t.Log("no search, no filter")
	{
		rulesets, err := apih.SearchRulesets(nil, nil)
		if err != nil {
			t.Fatalf("Expected no error, got '%v'", err)
		}

		actualType := reflect.TypeOf(rulesets)
		expectedType := "*[]api.Ruleset"
		if actualType.String() != expectedType {
			t.Fatalf("Expected %s, got %s", expectedType, actualType.String())
		}
	}

	t.Log("search, no filter")
	{
		rulesets, err := apih.SearchRulesets(&search, nil)
		if err != nil {
			t.Fatalf("Expected no error, got '%v'", err)
		}

		actualType := reflect.TypeOf(rulesets)
		expectedType := "*[]api.Ruleset"
		if actualType.String() != expectedType {
			t.Fatalf("Expected %s, got %s", expectedType, actualType.String())
		}
	}

	t.Log("no search, filter")
	{
		rulesets, err := apih.SearchRulesets(nil, &filter)
		if err != nil {
			t.Fatalf("Expected no error, got '%v'", err)
		}

		actualType := reflect.TypeOf(rulesets)
		expectedType := "*[]api.Ruleset"
		if actualType.String() != expectedType {
			t.Fatalf("Expected %s, got %s", expectedType, actualType.String())
		}
	}

	t.Log("search, filter")
	{
		rulesets, err := apih.SearchRulesets(&search, &filter)
		if err != nil {
			t.Fatalf("Expected no error, got '%v'", err)
		}

		actualType := reflect.TypeOf(rulesets)
		expectedType := "*[]api.Ruleset"
		if actualType.String() != expectedType {
			t.Fatalf("Expected %s, got %s", expectedType, actualType.String())
		}
	}
}
