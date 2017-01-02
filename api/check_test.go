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
		path := r.URL.Path
		if path == "/check/1234" {
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
				w.WriteHeader(404)
				fmt.Fprintln(w, fmt.Sprintf("not found: %s %s", r.Method, path))
			}
		} else if path == "/check" {
			switch r.Method {
			case "GET": // search or filter
				var c []Check
				if strings.Contains(r.URL.String(), "f_name=test") {
					c = []Check{}
				} else if strings.Contains(r.URL.String(), "search=test") {
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

func TestFetchCheck(t *testing.T) {
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

	t.Log("invalid CID")
	{
		cid := "/invalid"
		expectedError := errors.New("Invalid check CID [/invalid]")
		_, err = apih.FetchCheck(CIDType(&cid))
		if err == nil {
			t.Fatalf("Expected error")
		}
		if err.Error() != expectedError.Error() {
			t.Fatalf("Expected %+v got '%+v'", expectedError, err)
		}
	}

	t.Log("valid CID")
	{
		cid := CIDType(&testCheck.CID)
		check, err := apih.FetchCheck(cid)
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

/*
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
*/

func TestSearchChecks(t *testing.T) {
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

	t.Log("no search, no filter")
	{
		clusters, err := apih.SearchChecks(nil, nil)
		if err != nil {
			t.Fatalf("Expected no error, got '%v'", err)
		}

		actualType := reflect.TypeOf(clusters)
		expectedType := "*[]api.Check"
		if actualType.String() != expectedType {
			t.Fatalf("Expected %s, got %s", expectedType, actualType.String())
		}
	}

	t.Log("search criteria")
	{
		search := SearchQueryType("test")
		clusters, err := apih.SearchChecks(&search, nil)
		if err != nil {
			t.Fatalf("Expected no error, got '%v'", err)
		}

		actualType := reflect.TypeOf(clusters)
		expectedType := "*[]api.Check"
		if actualType.String() != expectedType {
			t.Fatalf("Expected %s, got %s", expectedType, actualType.String())
		}
	}

	t.Log("filter criteria")
	{
		filter := SearchFilterType{"f_name": []string{"test"}}
		clusters, err := apih.SearchChecks(nil, &filter)
		if err != nil {
			t.Fatalf("Expected no error, got '%v'", err)
		}

		actualType := reflect.TypeOf(clusters)
		expectedType := "*[]api.Check"
		if actualType.String() != expectedType {
			t.Fatalf("Expected %s, got %s", expectedType, actualType.String())
		}
	}
}
