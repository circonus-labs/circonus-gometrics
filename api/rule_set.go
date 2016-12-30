// Copyright 2016 Circonus, Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Ruleset API support - Fetch, Create, Delete, Search, and Update
// See: https://login.circonus.com/resources/api/calls/rule_set

package api

import (
	"encoding/json"
	"fmt"
	"net/url"
	"regexp"
)

// RulesetRule defines a ruleset rule
type RulesetRule struct {
	Criteria          string `json:"criteria"`
	Severity          int    `json:"severity"`
	Value             string `json:"value"`
	WindowingDuration int    `json:"windowing_duration,omitempty"`
	WindowingFunction string `json:"windowing_function,omitempty"`
	Wait              int    `json:"wait,omitempty"`
}

// Ruleset defines a ruleset
type Ruleset struct {
	CID           string           `json:"_cid,omitempty"`
	CheckCID      string           `json:"check"`
	ContactGroups map[int][]string `json:"contact_groups"`
	Derive        string           `json:"derive,omitempty"`
	Link          string           `json:"link"`
	MetricName    string           `json:"metric_name"`
	MetricTags    []string         `json:"metric_tags"`
	MetricType    string           `json:"metric_type"`
	Notes         string           `json:"notes"`
	Parent        string           `json:"parent,omitempty"`
	Rules         []RulesetRule    `json:"rules"`
	Tags          []string         `json:"tags"`
}

const (
	baseRulesetPath = "/ruleset"
	rulesetCIDRegex = "^" + baseRulesetPath + "/[0-9]+_.+$"
)

// FetchRuleset retrieves a ruleset definition
func (a *API) FetchRuleset(cid CIDType) (*Ruleset, error) {
	if matched, err := regexp.MatchString(rulesetCIDRegex, string(*cid)); err != nil {
		return nil, err
	} else if !matched {
		return nil, fmt.Errorf("Invalid ruleset CID %v", *cid)
	}

	result, err := a.Get(string(*cid))
	if err != nil {
		return nil, err
	}

	ruleset := new(Ruleset)
	if err := json.Unmarshal(result, ruleset); err != nil {
		return nil, err
	}

	return ruleset, nil
}

// FetchRulesets retrieves all rulesets
func (a *API) FetchRulesets() ([]Ruleset, error) {
	result, err := a.Get(baseRulesetPath)
	if err != nil {
		return nil, err
	}

	var rulesets []Ruleset
	if err := json.Unmarshal(result, &rulesets); err != nil {
		return nil, err
	}

	return rulesets, nil
}

// UpdateRuleset update ruleset definition
func (a *API) UpdateRuleset(config *Ruleset) (*Ruleset, error) {
	if matched, err := regexp.MatchString(rulesetCIDRegex, string(config.CID)); err != nil {
		return nil, err
	} else if !matched {
		return nil, fmt.Errorf("Invalid ruleset CID %v", config.CID)
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

	ruleset := &Ruleset{}
	if err := json.Unmarshal(resp, ruleset); err != nil {
		return nil, err
	}

	return ruleset, nil
}

// CreateRuleset create a new ruleset
func (a *API) CreateRuleset(config *Ruleset) (*Ruleset, error) {
	reqURL := url.URL{
		Path: baseRulesetPath,
	}

	cfg, err := json.Marshal(config)
	if err != nil {
		return nil, err
	}

	resp, err := a.Post(reqURL.String(), cfg)
	if err != nil {
		return nil, err
	}

	ruleset := &Ruleset{}
	if err := json.Unmarshal(resp, ruleset); err != nil {
		return nil, err
	}

	return ruleset, nil
}

// DeleteRuleset delete a ruleset
func (a *API) DeleteRuleset(bundle *Ruleset) (bool, error) {
	cid := CIDType(&bundle.CID)
	return a.DeleteRulesetByCID(cid)
}

// DeleteRulesetByCID delete a ruleset by cid
func (a *API) DeleteRulesetByCID(cid CIDType) (bool, error) {
	if matched, err := regexp.MatchString(rulesetCIDRegex, string(*cid)); err != nil {
		return false, err
	} else if !matched {
		return false, fmt.Errorf("Invalid ruleset CID %v", *cid)
	}

	reqURL := url.URL{
		Path: string(*cid),
	}

	_, err := a.Delete(reqURL.String())
	if err != nil {
		return false, err
	}

	return true, nil
}

// RulesetSearch returns list of rulesets matching a search query and/or filter
//    - a search query (see: https://login.circonus.com/resources/api#searching)
//    - a filter (see: https://login.circonus.com/resources/api#filtering)
func (a *API) RulesetSearch(searchCriteria SearchQueryType, filterCriteria map[string]string) ([]Ruleset, error) {

	if searchCriteria == "" && len(filterCriteria) == 0 {
		return a.FetchRulesets()
	}

	reqURL := url.URL{
		Path: baseRulesetPath,
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

	var results []Ruleset
	if err := json.Unmarshal(resp, &results); err != nil {
		return nil, err
	}

	return results, nil
}
