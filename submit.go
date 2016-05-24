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

func submit(output map[string]interface{}) {
	str, err := json.Marshal(output)
	if err == nil {
		trapCall(str)
	}
}

func trapCall(payload []byte) (int, error) {
	tr := &http.Transport{
		TLSClientConfig:    &tls.Config{RootCAs: rootCA},
		DisableCompression: true,
	}
	client := &http.Client{Transport: tr}
	req, err := http.NewRequest("POST", submissionUrl, bytes.NewBuffer(payload))
	if err != nil {
		return 0, err
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("X-Circonus-Auth-Token", authtoken)
	req.Header.Add("X-Circonus-App-Name", "circonus-cip")
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
