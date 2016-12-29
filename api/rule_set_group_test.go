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
	testRulesetGroup = RulesetGroup{
		CID: "/rule_set_group/1234",
		ContactGroups: map[int][]string{
			1: []string{"/contact_group/1234", "/contact_group/5678"},
			2: []string{"/contact_group/1234"},
			3: []string{"/contact_group/1234"},
			4: []string{},
			5: []string{},
		},
		Formulas: []RulesetGroupFormula{
			RulesetGroupFormula{
				Expression:    "(A and B) and not C",
				RaiseSeverity: 2,
				Wait:          0,
			},
			RulesetGroupFormula{
				Expression:    "3",
				RaiseSeverity: 1,
				Wait:          5,
			},
		},
		Name: "Multiple webservers gone bad",
		RulesetConditions: []RulesetGroupCondition{
			RulesetGroupCondition{
				MatchingSeverities: []string{"1", "2"},
				RulesetCID:         "/rule_set/1234_tt_firstbyte",
			},
			RulesetGroupCondition{
				MatchingSeverities: []string{"1", "2"},
				RulesetCID:         "/rule_set/5678_tt_firstbyte",
			},
			RulesetGroupCondition{
				MatchingSeverities: []string{"1", "2"},
				RulesetCID:         "/rule_set/9012_tt_firstbyte",
			},
		},
	}
)

func testRulesetGroupServer() *httptest.Server {
	f := func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if path == "/rule_set_group/1234" {
			switch r.Method {
			case "GET":
				ret, err := json.Marshal(testRulesetGroup)
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
		} else if path == "/rule_set_group" {
			switch r.Method {
			case "GET":
				c := []RulesetGroup{testRulesetGroup}
				ret, err := json.Marshal(c)
				if err != nil {
					panic(err)
				}
				w.WriteHeader(200)
				w.Header().Set("Content-Type", "application/json")
				fmt.Fprintln(w, string(ret))
			case "POST":
				defer r.Body.Close()
				_, err := ioutil.ReadAll(r.Body)
				if err != nil {
					panic(err)
				}
				ret, err := json.Marshal(testRulesetGroup)
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

func TestFetchRulesetGroup(t *testing.T) {
	server := testRulesetGroupServer()
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
		expectedError := errors.New("Invalid rule set group CID ")
		_, err := apih.FetchRulesetGroup(CIDType(""))
		if err == nil {
			t.Fatalf("Expected error")
		}
		if err.Error() != expectedError.Error() {
			t.Fatalf("Expected %+v got '%+v'", expectedError, err)
		}
	}

	t.Log("with valid CID")
	{
		cid := CIDType("/rule_set_group/1234")
		ruleset, err := apih.FetchRulesetGroup(cid)
		if err != nil {
			t.Fatalf("Expected no error, got '%v'", err)
		}

		actualType := reflect.TypeOf(ruleset)
		expectedType := "*api.RulesetGroup"
		if actualType.String() != expectedType {
			t.Fatalf("Expected %s, got %s", expectedType, actualType.String())
		}

		if ruleset.CID != testRulesetGroup.CID {
			t.Fatalf("CIDs do not match: %+v != %+v\n", ruleset, testRulesetGroup)
		}
	}

	t.Log("with invalid CID")
	{
		expectedError := errors.New("Invalid rule set group CID /invalid")
		_, err := apih.FetchRulesetGroup(CIDType("/invalid"))
		if err == nil {
			t.Fatalf("Expected error")
		}
		if err.Error() != expectedError.Error() {
			t.Fatalf("Expected %+v got '%+v'", expectedError, err)
		}
	}
}

func TestFetchRulesetGroups(t *testing.T) {
	server := testRulesetGroupServer()
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

	rulesets, err := apih.FetchRulesetGroups()
	if err != nil {
		t.Fatalf("Expected no error, got '%v'", err)
	}

	actualType := reflect.TypeOf(rulesets)
	expectedType := "[]api.RulesetGroup"
	if actualType.String() != expectedType {
		t.Fatalf("Expected %s, got %s", expectedType, actualType.String())
	}

}

func TestCreateRulesetGroup(t *testing.T) {
	server := testRulesetGroupServer()
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

	ruleset, err := apih.CreateRulesetGroup(&testRulesetGroup)
	if err != nil {
		t.Fatalf("Expected no error, got '%v'", err)
	}

	actualType := reflect.TypeOf(ruleset)
	expectedType := "*api.RulesetGroup"
	if actualType.String() != expectedType {
		t.Fatalf("Expected %s, got %s", expectedType, actualType.String())
	}
}

func TestUpdateRulesetGroup(t *testing.T) {
	server := testRulesetGroupServer()
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

	t.Log("valid RulesetGroup")
	{
		ruleset, err := apih.UpdateRulesetGroup(&testRulesetGroup)
		if err != nil {
			t.Fatalf("Expected no error, got '%v'", err)
		}

		actualType := reflect.TypeOf(ruleset)
		expectedType := "*api.RulesetGroup"
		if actualType.String() != expectedType {
			t.Fatalf("Expected %s, got %s", expectedType, actualType.String())
		}
	}

	t.Log("Test with invalid CID")
	{
		expectedError := errors.New("Invalid rule set group CID /invalid")
		x := &RulesetGroup{CID: "/invalid"}
		_, err := apih.UpdateRulesetGroup(x)
		if err == nil {
			t.Fatal("Expected an error")
		}
		if err.Error() != expectedError.Error() {
			t.Fatalf("Expected %+v got '%+v'", expectedError, err)
		}
	}
}

func TestDeleteRulesetGroup(t *testing.T) {
	server := testRulesetGroupServer()
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

	t.Log("valid RulesetGroup")
	{
		_, err := apih.DeleteRulesetGroup(&testRulesetGroup)
		if err != nil {
			t.Fatalf("Expected no error, got '%v'", err)
		}
	}

	t.Log("Test with invalid CID")
	{
		expectedError := errors.New("Invalid rule set group CID /invalid")
		x := &RulesetGroup{CID: "/invalid"}
		_, err := apih.UpdateRulesetGroup(x)
		if err == nil {
			t.Fatal("Expected an error")
		}
		if err.Error() != expectedError.Error() {
			t.Fatalf("Expected %+v got '%+v'", expectedError, err)
		}
	}
}
