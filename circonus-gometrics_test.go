// Copyright 2016 Circonus, Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package circonusgometrics

import (
	"errors"
	"testing"
)

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
	cfg.CheckManager.Check.SubmissionURL = "http://127.0.0.1:56104"

	cm, err := NewCirconusMetrics(cfg)
	if err != nil {
		t.Errorf("Expected no error, got '%v'", err)
	}

	trap, err := cm.check.GetTrap()
	if err != nil {
		t.Errorf("Expected no error, got '%v'", err)
	}

	if trap.URL.String() != cfg.CheckManager.Check.SubmissionURL {
		t.Errorf("Expected '%s' == '%s'", trap.URL.String(), cfg.CheckManager.Check.SubmissionURL)
	}

	if trap.TLS != nil {
		t.Errorf("Expected nil found %#v", trap.TLS)
	}
}

func TestNewCirconusMetricsHttpsUrlNoToken(t *testing.T) {
	t.Log("Testing correct return with Submission URL (https) and no API Token supplied")

	cfg := &Config{}
	cfg.CheckManager.Check.SubmissionURL = "https://127.0.0.1/v2"

	cm, err := NewCirconusMetrics(cfg)
	if err != nil {
		t.Errorf("Expected no error, got '%v'", err)
	}

	trap, err := cm.check.GetTrap()
	if err != nil {
		t.Errorf("Expected no error, got '%v'", err)
	}

	if trap.URL.String() != cfg.CheckManager.Check.SubmissionURL {
		t.Errorf("Expected '%s' == '%s'", trap.URL.String(), cfg.CheckManager.Check.SubmissionURL)
	}

	if trap.TLS == nil {
		t.Errorf("Expected a x509 cert pool, found nil")
	}
}
