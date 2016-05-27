package circonusgometrics

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
)

func (m *CirconusMetrics) submit(output map[string]interface{}) {
	str, err := json.Marshal(output)
	if err != nil {
        m.Log.Printf("Error marshling output %+v", err)
        return
    }

	numStats, err := m.trapCall(str)
    if err != nil {
        m.Log.Printf("Error sending metrics to %s %+v\n", m.trapUrl, err)
    }
    if m.Debug {
        m.Log.Printf("%d stats sent to %s\n", numStats, m.trapUrl)
    }
}

func (m *CirconusMetrics) trapCall(payload []byte) (int, error) {
	tr := &http.Transport{
		TLSClientConfig:    &tls.Config{RootCAs: m.certPool}, // add server name from submission_url or lookup of IP in broker list. ServerName: "noit3.dev.circonus.net"},
		DisableCompression: true,
	}
	client := &http.Client{Transport: tr}
	req, err := http.NewRequest("POST", m.trapUrl, bytes.NewBuffer(payload))
	if err != nil {
		return 0, err
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("X-Circonus-Auth-Token", m.ApiToken)
	req.Header.Add("X-Circonus-App-Name", m.ApiApp)
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	var response map[string]interface{}
	json.Unmarshal(body, &response)
	if resp.StatusCode != 200 {
		return 0, errors.New("bad response code: " + strconv.Itoa(resp.StatusCode))
	}
	switch v := response["stats"].(type) {
	case float64:
		return int(v), nil
	case int:
		return v, nil
	default:
	}
	return 0, errors.New("bad response type")
}
