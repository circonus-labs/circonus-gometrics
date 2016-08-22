// Copyright 2016 Circonus, Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package checkmgr

import (
	"os"
	"testing"

	"github.com/circonus-labs/circonus-gometrics/api"
)

func TestLoadCertNoToken(t *testing.T) {
	t.Log("Testing cert load w/o API token")

	cm := &CheckManager{}
	cm.enabled = false

	cm.loadCACert()

	if cm.certPool == nil {
		t.Errorf("Expected cert pool to be initialized, still nil.")
	}

	subjs := cm.certPool.Subjects()
	if len(subjs) == 0 {
		t.Errorf("Expected > 0 certs in pool")
	}
}

func TestLoadCertWithToken(t *testing.T) {
	if os.Getenv("CIRCONUS_API_TOKEN") == "" {
		t.Skip("skipping test; $CIRCONUS_API_TOKEN not set")
	}

	t.Log("Testing cert load with API token")

	cm := &CheckManager{}
	ac := &api.Config{}
	ac.TokenKey = os.Getenv("CIRCONUS_API_TOKEN")
	apih, err := api.NewAPI(ac)
	if err != nil {
		t.Errorf("Expected no error, got '%v'", err)
	}
	cm.apih = apih
	cm.enabled = true

	cm.loadCACert()

	if cm.certPool == nil {
		t.Errorf("Expected cert pool to be initialized, still nil.")
	}

	subjs := cm.certPool.Subjects()
	if len(subjs) == 0 {
		t.Errorf("Expected > 0 certs in pool")
	}
}
