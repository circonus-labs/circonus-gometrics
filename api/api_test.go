package api

import (
	"errors"
	"os"
	"strings"
	"testing"
)

func TestNewApiInvalidConfig(t *testing.T) {
	t.Log("Testing correct error return when no API config supplied")

	expectedError := errors.New("Invalid API configuration (nil).")

	_, err := NewApi(nil)

	if err == nil || err.Error() != expectedError.Error() {
		t.Errorf("Expected an '%#v' error, got '%#v'", expectedError, err)
	}

}

func TestNewApiNoToken(t *testing.T) {
	t.Log("Testing correct error return when no API Token supplied")

	expectedError := errors.New("API Token is required.")

	ac := &Config{}
	ac.Token = TokenConfig{
		Key: "",
		App: os.Getenv("CIRCONUS_API_APP"),
	}
	_, err := NewApi(ac)

	if err == nil || err.Error() != expectedError.Error() {
		t.Errorf("Expected an '%#v' error, got '%#v'", expectedError, err)
	}

}

func TestNewApiInvalidToken(t *testing.T) {
	t.Log("Testing correct error return when INVALID API Token supplied")

	ac := &Config{}
	ac.Token = TokenConfig{
		Key: "abc-123",
		App: os.Getenv("CIRCONUS_API_APP"),
	}
	apih, err := NewApi(ac)
	if err != nil {
		t.Errorf("Expected no error, got '%v'", err)
	}

	_, err = apih.Get("/user/current")
	if err == nil || !strings.Contains(err.Error(), "The authentication token you supplied is invalid") {
		t.Errorf("Expected an error containing 'The authentication token you supplied is invalid' error, got '%#v'", err)
	}

}

func TestApiGetUser(t *testing.T) {
	if os.Getenv("CIRCONUS_API_TOKEN") == "" {
		t.Skip("skipping test; $CIRCONUS_API_TOKEN not set")
	}

	t.Log("Testing correct API call to /user/current [defaults]")

	ac := &Config{}
	ac.Token = TokenConfig{
		Key: os.Getenv("CIRCONUS_API_TOKEN"),
		App: os.Getenv("CIRCONUS_API_APP"),
	}
	apih, err := NewApi(ac)
	if err != nil {
		t.Errorf("Expected no error, got '%v'", err)
	}

	if _, err := apih.Get("/user/current"); err != nil {
		t.Errorf("Expected no error, got '%v'", err)
	}

}

func TestApiGetUser2(t *testing.T) {
	if os.Getenv("CIRCONUS_API_TOKEN") == "" {
		t.Skip("skipping test; $CIRCONUS_API_TOKEN not set")
	}

	t.Log("Testing correct API call to /user/current [url=hostname]")

	ac := &Config{}
	ac.Token = TokenConfig{
		Key: os.Getenv("CIRCONUS_API_TOKEN"),
		App: os.Getenv("CIRCONUS_API_APP"),
	}
	ac.Url = "api.circonus.com"
	apih, err := NewApi(ac)
	if err != nil {
		t.Errorf("Expected no error, got '%v'", err)
	}

	if _, err := apih.Get("/user/current"); err != nil {
		t.Errorf("Expected no error, got '%v'", err)
	}

}
