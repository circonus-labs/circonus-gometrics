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
	"strings"
	"testing"

	"github.com/circonus-labs/circonus-gometrics/api/config"
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
		Config:             map[config.Key]string{},
		Brokers:            []string{"/broker/1234"},
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
		path := r.URL.Path
		if path == "/check_bundle/1234" {
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
				w.WriteHeader(404)
				fmt.Fprintln(w, fmt.Sprintf("not found: %s %s", r.Method, path))
			}
		} else if path == "/check_bundle" {
			switch r.Method {
			case "GET":
				var c []CheckBundle
				if strings.Contains(r.URL.String(), "search=test1") {
					c = []CheckBundle{}
				} else if strings.Contains(r.URL.String(), "f__tags_has=cat%3Atag") {
					c = []CheckBundle{testCheckBundle, testCheckBundle}
				} else if strings.Contains(r.URL.String(), "search=HTTPTrap") {
					c = []CheckBundle{testCheckBundle, testCheckBundle}
				} else if strings.Contains(r.URL.String(), "search=notfound") {
					c = []CheckBundle{}
				} else if strings.Contains(r.URL.String(), "f__tags_has=Found&search=Found") {
					c = []CheckBundle{testCheckBundle, testCheckBundle}
				} else if strings.Contains(r.URL.String(), "f__tags_has=NotFound&search=NotFound") {
					c = []CheckBundle{}
				} else {
					c = []CheckBundle{testCheckBundle}
				}

				ret, err := json.Marshal(c)
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

func TestFetchCheckBundle(t *testing.T) {
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

	t.Log("Testing invalid CID")
	{
		cid := CIDType("/invalid")
		expectedError := errors.New("Invalid check bundle CID [/invalid]")
		_, err := apih.FetchCheckBundle(&cid)
		if err == nil {
			t.Fatalf("Expected error")
		}
		if err.Error() != expectedError.Error() {
			t.Fatalf("Expected %+v got '%+v'", expectedError, err)
		}

	}

	t.Log("Testing valid CID")
	{
		cid := CIDType(testCheckBundle.CID)
		bundle, err := apih.FetchCheckBundle(&cid)
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
		bundles, err := apih.SearchCheckBundles(nil, nil)
		if err != nil {
			t.Fatalf("Expected no error, got '%v'", err)
		}

		actualType := reflect.TypeOf(bundles)
		expectedType := "*[]api.CheckBundle"
		if actualType.String() != expectedType {
			t.Fatalf("Expected %s, got %s", expectedType, actualType.String())
		}
	}

	t.Log("Testing with search criteria")
	{
		search := SearchQueryType("test1")
		bundles, err := apih.SearchCheckBundles(&search, nil)
		if err != nil {
			t.Fatalf("Expected no error, got '%v'", err)
		}

		actualType := reflect.TypeOf(bundles)
		expectedType := "*[]api.CheckBundle"
		if actualType.String() != expectedType {
			t.Fatalf("Expected %s, got %s", expectedType, actualType.String())
		}
	}

	t.Log("Testing with search and filter criteria")
	{
		search := SearchQueryType("test")
		filter := map[string]string{"f_notes": "foo"}
		bundles, err := apih.SearchCheckBundles(&search, &filter)
		if err != nil {
			t.Fatalf("Expected no error, got '%v'", err)
		}

		actualType := reflect.TypeOf(bundles)
		expectedType := "*[]api.CheckBundle"
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
	expectedError := errors.New("Invalid check bundle CID [/invalid]")
	x := &CheckBundle{CID: "/invalid"}
	_, err = apih.UpdateCheckBundle(x)
	if err == nil {
		t.Fatal("Expected an error")
	}
	if err.Error() != expectedError.Error() {
		t.Fatalf("Expected %+v got '%+v'", expectedError, err)
	}
}

func TestDeleteCheckBundleByCID(t *testing.T) {
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

	t.Log("Testing invalid CID")
	{
		cid := CIDType("/invalid")
		expectedError := errors.New("Invalid check bundle CID [/invalid]")
		_, err := apih.DeleteCheckBundleByCID(&cid)
		if err == nil {
			t.Fatalf("Expected error")
		}
		if err.Error() != expectedError.Error() {
			t.Fatalf("Expected %+v got '%+v'", expectedError, err)
		}
	}

	t.Log("Testing valid CID")
	{
		cid := CIDType(testCheckBundle.CID)
		success, err := apih.DeleteCheckBundleByCID(&cid)
		if err != nil {
			t.Fatalf("Expected no error, got '%v'", err)
		}

		if !success {
			t.Fatalf("Expected success to be true")
		}
	}
}

func TestDeleteCheckBundle(t *testing.T) {
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

	success, err := apih.DeleteCheckBundle(&testCheckBundle)
	if err != nil {
		t.Fatalf("Expected no error, got '%v'", err)
	}

	if !success {
		t.Fatalf("Expected success to be true")
	}
}
