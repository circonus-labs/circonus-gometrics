package circonusgometrics

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func (m *CirconusMetrics) apiCall(reqPath string) map[string]interface{} {
	url := fmt.Sprintf("https://%s%s", m.ApiHost, reqPath)
	if m.Debug {
		m.Log.Printf("Calling %s", url)
	}
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		m.Log.Printf("%+v", err)
		return nil
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("X-Circonus-Auth-Token", m.ApiToken)
	req.Header.Add("X-Circonus-App-Name", m.ApiApp)
	resp, err := client.Do(req)
	if err != nil {
		m.Log.Printf("Error fetching %s: %s\n", url, err)
		return nil
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	var response map[string]interface{}
	json.Unmarshal(body, &response)
	if resp.StatusCode != 200 {
		m.Log.Printf("response: %v\n", response)
		return nil
	}
	return response
}
