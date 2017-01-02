// Copyright 2016 Circonus, Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Check bundle API support - Fetch, Create, Update, Delete, and Search
// See: https://login.circonus.com/resources/api/calls/check_bundle

package api

import (
	"encoding/json"
	"fmt"
	"net/url"
	"regexp"

	"github.com/circonus-labs/circonus-gometrics/api/config"
)

// CheckBundleMetric individual metric configuration
type CheckBundleMetric struct {
	Name   string   `json:"name"`
	Type   string   `json:"type"`
	Units  string   `json:"units"`
	Status string   `json:"status"`
	Tags   []string `json:"tags"`
}

// CheckBundleConfig contains the check type specific configuration settings
// as k/v pairs (see https://login.circonus.com/resources/api/calls/check_bundle
// for the specific settings available for each distinct check type)
type CheckBundleConfig map[config.Key]string

// CheckBundle definition
type CheckBundle struct {
	CheckUUIDs         []string            `json:"_check_uuids,omitempty"`
	Checks             []string            `json:"_checks,omitempty"`
	CID                string              `json:"_cid,omitempty"`
	Created            int                 `json:"_created,omitempty"`
	LastModified       int                 `json:"_last_modified,omitempty"`
	LastModifedBy      string              `json:"_last_modifed_by,omitempty"`
	ReverseConnectURLs []string            `json:"_reverse_connection_urls,omitempty"`
	Brokers            []string            `json:"brokers"`
	Config             CheckBundleConfig   `json:"config,omitempty"`
	DisplayName        string              `json:"display_name"`
	Metrics            []CheckBundleMetric `json:"metrics"`
	MetricLimit        int                 `json:"metric_limit,omitempty"`
	Notes              string              `json:"notes,omitempty"`
	Period             int                 `json:"period,omitempty"`
	Status             string              `json:"status,omitempty"`
	Tags               []string            `json:"tags,omitempty"`
	Target             string              `json:"target"`
	Timeout            float64             `json:"timeout,omitempty"`
	Type               string              `json:"type"`
}

const (
	baseCheckBundlePath = "/check_bundle"
	checkBundleCIDRegex = "^" + baseCheckBundlePath + "/[0-9]+$"
)

// NewCheckBundle returns a check bundle with defaults
func (a *API) NewCheckBundle() *CheckBundle {
	return &CheckBundle{
		MetricLimit: config.DefaultCheckBundleMetricLimit,
		Period:      config.DefaultCheckBundlePeriod,
		Timeout:     config.DefaultCheckBundleTimeout,
		Status:      config.DefaultCheckBundleStatus,
	}
}

// FetchCheckBundle fetch a check bundle configuration by cid
func (a *API) FetchCheckBundle(cid CIDType) (*CheckBundle, error) {
	if cid == nil || *cid == "" {
		return nil, fmt.Errorf("Invalid check bundle CID [none]")
	}

	bundleCID := string(*cid)

	matched, err := regexp.MatchString(checkBundleCIDRegex, bundleCID)
	if err != nil {
		return nil, err
	}
	if !matched {
		return nil, fmt.Errorf("Invalid check bundle CID [%v]", bundleCID)
	}

	result, err := a.Get(bundleCID)
	if err != nil {
		return nil, err
	}

	checkBundle := &CheckBundle{}
	if err := json.Unmarshal(result, checkBundle); err != nil {
		return nil, err
	}

	return checkBundle, nil
}

// FetchCheckBundles fetch a check bundle configurations
func (a *API) FetchCheckBundles() (*[]CheckBundle, error) {
	result, err := a.Get(baseCheckBundlePath)
	if err != nil {
		return nil, err
	}

	var checkBundles []CheckBundle
	if err := json.Unmarshal(result, &checkBundles); err != nil {
		return nil, err
	}

	return &checkBundles, nil
}

// UpdateCheckBundle updates a check bundle configuration
func (a *API) UpdateCheckBundle(config *CheckBundle) (*CheckBundle, error) {
	if config == nil {
		return nil, fmt.Errorf("Invalid check bundle config [nil]")
	}

	bundleCID := string(config.CID)

	if matched, err := regexp.MatchString(checkBundleCIDRegex, bundleCID); err != nil {
		return nil, err
	} else if !matched {
		return nil, fmt.Errorf("Invalid check bundle CID [%s]", bundleCID)
	}

	cfg, err := json.Marshal(config)
	if err != nil {
		return nil, err
	}

	result, err := a.Put(bundleCID, cfg)
	if err != nil {
		return nil, err
	}

	checkBundle := &CheckBundle{}
	if err := json.Unmarshal(result, checkBundle); err != nil {
		return nil, err
	}

	return checkBundle, nil
}

// CreateCheckBundle create a new check bundle (check)
func (a *API) CreateCheckBundle(config *CheckBundle) (*CheckBundle, error) {
	if config == nil {
		return nil, fmt.Errorf("Invalid check bundle config [nil]")
	}

	cfg, err := json.Marshal(config)
	if err != nil {
		return nil, err
	}

	result, err := a.Post(baseCheckBundlePath, cfg)
	if err != nil {
		return nil, err
	}

	checkBundle := &CheckBundle{}
	if err := json.Unmarshal(result, checkBundle); err != nil {
		return nil, err
	}

	return checkBundle, nil
}

// DeleteCheckBundle delete a check bundle
func (a *API) DeleteCheckBundle(config *CheckBundle) (bool, error) {
	if config == nil {
		return false, fmt.Errorf("Invalid check bundle config [nil]")
	}
	return a.DeleteCheckBundleByCID(CIDType(&config.CID))
}

// DeleteCheckBundleByCID delete a check bundle by cid
func (a *API) DeleteCheckBundleByCID(cid CIDType) (bool, error) {

	if cid == nil || *cid == "" {
		return false, fmt.Errorf("Invalid check bundle CID [none]")
	}

	bundleCID := string(*cid)

	if matched, err := regexp.MatchString(checkBundleCIDRegex, bundleCID); err != nil {
		return false, err
	} else if !matched {
		return false, fmt.Errorf("Invalid check bundle CID [%v]", bundleCID)
	}

	_, err := a.Delete(bundleCID)
	if err != nil {
		return false, err
	}

	return true, nil
}

// SearchCheckBundles returns list of annotations matching a search query and/or filter
//    - a search query (see: https://login.circonus.com/resources/api#searching)
//    - a filter (see: https://login.circonus.com/resources/api#filtering)
func (a *API) SearchCheckBundles(searchCriteria *SearchQueryType, filterCriteria *map[string][]string) (*[]CheckBundle, error) {

	q := url.Values{}

	if searchCriteria != nil && *searchCriteria != "" {
		q.Set("search", string(*searchCriteria))
	}

	if filterCriteria != nil && len(*filterCriteria) > 0 {
		for filter, criteria := range *filterCriteria {
			for _, val := range criteria {
				q.Add(filter, val)
			}
		}
	}

	if q.Encode() == "" {
		return a.FetchCheckBundles()
	}

	reqURL := url.URL{
		Path:     baseCheckBundlePath,
		RawQuery: q.Encode(),
	}

	resp, err := a.Get(reqURL.String())
	if err != nil {
		return nil, fmt.Errorf("[ERROR] API call error %+v", err)
	}

	var results []CheckBundle
	if err := json.Unmarshal(resp, &results); err != nil {
		return nil, err
	}

	return &results, nil
}
