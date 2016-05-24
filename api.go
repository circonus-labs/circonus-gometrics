package circonusgometrics

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func apiCall(url string) map[string]interface{} {
	client := &http.Client{}
	req, err := http.NewRequest("GET", strings.Join([]string{"https://", apiUrl, url}, ""), nil)
	if err != nil {
		return nil
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("X-Circonus-Auth-Token", authtoken)
	req.Header.Add("X-Circonus-App-Name", "circonus-cip")
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error fetching %s: %s\n", url, err)
		return nil
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	var response map[string]interface{}
	json.Unmarshal(body, &response)
	if resp.StatusCode != 200 {
		log.Printf("response: %v\n", response)
		return nil
	}
	return response
}
