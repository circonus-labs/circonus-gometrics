// Copyright 2016 Circonus, Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package api

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func retryServer() *httptest.Server {
	f := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, "blah blah blah")
	}

	return httptest.NewServer(http.HandlerFunc(f))
}

func TestNewAPI(t *testing.T) {
	t.Log("Testing correct error return when no API config supplied")

	expectedError := errors.New("Invalid API configuration (nil)")

	_, err := NewAPI(nil)

	if err == nil || err.Error() != expectedError.Error() {
		t.Errorf("Expected an '%#v' error, got '%#v'", expectedError, err)
	}

	t.Log("Testing correct error return when no API Token supplied")

	expectedError = errors.New("API Token is required")

	ac := &Config{}
	ac.TokenKey = ""
	_, err = NewAPI(ac)

	if err == nil || err.Error() != expectedError.Error() {
		t.Errorf("Expected an '%#v' error, got '%#v'", expectedError, err)
	}

	t.Log("Testing correct error return when INVALID API Token supplied")

	ac = &Config{}
	ac.TokenKey = "abc-123"
	apih, err := NewAPI(ac)
	if err != nil {
		t.Errorf("Expected no error, got '%v'", err)
	}

	_, err = apih.Get("/user/current")
	if err == nil || !strings.Contains(err.Error(), "The authentication token you supplied is invalid") {
		t.Errorf("Expected an error containing 'The authentication token you supplied is invalid' error, got '%#v'", err)
	}
}

func TestApiCall(t *testing.T) {
	server := retryServer()
	defer server.Close()

	t.Log("Testing correct error return when API call fails retries")

	ac := &Config{}
	ac.TokenKey = "foo"
	ac.TokenApp = "bar"
	ac.URL = server.URL

	apih, err := NewAPI(ac)
	if err != nil {
		t.Errorf("Expected no error, got '%+v'", err)
	}

	resp, err := apih.apiCall("GET", "/retrytest", nil)
	if err == nil {
		t.Errorf("Expected error, got '%+v'", resp)
	}

	expected := fmt.Sprintf("[ERROR] fetching: GET %s/retrytest giving up after 4 attempts - last HTTP error: 500 blah blah blah\n", server.URL)
	if err.Error() != expected {
		t.Errorf("Expected\n'%s'\ngot\n'%s'\n", expected, err)
	}
}

func TestApiGet(t *testing.T) {
	if os.Getenv("CIRCONUS_API_TOKEN") == "" {
		t.Skip("skipping test; $CIRCONUS_API_TOKEN not set")
	}

	t.Log("Testing correct API call to /user/current [defaults]")

	ac := &Config{}
	ac.TokenKey = os.Getenv("CIRCONUS_API_TOKEN")
	ac.TokenApp = os.Getenv("CIRCONUS_API_APP")

	apih, err := NewAPI(ac)
	if err != nil {
		t.Errorf("Expected no error, got '%v'", err)
	}

	if _, err = apih.Get("/user/current"); err != nil {
		t.Errorf("Expected no error, got '%v'", err)
	}

	t.Log("Testing correct API call to /user/current [url=hostname]")

	ac = &Config{}
	ac.TokenKey = os.Getenv("CIRCONUS_API_TOKEN")
	ac.TokenApp = os.Getenv("CIRCONUS_API_APP")
	ac.URL = "api.circonus.com"
	apih, err = NewAPI(ac)
	if err != nil {
		t.Errorf("Expected no error, got '%v'", err)
	}

	if _, err := apih.Get("/user/current"); err != nil {
		t.Errorf("Expected no error, got '%v'", err)
	}

}
