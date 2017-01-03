// Copyright 2016 Circonus, Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Acknowledgement API support - Fetch, Create, Update, Delete*, and Search
// See: https://login.circonus.com/resources/api/calls/acknowledgement
// *  : delete (cancel) by updating with AcknowledgedUntil set to 0

package api

import (
	"encoding/json"
	"fmt"
	"net/url"
	"regexp"
)

// Acknowledgement defines a acknowledgement
type Acknowledgement struct {
	CID               string      `json:"_cid,omitempty"`
	AcknowledgedBy    string      `json:"_acknowledged_by,omitempty"`
	AcknowledgedOn    uint        `json:"_acknowledged_on,omitempty"`
	Active            bool        `json:"_active,omitempty"`
	LastModified      uint        `json:"_last_modified,omitempty"`
	LastModifiedBy    string      `json:"_last_modified_by,omitempty"`
	AcknowledgedUntil interface{} `json:"acknowledged_until,omitempty"`
	AlertCID          string      `json:"alert,omitempty"`
	Notes             string      `json:"notes,omitempty"`
}

const (
	baseAcknowledgementPath = "/acknowledgement"
	acknowledgementCIDRegex = "^" + baseAcknowledgementPath + "/[0-9]+$"
)

// FetchAcknowledgement retrieves a acknowledgement definition
func (a *API) FetchAcknowledgement(cid CIDType) (*Acknowledgement, error) {
	if cid == nil || *cid == "" {
		return nil, fmt.Errorf("Invalid acknowledgement CID [none]")
	}

	acknowledgementCID := string(*cid)

	matched, err := regexp.MatchString(acknowledgementCIDRegex, acknowledgementCID)
	if err != nil {
		return nil, err
	}
	if !matched {
		return nil, fmt.Errorf("Invalid acknowledgement CID [%s]", acknowledgementCID)
	}

	result, err := a.Get(acknowledgementCID)
	if err != nil {
		return nil, err
	}

	acknowledgement := &Acknowledgement{}
	if err := json.Unmarshal(result, acknowledgement); err != nil {
		return nil, err
	}

	return acknowledgement, nil
}

// FetchAcknowledgements retrieves all acknowledgements
func (a *API) FetchAcknowledgements() (*[]Acknowledgement, error) {
	result, err := a.Get(baseAcknowledgementPath)
	if err != nil {
		return nil, err
	}

	var acknowledgements []Acknowledgement
	if err := json.Unmarshal(result, &acknowledgements); err != nil {
		return nil, err
	}

	return &acknowledgements, nil
}

// UpdateAcknowledgement update acknowledgement definition
func (a *API) UpdateAcknowledgement(config *Acknowledgement) (*Acknowledgement, error) {
	if config == nil {
		return nil, fmt.Errorf("Invalid acknowledgement config [nil]")
	}

	acknowledgementCID := string(config.CID)

	matched, err := regexp.MatchString(acknowledgementCIDRegex, acknowledgementCID)
	if err != nil {
		return nil, err
	}
	if !matched {
		return nil, fmt.Errorf("Invalid acknowledgement CID [%s]", acknowledgementCID)
	}

	cfg, err := json.Marshal(config)
	if err != nil {
		return nil, err
	}

	result, err := a.Put(acknowledgementCID, cfg)
	if err != nil {
		return nil, err
	}

	acknowledgement := &Acknowledgement{}
	if err := json.Unmarshal(result, acknowledgement); err != nil {
		return nil, err
	}

	return acknowledgement, nil
}

// CreateAcknowledgement create a new acknowledgement
func (a *API) CreateAcknowledgement(config *Acknowledgement) (*Acknowledgement, error) {
	if config == nil {
		return nil, fmt.Errorf("Invalid acknowledgement config [nil]")
	}

	cfg, err := json.Marshal(config)
	if err != nil {
		return nil, err
	}

	result, err := a.Post(baseAcknowledgementPath, cfg)
	if err != nil {
		return nil, err
	}

	acknowledgement := &Acknowledgement{}
	if err := json.Unmarshal(result, acknowledgement); err != nil {
		return nil, err
	}

	return acknowledgement, nil
}

// SearchAcknowledgements returns list of acknowledgements matching a search query and/or filter
//    - a search query (see: https://login.circonus.com/resources/api#searching)
//    - a filter (see: https://login.circonus.com/resources/api#filtering)
func (a *API) SearchAcknowledgements(searchCriteria *SearchQueryType, filterCriteria *SearchFilterType) (*[]Acknowledgement, error) {
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
		return a.FetchAcknowledgements()
	}

	reqURL := url.URL{
		Path:     baseAcknowledgementPath,
		RawQuery: q.Encode(),
	}

	result, err := a.Get(reqURL.String())
	if err != nil {
		return nil, fmt.Errorf("[ERROR] API call error %+v", err)
	}

	var acknowledgements []Acknowledgement
	if err := json.Unmarshal(result, &acknowledgements); err != nil {
		return nil, err
	}

	return &acknowledgements, nil
}