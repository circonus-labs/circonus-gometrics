// Copyright 2016 Circonus, Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package checkmgr

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/circonus-labs/circonus-gometrics/api"
)

var (
	testCheckBundle = api.CheckBundle{
		CheckUUIDs:         []string{"abc123-a1b2-c3d4-e5f6-123abc"},
		Checks:             []string{"/check/1234"},
		CID:                "/check_bundle/1234",
		Created:            0,
		LastModified:       0,
		LastModifedBy:      "",
		ReverseConnectURLs: []string{""},
		Brokers:            []string{"/broker/1234"},
		Config:             api.CheckBundleConfig{},
		DisplayName:        "test check",
		Metrics: []api.CheckBundleMetric{
			api.CheckBundleMetric{
				Name:   "elmo",
				Type:   "numeric",
				Status: "active",
			},
		},
		MetricLimit: 0,
		Notes:       "",
		Period:      60,
		Status:      "active",
		Target:      "127.0.0.1",
		Timeout:     10,
		Type:        "httptrap",
		Tags:        []string{},
	}
)

func testCheckBundleServer() *httptest.Server {
	f := func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/check_bundle/1234": // handle GET/PUT/DELETE
			switch r.Method {
			case "PUT": // update
				defer r.Body.Close()
				b, err := ioutil.ReadAll(r.Body)
				if err != nil {
					panic(err)
				}
				w.WriteHeader(200)
				w.Header().Set("Content-Type", "application/json")
				fmt.Fprintln(w, string(b))
			case "GET": // get by id/cid
				ret, err := json.Marshal(testCheckBundle)
				if err != nil {
					panic(err)
				}
				w.WriteHeader(200)
				w.Header().Set("Content-Type", "application/json")
				fmt.Fprintln(w, string(ret))
			case "DELETE": // delete
				w.WriteHeader(200)
				fmt.Fprintln(w, "")
			default:
				w.WriteHeader(500)
				fmt.Fprintln(w, "unsupported method")
			}
		case "/check_bundle":
			switch r.Method {
			case "GET": // search
				r := []api.CheckBundle{testCheckBundle}
				ret, err := json.Marshal(r)
				if err != nil {
					panic(err)
				}
				w.WriteHeader(200)
				w.Header().Set("Content-Type", "application/json")
				fmt.Fprintln(w, string(ret))
			case "POST": // create
				defer r.Body.Close()
				b, err := ioutil.ReadAll(r.Body)
				if err != nil {
					panic(err)
				}
				w.WriteHeader(200)
				w.Header().Set("Content-Type", "application/json")
				fmt.Fprintln(w, string(b))
			default:
				w.WriteHeader(500)
				fmt.Fprintln(w, "unsupported method")
			}
		default:
			msg := fmt.Sprintf("unsupported path %s", r.URL.Path)
			w.WriteHeader(500)
			fmt.Fprintln(w, msg)
		}
	}

	return httptest.NewServer(http.HandlerFunc(f))
}

func TestUpdateCheck(t *testing.T) {
	server := testCheckBundleServer()
	defer server.Close()

	cm := &CheckManager{
		enabled: true,
	}
	ac := &api.Config{
		TokenApp: "abcd",
		TokenKey: "1234",
		URL:      server.URL,
	}
	apih, err := api.NewAPI(ac)
	if err != nil {
		t.Errorf("Expected no error, got '%v'", err)
	}
	cm.apih = apih

	newMetrics := make(map[string]*api.CheckBundleMetric)

	t.Log("check manager disabled")
	{
		cm.enabled = false
		cm.UpdateCheck(newMetrics)
	}

	t.Log("no check bundle")
	{
		cm.enabled = true
		cm.checkBundle = nil
		cm.UpdateCheck(newMetrics)
	}

	t.Log("nothing to update (!force metrics, 0 metrics, 0 tags)")
	{
		cm.enabled = true
		cm.checkBundle = &testCheckBundle
		cm.forceCheckUpdate = false
		cm.UpdateCheck(newMetrics)
	}

	newMetrics["test`metric"] = &api.CheckBundleMetric{
		Name:   "test`metric",
		Type:   "numeric",
		Status: "active",
	}

	t.Log("new metric")
	{
		cm.enabled = true
		cm.checkBundle = &testCheckBundle
		cm.forceCheckUpdate = false
		cm.UpdateCheck(newMetrics)
	}

	cm.metricTags = make(map[string][]string)
	cm.metricTags["elmo"] = []string{"cat:tag"}

	t.Log("metric tag")
	{
		cm.enabled = true
		cm.checkBundle = &testCheckBundle
		cm.forceCheckUpdate = false
		cm.UpdateCheck(newMetrics)
	}

}
