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
	testAnnotation = Annotation{
		CID:            "/annotation/1234",
		Created:        1483033102,
		LastModified:   1483033102,
		LastModifiedBy: "/user/1234",
		Start:          1483033100,
		Stop:           1483033102,
		Category:       "foo",
		Title:          "Foo Bar Baz",
	}
)

func testAnnotationServer() *httptest.Server {
	f := func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if path == "/annotation/1234" {
			switch r.Method {
			case "GET":
				ret, err := json.Marshal(testAnnotation)
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
		} else if path == "/annotation" {
			switch r.Method {
			case "GET":
				reqURL := r.URL.String()
				var c []Annotation
				if reqURL == "/annotation?search=%28category%3D%22updates%22%29" {
					c = []Annotation{testAnnotation}
				} else if reqURL == "/annotation?f__created_gt=1483639916" {
					c = []Annotation{testAnnotation}
				} else if reqURL == "/annotation?f__created_gt=1483639916&search=%28category%3D%22updates%22%29" {
					c = []Annotation{testAnnotation}
				} else if reqURL == "/annotation" {
					c = []Annotation{testAnnotation}
				} else {
					c = []Annotation{}
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
				ret, err := json.Marshal(testAnnotation)
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

func TestFetchAnnotation(t *testing.T) {
	server := testAnnotationServer()
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
		expectedError := errors.New("Invalid annotation CID [none]")
		_, err := apih.FetchAnnotation(CIDType(&cid))
		if err == nil {
			t.Fatalf("Expected error")
		}
		if err.Error() != expectedError.Error() {
			t.Fatalf("Expected %+v got '%+v'", expectedError, err)
		}
	}

	t.Log("with valid CID")
	{
		cid := "/annotation/1234"
		annotation, err := apih.FetchAnnotation(CIDType(&cid))
		if err != nil {
			t.Fatalf("Expected no error, got '%v'", err)
		}

		actualType := reflect.TypeOf(annotation)
		expectedType := "*api.Annotation"
		if actualType.String() != expectedType {
			t.Fatalf("Expected %s, got %s", expectedType, actualType.String())
		}

		if annotation.CID != testAnnotation.CID {
			t.Fatalf("CIDs do not match: %+v != %+v\n", annotation, testAnnotation)
		}
	}

	t.Log("with invalid CID")
	{
		cid := "/invalid"
		expectedError := errors.New("Invalid annotation CID [/invalid]")
		_, err := apih.FetchAnnotation(CIDType(&cid))
		if err == nil {
			t.Fatalf("Expected error")
		}
		if err.Error() != expectedError.Error() {
			t.Fatalf("Expected %+v got '%+v'", expectedError, err)
		}
	}
}

func TestFetchAnnotations(t *testing.T) {
	server := testAnnotationServer()
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

	annotations, err := apih.FetchAnnotations()
	if err != nil {
		t.Fatalf("Expected no error, got '%v'", err)
	}

	actualType := reflect.TypeOf(annotations)
	expectedType := "*[]api.Annotation"
	if actualType.String() != expectedType {
		t.Fatalf("Expected %s, got %s", expectedType, actualType.String())
	}

}

func TestCreateAnnotation(t *testing.T) {
	server := testAnnotationServer()
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

	annotation, err := apih.CreateAnnotation(&testAnnotation)
	if err != nil {
		t.Fatalf("Expected no error, got '%v'", err)
	}

	actualType := reflect.TypeOf(annotation)
	expectedType := "*api.Annotation"
	if actualType.String() != expectedType {
		t.Fatalf("Expected %s, got %s", expectedType, actualType.String())
	}
}

func TestUpdateAnnotation(t *testing.T) {
	server := testAnnotationServer()
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

	t.Log("valid Annotation")
	{
		annotation, err := apih.UpdateAnnotation(&testAnnotation)
		if err != nil {
			t.Fatalf("Expected no error, got '%v'", err)
		}

		actualType := reflect.TypeOf(annotation)
		expectedType := "*api.Annotation"
		if actualType.String() != expectedType {
			t.Fatalf("Expected %s, got %s", expectedType, actualType.String())
		}
	}

	t.Log("Test with invalid CID")
	{
		expectedError := errors.New("Invalid annotation CID [/invalid]")
		x := &Annotation{CID: "/invalid"}
		_, err := apih.UpdateAnnotation(x)
		if err == nil {
			t.Fatal("Expected an error")
		}
		if err.Error() != expectedError.Error() {
			t.Fatalf("Expected %+v got '%+v'", expectedError, err)
		}
	}
}

func TestDeleteAnnotation(t *testing.T) {
	server := testAnnotationServer()
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

	t.Log("valid Annotation")
	{
		_, err := apih.DeleteAnnotation(&testAnnotation)
		if err != nil {
			t.Fatalf("Expected no error, got '%v'", err)
		}
	}

	t.Log("Test with invalid CID")
	{
		expectedError := errors.New("Invalid annotation CID [/invalid]")
		x := &Annotation{CID: "/invalid"}
		_, err := apih.UpdateAnnotation(x)
		if err == nil {
			t.Fatal("Expected an error")
		}
		if err.Error() != expectedError.Error() {
			t.Fatalf("Expected %+v got '%+v'", expectedError, err)
		}
	}
}

func TestSearchAnnotations(t *testing.T) {
	server := testAnnotationServer()
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

	t.Log("no search, no filter")
	{
		annotations, err := apih.SearchAnnotations(nil, nil)
		if err != nil {
			t.Fatalf("Expected no error, got '%v'", err)
		}

		actualType := reflect.TypeOf(annotations)
		expectedType := "*[]api.Annotation"
		if actualType.String() != expectedType {
			t.Fatalf("Expected %s, got %s", expectedType, actualType.String())
		}
	}

	t.Log("search, no filter")
	{
		search := SearchQueryType(`(category="updates")`)
		annotations, err := apih.SearchAnnotations(&search, nil)
		if err != nil {
			t.Fatalf("Expected no error, got '%v'", err)
		}

		actualType := reflect.TypeOf(annotations)
		expectedType := "*[]api.Annotation"
		if actualType.String() != expectedType {
			t.Fatalf("Expected %s, got %s", expectedType, actualType.String())
		}
	}

	t.Log("no search, filter")
	{
		filter := SearchFilterType(map[string][]string{"f__created_gt": []string{"1483639916"}})
		annotations, err := apih.SearchAnnotations(nil, &filter)
		if err != nil {
			t.Fatalf("Expected no error, got '%v'", err)
		}

		actualType := reflect.TypeOf(annotations)
		expectedType := "*[]api.Annotation"
		if actualType.String() != expectedType {
			t.Fatalf("Expected %s, got %s", expectedType, actualType.String())
		}
	}

	t.Log("search, filter")
	{
		search := SearchQueryType(`(category="updates")`)
		filter := SearchFilterType(map[string][]string{"f__created_gt": []string{"1483639916"}})
		annotations, err := apih.SearchAnnotations(&search, &filter)
		if err != nil {
			t.Fatalf("Expected no error, got '%v'", err)
		}

		actualType := reflect.TypeOf(annotations)
		expectedType := "*[]api.Annotation"
		if actualType.String() != expectedType {
			t.Fatalf("Expected %s, got %s", expectedType, actualType.String())
		}
	}
}
