// Copyright 2016 Circonus, Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package api

import (
	"errors"
	"os"
	"strings"
	"testing"
)

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

	if _, err := apih.Get("/user/current"); err != nil {
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
