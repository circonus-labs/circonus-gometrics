package circonusgometrics

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/hashicorp/go-retryablehttp"
)

func (m *CirconusMetrics) apiGetBrokers() (*[]Broker, error) {
	response, err := m.apiCall("GET", "/v2/broker", nil)
	if err != nil {
		return nil, err
	}

	brokerList := new([]Broker)
	err = json.Unmarshal(response, brokerList)
	if err != nil {
		return nil, err
	}

	return brokerList, nil
}

func (m *CirconusMetrics) apiGetCert() ([]byte, error) {
	response, err := m.apiCall("GET", "/v2/pki/ca.crt", nil)
	if err != nil {
		return nil, err
	}

	cadata := new(CACert)
	err = json.Unmarshal(response, cadata)
	if err != nil {
		return nil, err
	}

	if cadata.Contents == "" {
		return nil, fmt.Errorf("Unable to find ca cert %+v", cadata)
	}
	return []byte(cadata.Contents), nil
}

func (m *CirconusMetrics) apiCreateCheck(checkConfig []byte) (*Check, error) {
	response, err := m.apiCall("POST", "/v2/check_bundle", checkConfig)
	if err != nil {
		return nil, err
	}

	check := new(Check)
	err = json.Unmarshal(response, check)
	if err != nil {
		return nil, err
	}

	return check, nil

}

func (m *CirconusMetrics) apiCheckSearch() (*Check, error) {
	searchCriteria := fmt.Sprintf("(active:1)(host:\"%s\")(type:\"%s\")(tags:%s)", m.InstanceId, checkType, circonusSearchTag)
	apiPath := fmt.Sprintf("/v2/check_bundle?search=%s", searchCriteria)

	response, err := m.apiCall("GET", apiPath, nil)
	if err != nil {
		return nil, fmt.Errorf("API call error %+v", response)
	}

	var results []Check
	err = json.Unmarshal(response, &results)
	if err != nil {
		return nil, fmt.Errorf("Error parsing JSON response %+v", err)
	}

	if len(results) == 0 {
		return nil, nil
		//return nil, fmt.Errorf("No checks found")
	}

	numActive := 0
	checkId := -1

	for idx, check := range results {
		if check.Status == "active" {
			numActive++
			checkId = idx
		}
	}

	if numActive > 1 {
		return nil, fmt.Errorf("multiple possibilities >1 check matches criteria target=%s, type=%s, with tag %s\n", m.InstanceId, checkType, circonusSearchTag)
	}

	return &results[checkId], nil
}

func (m *CirconusMetrics) apiCall(reqMethod string, reqPath string, data []byte) ([]byte, error) {
	dataReader := bytes.NewReader(data)

	url := fmt.Sprintf("https://%s%s", m.ApiHost, reqPath)

	req, err := retryablehttp.NewRequest(reqMethod, url, dataReader)
	if err != nil {
		return nil, fmt.Errorf("Error making API request to %s %+v", url, err)
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("X-Circonus-Auth-Token", m.ApiToken)
	req.Header.Add("X-Circonus-App-Name", m.ApiApp)

	client := retryablehttp.NewClient()
	client.RetryWaitMin = 10 * time.Millisecond
	client.RetryWaitMax = 50 * time.Millisecond
	client.RetryMax = 3
	// silence the debug messages (if consul or go-metrics offers a standard logging, use that instead)
	if m.Debug {
		client.Logger = m.Log //log.New(ioutil.Discard, "", log.LstdFlags)
	} else {
		client.Logger = log.New(ioutil.Discard, "", log.LstdFlags)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Error fetching %s: %s\n", url, err)
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
