package circonusgometrics

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/hashicorp/go-retryablehttp"
)

// Call Circonus API
func (m *CirconusMetrics) apiCall(reqMethod string, reqPath string, data []byte) ([]byte, error) {
	dataReader := bytes.NewReader(data)

	// default to SSL
	proto := "https://"
	// allow override with explict "http://" in ApiHost
	if m.ApiHost[0:4] == "http" {
		proto = ""
	}

	url := fmt.Sprintf("%s%s%s", proto, m.ApiHost, reqPath)

	req, err := retryablehttp.NewRequest(reqMethod, url, dataReader)
	if err != nil {
		return nil, fmt.Errorf("Error creating API request: %s %+v", url, err)
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("X-Circonus-Auth-Token", m.ApiToken)
	req.Header.Add("X-Circonus-App-Name", m.ApiApp)

	client := retryablehttp.NewClient()
	client.RetryWaitMin = 10 * time.Millisecond
	client.RetryWaitMax = 50 * time.Millisecond
	client.RetryMax = 3
	client.Logger = m.Log

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Error reading body %+v", err)
	}

	if resp.StatusCode != 200 {
		m.Log.Printf("%+v\n", string(body))

		var response map[string]interface{}
		json.Unmarshal(body, &response)
		if err != nil {
			return nil, fmt.Errorf("Error parsing JSON response %+v", err)
		}
		return nil, fmt.Errorf("Error API response code %+v", response)
	}

	return body, nil
}
