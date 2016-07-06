package api

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/hashicorp/go-retryablehttp"
)

const (
	// a few sensible defaults
	defaultApiHost = "api.circonus.com"
	defaultApiApp  = "circonus-gometrics"
)

type Config struct {
	Host     string
	TokenKey string
	TokenApp string
	Log      *log.Logger
	Debug    bool
}

type Api struct {
	host  string
	token string
	app   string
	proto string
	Debug bool
	Log   *log.Logger
}

func NewApi(ac *Config) (*Api, error) {
	if ac.TokenKey == "" {
		return nil, errors.New("API Token is required.")
	}

	token := ac.TokenKey

	host := defaultApiHost
	if ac.Host != "" {
		host = ac.Host
	}

	app := defaultApiApp
	if ac.TokenApp != "" {
		app = ac.TokenApp
	}

	// allow override with explict "http://" in ApiHost
	proto := "https://"
	if host[0:5] == "http:" {
		proto = ""
	}

	a := &Api{host, token, app, proto, ac.Debug, ac.Log}

	if a.Log == nil {
		if a.Debug {
			a.Log = log.New(os.Stderr, "", log.LstdFlags)
		} else {
			a.Log = log.New(ioutil.Discard, "", log.LstdFlags)
		}
	}

	return a, nil
}

// Call Circonus API
func (a *Api) apiCall(reqMethod string, reqPath string, data []byte) ([]byte, error) {
	dataReader := bytes.NewReader(data)

	url := fmt.Sprintf("%s%s%s", a.proto, a.host, reqPath)

	req, err := retryablehttp.NewRequest(reqMethod, url, dataReader)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] creating API request: %s %+v", url, err)
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("X-Circonus-Auth-Token", a.token)
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
		standard_req, _ := http.NewRequest(reqMethod, url, dataReader)
		standard_req.Header.Add("Accept", "application/json")
		standard_req.Header.Add("X-Circonus-Auth-Token", a.token)
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
		return nil, fmt.Errorf("[ERROR] fetching %s: %s", url, err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] reading body %+v", err)
	}

	if resp.StatusCode != 200 {
		msg := fmt.Sprintf("API response code %d: %s", resp.StatusCode, string(body))
		if a.Debug {
			a.Log.Printf("[DEBUG] %s\n", msg)
		}

		return nil, fmt.Errorf("[ERROR] %s", msg)
	}

	return body, nil
}
