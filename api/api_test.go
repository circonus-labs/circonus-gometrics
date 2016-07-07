package api

import (
	"errors"
	"os"
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
		Key: "", //os.Getenv("CIRCONUS_API_TOKEN"),
		App: os.Getenv("CIRCONUS_API_APP"),
	}
	_, err := NewApi(ac)

	if err == nil || err.Error() != expectedError.Error() {
		t.Errorf("Expected an '%#v' error, got '%#v'", expectedError, err)
	}

}

func TestApiGetUser(t *testing.T) {
	t.Log("Testing correct API call to /user/current")

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
