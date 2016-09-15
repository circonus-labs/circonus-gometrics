// Copyright 2016 Circonus, Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package api

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
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

	if os.Getenv("CIRC_API_TEST_DELETED_CHECK_ID") == "" {
		t.Skip("skipping test; $CIRC_API_TEST_DELETED_CHECK_ID not set")
	}

	t.Log("Testing correct return from API call with DELETED check ID")

	ac := &Config{}
	ac.TokenKey = os.Getenv("CIRCONUS_API_TOKEN")
	apih, err := NewAPI(ac)
	if err != nil {
		t.Errorf("Expected no error, got '%v'", err)
	}

	cid := os.Getenv("CIRC_API_TEST_DELETED_CHECK_ID")
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

func TestCheckSearch1(t *testing.T) {
	if os.Getenv("CIRCONUS_API_TOKEN") == "" {
		t.Skip("skipping test; $CIRCONUS_API_TOKEN not set")
	}
	if os.Getenv("CIRC_API_TEST_CHECK_TAG") == "" {
		t.Skip("skipping test; $CIRC_API_TEST_CHECK_TAG not set")
	}

	t.Log("Testing correct return from API call")

	ac := &Config{}
	ac.TokenKey = os.Getenv("CIRCONUS_API_TOKEN")

	apih, err := NewAPI(ac)
	if err != nil {
		t.Errorf("Expected no error, got '%v'", err)
	}

	tag := os.Getenv("CIRC_API_TEST_CHECK_TAG")
	if tag == "" {
		t.Fatal("Invalid check search tag (empty)")
	}

	searchQuery := SearchQueryType(fmt.Sprintf("(active:1)(tags:%s)", tag))

	checks, err := apih.CheckSearch(searchQuery)
	if err != nil {
		t.Fatalf("Expected no error, got '%v'", err)
	}

	actualType := reflect.TypeOf(checks)
	expectedType := "[]api.Check"
	if actualType.String() != expectedType {
		t.Errorf("Expected %s, got %s", expectedType, actualType.String())
	}

	if len(checks) < 1 {
		t.Fatalf("Expected at least 1 check returned, recieved %d", len(checks))
	}

	t.Logf("%d checks returned", len(checks))

}

func TestCheckSearch2(t *testing.T) {
	if os.Getenv("CIRCONUS_API_TOKEN") == "" {
		t.Skip("skipping test; $CIRCONUS_API_TOKEN not set")
	}
	if os.Getenv("CIRC_API_TEST_CHECK_MULTI_TAG") == "" {
		t.Skip("skipping test; $CIRC_API_TEST_CHECK_MULTI_TAG not set")
	}

	t.Log("Testing correct return from API call")

	ac := &Config{}
	ac.TokenKey = os.Getenv("CIRCONUS_API_TOKEN")

	apih, err := NewAPI(ac)
	if err != nil {
		t.Errorf("Expected no error, got '%v'", err)
	}

	tags := strings.Split(strings.Replace(os.Getenv("CIRC_API_TEST_CHECK_MULTI_TAG"), " ", "", -1), ",")
	if len(tags) == 0 {
		t.Fatal("Invalid check search tags (empty)")
	}

	searchQuery := SearchQueryType(fmt.Sprintf("(active:1)(tags:%s)", strings.Join(tags, ",")))

	checks, err := apih.CheckSearch(searchQuery)
	if err != nil {
		t.Fatalf("Expected no error, got '%v'", err)
	}

	actualType := reflect.TypeOf(checks)
	expectedType := "[]api.Check"
	if actualType.String() != expectedType {
		t.Errorf("Expected %s, got %s", expectedType, actualType.String())
	}

	if len(checks) < 1 {
		t.Fatalf("Expected at least 1 check returned, recieved %d", len(checks))
	}

	t.Logf("%d checks returned", len(checks))

}

func TestCheckFileterSearch(t *testing.T) {
	if os.Getenv("CIRCONUS_API_TOKEN") == "" {
		t.Skip("skipping test; $CIRCONUS_API_TOKEN not set")
	}
	if os.Getenv("CIRC_API_TEST_CHECK_FILTER") == "" {
		t.Skip("skipping test; $CIRC_API_TEST_CHECK_FILTER not set")
	}

	t.Log("Testing correct return from API call")

	ac := &Config{}
	ac.TokenKey = os.Getenv("CIRCONUS_API_TOKEN")

	apih, err := NewAPI(ac)
	if err != nil {
		t.Errorf("Expected no error, got '%v'", err)
	}

	filter := os.Getenv("CIRC_API_TEST_CHECK_FILTER")
	if filter == "" {
		t.Fatal("Invalid check search filter (empty)")
	}

	checks, err := apih.CheckFilterSearch(SearchFilterType(filter))
	if err != nil {
		t.Fatalf("Expected no error, got '%v'", err)
	}

	actualType := reflect.TypeOf(checks)
	expectedType := "[]api.Check"
	if actualType.String() != expectedType {
		t.Errorf("Expected %s, got %s", expectedType, actualType.String())
	}

	if len(checks) < 1 {
		t.Fatalf("Expected at least 1 check returned, recieved %d", len(checks))
	}

	t.Logf("%d checks returned", len(checks))

}
