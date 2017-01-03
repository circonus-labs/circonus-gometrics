// Copyright 2016 Circonus, Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Alert API support - Fetch and Search
// See: https://login.circonus.com/resources/api/calls/alert

package api

import (
	"encoding/json"
	"fmt"
	"net/url"
	"regexp"
)

// Alert defines a alert
type Alert struct {
	CID                string   `json:"_cid,omitempty"`
	AcknowledgementCID string   `json:"_acknowledgement,omitempty"`
	AlertURL           string   `json:"_alert_url,omitempty"`
	BrokerCID          string   `json:"_broker,omitempty"`
	CheckCID           string   `json:"_check,omitempty"`
	CheckName          string   `json:"_check_name,omitempty"`
	ClearedOn          uint     `json:"_cleared_on,omitempty"`
	ClearedValue       string   `json:"_cleared_value,omitempty"`
	Maintenance        []string `json:"_maintenance,omitempty"`
	MetricLinkURL      string   `json:"_metric_link,omitempty"`
	MetricName         string   `json:"_metric_name,omitempty"`
	MetricNotes        string   `json:"_metric_notes,omitempty"`
	OccurredOn         uint     `json:"_occurred_on,omitempty"`
	RuleSetCID         string   `json:"_rule_set,omitempty"`
	Severity           uint     `json:"_severity,omitempty"`
	Tags               []string `json:"_tags,omitempty"`
	Value              string   `json:"_value,omitempty"`
}

const (
	baseAlertPath = "/alert"
	alertCIDRegex = "^" + baseAlertPath + "/[0-9]+$"
)

// FetchAlert retrieves a alert definition
func (a *API) FetchAlert(cid CIDType) (*Alert, error) {
	if cid == nil || *cid == "" {
		return nil, fmt.Errorf("Invalid alert CID [none]")
	}

	alertCID := string(*cid)

	matched, err := regexp.MatchString(alertCIDRegex, alertCID)
	if err != nil {
		return nil, err
	}
	if !matched {
		return nil, fmt.Errorf("Invalid alert CID [%s]", alertCID)
	}

	result, err := a.Get(alertCID)
	if err != nil {
		return nil, err
	}

	alert := &Alert{}
	if err := json.Unmarshal(result, alert); err != nil {
		return nil, err
	}

	return alert, nil
}

// FetchAlerts retrieves all alerts
func (a *API) FetchAlerts() (*[]Alert, error) {
	result, err := a.Get(baseAlertPath)
	if err != nil {
		return nil, err
	}

	var alerts []Alert
	if err := json.Unmarshal(result, &alerts); err != nil {
		return nil, err
	}

	return &alerts, nil
}

// SearchAlerts returns list of alerts matching a search query and/or filter
//    - a search query (see: https://login.circonus.com/resources/api#searching)
//    - a filter (see: https://login.circonus.com/resources/api#filtering)
func (a *API) SearchAlerts(searchCriteria *SearchQueryType, filterCriteria *SearchFilterType) (*[]Alert, error) {
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
		return a.FetchAlerts()
	}

	reqURL := url.URL{
		Path:     baseAlertPath,
		RawQuery: q.Encode(),
	}

	result, err := a.Get(reqURL.String())
	if err != nil {
		return nil, fmt.Errorf("[ERROR] API call error %+v", err)
	}

	var alerts []Alert
	if err := json.Unmarshal(result, &alerts); err != nil {
		return nil, err
	}

	return &alerts, nil
}