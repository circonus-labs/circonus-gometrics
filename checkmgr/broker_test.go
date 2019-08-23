// Copyright 2016 Circonus, Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package checkmgr

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	apiclient "github.com/circonus-labs/go-apiclient"
)

var (
	invalidBroker = apiclient.Broker{
		CID:       "/broker/1",
		Longitude: nil,
		Latitude:  nil,
		Name:      "test broker",
		Tags:      []string{},
		Type:      "foo",
		Details: []apiclient.BrokerDetail{
			{
				CN:           "testbroker.example.com",
				ExternalHost: &[]string{"testbroker.example.com"}[0],
				ExternalPort: 43191,
				IP:           &[]string{"127.0.0.1"}[0],
				MinVer:       0,
				Modules:      []string{"a", "b", "c"},
				Port:         &[]uint16{43191}[0],
				Skew:         nil,
				Status:       "unprovisioned",
				Version:      nil,
			},
		},
	}

	noIPorHostBroker = apiclient.Broker{
		CID:       "/broker/2",
		Longitude: nil,
		Latitude:  nil,
		Name:      "no ip or external host broker",
		Tags:      []string{},
		Type:      "enterprise",
		Details: []apiclient.BrokerDetail{
			{
				CN:           "foobar",
				ExternalHost: nil,
				ExternalPort: 43191,
				IP:           nil,
				MinVer:       0,
				Modules:      []string{"httptrap"},
				Port:         &[]uint16{43191}[0],
				Skew:         nil,
				Status:       "active",
				Version:      nil,
			},
		},
	}

	validBroker = apiclient.Broker{
		CID:       "/broker/2",
		Longitude: nil,
		Latitude:  nil,
		Name:      "test broker",
		Tags:      []string{},
		Type:      "enterprise",
		Details: []apiclient.BrokerDetail{
			{
				CN:           "testbroker.example.com",
				ExternalHost: nil,
				ExternalPort: 43191,
				IP:           &[]string{"127.0.0.1"}[0],
				MinVer:       0,
				Modules:      []string{"httptrap"},
				Port:         &[]uint16{43191}[0],
				Skew:         nil,
				Status:       "active",
				Version:      nil,
			},
		},
	}

	validBrokerNonEnterprise = apiclient.Broker{
		CID:       "/broker/3",
		Longitude: nil,
		Latitude:  nil,
		Name:      "test broker",
		Tags:      []string{},
		Type:      "foo",
		Details: []apiclient.BrokerDetail{
			{
				CN:           "testbroker.example.com",
				ExternalHost: nil,
				ExternalPort: 43191,
				IP:           &[]string{"127.0.0.1"}[0],
				MinVer:       0,
				Modules:      []string{"httptrap"},
				Port:         &[]uint16{43191}[0],
				Skew:         nil,
				Status:       "active",
				Version:      nil,
			},
		},
	}
)

func testBrokerServer() *httptest.Server {
	f := func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/broker/1":
			switch r.Method {
			case "GET": // get by id/cid
				ret, err := json.Marshal(invalidBroker)
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
		case "/broker/2":
			switch r.Method {
			case "GET": // get by id/cid
				ret, err := json.Marshal(validBroker)
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
		case "/broker":
			switch r.Method {
			case "GET": // search or filter
				var c []apiclient.Broker
				switch {
				case strings.Contains(r.URL.String(), "f__tags_has=no%3Abroker"):
					c = []apiclient.Broker{}
				case strings.Contains(r.URL.String(), "f__tags_has=multi%3Abroker"):
					c = []apiclient.Broker{invalidBroker, invalidBroker}
				default:
					c = []apiclient.Broker{validBroker, validBrokerNonEnterprise}
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

func TestBrokerSupportsCheckType(t *testing.T) {
	detail := &apiclient.BrokerDetail{
		Modules: []string{"httptrap"},
	}

	cm := CheckManager{}

	t.Log("supports 'httptrap' check type?")
	{
		ok := cm.brokerSupportsCheckType("httptrap", detail)
		if !ok {
			t.Fatal("Expected OK")
		}
	}

	t.Log("supports 'foo' check type?")
	{
		ok := cm.brokerSupportsCheckType("foo", detail)
		if ok {
			t.Fatal("Expected not OK")
		}
	}
}

func TestGetBrokerCN(t *testing.T) {

	t.Log("URL with IP")
	{
		submissionURL := apiclient.URLType("http://127.0.0.1:43191/blah/blah/blah")
		cm := CheckManager{}

		_, err := cm.getBrokerCN(&validBroker, submissionURL)
		if err != nil {
			t.Fatalf("Expected no error, got %+v", err)
		}
	}

	t.Log("URL with FQDN")
	{
		submissionURL := apiclient.URLType("http://test.example.com:43191/blah/blah/blah")
		cm := CheckManager{}

		_, err := cm.getBrokerCN(&validBroker, submissionURL)
		if err != nil {
			t.Fatalf("Expected no error, got %+v", err)
		}
	}

	t.Log("URL with invalid IP")
	{
		submissionURL := apiclient.URLType("http://127.0.0.2:43191/blah/blah/blah")
		cm := CheckManager{}

		_, err := cm.getBrokerCN(&validBroker, submissionURL)
		if err == nil {
			t.Fatal("expected error")
		}
		if err.Error() != "error, unable to match URL host (127.0.0.2:43191) to Broker" {
			t.Fatalf("unexpected error (%s)", err)
		}
	}
}

func TestSelectBroker(t *testing.T) {
	server := testBrokerServer()
	defer server.Close()

	testURL, err := url.Parse(server.URL)
	if err != nil {
		t.Fatalf("Error parsing temporary url %v", err)
	}

	hostParts := strings.Split(testURL.Host, ":")
	hostPort, err := strconv.Atoi(hostParts[1])
	if err != nil {
		t.Fatalf("Error converting port to numeric %v", err)
	}

	validBroker.Details[0].ExternalHost = &hostParts[0]
	validBroker.Details[0].ExternalPort = uint16(hostPort)
	validBroker.Details[0].IP = &hostParts[0]
	validBroker.Details[0].Port = &[]uint16{uint16(hostPort)}[0]

	validBrokerNonEnterprise.Details[0].ExternalHost = &hostParts[0]
	validBrokerNonEnterprise.Details[0].ExternalPort = uint16(hostPort)
	validBrokerNonEnterprise.Details[0].IP = &hostParts[0]
	validBrokerNonEnterprise.Details[0].Port = &[]uint16{uint16(hostPort)}[0]

	t.Log("default broker selection")
	{
		cm := &CheckManager{
			checkType:             "httptrap",
			brokerMaxResponseTime: time.Duration(time.Millisecond * 500),
		}
		ac := &apiclient.Config{
			TokenApp: "abcd",
			TokenKey: "1234",
			URL:      server.URL,
		}
		apih, err := apiclient.New(ac)
		if err != nil {
			t.Errorf("Expected no error, got '%v'", err)
		}
		cm.apih = apih

		_, err = cm.selectBroker()
		if err != nil {
			t.Fatal("Expected no error")
		}
	}

	t.Log("tag, no brokers matching")
	{
		cm := &CheckManager{
			checkType:             "httptrap",
			brokerMaxResponseTime: time.Duration(time.Millisecond * 500),
			brokerSelectTag:       apiclient.TagType([]string{"no:broker"}),
		}
		ac := &apiclient.Config{
			TokenApp: "abcd",
			TokenKey: "1234",
			URL:      server.URL,
		}
		apih, err := apiclient.New(ac)
		if err != nil {
			t.Errorf("Expected no error, got '%v'", err)
		}
		cm.apih = apih

		expectedError := errors.New("zero brokers found")

		_, err = cm.selectBroker()
		if err == nil {
			t.Fatal("Expected an error")
		}
		if expectedError.Error() != err.Error() {
			t.Errorf("Expected %v got '%v'", expectedError, err)
		}
	}

	t.Log("multiple brokers with tag, none valid")
	{
		cm := &CheckManager{
			checkType:             "httptrap",
			brokerMaxResponseTime: time.Duration(time.Millisecond * 500),
			brokerSelectTag:       apiclient.TagType([]string{"multi:broker"}),
		}
		ac := &apiclient.Config{
			TokenApp: "abcd",
			TokenKey: "1234",
			URL:      server.URL,
		}
		apih, err := apiclient.NewAPI(ac)
		if err != nil {
			t.Errorf("Expected no error, got '%v'", err)
		}
		cm.apih = apih

		expectedError := errors.New("found 2 broker(s), zero are valid")
		_, err = cm.selectBroker()
		if err == nil {
			t.Fatalf("Expected an error")
		}

		if expectedError.Error() != err.Error() {
			t.Fatalf("Expected %v got '%v'", expectedError, err)
		}
	}

}

func TestIsValidBroker(t *testing.T) {
	cm := &CheckManager{
		Log:                   log.New(os.Stderr, "", log.LstdFlags),
		checkType:             "httptrap",
		brokerMaxResponseTime: time.Duration(time.Millisecond * 50),
	}

	broker := apiclient.Broker{
		CID:  "/broker/2",
		Name: "test broker",
		Type: "enterprise",
		Details: []apiclient.BrokerDetail{
			{
				CN:           "testbroker.example.com",
				ExternalHost: nil,
				ExternalPort: 43191,
				IP:           &[]string{"127.0.0.1"}[0],
				Modules:      []string{"httptrap"},
				Port:         &[]uint16{43191}[0],
				Status:       "unprovisioned",
			},
		},
	}

	t.Log("status unprovisioned")
	{
		if cm.isValidBroker(&broker) {
			t.Fatal("Expected invalid broker")
		}
	}

	t.Log("no ip or host")
	{
		if cm.isValidBroker(&noIPorHostBroker) {
			t.Fatal("Expected invalid broker")
		}
	}

	t.Log("does not have required module")
	{
		broker.Details[0].Modules = []string{"foo"}
		broker.Details[0].Status = "active"
		if cm.isValidBroker(&broker) {
			t.Fatal("Expected invalid broker")
		}
	}
}

func TestIsValidBrokerTimeout(t *testing.T) {
	if os.Getenv("CIRCONUS_BROKER_TEST_TIMEOUT") == "" {
		t.Skip("not testing timeouts, CIRCONUS_BROKER_TEST_TIMEOUT not set")
	}

	cm := &CheckManager{
		Log:                   log.New(os.Stderr, "", log.LstdFlags),
		checkType:             "httptrap",
		brokerMaxResponseTime: time.Duration(time.Millisecond * 50),
	}

	broker := apiclient.Broker{
		CID:  "/broker/2",
		Name: "test broker",
		Type: "enterprise",
		Details: []apiclient.BrokerDetail{
			{
				CN:           "testbroker.example.com",
				ExternalHost: nil,
				ExternalPort: 43191,
				IP:           &[]string{"127.0.0.1"}[0],
				Modules:      []string{"httptrap"},
				Port:         &[]uint16{43191}[0],
				Status:       "unprovisioned",
			},
		},
	}

	t.Log("unable to connect, broker.ExternalPort")
	{
		broker.Name = "test"
		broker.Details[0].Modules = []string{"httptrap"}
		broker.Details[0].Status = "active"
		if cm.isValidBroker(&broker) {
			t.Fatal("Expected invalid broker")
		}
	}

	t.Log("unable to connect, broker.Port")
	{
		broker.Name = "test"
		broker.Details[0].ExternalPort = 0
		broker.Details[0].Modules = []string{"httptrap"}
		broker.Details[0].Status = "active"
		if cm.isValidBroker(&broker) {
			t.Fatal("Expected invalid broker")
		}
	}

	t.Log("unable to connect, default port")
	{
		broker.Name = "test"
		broker.Details[0].ExternalPort = 0
		broker.Details[0].Port = &[]uint16{0}[0]
		broker.Details[0].Modules = []string{"httptrap"}
		broker.Details[0].Status = "active"
		if cm.isValidBroker(&broker) {
			t.Fatal("Expected invalid broker")
		}
	}
}

func TestGetBroker(t *testing.T) {
	server := testBrokerServer()
	defer server.Close()

	testURL, err := url.Parse(server.URL)
	if err != nil {
		t.Fatalf("Error parsing temporary url %v", err)
	}

	hostParts := strings.Split(testURL.Host, ":")
	hostPort, err := strconv.Atoi(hostParts[1])
	if err != nil {
		t.Fatalf("Error converting port to numeric %v", err)
	}

	validBroker.Details[0].ExternalHost = &hostParts[0]
	validBroker.Details[0].ExternalPort = uint16(hostPort)
	validBroker.Details[0].IP = &hostParts[0]
	validBroker.Details[0].Port = &[]uint16{uint16(hostPort)}[0]

	t.Log("invalid custom broker")
	{
		cm := &CheckManager{}
		ac := &apiclient.Config{
			TokenApp: "abcd",
			TokenKey: "1234",
			URL:      server.URL,
		}
		apih, err := apiclient.NewAPI(ac)
		if err != nil {
			t.Fatalf("unexpected error (%s)", err)
		}
		cm.apih = apih
		cm.brokerID = 1

		_, err = cm.getBroker()
		if err == nil || err.Error() != "error, designated broker 1 [test broker] is invalid (not active, does not support required check type, or connectivity issue)" {
			t.Fatalf("unexpected error (%s)", err)
		}
	}

	t.Log("valid custom broker")
	{

		cm := &CheckManager{
			checkType:             "httptrap",
			brokerMaxResponseTime: time.Duration(time.Millisecond * 500),
		}
		ac := &apiclient.Config{
			TokenApp: "abcd",
			TokenKey: "1234",
			URL:      server.URL,
		}

		apih, err := apiclient.NewAPI(ac)
		if err != nil {
			t.Fatalf("unexpected error (%s)", err)
		}
		cm.apih = apih
		cm.brokerID = 2

		_, err = cm.getBroker()
		if err != nil {
			t.Fatalf("unexpected error (%s)", err)
		}
	}

}
