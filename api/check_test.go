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
	testCheck = Check{
		CID:            "/check/1234",
		Active:         true,
		BrokerCID:      "/broker/1234",
		CheckBundleCID: "/check_bundle/1234",
		CheckUUID:      "abc123-a1b2-c3d4-e5f6-123abc",
		Details:        CheckDetails{SubmissionURL: "http://example.com/"},
	}
)

func testCheckServer() *httptest.Server {
	f := func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/check/1234": // handle GET/PUT/DELETE
			switch r.Method {
			case "GET": // get by id/cid
				ret, err := json.Marshal(testCheck)
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
		case "/check":
			switch r.Method {
			case "GET": // search or filter
				var c []Check
				if strings.Contains(r.URL.String(), "f__check_uuid=none") {
					c = []Check{}
				} else if strings.Contains(r.URL.String(), "f__check_uuid=multi") {
					c = []Check{testCheck, testCheck}
				} else {
					c = []Check{testCheck}
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

func TestFetchCheckByID(t *testing.T) {
	server := testCheckServer()
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

	checkID := IDType(id)

	check, err := apih.FetchCheckByID(checkID)
	if err != nil {
		t.Fatalf("Expected no error, got '%v'", err)
	}

	actualType := reflect.TypeOf(check)
	expectedType := "*api.Check"
	if actualType.String() != expectedType {
		t.Fatalf("Expected %s, got %s", expectedType, actualType.String())
	}

	if check.CID != testCheck.CID {
		t.Fatalf("CIDs do not match: %+v != %+v\n", check, testCheck)
	}
}

func TestFetchCheckByCID(t *testing.T) {
	server := testCheckServer()
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

	t.Log("Testing invalid CID")
	expectedError := errors.New("Invalid check CID /1234")
	_, err = apih.FetchCheckByCID("/1234")
	if err == nil {
		t.Fatalf("Expected error")
	}
	if err.Error() != expectedError.Error() {
		t.Fatalf("Expected %+v got '%+v'", expectedError, err)
	}

	t.Log("Testing valid CID")
	check, err := apih.FetchCheckByCID(CIDType(testCheck.CID))
	if err != nil {
		t.Fatalf("Expected no error, got '%v'", err)
	}

	actualType := reflect.TypeOf(check)
	expectedType := "*api.Check"
	if actualType.String() != expectedType {
		t.Fatalf("Expected %s, got %s", expectedType, actualType.String())
	}

	if check.CID != testCheck.CID {
		t.Fatalf("CIDs do not match: %+v != %+v\n", check, testCheck)
	}
}

func TestFetchCheckBySubmissionURL(t *testing.T) {
	server := testCheckServer()
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

	t.Log("Testing invalid URL (blank)")
	{
		expectedError := errors.New("[ERROR] Invalid submission URL (blank)")
		_, err = apih.FetchCheckBySubmissionURL("")
		if err == nil {
			t.Fatalf("Expected error")
		}
		if err.Error() != expectedError.Error() {
			t.Fatalf("Expected %+v got '%+v'", expectedError, err)
		}
	}

	t.Log("Testing invalid URL (bad format)")
	{
		expectedError := errors.New("parse http://example.com\\noplace$: invalid character \"\\\\\" in host name")
		_, err = apih.FetchCheckBySubmissionURL(URLType("http://example.com\\noplace$"))
		if err == nil {
			t.Fatalf("Expected error")
		}
		if err.Error() != expectedError.Error() {
			t.Fatalf("Expected %+v got '%+v'", expectedError, err)
		}
	}

	t.Log("Testing invalid URL (bad path)")
	{
		expectedError := errors.New("[ERROR] Invalid submission URL 'http://example.com/foo', unrecognized path")
		_, err = apih.FetchCheckBySubmissionURL(URLType("http://example.com/foo"))
		if err == nil {
			t.Fatalf("Expected error")
		}
		if err.Error() != expectedError.Error() {
			t.Fatalf("Expected %+v got '%+v'", expectedError, err)
		}
	}

	t.Log("Testing invalid URL (no uuid)")
	{
		expectedError := errors.New("[ERROR] Invalid submission URL 'http://example.com/module/httptrap/', UUID not where expected")
		_, err = apih.FetchCheckBySubmissionURL(URLType("http://example.com/module/httptrap/"))
		if err == nil {
			t.Fatalf("Expected error")
		}
		if err.Error() != expectedError.Error() {
			t.Fatalf("Expected %+v got '%+v'", expectedError, err)
		}
	}

	t.Log("Testing valid URL (0 checks returned)")
	{
		expectedError := errors.New("[ERROR] No checks found with UUID none")
		_, err := apih.FetchCheckBySubmissionURL(URLType("http://example.com/module/httptrap/none/boo"))
		if err == nil {
			t.Fatalf("Expected error")
		}
		if err.Error() != expectedError.Error() {
			t.Fatalf("Expected %+v got '%+v'", expectedError, err)
		}
	}

	t.Log("Testing valid URL (multiple checks returned)")
	{
		expectedError := errors.New("[ERROR] Multiple checks with same UUID multi")
		_, err := apih.FetchCheckBySubmissionURL(URLType("http://example.com/module/httptrap/multi/boo"))
		if err == nil {
			t.Fatalf("Expected error")
		}
		if err.Error() != expectedError.Error() {
			t.Fatalf("Expected %+v got '%+v'", expectedError, err)
		}
	}

	t.Log("Testing valid URL (1 check returned)")
	{
		check, err := apih.FetchCheckBySubmissionURL(URLType("http://example.com/module/httptrap/abc123-abc1-def2-ghi3-123abc/boo"))
		if err != nil {
			t.Fatalf("Expected no error, got '%v'", err)
		}

		actualType := reflect.TypeOf(check)
		expectedType := "*api.Check"
		if actualType.String() != expectedType {
			t.Fatalf("Expected %s, got %s", expectedType, actualType.String())
		}

		if check.CID != testCheck.CID {
			t.Fatalf("CIDs do not match: %+v != %+v\n", check, testCheck)
		}
	}
}

func TestCheckSearch(t *testing.T) {
	server := testCheckServer()
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
		clusters, err := apih.CheckSearch("")
		if err != nil {
			t.Fatalf("Expected no error, got '%v'", err)
		}

		actualType := reflect.TypeOf(clusters)
		expectedType := "[]api.Check"
		if actualType.String() != expectedType {
			t.Fatalf("Expected %s, got %s", expectedType, actualType.String())
		}
	}

	t.Log("Testing with search criteria")
	{
		clusters, err := apih.CheckSearch("test")
		if err != nil {
			t.Fatalf("Expected no error, got '%v'", err)
		}

		actualType := reflect.TypeOf(clusters)
		expectedType := "[]api.Check"
		if actualType.String() != expectedType {
			t.Fatalf("Expected %s, got %s", expectedType, actualType.String())
		}
	}
}

func TestCheckFilterSearch(t *testing.T) {
	server := testCheckServer()
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
		expectedError := errors.New("[ERROR] invalid filter supplied (blank)")
		_, err := apih.CheckFilterSearch("")
		if err == nil {
			t.Fatalf("Expected error")
		}
		if err.Error() != expectedError.Error() {
			t.Fatalf("Expected %v got '%v'", expectedError, err)
		}
	}

	t.Log("Testing with search criteria")
	{
		clusters, err := apih.CheckFilterSearch("f_name=test")
		if err != nil {
			t.Fatalf("Expected no error, got '%v'", err)
		}

		actualType := reflect.TypeOf(clusters)
		expectedType := "[]api.Check"
		if actualType.String() != expectedType {
			t.Fatalf("Expected %s, got %s", expectedType, actualType.String())
		}
	}
}

// Implicit tests:
//
// FetchCheckByCID is called by FetchCheckByID
// CheckSearch is called by FetchCheckBySubmissionURL
/*
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
*/
