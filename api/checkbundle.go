package api

// abstracted in preparation of separate circonus-api-go package

import (
	"encoding/json"
	"fmt"
)

type CheckBundleConfig struct {
	AsyncMetrics  bool   `json:"async_metrics"`
	Secret        string `json:"secret"`
	SubmissionUrl string `json:"submission_url"`
}

type CheckBundleMetric struct {
	Name   string `json:"name"`
	Type   string `json:"type"`
	Units  string `json:"units"`
	Status string `json:"status"`
}

type CheckBundle struct {
	CheckUUIDs         []string            `json:"_check_uuids,omitempty"`
	Checks             []string            `json:"_checks,omitempty"`
	Cid                string              `json:"_cid,omitempty"`
	Created            int                 `json:"_created,omitempty"`
	LastModified       int                 `json:"_last_modified,omitempty"`
	LastModifedBy      string              `json:"_last_modifed_by,omitempty"`
	ReverseConnectUrls []string            `json:"_reverse_connection_urls,omitempty"`
	Brokers            []string            `json:"brokers"`
	Config             CheckBundleConfig   `json:"config"`
	DisplayName        string              `json:"display_name"`
	Metrics            []CheckBundleMetric `json:"metrics"`
	MetricLimit        int                 `json:"metric_limit"`
	Notes              string              `json:"notes"`
	Period             int                 `json:"period"`
	Status             string              `json:"status"`
	Tags               []string            `json:"tags"`
	Target             string              `json:"target"`
	Timeout            int                 `json:"timeout"`
	Type               string              `json:"type"`
}

// Use Circonus API to retrieve a check bundle by ID
func (a *Api) FetchCheckBundleById(id int) (*CheckBundle, error) {
	cid := fmt.Sprintf("/check_bundle/%d", id)
	return a.FetchCheckBundleByCid(cid)
}

// Use Circonus API to retrieve a check bundle by CID
func (a *Api) FetchCheckBundleByCid(cid string) (*CheckBundle, error) {
	result, err := a.Get(cid)
	if err != nil {
		return nil, err
	}

	checkBundle := new(CheckBundle)
	json.Unmarshal(result, checkBundle)

	return checkBundle, nil
}

// Use Circonus API to search for a check bundle
func (a *Api) SearchCheckBundles(searchCriteria string) ([]CheckBundle, error) {
	apiPath := fmt.Sprintf("/check_bundle?search=%s", searchCriteria)

	response, err := a.Get(apiPath)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] API call error %+v", err)
	}

	var results []CheckBundle
	err = json.Unmarshal(response, &results)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Parsing JSON response %+v", err)
	}

	return results, nil
}

// Use Circonus API to create a check bundle
func (a *Api) CreateCheckBundle(config CheckBundle) (*CheckBundle, error) {
	cfgJson, err := json.Marshal(config)
	if err != nil {
		return nil, err
	}

	response, err := a.Post("/check_bundle", cfgJson)
	if err != nil {
		return nil, err
	}

	checkBundle := new(CheckBundle)
	err = json.Unmarshal(response, checkBundle)
	if err != nil {
		return nil, err
	}

	return checkBundle, nil
}

// Use Circonus API to update a check bundle
func (a *Api) UpdateCheckBundle(config *CheckBundle) (*CheckBundle, error) {
	if a.Debug {
		a.Log.Printf("[DEBUG] Updating check bundle with new metrics.")
	}

	cfgJson, err := json.Marshal(config)
	if err != nil {
		return nil, err
	}

	response, err := a.Put(config.Cid, cfgJson)
	if err != nil {
		return nil, err
	}

	checkBundle := new(CheckBundle)
	err = json.Unmarshal(response, checkBundle)
	if err != nil {
		return nil, err
	}

	return checkBundle, nil
}