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
	testUser = User{
		CID:       "/user/1234",
		Email:     "john@example.com",
		Firstname: "John",
		Lastname:  "Doe",
		ContactInfo: UserContactInfo{
			SMS:  "123-456-7890",
			XMPP: "foobar",
		},
	}
)

func testUserServer() *httptest.Server {
	f := func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if path == "/user/1234" || path == "/user/current" {
			switch r.Method {
			case "GET": // get by id/cid
				ret, err := json.Marshal(testUser)
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
			default:
				w.WriteHeader(404)
				fmt.Fprintln(w, "not found")
			}
		} else if path == "/user" {
			switch r.Method {
			case "GET":
				reqURL := r.URL.String()
				var c []User
				if reqURL == "/user?f_firstname=john&f_lastname=doe" {
					c = []User{testUser}
				} else if reqURL == "/user" {
					c = []User{testUser}
				} else {
					c = []User{}
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
			default:
				w.WriteHeader(404)
				fmt.Fprintln(w, "not found")
			}
		} else {
			w.WriteHeader(404)
			fmt.Fprintln(w, "not found")
		}
	}

	return httptest.NewServer(http.HandlerFunc(f))
}

func TestFetchUser(t *testing.T) {
	server := testUserServer()
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
		user, err := apih.FetchUser(CIDType(&cid))
		if err != nil {
			t.Fatalf("Expected no error, got '%v'", err)
		}

		actualType := reflect.TypeOf(user)
		expectedType := "*api.User"
		if actualType.String() != expectedType {
			t.Fatalf("Expected %s, got %s", expectedType, actualType.String())
		}
	}

	t.Log("with valid CID")
	{
		cid := "/user/1234"
		user, err := apih.FetchUser(CIDType(&cid))
		if err != nil {
			t.Fatalf("Expected no error, got '%v'", err)
		}

		actualType := reflect.TypeOf(user)
		expectedType := "*api.User"
		if actualType.String() != expectedType {
			t.Fatalf("Expected %s, got %s", expectedType, actualType.String())
		}

		if user.CID != testUser.CID {
			t.Fatalf("CIDs do not match: %+v != %+v\n", user, testUser)
		}
	}

	t.Log("with invalid CID")
	{
		cid := "/invalid"
		expectedError := errors.New("Invalid user CID [/invalid]")
		_, err := apih.FetchUser(CIDType(&cid))
		if err == nil {
			t.Fatalf("Expected error")
		}
		if err.Error() != expectedError.Error() {
			t.Fatalf("Expected %+v got '%+v'", expectedError, err)
		}
	}
}

func TestFetchUsers(t *testing.T) {
	server := testUserServer()
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

	users, err := apih.FetchUsers()
	if err != nil {
		t.Fatalf("Expected no error, got '%v'", err)
	}

	actualType := reflect.TypeOf(users)
	expectedType := "*[]api.User"
	if actualType.String() != expectedType {
		t.Fatalf("Expected %s, got %s", expectedType, actualType.String())
	}

}

func TestUpdateUser(t *testing.T) {
	server := testUserServer()
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

	t.Log("valid User")
	{
		User, err := apih.UpdateUser(&testUser)
		if err != nil {
			t.Fatalf("Expected no error, got '%v'", err)
		}

		actualType := reflect.TypeOf(User)
		expectedType := "*api.User"
		if actualType.String() != expectedType {
			t.Fatalf("Expected %s, got %s", expectedType, actualType.String())
		}
	}

	t.Log("Test with invalid CID")
	{
		expectedError := errors.New("Invalid user CID [/invalid]")
		x := &User{CID: "/invalid"}
		_, err := apih.UpdateUser(x)
		if err == nil {
			t.Fatal("Expected an error")
		}
		if err.Error() != expectedError.Error() {
			t.Fatalf("Expected %+v got '%+v'", expectedError, err)
		}
	}
}

func TestSearchUsers(t *testing.T) {
	server := testUserServer()
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

	filter := SearchFilterType(map[string][]string{"f_firstname": []string{"john"}, "f_lastname": []string{"doe"}})

	t.Log("no filter")
	{
		users, err := apih.SearchUsers(nil)
		if err != nil {
			t.Fatalf("Expected no error, got '%v'", err)
		}

		actualType := reflect.TypeOf(users)
		expectedType := "*[]api.User"
		if actualType.String() != expectedType {
			t.Fatalf("Expected %s, got %s", expectedType, actualType.String())
		}
	}

	t.Log("filter")
	{
		users, err := apih.SearchUsers(&filter)
		if err != nil {
			t.Fatalf("Expected no error, got '%v'", err)
		}

		actualType := reflect.TypeOf(users)
		expectedType := "*[]api.User"
		if actualType.String() != expectedType {
			t.Fatalf("Expected %s, got %s", expectedType, actualType.String())
		}
	}
}
