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
	"strconv"
	"testing"
)

var (
	testCheckBundle = CheckBundle{
		CheckUUIDs:         []string{"abc123-a1b2-c3d4-e5f6-123abc"},
		Checks:             []string{"/check/1234"},
		CID:                "/check_bundle/1234",
		Created:            0,
		LastModified:       0,
		LastModifedBy:      "",
		ReverseConnectURLs: []string{""},
		Brokers:            []string{"/broker/1234"},
		Config:             CheckBundleConfig{},
		DisplayName:        "test check",
		Metrics:            []CheckBundleMetric{},
		MetricLimit:        0,
		Notes:              "",
		Period:             60,
		Status:             "active",
		Target:             "127.0.0.1",
		Timeout:            10,
		Type:               "httptrap",
		Tags:               []string{},
	}
)

func testCheckBundleServer() *httptest.Server {
	f := func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/check_bundle/1234": // handle GET/PUT/DELETE
			switch r.Method {
			case "PUT": // update
				defer r.Body.Close()
				b, err := ioutil.ReadAll(r.Body)
				if err != nil {
					panic(err)
				}
				w.WriteHeader(200)
				w.Header().Set("Content-Type", "application/json")
				fmt.Fprintln(w, string(b))
			case "GET": // get by id/cid
				ret, err := json.Marshal(testCheckBundle)
				if err != nil {
					panic(err)
				}
				w.WriteHeader(200)
				w.Header().Set("Content-Type", "application/json")
				fmt.Fprintln(w, string(ret))
			case "DELETE": // delete
				w.WriteHeader(200)
				fmt.Fprintln(w, "")
			default:
				w.WriteHeader(500)
				fmt.Fprintln(w, "unsupported")
			}
		case "/check_bundle":
			switch r.Method {
			case "GET": // search
				r := []CheckBundle{testCheckBundle}
				ret, err := json.Marshal(r)
				if err != nil {
					panic(err)
				}
				w.WriteHeader(200)
				w.Header().Set("Content-Type", "application/json")
				fmt.Fprintln(w, string(ret))
			case "POST": // create
				defer r.Body.Close()
				b, err := ioutil.ReadAll(r.Body)
				if err != nil {
					panic(err)
				}
				w.WriteHeader(200)
				w.Header().Set("Content-Type", "application/json")
				fmt.Fprintln(w, string(b))
			default:
				w.WriteHeader(500)
				fmt.Fprintln(w, "unsupported")
			}
		default:
			w.WriteHeader(500)
			fmt.Fprintln(w, "unsupported")
		}
	}

	return httptest.NewServer(http.HandlerFunc(f))
}

func TestFetchCheckBundleByID(t *testing.T) {
	server := testCheckBundleServer()
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

	cid := "1234"
	id, err := strconv.Atoi(cid)
	if err != nil {
		t.Fatalf("Unable to convert id %s to int", cid)
	}

	cbID := IDType(id)

	bundle, err := apih.FetchCheckBundleByID(cbID)
	if err != nil {
		t.Fatalf("Expected no error, got '%v'", err)
	}

	actualType := reflect.TypeOf(bundle)
	expectedType := "*api.CheckBundle"
	if actualType.String() != expectedType {
		t.Fatalf("Expected %s, got %s", expectedType, actualType.String())
	}

	if bundle.CID != testCheckBundle.CID {
		t.Fatalf("CIDs do not match: %+v != %+v\n", bundle, testCheckBundle)
	}
}

func TestFetchCheckBundleByCID(t *testing.T) {
	server := testCheckBundleServer()
	defer server.Close()

	var apih *API
	var err error

	ac := &Config{
		TokenKey: "abc123",
		TokenApp: "test",
		URL:      server.URL,
	}
	apih, err = NewAPI(ac)
	if err != nil {
		t.Errorf("Expected no error, got '%v'", err)
	}

	t.Log("Testing invalid CID")
	expectedError := errors.New("Invalid check bundle CID /1234")
	_, err = apih.FetchCheckBundleByCID("/1234")
	if err == nil {
		t.Fatalf("Expected error")
	}
	if err.Error() != expectedError.Error() {
		t.Fatalf("Expected %+v got '%+v'", expectedError, err)
	}

	t.Log("Testing valid CID")
	bundle, err := apih.FetchCheckBundleByCID(CIDType(testCheckBundle.CID))
	if err != nil {
		t.Fatalf("Expected no error, got '%v'", err)
	}

	actualType := reflect.TypeOf(bundle)
	expectedType := "*api.CheckBundle"
	if actualType.String() != expectedType {
		t.Fatalf("Expected %s, got %s", expectedType, actualType.String())
	}

	if bundle.CID != testCheckBundle.CID {
		t.Fatalf("CIDs do not match: %+v != %+v\n", bundle, testCheckBundle)
	}
}

func TestCheckBundleSearch(t *testing.T) {
	server := testCheckBundleServer()
	defer server.Close()

	var apih *API
	var err error

	ac := &Config{
		TokenKey: "abc123",
		TokenApp: "test",
		URL:      server.URL,
	}
	apih, err = NewAPI(ac)
	if err != nil {
		t.Errorf("Expected no error, got '%v'", err)
	}

	t.Log("Testing w/o search criteria")
	{
		bundles, err := apih.CheckBundleSearch("", map[string]string{})
		if err != nil {
			t.Fatalf("Expected no error, got '%v'", err)
		}

		actualType := reflect.TypeOf(bundles)
		expectedType := "[]api.CheckBundle"
		if actualType.String() != expectedType {
			t.Fatalf("Expected %s, got %s", expectedType, actualType.String())
		}
	}

	t.Log("Testing with search criteria")
	{
		bundles, err := apih.CheckBundleSearch("test", map[string]string{"f_notes": "foo"})
		if err != nil {
			t.Fatalf("Expected no error, got '%v'", err)
		}

		actualType := reflect.TypeOf(bundles)
		expectedType := "[]api.CheckBundle"
		if actualType.String() != expectedType {
			t.Fatalf("Expected %s, got %s", expectedType, actualType.String())
		}
	}
}

func TestCreateCheckBundle(t *testing.T) {
	server := testCheckBundleServer()
	defer server.Close()

	var apih *API
	var err error

	ac := &Config{
		TokenKey: "abc123",
		TokenApp: "test",
		URL:      server.URL,
	}
	apih, err = NewAPI(ac)
	if err != nil {
		t.Errorf("Expected no error, got '%v'", err)
	}

	bundle, err := apih.CreateCheckBundle(&testCheckBundle)
	if err != nil {
		t.Fatalf("Expected no error, got '%v'", err)
	}

	actualType := reflect.TypeOf(bundle)
	expectedType := "*api.CheckBundle"
	if actualType.String() != expectedType {
		t.Fatalf("Expected %s, got %s", expectedType, actualType.String())
	}
}

func TestUpdateCheckBundle(t *testing.T) {
	server := testCheckBundleServer()
	defer server.Close()

	var apih *API
	var err error

	ac := &Config{
		TokenKey: "abc123",
		TokenApp: "test",
		URL:      server.URL,
	}
	apih, err = NewAPI(ac)
	if err != nil {
		t.Errorf("Expected no error, got '%v'", err)
	}

	bundle, err := apih.UpdateCheckBundle(&testCheckBundle)
	if err != nil {
		t.Fatalf("Expected no error, got '%v'", err)
	}

	actualType := reflect.TypeOf(bundle)
	expectedType := "*api.CheckBundle"
	if actualType.String() != expectedType {
		t.Fatalf("Expected %s, got %s", expectedType, actualType.String())
	}

	t.Log("Test with invalid CID")
	expectedError := errors.New("Invalid check bundle CID xxx")
	x := &CheckBundle{}
	x = &testCheckBundle
	x.CID = "xxx"
	_, err = apih.UpdateCheckBundle(x)
	if err == nil {
		t.Fatal("Expected an error")
	}
	if err.Error() != expectedError.Error() {
		t.Fatalf("Expected %+v got '%+v'", expectedError, err)
	}
}
