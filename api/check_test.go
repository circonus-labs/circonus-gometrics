// Copyright 2016 Circonus, Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package api

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"testing"
)

// Implicit tests:
//
// FetchCheckByCID is called by FetchCheckByID
// CheckSearch is called by FetchCheckBySubmissionURL

func TestFetchCheckByID1(t *testing.T) {
	if os.Getenv("CIRCONUS_API_TOKEN") == "" {
		t.Skip("skipping test; $CIRCONUS_API_TOKEN not set")
	}

	if os.Getenv("CIRC_API_TEST_CHECK_ID") == "" {
		t.Skip("skipping test; $CIRC_API_TEST_CHECK_ID not set")
	}

	t.Log("Testing correct return from API call")
	ac := &Config{}
	ac.TokenKey = os.Getenv("CIRCONUS_API_TOKEN")
	apih, err := NewAPI(ac)
	if err != nil {
		t.Errorf("Expected no error, got '%v'", err)
	}

	cid := os.Getenv("CIRC_API_TEST_CHECK_ID")
	if cid == "" {
		t.Fatal("Invalid check id (empty)")
	}

	id, err := strconv.Atoi(cid)
	if err != nil {
		t.Fatalf("Unable to convert check id %s to int", cid)
	}

	checkID := IDType(id)

	check, err := apih.FetchCheckByID(checkID)
	if err != nil {
		t.Fatalf("Expected no error, got '%v'", err)
	}

	actualType := reflect.TypeOf(check)
	expectedType := "*api.Check"
	if actualType.String() != expectedType {
		t.Errorf("Expected %s, got %s", expectedType, actualType.String())
	}

	expectedCid := fmt.Sprintf("/check/%s", strconv.Itoa(int(checkID)))
	if check.Cid != expectedCid {
		t.Fatalf("%s != %s", check.Cid, expectedCid)
	}

	t.Logf("Check returned %s %s", check.CheckUUID, check.Cid)

}

func TestFetchCheckByID2(t *testing.T) {
	if os.Getenv("CIRCONUS_API_TOKEN") == "" {
		t.Skip("skipping test; $CIRCONUS_API_TOKEN not set")
	}

	if os.Getenv("CIRCONUS_DELETED_CHECK_ID") == "" {
		t.Skip("skipping test; $CIRCONUS_DELETED_CHECK_ID not set")
	}

	t.Log("Testing correct return from API call with DELETED check ID")

	ac := &Config{}
	ac.TokenKey = os.Getenv("CIRCONUS_API_TOKEN")
	apih, err := NewAPI(ac)
	if err != nil {
		t.Errorf("Expected no error, got '%v'", err)
	}

	cid := os.Getenv("CIRCONUS_DELETED_CHECK_ID")
	if cid == "" {
		t.Fatal("Invalid check id (empty)")
	}

	id, err := strconv.Atoi(cid)
	if err != nil {
		t.Fatalf("Unable to convert check id %s to int", cid)
	}

	checkID := IDType(id)

	check, err := apih.FetchCheckByID(checkID)
	if err != nil {
		t.Fatalf("Expected no error, got '%s'", err)
	}
	if check.Active {
		t.Fatalf("Expected non-active check, got '%#v'", check)
	}
}

func TestFetchCheckBySubmissionURL(t *testing.T) {
	if os.Getenv("CIRCONUS_API_TOKEN") == "" {
		t.Skip("skipping test; $CIRCONUS_API_TOKEN not set")
	}
	if os.Getenv("CIRC_API_TEST_CHECK_SUBMISSION_URL") == "" {
		t.Skip("skipping test; $CIRC_API_TEST_CHECK_SUBMISSION_URL not set")
	}

	t.Log("Testing correct return from API call")

	ac := &Config{}
	ac.TokenKey = os.Getenv("CIRCONUS_API_TOKEN")
	apih, err := NewAPI(ac)
	if err != nil {
		t.Errorf("Expected no error, got '%v'", err)
	}

	url := os.Getenv("CIRC_API_TEST_CHECK_SUBMISSION_URL")
	if url == "" {
		t.Fatal("Invalid check submission url (empty)")
	}

	submissionURL := URLType(url)

	check, err := apih.FetchCheckBySubmissionURL(submissionURL)
	if err != nil {
		t.Fatalf("Expected no error, got '%v'", err)
	}

	actualType := reflect.TypeOf(check)
	expectedType := "*api.Check"
	if actualType.String() != expectedType {
		t.Errorf("Expected %s, got %s", expectedType, actualType.String())
	}

	expectedURL := string(submissionURL)
	if check.Details.SubmissionURL != expectedURL {
		t.Fatalf("%s != %s", check.Details.SubmissionURL, expectedURL)
	}

	t.Logf("Check returned %s %s", check.Cid, check.Details.SubmissionURL)
}
