package circonusgometrics

import (
	"errors"
	"os"
	"testing"
)

/*

because some of these tests interact with the circonus api directly
to create checks, the tests in question have environment variable
flags to gate whether they run or not.

*/

func TestNewCirconusMetricsInvalidConfig(t *testing.T) {
	t.Log("Testing correct error return when no config supplied")

	expectedError := errors.New("Invalid configuration (nil).")

	_, err := NewCirconusMetrics(nil)

	if err == nil || err.Error() != expectedError.Error() {
		t.Errorf("Expected an '%#v' error, got '%#v'", expectedError, err)
	}
}

func TestNewCirconusMetricsNoTokenNoUrl(t *testing.T) {
	t.Log("Testing correct error return when no API Token and no Submission URL supplied")

	expectedError := errors.New("Invalid check manager configuration (no API token AND no submission url).")

	cfg := &Config{}
	_, err := NewCirconusMetrics(cfg)

	if err == nil || err.Error() != expectedError.Error() {
		t.Errorf("Expected an '%#v' error, got '%#v'", expectedError, err)
	}
}

func TestNewCirconusMetricsHttpUrlNoToken(t *testing.T) {
	t.Log("Testing correct return with Submission URL (http) and no API Token supplied")

	cfg := &Config{}
	cfg.CheckManager.Check.SubmissionUrl = "http://127.0.0.1:56104"

	cm, err := NewCirconusMetrics(cfg)
	if err != nil {
		t.Errorf("Expected no error, got '%v'", err)
	}

	trap, err := cm.check.GetTrap()
	if err != nil {
		t.Errorf("Expected no error, got '%v'", err)
	}

	if trap.Url.String() != cfg.CheckManager.Check.SubmissionUrl {
		t.Errorf("Expected '%s' == '%s'", trap.Url.String(), cfg.CheckManager.Check.SubmissionUrl)
	}

	if trap.Tls != nil {
		t.Errorf("Expected nil found %#v", trap.Tls)
	}
}

func TestNewCirconusMetricsHttpsUrlNoToken(t *testing.T) {
	t.Log("Testing correct return with Submission URL (https) and no API Token supplied")

	cfg := &Config{}
	cfg.CheckManager.Check.SubmissionUrl = "https://127.0.0.1/v2"

	cm, err := NewCirconusMetrics(cfg)
	if err != nil {
		t.Errorf("Expected no error, got '%v'", err)
	}

	trap, err := cm.check.GetTrap()
	if err != nil {
		t.Errorf("Expected no error, got '%v'", err)
	}

	if trap.Url.String() != cfg.CheckManager.Check.SubmissionUrl {
		t.Errorf("Expected '%s' == '%s'", trap.Url.String(), cfg.CheckManager.Check.SubmissionUrl)
	}

	if trap.Tls == nil {
		t.Errorf("Expected a x509 cert pool, found nil")
	}
}

func TestNewCirconusMetrics1(t *testing.T) {
	// flag to indicate whether to do this test
	if os.Getenv("CIRCONUS_CGM_TEST1") == "" {
		t.Skip("skipping test; $CIRCONUS_CGM_TEST1 not set")
	}

	// note, this test expects to CREATE a check then, search (and find) the check.
	// ensure there is no existing check which would match the default search criteria
	// iow, remember to remove the check after running this test otherwise it
	// will test search+search not create+search

	if os.Getenv("CIRCONUS_API_TOKEN") == "" {
		t.Skip("skipping test; $CIRCONUS_API_TOKEN not set")
	}
	t.Log("Testing correct check creation and search with API Token only")

	cfg := &Config{}
	cfg.CheckManager.Api.Token.Key = os.Getenv("CIRCONUS_API_TOKEN")

	t.Log("Testing correct check creation - should create a check, if it doesn't exist")
	cm, err := NewCirconusMetrics(cfg)
	if err != nil {
		t.Errorf("Expected no error, got '%v'", err)
	}

	trap, err := cm.check.GetTrap()
	if err != nil {
		t.Errorf("Expected no error, got '%v'", err)
	}

	t.Log("Testing correct check search - should find the check created")
	cm2, err := NewCirconusMetrics(cfg)
	if err != nil {
		t.Errorf("Expected no error, got '%v'", err)
	}

	trap2, err := cm2.check.GetTrap()
	if err != nil {
		t.Errorf("Expected no error, got '%v'", err)
	}

	if trap.Url.String() != trap2.Url.String() {
		t.Errorf("Expected '%s' == '%s'", trap.Url.String(), trap2.Url.String())
	}
}

func TestNewCirconusMetrics2(t *testing.T) {
	// flag to indicate whether to do this test
	if os.Getenv("CIRCONUS_CGM_TEST2") == "" {
		t.Skip("skipping test; $CIRCONUS_CGM_TEST2 not set")
	}

	// note, this test expects to CREATE a check then, search (and find) the check.
	// ensure there is no existing check which would match the search criteria
	// iow, remember to remove the check after running this test otherwise it
	// will test search+search not create+search

	if os.Getenv("CIRCONUS_API_TOKEN") == "" {
		t.Skip("skipping test; $CIRCONUS_API_TOKEN not set")
	}
	t.Log("Testing correct check creation and search with API Token only")

	cfg := &Config{}
	cfg.CheckManager.Api.Token.Key = os.Getenv("CIRCONUS_API_TOKEN")
	cfg.CheckManager.Check.InstanceId = "cgmtest2:gotest"
	cfg.CheckManager.Check.SearchTag = "gotest:cgm"

	t.Log("Testing correct check creation - should create a check, if it doesn't exist")
	cm, err := NewCirconusMetrics(cfg)
	if err != nil {
		t.Errorf("Expected no error, got '%v'", err)
	}

	trap, err := cm.check.GetTrap()
	if err != nil {
		t.Errorf("Expected no error, got '%v'", err)
	}

	t.Log("Testing correct check search - should find the check created")
	cm2, err := NewCirconusMetrics(cfg)
	if err != nil {
		t.Errorf("Expected no error, got '%v'", err)
	}

	trap2, err := cm2.check.GetTrap()
	if err != nil {
		t.Errorf("Expected no error, got '%v'", err)
	}

	if trap.Url.String() != trap2.Url.String() {
		t.Errorf("Expected '%s' == '%s'", trap.Url.String(), trap2.Url.String())
	}

	t.Log(trap2.Url.String())
}
