// Copyright 2016 Circonus, Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package api

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/hashicorp/go-retryablehttp"
)

const (
	// a few sensible defaults
	defaultAPIURL = "https://api.circonus.com/v2"
	defaultAPIApp = "circonus-gometrics"
	minRetryWait  = 1 * time.Second
	maxRetryWait  = 15 * time.Second
	maxRetries    = 4 // equating to 1 + maxRetries total attempts
)

// TokenKeyType - Circonus API Token key
type TokenKeyType string

// TokenAppType - Circonus API Token app name
type TokenAppType string

// CIDType Circonus object cid
type CIDType *string

// IDType Circonus object id
type IDType int

// URLType submission url type
type URLType string

// SearchQueryType search query (see: https://login.circonus.com/resources/api#searching)
type SearchQueryType string

// SearchFilterType search filter (see: https://login.circonus.com/resources/api#filtering)
type SearchFilterType map[string][]string

// TagType search/select/custom tag(s) type
type TagType []string

// Config options for Circonus API
type Config struct {
	URL      string
	TokenKey string
	TokenApp string
	Log      *log.Logger
	Debug    bool
}

// API Circonus API
type API struct {
	apiURL                *url.URL
	key                   TokenKeyType
	app                   TokenAppType
	Debug                 bool
	Log                   *log.Logger
	useExponentialBackoff bool
}

// NewClient returns a new Circonus API (alias for New)
func NewClient(ac *Config) (*API, error) {
	return New(ac)
}

// NewAPI returns a new Circonus API (alias for New)
func NewAPI(ac *Config) (*API, error) {
	return New(ac)
}

// New returns a new Circonus API
func New(ac *Config) (*API, error) {

	if ac == nil {
		return nil, errors.New("Invalid API configuration (nil)")
	}

	key := TokenKeyType(ac.TokenKey)
	if key == "" {
		return nil, errors.New("API Token is required")
	}

	app := TokenAppType(ac.TokenApp)
	if app == "" {
		app = defaultAPIApp
	}

	au := string(ac.URL)
	if au == "" {
		au = defaultAPIURL
	}
	if !strings.Contains(au, "/") {
		// if just a hostname is passed, ASSume "https" and a path prefix of "/v2"
		au = fmt.Sprintf("https://%s/v2", ac.URL)
	}
	if last := len(au) - 1; last >= 0 && au[last] == '/' {
		// strip off trailing '/'
		au = au[:last]
	}
	apiURL, err := url.Parse(au)
	if err != nil {
		return nil, err
	}

	a := &API{apiURL, key, app, ac.Debug, ac.Log, false}

	a.Debug = ac.Debug
	a.Log = ac.Log
	if a.Debug && a.Log == nil {
		a.Log = log.New(os.Stderr, "", log.LstdFlags)
	}
	if a.Log == nil {
		a.Log = log.New(ioutil.Discard, "", log.LstdFlags)
	}

	return a, nil
}

// EnableExponentialBackoff enables use of exponential backoff for next API call(s)
// and use exponential backoff for all API calls until exponential backoff is disabled.
func (a *API) EnableExponentialBackoff() {
	a.useExponentialBackoff = true
}

// DisableExponentialBackoff disables use of exponential backoff. If a request using
// exponential backoff is currently running, it will stop using exponential backoff
// on its next iteration (if needed).
func (a *API) DisableExponentialBackoff() {
	a.useExponentialBackoff = false
}

// Get API request
func (a *API) Get(reqPath string) ([]byte, error) {
	return a.apiRequest("GET", reqPath, nil)
}

// Delete API request
func (a *API) Delete(reqPath string) ([]byte, error) {
	return a.apiRequest("DELETE", reqPath, nil)
}

// Post API request
func (a *API) Post(reqPath string, data []byte) ([]byte, error) {
	return a.apiRequest("POST", reqPath, data)
}

// Put API request
func (a *API) Put(reqPath string, data []byte) ([]byte, error) {
	return a.apiRequest("PUT", reqPath, data)
}

// apiRequest manages retry strategy for exponential backoffs
func (a *API) apiRequest(reqMethod string, reqPath string, data []byte) ([]byte, error) {
	backoffs := []uint{2, 4, 8, 16, 32}
	attempts := 0
	success := false

	var result []byte
	var err error

	for !success {
		result, err = a.apiCall(reqMethod, reqPath, data)
		if err == nil {
			success = true
		}

		// break and return error if not using exponential backoff
		if err != nil {
			if !a.useExponentialBackoff {
				break
			}
		}

		if !success {
			attempts++
			var wait uint
			if attempts >= len(backoffs) {
				wait = backoffs[len(backoffs)]
			} else {
				wait = backoffs[attempts]
			}
			a.Log.Printf("[WARN] %s, retrying in %d seconds.\n", err.Error(), wait)
			time.Sleep(time.Duration(wait) * time.Second)
		}
	}

	return result, err
}

// apiCall call Circonus API
func (a *API) apiCall(reqMethod string, reqPath string, data []byte) ([]byte, error) {
	reqURL := a.apiURL.String()

	if reqPath == "" {
		return nil, errors.New("Invalid URL path")
	}
	if reqPath[:1] != "/" {
		reqURL += "/"
	}
	if len(reqPath) >= 3 && reqPath[:3] == "/v2" {
		reqURL += reqPath[3:len(reqPath)]
	} else {
		reqURL += reqPath
	}

	// keep last HTTP error in the event of retry failure
	var lastHTTPError error
	retryPolicy := func(resp *http.Response, err error) (bool, error) {
		if err != nil {
			lastHTTPError = err
			return true, err
		}
		// Check the response code. We retry on 500-range responses to allow
		// the server time to recover, as 500's are typically not permanent
		// errors and may relate to outages on the server side. This will catch
		// invalid response codes as well, like 0 and 999.
		// Retry on 429 (rate limit) as well.
		if resp.StatusCode == 0 || // wtf?!
			resp.StatusCode >= 500 || // rutroh
			resp.StatusCode == 429 { // rate limit
			body, readErr := ioutil.ReadAll(resp.Body)
			if readErr != nil {
				lastHTTPError = fmt.Errorf("- last HTTP error: %d %+v", resp.StatusCode, readErr)
			} else {
				lastHTTPError = fmt.Errorf("- last HTTP error: %d %s", resp.StatusCode, string(body))
			}
			return true, nil
		}
		return false, nil
	}

	dataReader := bytes.NewReader(data)

	req, err := retryablehttp.NewRequest(reqMethod, reqURL, dataReader)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] creating API request: %s %+v", reqURL, err)
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("X-Circonus-Auth-Token", string(a.key))
	req.Header.Add("X-Circonus-App-Name", string(a.app))

	client := retryablehttp.NewClient()
	if a.useExponentialBackoff {
		// limit to one request if using exponential backoff
		client.RetryWaitMin = 1
		client.RetryWaitMax = 2
		client.RetryMax = 1
	} else {
		client.RetryWaitMin = minRetryWait
		client.RetryWaitMax = maxRetryWait
		client.RetryMax = maxRetries
	}

	// retryablehttp only groks log or no log
	if a.Debug {
		client.Logger = a.Log
	} else {
		client.Logger = log.New(ioutil.Discard, "", log.LstdFlags)
	}

	client.CheckRetry = retryPolicy

	resp, err := client.Do(req)
	if err != nil {
		if lastHTTPError != nil {
			return nil, lastHTTPError
		}
		return nil, fmt.Errorf("[ERROR] %s: %+v", reqURL, err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] reading response %+v", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		msg := fmt.Sprintf("API response code %d: %s", resp.StatusCode, string(body))
		if a.Debug {
			a.Log.Printf("[DEBUG] %s\n", msg)
		}

		return nil, fmt.Errorf("[ERROR] %s", msg)
	}

	return body, nil
}
