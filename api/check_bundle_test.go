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
// FetchCheckBundleByCID is called by FetchCheckBundleByID
//
// Excluded tests:
//
// CreateCheckBundle and UpdateCheckBundle because
// both calls change an account (creating and updating)
// checks/check bundles they are tested in checkmgr.

func TestFetchCheckBundleByID(t *testing.T) {
	if os.Getenv("CIRCONUS_API_TOKEN") == "" {
		t.Skip("skipping test; $CIRCONUS_API_TOKEN not set")
	}
	if os.Getenv("CIRC_API_TEST_CHECK_BUNDLE_ID") == "" {
		t.Skip("skipping test; $CIRC_API_TEST_CHECK_BUNDLE_ID not set")
	}

	t.Log("Testing correct return from API call")

	ac := &Config{}
	ac.TokenKey = os.Getenv("CIRCONUS_API_TOKEN")
	apih, err := NewAPI(ac)
	if err != nil {
		t.Errorf("Expected no error, got '%v'", err)
	}

	cid := os.Getenv("CIRC_API_TEST_CHECK_BUNDLE_ID")
	if cid == "" {
		t.Fatal("Invalid check bundle id (empty)")
	}

	id, err := strconv.Atoi(cid)
	if err != nil {
		t.Fatalf("Unable to convert check bundle id %s to int", cid)
	}

	checkBundleID := IDType(id)

	checkBundle, err := apih.FetchCheckBundleByID(checkBundleID)
	if err != nil {
		t.Fatalf("Expected no error, got '%v'", err)
	}

	actualType := reflect.TypeOf(checkBundle)
	expectedType := "*api.CheckBundle"
	if actualType.String() != expectedType {
		t.Errorf("Expected %s, got %s", expectedType, actualType.String())
	}

	expectedCid := fmt.Sprintf("/check_bundle/%s", strconv.Itoa(int(checkBundleID)))
	if checkBundle.Cid != expectedCid {
		t.Fatalf("%s != %s", checkBundle.Cid, expectedCid)
	}

	t.Logf("Check bundle returned %s %s", checkBundle.DisplayName, checkBundle.Cid)
}
