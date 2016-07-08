package checkmgr

import (
	"errors"
	"os"
	"testing"
)

func TestNewCheckManager1(t *testing.T) {
	t.Log("Testing correct error return when no Check Manager config supplied")

	expectedError := errors.New("Invalid Check Manager configuration (nil).")

	_, err := NewCheckManager(nil)

	if err == nil || err.Error() != expectedError.Error() {
		t.Errorf("Expected an '%#v' error, got '%#v'", expectedError, err)
	}

}

func TestNewCheckManager2(t *testing.T) {
	t.Log("Testing correct error return when no API Token and no Submission URL supplied")

	expectedError := errors.New("Invalid check manager configuration (no API token AND no submission url).")

	cfg := &Config{}
	_, err := NewCheckManager(cfg)

	if err == nil || err.Error() != expectedError.Error() {
		t.Errorf("Expected an '%#v' error, got '%#v'", expectedError, err)
	}

}

func TestNewCheckManager3(t *testing.T) {
	t.Log("Testing correct return with Submission URL (http) and no API Token supplied")

	cfg := &Config{}
	cfg.Check.SubmissionUrl = "http://127.0.0.1:56104"

	cm, err := NewCheckManager(cfg)
	if err != nil {
		t.Errorf("Expected no error, got '%v'", err)
	}

	trap, err := cm.GetTrap()
	if err != nil {
		t.Errorf("Expected no error, got '%v'", err)
	}

	if trap.Url.String() != cfg.Check.SubmissionUrl {
		t.Errorf("Expected '%s' == '%s'", trap.Url.String(), cfg.Check.SubmissionUrl)
	}

	if trap.Tls != nil {
		t.Errorf("Expected nil found %#v", trap.Tls)
	}
}

func TestNewCheckManager4(t *testing.T) {
	t.Log("Testing correct return with Submission URL (https) and no API Token supplied")

	cfg := &Config{}
	cfg.Check.SubmissionUrl = "https://127.0.0.1/v2"

	cm, err := NewCheckManager(cfg)
	if err != nil {
		t.Fatalf("Expected no error, got '%v'", err)
	}

	trap, err := cm.GetTrap()
	if err != nil {
		t.Fatalf("Expected no error, got '%v'", err)
	}

	if trap.Url.String() != cfg.Check.SubmissionUrl {
		t.Fatalf("Expected '%s' == '%s'", trap.Url.String(), cfg.Check.SubmissionUrl)
	}

	if trap.Tls == nil {
		t.Fatalf("Expected a x509 cert pool, found nil")
	}
}

func TestNewCheckManager5(t *testing.T) {
	if os.Getenv("CIRCONUS_API_TOKEN") == "" {
		t.Skip("skipping test; $CIRCONUS_API_TOKEN not set")
	}

    t.Log("Testing custom/specific Broker ID (invalid)")

	cfg := &Config{}
	cfg.Api.Token.Key = os.Getenv("CIRCONUS_API_TOKEN")
    // chicago broker, does not support httptrap
    cfg.Broker.Id = 275

    expectedError := errors.New("[ERROR] designated broker 275 [Chicago, IL, US] is invalid (not active, does not support required check type, or connectivity issue).")

	_, err := NewCheckManager(cfg)

    if err == nil || err.Error() != expectedError.Error() {
        t.Errorf("Expected an '%#v' error, got '%#v'", expectedError, err)
    }
}

func TestNewCheckManager6(t *testing.T) {
	// flag to indicate whether to do this test
	if os.Getenv("CIRCONUS_CGM_CMTEST6") == "" {
		t.Skip("skipping test; $CIRCONUS_CGM_CMTEST6 not set")
	}

	// !!IMPORTANT!! this test is DESTRUCTIVE it will DELETE the check bundle
	//
	// this test expects to CREATE a check then, search (and find) the check.
	//
	// ensure there is no existing check which would match the default search criteria
	// it *will* be deleted at the end of this test...
	//
	// the default InstanceId is "os.hostname():program name" e.g. testhost1:checkmgr.test
	// the default SearchTag is "service:program name" e.g. service:checkmgr.test

	if os.Getenv("CIRCONUS_API_TOKEN") == "" {
		t.Skip("skipping test; $CIRCONUS_API_TOKEN not set")
	}

	t.Log("Testing correct check creation and search with API Token only")

	cfg := &Config{}
	cfg.Api.Token.Key = os.Getenv("CIRCONUS_API_TOKEN")

	t.Log("Testing correct check creation - should create a check, if it doesn't exist")
	cm, err := NewCheckManager(cfg)
	if err != nil {
		t.Fatalf("Expected no error, got '%v'", err)
	}

    t.Log("Getting Trap from cm instance")
    trap, err := cm.GetTrap()
	if err != nil {
		t.Fatalf("Expected no error, got '%v'", err)
	}

	t.Log("Testing correct check search - should find the check created")
	cm2, err := NewCheckManager(cfg)
	if err != nil {
		t.Fatalf("Expected no error, got '%v'", err)
	}

    t.Log("Getting Trap from cm2 instance")
	trap2, err := cm2.GetTrap()
	if err != nil {
		t.Fatalf("Expected no error, got '%v'", err)
	}

    t.Log("Comparing Trap URLs")
	if trap.Url.String() != trap2.Url.String() {
		t.Fatalf("Expected '%s' == '%s'", trap.Url.String(), trap2.Url.String())
	}

	t.Logf("Deleting %s %s", cm2.checkBundle.Cid, cm2.checkBundle.DisplayName)
	_, err = cm2.apih.Delete(cm2.checkBundle.Cid)
	if err != nil {
		t.Fatalf("Expected no error, got '%v'", err)
	}
}
