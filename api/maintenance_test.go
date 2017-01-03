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
	testMaintenance = Maintenance{
		CID:        "/maintenance/1234",
		Item:       "/check/1234",
		Notes:      "upgrading blah",
		Severities: []string{"1", "2", "3", "4", "5"},
		Start:      1483033100,
		Stop:       1483033102,
		Type:       "check",
		Tags:       []string{"cat:tag"},
	}
)

func testMaintenanceServer() *httptest.Server {
	f := func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if path == "/maintenance/1234" {
			switch r.Method {
			case "GET":
				ret, err := json.Marshal(testMaintenance)
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
		} else if path == "/maintenance" {
			switch r.Method {
			case "GET":
				c := []Maintenance{testMaintenance}
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
				ret, err := json.Marshal(testMaintenance)
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

func TestFetchMaintenanceWindow(t *testing.T) {
	server := testMaintenanceServer()
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
		expectedError := errors.New("Invalid maintenance window CID [none]")
		_, err := apih.FetchMaintenanceWindow(CIDType(&cid))
		if err == nil {
			t.Fatalf("Expected error")
		}
		if err.Error() != expectedError.Error() {
			t.Fatalf("Expected %+v got '%+v'", expectedError, err)
		}
	}

	t.Log("with valid CID")
	{
		cid := "/maintenance/1234"
		maintenance, err := apih.FetchMaintenanceWindow(CIDType(&cid))
		if err != nil {
			t.Fatalf("Expected no error, got '%v'", err)
		}

		actualType := reflect.TypeOf(maintenance)
		expectedType := "*api.Maintenance"
		if actualType.String() != expectedType {
			t.Fatalf("Expected %s, got %s", expectedType, actualType.String())
		}

		if maintenance.CID != testMaintenance.CID {
			t.Fatalf("CIDs do not match: %+v != %+v\n", maintenance, testMaintenance)
		}
	}

	t.Log("with invalid CID")
	{
		cid := "/invalid"
		expectedError := errors.New("Invalid maintenance window CID [/invalid]")
		_, err := apih.FetchMaintenanceWindow(CIDType(&cid))
		if err == nil {
			t.Fatalf("Expected error")
		}
		if err.Error() != expectedError.Error() {
			t.Fatalf("Expected %+v got '%+v'", expectedError, err)
		}
	}
}

func TestFetchMaintenanceWindows(t *testing.T) {
	server := testMaintenanceServer()
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

	maintenances, err := apih.FetchMaintenanceWindows()
	if err != nil {
		t.Fatalf("Expected no error, got '%v'", err)
	}

	actualType := reflect.TypeOf(maintenances)
	expectedType := "*[]api.Maintenance"
	if actualType.String() != expectedType {
		t.Fatalf("Expected %s, got %s", expectedType, actualType.String())
	}

}

func TestCreateMaintenanceWindow(t *testing.T) {
	server := testMaintenanceServer()
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

	maintenance, err := apih.CreateMaintenanceWindow(&testMaintenance)
	if err != nil {
		t.Fatalf("Expected no error, got '%v'", err)
	}

	actualType := reflect.TypeOf(maintenance)
	expectedType := "*api.Maintenance"
	if actualType.String() != expectedType {
		t.Fatalf("Expected %s, got %s", expectedType, actualType.String())
	}
}

func TestUpdateMaintenanceWindow(t *testing.T) {
	server := testMaintenanceServer()
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

	t.Log("valid Maintenance")
	{
		maintenance, err := apih.UpdateMaintenanceWindow(&testMaintenance)
		if err != nil {
			t.Fatalf("Expected no error, got '%v'", err)
		}

		actualType := reflect.TypeOf(maintenance)
		expectedType := "*api.Maintenance"
		if actualType.String() != expectedType {
			t.Fatalf("Expected %s, got %s", expectedType, actualType.String())
		}
	}

	t.Log("Test with invalid CID")
	{
		expectedError := errors.New("Invalid maintenance window CID [/invalid]")
		x := &Maintenance{CID: "/invalid"}
		_, err := apih.UpdateMaintenanceWindow(x)
		if err == nil {
			t.Fatal("Expected an error")
		}
		if err.Error() != expectedError.Error() {
			t.Fatalf("Expected %+v got '%+v'", expectedError, err)
		}
	}
}

func TestDeleteMaintenanceWindow(t *testing.T) {
	server := testMaintenanceServer()
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

	t.Log("valid Maintenance")
	{
		_, err := apih.DeleteMaintenanceWindow(&testMaintenance)
		if err != nil {
			t.Fatalf("Expected no error, got '%v'", err)
		}
	}

	t.Log("Test with invalid CID")
	{
		expectedError := errors.New("Invalid maintenance window CID [/invalid]")
		x := &Maintenance{CID: "/invalid"}
		_, err := apih.UpdateMaintenanceWindow(x)
		if err == nil {
			t.Fatal("Expected an error")
		}
		if err.Error() != expectedError.Error() {
			t.Fatalf("Expected %+v got '%+v'", expectedError, err)
		}
	}
}
