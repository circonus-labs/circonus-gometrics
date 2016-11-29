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
	"strconv"
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
				if strings.Contains(r.URL.String(), "f__check_uuid=none") {
					c = []Broker{}
				} else if strings.Contains(r.URL.String(), "f__check_uuid=multi") {
					c = []Broker{testBroker, testBroker}
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
	apih, err := NewAPI(ac)
	if err != nil {
		t.Errorf("Expected no error, got '%v'", err)
	}

	cid := "1234"
	id, err := strconv.Atoi(cid)
	if err != nil {
		t.Fatalf("Unable to convert id %s to int", cid)
	}

	brokerID := IDType(id)

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

func TestFetchBrokerByCID(t *testing.T) {
	server := testBrokerServer()
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

	t.Log("invalid CID")
	{
		expectedError := errors.New("Invalid broker CID /1234")
		_, err := apih.FetchBrokerByCID("/1234")
		if err == nil {
			t.Fatalf("Expected error")
		}
		if err.Error() != expectedError.Error() {
			t.Fatalf("Expected %+v got '%+v'", expectedError, err)
		}
	}

	t.Log("valid CID")
	{
		broker, err := apih.FetchBrokerByCID(CIDType(testBroker.CID))
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

func TestFetchBrokerList(t *testing.T) {
	server := testBrokerServer()
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

	_, err = apih.FetchBrokerList()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestFetchBrokerListByTag(t *testing.T) {
	server := testBrokerServer()
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

	_, err = apih.FetchBrokerListByTag(TagType([]string{"cat:tag"}))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestBrokerSearch(t *testing.T) {
	server := testBrokerServer()
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

	_, err = apih.BrokerSearch("foo")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}
