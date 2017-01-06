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
	"os"
	"reflect"
	"testing"
)

var (
	testDashboard = Dashboard{}
)

func init() {
	file := "dashboard-example.json"
	fi, err := os.Open(file)
	if err != nil {
		fmt.Printf("Error opening example json file (%s)\n", file)
		os.Exit(1)
	}
	defer fi.Close()

	dec := json.NewDecoder(fi)
	dec.Decode(&testDashboard)
}

func testDashboardServer() *httptest.Server {
	f := func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if path == "/dashboard/1234" {
			switch r.Method {
			case "GET":
				ret, err := json.Marshal(testDashboard)
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
		} else if path == "/dashboard" {
			switch r.Method {
			case "GET":
				reqURL := r.URL.String()
				var c []Dashboard
				if reqURL == "/dashboard?search=my+dashboard" {
					c = []Dashboard{testDashboard}
				} else if reqURL == "/dashboard?f__created_gt=1483639916" {
					c = []Dashboard{testDashboard}
				} else if reqURL == "/dashboard?f__created_gt=1483639916&search=my+dashboard" {
					c = []Dashboard{testDashboard}
				} else if reqURL == "/dashboard" {
					c = []Dashboard{testDashboard}
				} else {
					c = []Dashboard{}
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
				ret, err := json.Marshal(testDashboard)
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

func TestFetchDashboard(t *testing.T) {
	server := testDashboardServer()
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
		expectedError := errors.New("Invalid dashboard CID [none]")
		_, err := apih.FetchDashboard(CIDType(&cid))
		if err == nil {
			t.Fatalf("Expected error")
		}
		if err.Error() != expectedError.Error() {
			t.Fatalf("Expected %+v got '%+v'", expectedError, err)
		}
	}

	t.Log("with valid CID")
	{
		cid := "/dashboard/1234"
		dashboard, err := apih.FetchDashboard(CIDType(&cid))
		if err != nil {
			t.Fatalf("Expected no error, got '%v'", err)
		}

		actualType := reflect.TypeOf(dashboard)
		expectedType := "*api.Dashboard"
		if actualType.String() != expectedType {
			t.Fatalf("Expected %s, got %s", expectedType, actualType.String())
		}

		if dashboard.CID != testDashboard.CID {
			t.Fatalf("CIDs do not match: %+v != %+v\n", dashboard, testDashboard)
		}
	}

	t.Log("with invalid CID")
	{
		cid := "/invalid"
		expectedError := errors.New("Invalid dashboard CID [/invalid]")
		_, err := apih.FetchDashboard(CIDType(&cid))
		if err == nil {
			t.Fatalf("Expected error")
		}
		if err.Error() != expectedError.Error() {
			t.Fatalf("Expected %+v got '%+v'", expectedError, err)
		}
	}
}

func TestFetchDashboards(t *testing.T) {
	server := testDashboardServer()
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

	dashboards, err := apih.FetchDashboards()
	if err != nil {
		t.Fatalf("Expected no error, got '%v'", err)
	}

	actualType := reflect.TypeOf(dashboards)
	expectedType := "*[]api.Dashboard"
	if actualType.String() != expectedType {
		t.Fatalf("Expected %s, got %s", expectedType, actualType.String())
	}

}

func TestCreateDashboard(t *testing.T) {
	server := testDashboardServer()
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

	dashboard, err := apih.CreateDashboard(&testDashboard)
	if err != nil {
		t.Fatalf("Expected no error, got '%v'", err)
	}

	actualType := reflect.TypeOf(dashboard)
	expectedType := "*api.Dashboard"
	if actualType.String() != expectedType {
		t.Fatalf("Expected %s, got %s", expectedType, actualType.String())
	}
}

func TestUpdateDashboard(t *testing.T) {
	server := testDashboardServer()
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

	t.Log("valid Dashboard")
	{
		dashboard, err := apih.UpdateDashboard(&testDashboard)
		if err != nil {
			t.Fatalf("Expected no error, got '%v'", err)
		}

		actualType := reflect.TypeOf(dashboard)
		expectedType := "*api.Dashboard"
		if actualType.String() != expectedType {
			t.Fatalf("Expected %s, got %s", expectedType, actualType.String())
		}
	}

	t.Log("Test with invalid CID")
	{
		expectedError := errors.New("Invalid dashboard CID [/invalid]")
		x := &Dashboard{CID: "/invalid"}
		_, err := apih.UpdateDashboard(x)
		if err == nil {
			t.Fatal("Expected an error")
		}
		if err.Error() != expectedError.Error() {
			t.Fatalf("Expected %+v got '%+v'", expectedError, err)
		}
	}
}

func TestDeleteDashboard(t *testing.T) {
	server := testDashboardServer()
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

	t.Log("valid Dashboard")
	{
		_, err := apih.DeleteDashboard(&testDashboard)
		if err != nil {
			t.Fatalf("Expected no error, got '%v'", err)
		}
	}

	t.Log("Test with invalid CID")
	{
		expectedError := errors.New("Invalid dashboard CID [/invalid]")
		x := &Dashboard{CID: "/invalid"}
		_, err := apih.UpdateDashboard(x)
		if err == nil {
			t.Fatal("Expected an error")
		}
		if err.Error() != expectedError.Error() {
			t.Fatalf("Expected %+v got '%+v'", expectedError, err)
		}
	}
}

func TestSearchDashboards(t *testing.T) {
	server := testDashboardServer()
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

	search := SearchQueryType("my dashboard")
	filter := SearchFilterType(map[string][]string{"f__created_gt": []string{"1483639916"}})

	t.Log("no search, no filter")
	{
		dashboards, err := apih.SearchDashboards(nil, nil)
		if err != nil {
			t.Fatalf("Expected no error, got '%v'", err)
		}

		actualType := reflect.TypeOf(dashboards)
		expectedType := "*[]api.Dashboard"
		if actualType.String() != expectedType {
			t.Fatalf("Expected %s, got %s", expectedType, actualType.String())
		}
	}

	t.Log("search, no filter")
	{
		dashboards, err := apih.SearchDashboards(&search, nil)
		if err != nil {
			t.Fatalf("Expected no error, got '%v'", err)
		}

		actualType := reflect.TypeOf(dashboards)
		expectedType := "*[]api.Dashboard"
		if actualType.String() != expectedType {
			t.Fatalf("Expected %s, got %s", expectedType, actualType.String())
		}
	}

	t.Log("no search, filter")
	{
		dashboards, err := apih.SearchDashboards(nil, &filter)
		if err != nil {
			t.Fatalf("Expected no error, got '%v'", err)
		}

		actualType := reflect.TypeOf(dashboards)
		expectedType := "*[]api.Dashboard"
		if actualType.String() != expectedType {
			t.Fatalf("Expected %s, got %s", expectedType, actualType.String())
		}
	}

	t.Log("search, filter")
	{
		dashboards, err := apih.SearchDashboards(&search, &filter)
		if err != nil {
			t.Fatalf("Expected no error, got '%v'", err)
		}

		actualType := reflect.TypeOf(dashboards)
		expectedType := "*[]api.Dashboard"
		if actualType.String() != expectedType {
			t.Fatalf("Expected %s, got %s", expectedType, actualType.String())
		}
	}
}
