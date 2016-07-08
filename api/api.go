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
	"time"

	"github.com/hashicorp/go-retryablehttp"
)

const (
	// a few sensible defaults
	defaultApiUrl = "https://api.circonus.com/v2"
	defaultApiApp = "circonus-gometrics"
)

type TokenConfig struct {
	Key string
	App string
}

type Config struct {
	Url   string
	Token TokenConfig
	Log   *log.Logger
	Debug bool
}

type Api struct {
	apiUrl *url.URL
	key    string
	app    string
	Debug  bool
	Log    *log.Logger
}

// New Circonus API handle
func NewApi(ac *Config) (*Api, error) {

	if ac == nil {
		return nil, errors.New("Invalid API configuration (nil).")
	}

	key := ac.Token.Key
	if key == "" {
		return nil, errors.New("API Token is required.")
	}

	app := ac.Token.App
	if app == "" {
		app = defaultApiApp
	}

	api_url := ac.Url
	if api_url == "" {
		api_url = defaultApiUrl
	}
	if last := len(api_url) - 1; last >= 0 && api_url[last] == '/' {
		api_url = api_url[:last]
	}
	apiUrl, err := url.Parse(api_url)
	if err != nil {
		return nil, err
	}

	a := &Api{apiUrl, key, app, ac.Debug, ac.Log}

	if a.Log == nil {
		if a.Debug {
			a.Log = log.New(os.Stderr, "", log.LstdFlags)
		} else {
			a.Log = log.New(ioutil.Discard, "", log.LstdFlags)
		}
	}

	return a, nil
}

// API GET request
func (a *Api) Get(reqPath string) ([]byte, error) {
	return a.apiCall("GET", reqPath, nil)
}

// API DELETE request
func (a *Api) Delete(reqPath string) ([]byte, error) {
	return a.apiCall("DELETE", reqPath, nil)
}

// API Post request
func (a *Api) Post(reqPath string, data []byte) ([]byte, error) {
	return a.apiCall("POST", reqPath, data)
}

// API PUT request
func (a *Api) Put(reqPath string, data []byte) ([]byte, error) {
	return a.apiCall("PUT", reqPath, data)
}

// Call Circonus API
func (a *Api) apiCall(reqMethod string, reqPath string, data []byte) ([]byte, error) {
	dataReader := bytes.NewReader(data)
	reqUrl := a.apiUrl.String()

	if reqPath[:1] != "/" {
		reqUrl += "/"
	}
	if reqPath[:3] == "/v2" {
		reqUrl += reqPath[3:len(reqPath)]
	} else {
		reqUrl += reqPath
	}

	req, err := retryablehttp.NewRequest(reqMethod, reqUrl, dataReader)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] creating API request: %s %+v", reqUrl, err)
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("X-Circonus-Auth-Token", a.key)
	req.Header.Add("X-Circonus-App-Name", a.app)

	client := retryablehttp.NewClient()
	client.RetryWaitMin = 10 * time.Millisecond
	client.RetryWaitMax = 50 * time.Millisecond
	client.RetryMax = 3
	client.Logger = a.Log

	resp, err := client.Do(req)
	if err != nil {
		standard_client := &http.Client{}
		dataReader.Seek(0, 0)
		standard_req, _ := http.NewRequest(reqMethod, reqUrl, dataReader)
		standard_req.Header.Add("Accept", "application/json")
		standard_req.Header.Add("X-Circonus-Auth-Token", a.key)
		standard_req.Header.Add("X-Circonus-App-Name", a.app)
		resp, err := standard_client.Do(standard_req)
		if resp != nil && resp.Body != nil {
			defer resp.Body.Close()
			body, _ := ioutil.ReadAll(resp.Body)
			if a.Debug {
				a.Log.Printf("[DEBUG] %v\n", string(body))
			}
			return nil, fmt.Errorf("[ERROR] %s", string(body))
		}
		return nil, fmt.Errorf("[ERROR] fetching %s: %s", reqUrl, err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] reading body %+v", err)
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
