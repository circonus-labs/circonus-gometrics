package api

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"testing"
)

func TestFetchBrokerByID(t *testing.T) {
	if os.Getenv("CIRCONUS_API_TOKEN") == "" {
		t.Skip("skipping test; $CIRCONUS_API_TOKEN not set")
	}
	if os.Getenv("CIRC_API_TEST_BROKER_ID") == "" {
		t.Skip("skipping test; $CIRC_API_TEST_BROKER_ID not set")
	}

	t.Log("Testing correct return from API call")

	ac := &Config{}
	ac.TokenKey = os.Getenv("CIRCONUS_API_TOKEN")
	apih, err := NewAPI(ac)
	if err != nil {
		t.Errorf("Expected no error, got '%v'", err)
	}

	bid := os.Getenv("CIRC_API_TEST_BROKER_ID")
	if bid == "" {
		t.Fatal("Invalid broker id (empty)")
	}

	id, err := strconv.Atoi(bid)
	if err != nil {
		t.Fatalf("Unable to convert broker id %s to int", bid)
	}

	brokerID := BrokerIDType(id)

	broker, err := apih.FetchBrokerByID(brokerID)
	if err != nil {
		t.Fatalf("Expected no error, got '%v'", err)
	}

	actualType := reflect.TypeOf(broker)
	expectedType := "*api.Broker"
	if actualType.String() != expectedType {
		t.Errorf("Expected %s, got %s", expectedType, actualType.String())
	}

	expectedCid := fmt.Sprintf("/broker/%s", strconv.Itoa(int(brokerID)))
	if broker.Cid != expectedCid {
		t.Fatalf("%s != %s", broker.Cid, expectedCid)
	}

	t.Logf("Broker returned %s %s", broker.Name, broker.Cid)
}

func TestFetchBrokerByCID(t *testing.T) {
	if os.Getenv("CIRCONUS_API_TOKEN") == "" {
		t.Skip("skipping test; $CIRCONUS_API_TOKEN not set")
	}
	if os.Getenv("CIRC_API_TEST_BROKER_CID") == "" {
		t.Skip("skipping test; $CIRC_API_TEST_BROKER_CID not set")
	}

	t.Log("Testing correct return from API call")

	ac := &Config{}
	ac.TokenKey = os.Getenv("CIRCONUS_API_TOKEN")
	apih, err := NewAPI(ac)
	if err != nil {
		t.Errorf("Expected no error, got '%v'", err)
	}

	cid := os.Getenv("CIRC_API_TEST_BROKER_CID")
	if cid == "" {
		t.Fatal("Invalid broker cid (empty)")
	}

	brokerCid := BrokerCIDType(cid)

	broker, err := apih.FetchBrokerByCID(brokerCid)
	if err != nil {
		t.Fatalf("Expected no error, got '%v'", err)
	}

	actualType := reflect.TypeOf(broker)
	expectedType := "*api.Broker"
	if actualType.String() != expectedType {
		t.Errorf("Expected %s, got %s", expectedType, actualType.String())
	}

	expectedCid := string(brokerCid)
	if broker.Cid != expectedCid {
		t.Fatalf("%s != %s", broker.Cid, expectedCid)
	}

	t.Logf("Broker returned %s %s", broker.Name, broker.Cid)
}

func TestFetchBrokerListByTag(t *testing.T) {
	if os.Getenv("CIRCONUS_API_TOKEN") == "" {
		t.Skip("skipping test; $CIRCONUS_API_TOKEN not set")
	}
	if os.Getenv("CIRC_API_TEST_BROKER_TAG") == "" {
		t.Skip("skipping test; $CIRC_API_TEST_BROKER_TAG not set")
	}

	t.Log("Testing correct return from API call")

	ac := &Config{}
	ac.TokenKey = os.Getenv("CIRCONUS_API_TOKEN")
	apih, err := NewAPI(ac)
	if err != nil {
		t.Errorf("Expected no error, got '%v'", err)
	}

	tag := os.Getenv("CIRC_API_TEST_BROKER_TAG")
	if tag == "" {
		t.Fatal("Invalid broker tag (empty)")
	}

	selectTag := BrokerSearchTagType(tag)

	brokers, err := apih.FetchBrokerListByTag(selectTag)
	if err != nil {
		t.Fatalf("Expected no error, got '%v'", err)
	}

	actualType := reflect.TypeOf(brokers)
	expectedType := "[]api.Broker"
	if actualType.String() != expectedType {
		t.Errorf("Expected %s, got %s", expectedType, actualType.String())
	}

	t.Logf("%d brokers returned", len(brokers))
}

func TestFetchBrokerList(t *testing.T) {
	if os.Getenv("CIRCONUS_API_TOKEN") == "" {
		t.Skip("skipping test; $CIRCONUS_API_TOKEN not set")
	}

	t.Log("Testing correct return from API call")

	ac := &Config{}
	ac.TokenKey = os.Getenv("CIRCONUS_API_TOKEN")
	apih, err := NewAPI(ac)
	if err != nil {
		t.Errorf("Expected no error, got '%v'", err)
	}

	brokers, err := apih.FetchBrokerList()
	if err != nil {
		t.Fatalf("Expected no error, got '%v'", err)
	}

	actualType := reflect.TypeOf(brokers)
	expectedType := "[]api.Broker"
	if actualType.String() != expectedType {
		t.Errorf("Expected %s, got %s", expectedType, actualType.String())
	}

	t.Logf("%d brokers returned", len(brokers))
}
