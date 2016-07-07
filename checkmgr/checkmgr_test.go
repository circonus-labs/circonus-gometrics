package checkmgr

import (
	"errors"
	"testing"
)

func TestNewCheckManagerInvalidConfig(t *testing.T) {
	t.Log("Testing correct error return when no Check Manager config supplied")

	expectedError := errors.New("Invalid Check Manager configuration (nil).")

	_, err := NewCheckManager(nil)

	if err == nil || err.Error() != expectedError.Error() {
		t.Errorf("Expected an '%#v' error, got '%#v'", expectedError, err)
	}

}

func TestNewCheckManagerNoTokenNoUrl(t *testing.T) {
	t.Log("Testing correct error return when no API Token and no Submission URL supplied")

	expectedError := errors.New("Invalid check manager configuration (no API token AND no submission url).")

	cfg := &Config{}
	_, err := NewCheckManager(cfg)

	if err == nil || err.Error() != expectedError.Error() {
		t.Errorf("Expected an '%#v' error, got '%#v'", expectedError, err)
	}

}

func TestNewCheckManagerUrlNoToken(t *testing.T) {
	t.Log("Testing correct return with Submission URL and no API Token supplied")

	cfg := &Config{}
	cfg.Check.SubmissionUrl = "http://127.0.0.1:56104"

	cm, err := NewCheckManager(cfg)
	if err != nil {
		t.Errorf("Expected no error, got '%v'", err)
	}

	trapUrl, err := cm.GetTrapUrl()
	if err != nil {
		t.Errorf("Expected no error, got '%v'", err)
	}
	if trapUrl != cfg.Check.SubmissionUrl {
		t.Errorf("Expected '%s' == '%s'", trapUrl, cfg.Check.SubmissionUrl)
	}

}
