package circonusgometrics

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/hashicorp/go-retryablehttp"
)

func (m *CirconusMetrics) getTrapUrl() (string, error) {
	if m.TrapUrl != "" {
		return m.TrapUrl, nil
	}

	searchCriteria := fmt.Sprintf("(active:1)(host:\"%s\")(type:\"%s\")(tags:%s)", m.InstanceId, checkType, circonusSearchTag)
	apiUrl := fmt.Sprintf("https://%s/v2/check?search=%s", m.ApiHost, searchCriteria)
	req, err := retryablehttp.NewRequest("GET", apiUrl, nil)
	if err != nil {
		return "", err
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
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != 200 {
		var errResp map[string]interface{}
		json.Unmarshal(body, &errResp)
		return "", fmt.Errorf("API error response: %v\n", errResp)
	}

	var response []map[string]interface{}
	json.Unmarshal(body, &response)

	if len(response) == 0 {
		return m.createCheck()
	}

	if len(response) > 1 {
		return "", fmt.Errorf("multiple possibilities >1 check matches criteria target=%s, type=httptrap, with tag %s\n", m.InstanceId, circonusSearchTag)
	}

	active := response[0]["_active"].(bool)
	url := response[0]["_details"].(map[string]interface{})["submission_url"]

	if active && url != nil {
		return url.(string), nil
	}

	return "", fmt.Errorf("No *active* check found.")
}

type CheckConfig struct {
	AsyncMetrics bool   `json:"async_metrics"`
	Secret       string `json:"secret"`
}

type CheckMetric struct {
	Status string `json:"status"`
	Name   string `json:"name"`
	Type   string `json:"type"`
	Units  string `json:"units"`
}

type Check struct {
	Brokers     []string      `json:"brokers"`
	Config      CheckConfig   `json:"config"`
	DisplayName string        `json:"display_name"`
	MetricLimit int           `json:"metric_limit"`
	Metrics     []CheckMetric `json:"metrics"`
	Notes       string        `json:"notes"`
	Period      int           `json:"period"`
	Status      string        `json:"status"`
	Tags        []string      `json:"tags"`
	Target      string        `json:"target"`
	Timeout     int           `json:"timeout"`
	Type        string        `json:"type"`
}

func makeSecret() (string, error) {
	hash := sha256.New()
	x := make([]byte, 2048)
	if _, err := rand.Read(x); err != nil {
		return "", err
	}
	hash.Write(x)
	return hex.EncodeToString(hash.Sum(nil))[0:16], nil
}

func (m *CirconusMetrics) getBrokerGroupId() int {
	bgi := m.BrokerGroupId
	if bgi == 0 {
		bgi = circonusTrapBroker
		// add lookup and leverage selection method unless we can add a method for the account admin to designate a "default" broker (maybe an overall then potentially by check type so that
		// if POST check_bundle {..., brokers: [], ...} is sent, the default broker for the checktype or the generic dfeault broker would be used)
	}
	return bgi
}

func (m *CirconusMetrics) createCheck() (string, error) {
	checkSecret := m.CheckSecret
	if checkSecret == "" {
		secret, err := makeSecret()
		if err != nil {
			secret = "myS3cr3t"
		}
		checkSecret = secret
	}

	config := &Check{
		Brokers:     []string{fmt.Sprintf("/broker/%d", m.getBrokerGroupId())},
		Config:      CheckConfig{AsyncMetrics: true, Secret: checkSecret},
		DisplayName: fmt.Sprintf("%s /%s", m.InstanceId, checkType),
		MetricLimit: 0,
		Notes:       "",
		Period:      60,
		Status:      "active",
		Tags:        append([]string{circonusSearchTag}, m.Tags...),
		Target:      m.InstanceId,
		Timeout:     10,
		Type:        checkType,
	}

	x, err := json.Marshal(config)
	if err != nil {
		m.Log.Fatalf("%+v\n", err)
	}

	if m.Debug {
		m.Log.Printf(string(x))
	}

	return "blah blah", nil
}

/*
func (m *CirconusMetrics) getTrapUrl(id int) {
	url := strings.Join([]string{"/v2/check/", strconv.Itoa(id)}, "")
	checkDetails := m.apiCall(url)
	details, ok := checkDetails["_details"]
	if !ok {
		log.Printf("Cannot find submission URL at %s\n", url)
		return
	}
	dmap := details.(map[string]interface{})
	val := dmap["submission_url"]
	m.trapUrl = val.(string)
}
*/
