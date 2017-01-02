// Copyright 2016 Circonus, Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package api

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func callServer() *httptest.Server {
	f := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, r.Method)
	}

	return httptest.NewServer(http.HandlerFunc(f))
}

// func retryServer() *httptest.Server {
// 	f := func(w http.ResponseWriter, r *http.Request) {
// 		w.WriteHeader(500)
// 		w.Header().Set("Content-Type", "application/json")
// 		fmt.Fprintln(w, "blah blah blah")
// 	}
//
// 	return httptest.NewServer(http.HandlerFunc(f))
// }

func TestNew(t *testing.T) {
	var expectedError error

	t.Log("Testing correct error return when no API config supplied")
	expectedError = errors.New("Invalid API configuration (nil)")
	_, err := New(nil)
	if err == nil {
		t.Error("Expected an error")
	}
	if err.Error() != expectedError.Error() {
		t.Errorf("Expected an '%#v' error, got '%#v'", expectedError, err)
	}

	t.Log("Testing correct error return when no API Token supplied")
	expectedError = errors.New("API Token is required")
	ac := &Config{}
	_, err = New(ac)
	if err == nil {
		t.Error("Expected an error")
	}
	if err.Error() != expectedError.Error() {
		t.Errorf("Expected an '%#v' error, got '%#v'", expectedError, err)
	}

	t.Log("Testing correct return when an API Token is supplied")
	ac = &Config{
		TokenKey: "abc123",
	}
	_, err = New(ac)
	if err != nil {
		t.Errorf("Expected no error, got '%v'", err)
	}

	t.Log("Testing correct return when an API Token and App are supplied")
	ac = &Config{
		TokenKey: "abc123",
		TokenApp: "someapp",
	}
	_, err = New(ac)
	if err != nil {
		t.Errorf("Expected no error, got '%v'", err)
	}

	t.Log("Testing correct return when an API Token, App, and URL (host) are supplied")
	ac = &Config{
		TokenKey: "abc123",
		TokenApp: "someapp",
		URL:      "something.somewhere.com",
	}
	_, err = New(ac)
	if err != nil {
		t.Errorf("Expected no error, got '%v'", err)
	}

	t.Log("Testing correct return when an API Token, App, and URL (w/trailing '/') are supplied")
	ac = &Config{
		TokenKey: "abc123",
		TokenApp: "someapp",
		URL:      "something.somewhere.com/somepath/",
	}
	_, err = New(ac)
	if err != nil {
		t.Errorf("Expected no error, got '%v'", err)
	}

	t.Log("Testing correct return when an API Token, App, and [invalid] URL are supplied")
	expectedError = errors.New("parse http://something.somewhere.com\\somepath$: invalid character \"\\\\\" in host name")
	ac = &Config{
		TokenKey: "abc123",
		TokenApp: "someapp",
		URL:      "http://something.somewhere.com\\somepath$",
	}
	_, err = New(ac)
	if err == nil {
		t.Error("Expected an error")
	}
	if err.Error() != expectedError.Error() {
		t.Errorf("Expected an '%#v' error, got '%#v'", expectedError, err)
	}

	t.Log("Testing correct return when an Debug true but no log.Logger are supplied")
	ac = &Config{
		TokenKey: "abc123",
		Debug:    true,
	}
	_, err = New(ac)
	if err != nil {
		t.Errorf("Expected no error, got '%v'", err)
	}
}

func TestApiCall(t *testing.T) {
	server := callServer()
	defer server.Close()

	ac := &Config{
		TokenKey: "foo",
		TokenApp: "bar",
		URL:      server.URL,
	}

	apih, err := NewAPI(ac)
	if err != nil {
		t.Errorf("Expected no error, got '%+v'", err)
	}

	t.Log("Testing invalid URL path")
	{
		_, err := apih.apiCall("GET", "", nil)
		expectedError := errors.New("Invalid URL path")
		if err == nil {
			t.Errorf("Expected error")
		}
		if err.Error() != expectedError.Error() {
			t.Errorf("Expected %+v go '%+v'", expectedError, err)
		}
	}

	t.Log("Testing URL path fixup, prefix '/'")
	{
		call := "GET"
		resp, err := apih.apiCall(call, "nothing", nil)
		if err != nil {
			t.Errorf("Expected no error, got '%+v'", resp)
		}
		expected := fmt.Sprintf("%s\n", call)
		if string(resp) != expected {
			t.Errorf("Expected\n'%s'\ngot\n'%s'\n", expected, resp)
		}
	}

	t.Log("Testing URL path fixup, remove '/v2' prefix")
	{
		call := "GET"
		resp, err := apih.apiCall(call, "/v2/nothing", nil)
		if err != nil {
			t.Errorf("Expected no error, got '%+v'", resp)
		}
		expected := fmt.Sprintf("%s\n", call)
		if string(resp) != expected {
			t.Errorf("Expected\n'%s'\ngot\n'%s'\n", expected, resp)
		}
	}

	calls := []string{"GET", "PUT", "POST", "DELETE"}
	for _, call := range calls {
		t.Logf("Testing %s call", call)
		resp, err := apih.apiCall(call, "/", nil)
		if err != nil {
			t.Errorf("Expected no error, got '%+v'", resp)
		}

		expected := fmt.Sprintf("%s\n", call)
		if string(resp) != expected {
			t.Errorf("Expected\n'%s'\ngot\n'%s'\n", expected, resp)
		}
	}

}

func TestApiGet(t *testing.T) {
	server := callServer()
	defer server.Close()

	ac := &Config{
		TokenKey: "foo",
		TokenApp: "bar",
		URL:      server.URL,
	}

	apih, err := NewAPI(ac)
	if err != nil {
		t.Errorf("Expected no error, got '%+v'", err)
	}

	resp, err := apih.Get("/")

	if err != nil {
		t.Errorf("Expected no error, got '%+v'", resp)
	}

	expected := "GET\n"
	if string(resp) != expected {
		t.Errorf("Expected\n'%s'\ngot\n'%s'\n", expected, resp)
	}

}

func TestApiPut(t *testing.T) {
	server := callServer()
	defer server.Close()

	ac := &Config{
		TokenKey: "foo",
		TokenApp: "bar",
		URL:      server.URL,
	}

	apih, err := NewAPI(ac)
	if err != nil {
		t.Errorf("Expected no error, got '%+v'", err)
	}

	resp, err := apih.Put("/", nil)

	if err != nil {
		t.Errorf("Expected no error, got '%+v'", resp)
	}

	expected := "PUT\n"
	if string(resp) != expected {
		t.Errorf("Expected\n'%s'\ngot\n'%s'\n", expected, resp)
	}

}

func TestApiPost(t *testing.T) {
	server := callServer()
	defer server.Close()

	ac := &Config{
		TokenKey: "foo",
		TokenApp: "bar",
		URL:      server.URL,
	}

	apih, err := NewAPI(ac)
	if err != nil {
		t.Errorf("Expected no error, got '%+v'", err)
	}

	resp, err := apih.Post("/", nil)

	if err != nil {
		t.Errorf("Expected no error, got '%+v'", resp)
	}

	expected := "POST\n"
	if string(resp) != expected {
		t.Errorf("Expected\n'%s'\ngot\n'%s'\n", expected, resp)
	}

}

func TestApiDelete(t *testing.T) {
	server := callServer()
	defer server.Close()

	ac := &Config{
		TokenKey: "foo",
		TokenApp: "bar",
		URL:      server.URL,
	}

	apih, err := NewAPI(ac)
	if err != nil {
		t.Errorf("Expected no error, got '%+v'", err)
	}

	resp, err := apih.Delete("/")

	if err != nil {
		t.Errorf("Expected no error, got '%+v'", resp)
	}

	expected := "DELETE\n"
	if string(resp) != expected {
		t.Errorf("Expected\n'%s'\ngot\n'%s'\n", expected, resp)
	}

}
