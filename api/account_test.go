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
	testAccount = Account{
		CID: "/account/1234",
		ContactGroups: []string{
			"/contact_group/1701",
			"/contact_group/3141",
		},
		OwnerCID: "/user/42",
		Usage: []AccountLimit{
			AccountLimit{
				Limit: 50,
				Type:  "Host",
				Used:  7,
			},
		},
		Address1:    "Hooper's Store",
		Address2:    "Sesame Street",
		CCEmail:     "accounts_payable@yourdomain.com",
		City:        "New York City",
		Country:     "US",
		Description: "Hooper's Store Account",
		Invites: []AccountInvite{
			AccountInvite{
				Email: "alan@example.com",
				Role:  "Admin",
			},
			AccountInvite{
				Email: "chris.robinson@example.com",
				Role:  "Normal",
			},
		},
		Name:      "hoopers-store",
		StateProv: "NY",
		Timezone:  "America/New_York",
		Users: []AccountUser{
			AccountUser{
				Role:    "Admin",
				UserCID: "/user/42",
			},
		},
	}
)

func testAccountServer() *httptest.Server {
	f := func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if path == "/account/1234" || path == "/account/current" {
			switch r.Method {
			case "GET": // get by id/cid
				ret, err := json.Marshal(testAccount)
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
		} else if path == "/account" {
			switch r.Method {
			case "GET":
				c := []Account{testAccount}
				ret, err := json.Marshal(c)
				if err != nil {
					panic(err)
				}
				w.WriteHeader(200)
				w.Header().Set("Content-Type", "application/json")
				fmt.Fprintln(w, string(ret))
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

func TestFetchAccount(t *testing.T) {
	server := testAccountServer()
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

	t.Log("without CID (nil)")
	{
		account, err := apih.FetchAccount(nil)
		if err != nil {
			t.Fatalf("Expected no error, got '%v'", err)
		}

		actualType := reflect.TypeOf(account)
		expectedType := "*api.Account"
		if actualType.String() != expectedType {
			t.Fatalf("Expected %s, got %s", expectedType, actualType.String())
		}
	}

	t.Log("without CID (\"\")")
	{
		cid := CIDType("")
		account, err := apih.FetchAccount(&cid)
		if err != nil {
			t.Fatalf("Expected no error, got '%v'", err)
		}

		actualType := reflect.TypeOf(account)
		expectedType := "*api.Account"
		if actualType.String() != expectedType {
			t.Fatalf("Expected %s, got %s", expectedType, actualType.String())
		}
	}

	t.Log("with valid CID")
	{
		cid := CIDType("/account/1234")
		account, err := apih.FetchAccount(&cid)
		if err != nil {
			t.Fatalf("Expected no error, got '%v'", err)
		}

		actualType := reflect.TypeOf(account)
		expectedType := "*api.Account"
		if actualType.String() != expectedType {
			t.Fatalf("Expected %s, got %s", expectedType, actualType.String())
		}

		if account.CID != testAccount.CID {
			t.Fatalf("CIDs do not match: %+v != %+v\n", account, testAccount)
		}
	}

	t.Log("with invalid CID")
	{
		cid := CIDType("/invalid")
		expectedError := errors.New("Invalid account CID [/invalid]")
		_, err := apih.FetchAccount(&cid)
		if err == nil {
			t.Fatalf("Expected error")
		}
		if err.Error() != expectedError.Error() {
			t.Fatalf("Expected %+v got '%+v'", expectedError, err)
		}
	}
}

func TestUpdateAccount(t *testing.T) {
	server := testAccountServer()
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

	t.Log("valid account")
	{
		account, err := apih.UpdateAccount(&testAccount)
		if err != nil {
			t.Fatalf("Expected no error, got '%v'", err)
		}

		actualType := reflect.TypeOf(account)
		expectedType := "*api.Account"
		if actualType.String() != expectedType {
			t.Fatalf("Expected %s, got %s", expectedType, actualType.String())
		}
	}

	t.Log("Test with invalid CID")
	{
		expectedError := errors.New("Invalid account CID [/invalid]")
		x := &Account{CID: "/invalid"}
		_, err := apih.UpdateAccount(x)
		if err == nil {
			t.Fatal("Expected an error")
		}
		if err.Error() != expectedError.Error() {
			t.Fatalf("Expected %+v got '%+v'", expectedError, err)
		}
	}
}
