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
		ContactGroups: map[uint8][]string{
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
				reqURL := r.URL.String()
				var c []RulesetGroup
				if reqURL == "/rule_set_group?search=web+requests" {
					c = []RulesetGroup{testRulesetGroup}
				} else if reqURL == "/rule_set_group?f_tags_has=location%3Aconus" {
					c = []RulesetGroup{testRulesetGroup}
				} else if reqURL == "/rule_set_group?f_tags_has=location%3Aconus&search=web+requests" {
					c = []RulesetGroup{testRulesetGroup}
				} else if reqURL == "/rule_set_group" {
					c = []RulesetGroup{testRulesetGroup}
				} else {
					c = []RulesetGroup{}
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

func TestNewRulesetGroup(t *testing.T) {
	bundle := NewRulesetGroup()
	actualType := reflect.TypeOf(bundle)
	expectedType := "*api.RulesetGroup"
	if actualType.String() != expectedType {
		t.Fatalf("Expected %s, got %s", expectedType, actualType.String())
	}
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
		cid := ""
		expectedError := errors.New("Invalid rule set group CID [none]")
		_, err := apih.FetchRulesetGroup(CIDType(&cid))
		if err == nil {
			t.Fatalf("Expected error")
		}
		if err.Error() != expectedError.Error() {
			t.Fatalf("Expected %+v got '%+v'", expectedError, err)
		}
	}

	t.Log("with valid CID")
	{
		cid := "/rule_set_group/1234"
		ruleset, err := apih.FetchRulesetGroup(CIDType(&cid))
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
		cid := "/invalid"
		expectedError := errors.New("Invalid rule set group CID [/invalid]")
		_, err := apih.FetchRulesetGroup(CIDType(&cid))
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
	expectedType := "*[]api.RulesetGroup"
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
		expectedError := errors.New("Invalid rule set group CID [/invalid]")
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
		expectedError := errors.New("Invalid rule set group CID [/invalid]")
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

func TestSearchRulesetGroups(t *testing.T) {
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

	search := SearchQueryType("web requests")
	filter := SearchFilterType(map[string][]string{"f_tags_has": []string{"location:conus"}})

	t.Log("no search, no filter")
	{
		groups, err := apih.SearchRulesetGroups(nil, nil)
		if err != nil {
			t.Fatalf("Expected no error, got '%v'", err)
		}

		actualType := reflect.TypeOf(groups)
		expectedType := "*[]api.RulesetGroup"
		if actualType.String() != expectedType {
			t.Fatalf("Expected %s, got %s", expectedType, actualType.String())
		}
	}

	t.Log("search, no filter")
	{
		groups, err := apih.SearchRulesetGroups(&search, nil)
		if err != nil {
			t.Fatalf("Expected no error, got '%v'", err)
		}

		actualType := reflect.TypeOf(groups)
		expectedType := "*[]api.RulesetGroup"
		if actualType.String() != expectedType {
			t.Fatalf("Expected %s, got %s", expectedType, actualType.String())
		}
	}

	t.Log("no search, filter")
	{
		groups, err := apih.SearchRulesetGroups(nil, &filter)
		if err != nil {
			t.Fatalf("Expected no error, got '%v'", err)
		}

		actualType := reflect.TypeOf(groups)
		expectedType := "*[]api.RulesetGroup"
		if actualType.String() != expectedType {
			t.Fatalf("Expected %s, got %s", expectedType, actualType.String())
		}
	}

	t.Log("search, filter")
	{
		groups, err := apih.SearchRulesetGroups(&search, &filter)
		if err != nil {
			t.Fatalf("Expected no error, got '%v'", err)
		}

		actualType := reflect.TypeOf(groups)
		expectedType := "*[]api.RulesetGroup"
		if actualType.String() != expectedType {
			t.Fatalf("Expected %s, got %s", expectedType, actualType.String())
		}
	}
}
