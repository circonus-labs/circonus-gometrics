package checkmgr

import (
	"errors"
	"os"
    "fmt"
	"testing"
    "time"
    "reflect"

    "github.com/circonus-labs/circonus-gometrics/api"
)

func TestGetBrokerInvalid(t *testing.T) {
    if os.Getenv("CIRCONUS_API_TOKEN") == "" {
        t.Skip("skipping test; $CIRCONUS_API_TOKEN not set")
    }

	t.Log("Testing invalid custom broker id (275/chicago-no httptrap)")

    cm := &CheckManager{}
    ac := &api.Config{}
    ac.Token = api.TokenConfig{
        Key: os.Getenv("CIRCONUS_API_TOKEN"),
        App: "circonus-gometrics",
    }
    apih, err := api.NewApi(ac)
    if err != nil {
        t.Errorf("Expected no error, got '%v'", err)
    }
    cm.apih = apih
    cm.brokerId = 275

    expectedError := errors.New("[ERROR] designated broker 275 [Chicago, IL, US] is invalid (not active, does not support required check type, or connectivity issue).")

    _, err = cm.getBroker()
    if err == nil || err.Error() != expectedError.Error() {
        t.Errorf("Expected an '%#v' error, got '%#v'", expectedError, err)
    }
}


func TestGetBrokerValid(t *testing.T) {
    if os.Getenv("CIRCONUS_API_TOKEN") == "" {
        t.Skip("skipping test; $CIRCONUS_API_TOKEN not set")
    }

	t.Log("Testing valid custom broker id (35/httptrap)")

    cm := &CheckManager{}
    ac := &api.Config{}
    ac.Token = api.TokenConfig{
        Key: os.Getenv("CIRCONUS_API_TOKEN"),
        App: "circonus-gometrics",
    }
    apih, err := api.NewApi(ac)
    if err != nil {
        t.Errorf("Expected no error, got '%v'", err)
    }
    cm.apih = apih
    cm.brokerId = 35
    cm.brokerMaxResponseTime = 5000 * time.Millisecond
    cm.checkType = "httptrap"

    broker, err := cm.getBroker()
    if err != nil {
        t.Fatalf("Expected no error, got '%v'", err)
    }

    expectedCid := fmt.Sprintf("/broker/%d", cm.brokerId)
    if broker.Cid !=  expectedCid {
        t.Fatalf("%s != %s", broker.Cid, expectedCid)
    }
}

func TestGetBrokerSelection(t *testing.T) {
    if os.Getenv("CIRCONUS_API_TOKEN") == "" {
        t.Skip("skipping test; $CIRCONUS_API_TOKEN not set")
    }

	t.Log("Testing broker selection")

    cm := &CheckManager{}
    ac := &api.Config{}
    ac.Token = api.TokenConfig{
        Key: os.Getenv("CIRCONUS_API_TOKEN"),
        App: "circonus-gometrics",
    }
    apih, err := api.NewApi(ac)
    if err != nil {
        t.Errorf("Expected no error, got '%v'", err)
    }
    cm.apih = apih
    cm.brokerMaxResponseTime = 5000 * time.Millisecond
    cm.checkType = "httptrap"

    broker, err := cm.getBroker()
    if err != nil {
        t.Fatalf("Expected no error, got '%v'", err)
    }

    actualType := reflect.TypeOf(broker)
    expectedType := "*api.Broker"
    if actualType.String() != expectedType {
        t.Errorf("Expected *api.Broker, got %s", actualType.String())
    }

    if broker.Cid[:8] != "/broker/" {
        t.Errorf("Expected cid to start with '/broker/', found: %s", broker.Cid)
    }

    t.Logf("Selected broker %s %s", broker.Name, broker.Cid)
}
