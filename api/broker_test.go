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
	"strings"
	"testing"
)

var (
	testBroker = Broker{
		CID:       "/broker/1234",
		Longitude: "",
		Latitude:  "",
		Name:      "test broker",
		Tags:      []string{},
		Type:      "enterprise",
		Details: []BrokerDetail{
			BrokerDetail{
				CN:           "testbroker.example.com",
				ExternalHost: "testbroker.example.com",
				ExternalPort: 43191,
				IP:           "127.0.0.1",
				MinVer:       0,
				Modules:      []string{"a", "b", "c"},
				Port:         43191,
				Skew:         "",
				Status:       "active",
				Version:      1,
			},
		},
	}
)

func testBrokerServer() *httptest.Server {
	f := func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/broker/1234": // handle GET/PUT/DELETE
			switch r.Method {
			case "GET": // get by id/cid
				ret, err := json.Marshal(testBroker)
				if err != nil {
					panic(err)
				}
				w.WriteHeader(200)
				w.Header().Set("Content-Type", "application/json")
				fmt.Fprintln(w, string(ret))
			default:
				w.WriteHeader(500)
				fmt.Fprintln(w, "unsupported")
			}
		case "/broker":
			switch r.Method {
			case "GET": // search or filter
				var c []Broker
				if strings.Contains(r.URL.String(), "f__tags_has=cat%3Anot_found") {
					c = []Broker{}
				} else if strings.Contains(r.URL.String(), "f__tags_has=cat%3Atag") {
					c = []Broker{testBroker, testBroker}
				} else if strings.Contains(r.URL.String(), "search=HTTPTrap") {
					c = []Broker{testBroker, testBroker}
				} else if strings.Contains(r.URL.String(), "search=notfound") {
					c = []Broker{}
				} else if strings.Contains(r.URL.String(), "f__tags_has=Found&search=Found") {
					c = []Broker{testBroker, testBroker}
				} else if strings.Contains(r.URL.String(), "f__tags_has=NotFound&search=NotFound") {
					c = []Broker{}
				} else {
					c = []Broker{testBroker}
				}

				ret, err := json.Marshal(c)
				if err != nil {
					panic(err)
				}
				w.WriteHeader(200)
				w.Header().Set("Content-Type", "application/json")
				fmt.Fprintln(w, string(ret))
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

func TestFetchBrokerByID(t *testing.T) {
	server := testBrokerServer()
	defer server.Close()

	ac := &Config{
		TokenKey: "abc123",
		TokenApp: "test",
		URL:      server.URL,
	}
	apih, err := New(ac)
	if err != nil {
		t.Errorf("Expected no error, got '%v'", err)
	}

	t.Log("valid ID")
	{
		brokerID := IDType(1234)

		broker, err := apih.FetchBrokerByID(brokerID)
		if err != nil {
			t.Fatalf("Expected no error, got '%v'", err)
		}

		actualType := reflect.TypeOf(broker)
		expectedType := "*api.Broker"
		if actualType.String() != expectedType {
			t.Fatalf("Expected %s, got %s", expectedType, actualType.String())
		}

		if broker.CID != testBroker.CID {
			t.Fatalf("CIDs do not match: %+v != %+v\n", broker, testBroker)
		}
	}

	t.Log("invalid ID")
	{
		brokerID := IDType(-1)

		expectedError := errors.New("Invalid broker ID [-1]")
		_, err := apih.FetchBrokerByID(brokerID)
		if err == nil {
			t.Fatalf("Expected error")
		}
		if err.Error() != expectedError.Error() {
			t.Fatalf("Expected %+v got '%+v'", expectedError, err)
		}
	}
}

func TestFetchBroker(t *testing.T) {
	server := testBrokerServer()
	defer server.Close()

	var apih *API
	var err error

	ac := &Config{
		TokenKey: "abc123",
		TokenApp: "test",
		URL:      server.URL,
	}
	apih, err = New(ac)
	if err != nil {
		t.Errorf("Expected no error, got '%v'", err)
	}

	t.Log("invalid CID [nil]")
	{
		expectedError := errors.New("Invalid broker CID [none]")
		_, err := apih.FetchBroker(nil)
		if err == nil {
			t.Fatalf("Expected error")
		}
		if err.Error() != expectedError.Error() {
			t.Fatalf("Expected %+v got '%+v'", expectedError, err)
		}
	}

	t.Log("invalid CID [/invalid]")
	{
		cid := CIDType("/invalid")
		expectedError := errors.New("Invalid broker CID [/invalid]")
		_, err := apih.FetchBroker(&cid)
		if err == nil {
			t.Fatalf("Expected error")
		}
		if err.Error() != expectedError.Error() {
			t.Fatalf("Expected %+v got '%+v'", expectedError, err)
		}
	}

	t.Log("valid CID")
	{
		cid := CIDType(testBroker.CID)
		broker, err := apih.FetchBroker(&cid)
		if err != nil {
			t.Fatalf("Expected no error, got '%v'", err)
		}

		actualType := reflect.TypeOf(broker)
		expectedType := "*api.Broker"
		if actualType.String() != expectedType {
			t.Fatalf("Expected %s, got %s", expectedType, actualType.String())
		}

		if broker.CID != testBroker.CID {
			t.Fatalf("CIDs do not match: %+v != %+v\n", broker, testBroker)
		}
	}
}

func TestFetchBrokers(t *testing.T) {
	server := testBrokerServer()
	defer server.Close()

	ac := &Config{
		TokenKey: "abc123",
		TokenApp: "test",
		URL:      server.URL,
	}
	apih, err := New(ac)
	if err != nil {
		t.Errorf("Expected no error, got '%v'", err)
	}

	brokers, err := apih.FetchBrokers()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	actualType := reflect.TypeOf(brokers)
	expectedType := "*[]api.Broker"
	if actualType.String() != expectedType {
		t.Fatalf("Expected %s, got %s", expectedType, actualType.String())
	}
}

func TestFetchBrokersByTag(t *testing.T) {
	server := testBrokerServer()
	defer server.Close()

	ac := &Config{
		TokenKey: "abc123",
		TokenApp: "test",
		URL:      server.URL,
	}
	apih, err := New(ac)
	if err != nil {
		t.Errorf("Expected no error, got '%v'", err)
	}

	t.Log("no tags")
	{
		brokers, err := apih.FetchBrokersByTag(TagType([]string{}))
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		actualType := reflect.TypeOf(brokers)
		expectedType := "*[]api.Broker"
		if actualType.String() != expectedType {
			t.Fatalf("Expected %s, got %s", expectedType, actualType.String())
		}
	}

	t.Log("tag(s) [found]")
	{
		brokers, err := apih.FetchBrokersByTag(TagType([]string{"cat:tag"}))
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		actualType := reflect.TypeOf(brokers)
		expectedType := "*[]api.Broker"
		if actualType.String() != expectedType {
			t.Fatalf("Expected %s, got %s", expectedType, actualType.String())
		}
		if len(*brokers) == 0 {
			t.Fatal("Expected >0 brokers, got 0")
		}
	}

	t.Log("tag(s) [not found]")
	{
		brokers, err := apih.FetchBrokersByTag(TagType([]string{"cat:not_found"}))
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		actualType := reflect.TypeOf(brokers)
		expectedType := "*[]api.Broker"
		if actualType.String() != expectedType {
			t.Fatalf("Expected %s, got %s", expectedType, actualType.String())
		}
		if len(*brokers) != 0 {
			t.Fatalf("Expected 0 brokers, got %d", len(*brokers))
		}
	}
}

func TestSearchBrokers(t *testing.T) {
	server := testBrokerServer()
	defer server.Close()

	ac := &Config{
		TokenKey: "abc123",
		TokenApp: "test",
		URL:      server.URL,
	}
	apih, err := New(ac)
	if err != nil {
		t.Errorf("Expected no error, got '%v'", err)
	}

	t.Log("search [nil search, nil filter]")
	{
		brokers, err := apih.SearchBrokers(nil, nil)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		actualType := reflect.TypeOf(brokers)
		expectedType := "*[]api.Broker"
		if actualType.String() != expectedType {
			t.Fatalf("Expected %s, got %s", expectedType, actualType.String())
		}
	}

	t.Log("search [search, nil filter, found]")
	{
		search := SearchQueryType("HTTPTrap")
		brokers, err := apih.SearchBrokers(&search, nil)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		actualType := reflect.TypeOf(brokers)
		expectedType := "*[]api.Broker"
		if actualType.String() != expectedType {
			t.Fatalf("Expected %s, got %s", expectedType, actualType.String())
		}
	}

	t.Log("search [search, nil filter, not found]")
	{
		search := SearchQueryType("notfound")
		brokers, err := apih.SearchBrokers(&search, nil)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		actualType := reflect.TypeOf(brokers)
		expectedType := "*[]api.Broker"
		if actualType.String() != expectedType {
			t.Fatalf("Expected %s, got %s", expectedType, actualType.String())
		}

		if len(*brokers) > 0 {
			t.Fatalf("Expected 0 got %d", len(*brokers))
		}
	}

	t.Log("search [search, filter, found]")
	{
		filter := map[string]string{"f__tags_has": "Found"}
		search := SearchQueryType("Found")
		brokers, err := apih.SearchBrokers(&search, &filter)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		actualType := reflect.TypeOf(brokers)
		expectedType := "*[]api.Broker"
		if actualType.String() != expectedType {
			t.Fatalf("Expected %s, got %s", expectedType, actualType.String())
		}
	}

	t.Log("search [search, filter, not found]")
	{
		filter := map[string]string{"f__tags_has": "NotFound"}
		search := SearchQueryType("NotFound")
		brokers, err := apih.SearchBrokers(&search, &filter)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		actualType := reflect.TypeOf(brokers)
		expectedType := "*[]api.Broker"
		if actualType.String() != expectedType {
			t.Fatalf("Expected %s, got %s", expectedType, actualType.String())
		}
		if len(*brokers) > 0 {
			t.Fatalf("Expected 0 got %d", len(*brokers))
		}
	}

}
