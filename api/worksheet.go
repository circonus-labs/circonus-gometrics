// Copyright 2016 Circonus, Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Worksheet API support - Fetch, Create, Delete, and Update
// See: https://login.circonus.com/resources/api/calls/worksheet

package api

import (
	"encoding/json"
	"fmt"
	"net/url"
	"regexp"
)

// WorksheetGraph defines a worksheet cid to be include in the worksheet
type WorksheetGraph struct {
	GraphCID string `json:"graph"`
}

// WorksheetSmartQuery defines a query to include multiple worksheets
type WorksheetSmartQuery struct {
	Name  string   `json:"name"`
	Query string   `json:"query"`
	Order []string `json:"order"`
}

// Worksheet defines a worksheet
type Worksheet struct {
	CID          string                `json:"_cid,omitempty"`
	Description  string                `json:"description"`
	Favorite     bool                  `json:"favorite"`
	Graphs       []WorksheetGraph      `json:"worksheets,omitempty"`
	Notes        string                `json:"notes"`
	SmartQueries []WorksheetSmartQuery `json:"smart_queries,omitempty"`
	Tags         []string              `json:"tags"`
	Title        string                `json:"title"`
}

const (
	baseWorksheetPath = "/worksheet"
	worksheetCIDRegex = "^" + baseWorksheetPath + "/[[:xdigit:]]{8}-[[:xdigit:]]{4}-[[:xdigit:]]{4}-[[:xdigit:]]{4}-[[:xdigit:]]{8,12}$"
)

// FetchWorksheet retrieves a worksheet definition
func (a *API) FetchWorksheet(cid CIDType) (*Worksheet, error) {
	if matched, err := regexp.MatchString(worksheetCIDRegex, string(cid)); err != nil {
		return nil, err
	} else if !matched {
		return nil, fmt.Errorf("Invalid worksheet CID %v", cid)
	}

	result, err := a.Get(string(cid))
	if err != nil {
		return nil, err
	}

	worksheet := new(Worksheet)
	if err := json.Unmarshal(result, worksheet); err != nil {
		return nil, err
	}

	return worksheet, nil
}

// FetchWorksheets retrieves all worksheets
func (a *API) FetchWorksheets() ([]Worksheet, error) {
	result, err := a.Get(baseWorksheetPath)
	if err != nil {
		return nil, err
	}

	var worksheets []Worksheet
	if err := json.Unmarshal(result, &worksheets); err != nil {
		return nil, err
	}

	return worksheets, nil
}

// UpdateWorksheet update worksheet definition
func (a *API) UpdateWorksheet(config *Worksheet) (*Worksheet, error) {
	if matched, err := regexp.MatchString(worksheetCIDRegex, string(config.CID)); err != nil {
		return nil, err
	} else if !matched {
		return nil, fmt.Errorf("Invalid worksheet CID %v", config.CID)
	}

	reqURL := url.URL{
		Path: config.CID,
	}

	cfg, err := json.Marshal(config)
	if err != nil {
		return nil, err
	}

	resp, err := a.Put(reqURL.String(), cfg)
	if err != nil {
		return nil, err
	}

	worksheet := &Worksheet{}
	if err := json.Unmarshal(resp, worksheet); err != nil {
		return nil, err
	}

	return worksheet, nil
}

// CreateWorksheet create a new worksheet
func (a *API) CreateWorksheet(config *Worksheet) (*Worksheet, error) {
	reqURL := url.URL{
		Path: baseWorksheetPath,
	}

	cfg, err := json.Marshal(config)
	if err != nil {
		return nil, err
	}

	resp, err := a.Post(reqURL.String(), cfg)
	if err != nil {
		return nil, err
	}

	worksheet := &Worksheet{}
	if err := json.Unmarshal(resp, worksheet); err != nil {
		return nil, err
	}

	return worksheet, nil
}

// DeleteWorksheet delete a worksheet
func (a *API) DeleteWorksheet(bundle *Worksheet) (bool, error) {
	cid := CIDType(bundle.CID)
	return a.DeleteWorksheetByCID(cid)
}

// DeleteWorksheetByCID delete a worksheet by cid
func (a *API) DeleteWorksheetByCID(cid CIDType) (bool, error) {
	if matched, err := regexp.MatchString(worksheetCIDRegex, string(cid)); err != nil {
		return false, err
	} else if !matched {
		return false, fmt.Errorf("Invalid worksheet CID %v", cid)
	}

	reqURL := url.URL{
		Path: string(cid),
	}

	_, err := a.Delete(reqURL.String())
	if err != nil {
		return false, err
	}

	return true, nil
}

// WorksheetSearch returns list of worksheets matching a search query and/or filter
//    - a search query (see: https://login.circonus.com/resources/api#searching)
//    - a filter (see: https://login.circonus.com/resources/api#filtering)
func (a *API) WorksheetSearch(searchCriteria SearchQueryType, filterCriteria map[string]string) ([]Worksheet, error) {

	if searchCriteria == "" && len(filterCriteria) == 0 {
		return a.FetchWorksheets()
	}

	reqURL := url.URL{
		Path: baseWorksheetPath,
	}

	q := url.Values{}

	if searchCriteria != "" {
		q.Set("search", string(searchCriteria))
	}

	if len(filterCriteria) > 0 {
		for filter, criteria := range filterCriteria {
			q.Set(filter, criteria)
		}
	}

	reqURL.RawQuery = q.Encode()

	resp, err := a.Get(reqURL.String())
	if err != nil {
		return nil, fmt.Errorf("[ERROR] API call error %+v", err)
	}

	var results []Worksheet
	if err := json.Unmarshal(resp, &results); err != nil {
		return nil, err
	}

	return results, nil
}
