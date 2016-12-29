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
	testGraph = Graph{
		CID:        "/graph/01234567-89ab-cdef-0123-456789abcdef",
		AccessKeys: []GraphAccessKey{},
		Composites: []GraphComposite{
			GraphComposite{
				Axis:        "l",
				Color:       "#000000",
				DataFormula: "=A-B",
				Hidden:      false,
				Name:        "Time After First Byte",
			},
		},
		Datapoints: []GraphDatapoint{
			GraphDatapoint{
				Axis:        "l",
				CheckID:     1234,
				Color:       "#ff0000",
				DataFormula: "=VAL/1000",
				Derive:      "gauge",
				Hidden:      false,
				MetricName:  "duration",
				MetricType:  "numeric",
				Name:        "Total Request Time",
			},
			GraphDatapoint{
				Axis:        "l",
				CheckID:     2345,
				Color:       "#00ff00",
				DataFormula: "=VAL/1000",
				Derive:      "gauge",
				Hidden:      false,
				MetricName:  "tt_firstbyte",
				MetricType:  "numeric",
				Name:        "Time Till First Byte",
			},
		},
		Description: "Time to first byte verses time to whole thing",
		LineStyle:   "interpolated",
		LogLeftY:    10,
		Notes:       "This graph shows just the main webserver",
		Style:       "line",
		Tags:        []string{"datacenter:primary"},
		Title:       "Slow Webserver",
	}
)

func testGraphServer() *httptest.Server {
	f := func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if path == "/graph/01234567-89ab-cdef-0123-456789abcdef" {
			switch r.Method {
			case "GET":
				ret, err := json.Marshal(testGraph)
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
		} else if path == "/graph" {
			switch r.Method {
			case "GET":
				c := []Graph{testGraph}
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
				ret, err := json.Marshal(testGraph)
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

func TestFetchGraph(t *testing.T) {
	server := testGraphServer()
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
		expectedError := errors.New("Invalid graph CID ")
		_, err := apih.FetchGraph(CIDType(""))
		if err == nil {
			t.Fatalf("Expected error")
		}
		if err.Error() != expectedError.Error() {
			t.Fatalf("Expected %+v got '%+v'", expectedError, err)
		}
	}

	t.Log("with valid CID")
	{
		cid := CIDType("/graph/01234567-89ab-cdef-0123-456789abcdef")
		graph, err := apih.FetchGraph(cid)
		if err != nil {
			t.Fatalf("Expected no error, got '%v'", err)
		}

		actualType := reflect.TypeOf(graph)
		expectedType := "*api.Graph"
		if actualType.String() != expectedType {
			t.Fatalf("Expected %s, got %s", expectedType, actualType.String())
		}

		if graph.CID != testGraph.CID {
			t.Fatalf("CIDs do not match: %+v != %+v\n", graph, testGraph)
		}
	}

	t.Log("with invalid CID")
	{
		expectedError := errors.New("Invalid graph CID /invalid")
		_, err := apih.FetchGraph(CIDType("/invalid"))
		if err == nil {
			t.Fatalf("Expected error")
		}
		if err.Error() != expectedError.Error() {
			t.Fatalf("Expected %+v got '%+v'", expectedError, err)
		}
	}
}

func TestFetchGraphs(t *testing.T) {
	server := testGraphServer()
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

	graphs, err := apih.FetchGraphs()
	if err != nil {
		t.Fatalf("Expected no error, got '%v'", err)
	}

	actualType := reflect.TypeOf(graphs)
	expectedType := "[]api.Graph"
	if actualType.String() != expectedType {
		t.Fatalf("Expected %s, got %s", expectedType, actualType.String())
	}

}

func TestCreateGraph(t *testing.T) {
	server := testGraphServer()
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

	graph, err := apih.CreateGraph(&testGraph)
	if err != nil {
		t.Fatalf("Expected no error, got '%v'", err)
	}

	actualType := reflect.TypeOf(graph)
	expectedType := "*api.Graph"
	if actualType.String() != expectedType {
		t.Fatalf("Expected %s, got %s", expectedType, actualType.String())
	}
}

func TestUpdateGraph(t *testing.T) {
	server := testGraphServer()
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

	t.Log("valid Graph")
	{
		graph, err := apih.UpdateGraph(&testGraph)
		if err != nil {
			t.Fatalf("Expected no error, got '%v'", err)
		}

		actualType := reflect.TypeOf(graph)
		expectedType := "*api.Graph"
		if actualType.String() != expectedType {
			t.Fatalf("Expected %s, got %s", expectedType, actualType.String())
		}
	}

	t.Log("Test with invalid CID")
	{
		expectedError := errors.New("Invalid graph CID /invalid")
		x := &Graph{CID: "/invalid"}
		_, err := apih.UpdateGraph(x)
		if err == nil {
			t.Fatal("Expected an error")
		}
		if err.Error() != expectedError.Error() {
			t.Fatalf("Expected %+v got '%+v'", expectedError, err)
		}
	}
}

func TestDeleteGraph(t *testing.T) {
	server := testGraphServer()
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

	t.Log("valid Graph")
	{
		_, err := apih.DeleteGraph(&testGraph)
		if err != nil {
			t.Fatalf("Expected no error, got '%v'", err)
		}
	}

	t.Log("Test with invalid CID")
	{
		expectedError := errors.New("Invalid graph CID /invalid")
		x := &Graph{CID: "/invalid"}
		_, err := apih.UpdateGraph(x)
		if err == nil {
			t.Fatal("Expected an error")
		}
		if err.Error() != expectedError.Error() {
			t.Fatalf("Expected %+v got '%+v'", expectedError, err)
		}
	}
}
