// Copyright 2016 Circonus, Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// RulesetGroup API support - Fetch, Create, Delete, Search, and Update
// See: https://login.circonus.com/resources/api/calls/rule_set_group

package api

import (
	"encoding/json"
	"fmt"
	"net/url"
	"regexp"
)

// RulesetGroupRule defines a rulesetGroup rule
type RulesetGroupRule struct {
	Criteria          string `json:"criteria"`
	Severity          int    `json:"severity"`
	Value             string `json:"value"`
	WindowingDuration int    `json:"windowing_duration,omitempty"`
	WindowingFunction string `json:"windowing_function,omitempty"`
	Wait              int    `json:"wait,omitempty"`
}

// RulesetGroupFormula defines a formula for raising alerts
type RulesetGroupFormula struct {
	Expression    string `json:"expression"`
	RaiseSeverity int    `json:"raise_severity"`
	Wait          int    `json:"wait"`
}

// RulesetGroupCondition defines conditions for raising alerts
type RulesetGroupCondition struct {
	MatchingSeverities []string `json:"matching_serverities"`
	RulesetCID         string   `json:"rule_set"`
}

// RulesetGroup defines a ruleset group
type RulesetGroup struct {
	CID               string                  `json:"_cid,omitempty"`
	ContactGroups     map[int][]string        `json:"contact_groups"`
	Formulas          []RulesetGroupFormula   `json:"formulas"`
	Name              string                  `json:"name"`
	RulesetConditions []RulesetGroupCondition `json:"rule_set_conditions"`
	Tags              []string                `json:"tags"`
}

const (
	baseRulesetGroupPath = "/rule_set_group"
	rulesetGroupCIDRegex = "^" + baseRulesetGroupPath + "/[0-9]+$"
)

// FetchRulesetGroup retrieves a rulesetGroup definition
func (a *API) FetchRulesetGroup(cid CIDType) (*RulesetGroup, error) {
	if matched, err := regexp.MatchString(rulesetGroupCIDRegex, string(*cid)); err != nil {
		return nil, err
	} else if !matched {
		return nil, fmt.Errorf("Invalid rule set group CID %v", *cid)
	}

	result, err := a.Get(string(*cid))
	if err != nil {
		return nil, err
	}

	rulesetGroup := new(RulesetGroup)
	if err := json.Unmarshal(result, rulesetGroup); err != nil {
		return nil, err
	}

	return rulesetGroup, nil
}

// FetchRulesetGroups retrieves all rulesetGroups
func (a *API) FetchRulesetGroups() ([]RulesetGroup, error) {
	result, err := a.Get(baseRulesetGroupPath)
	if err != nil {
		return nil, err
	}

	var rulesetGroups []RulesetGroup
	if err := json.Unmarshal(result, &rulesetGroups); err != nil {
		return nil, err
	}

	return rulesetGroups, nil
}

// UpdateRulesetGroup update rulesetGroup definition
func (a *API) UpdateRulesetGroup(config *RulesetGroup) (*RulesetGroup, error) {
	if matched, err := regexp.MatchString(rulesetGroupCIDRegex, string(config.CID)); err != nil {
		return nil, err
	} else if !matched {
		return nil, fmt.Errorf("Invalid rule set group CID %v", config.CID)
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

	rulesetGroup := &RulesetGroup{}
	if err := json.Unmarshal(resp, rulesetGroup); err != nil {
		return nil, err
	}

	return rulesetGroup, nil
}

// CreateRulesetGroup create a new rulesetGroup
func (a *API) CreateRulesetGroup(config *RulesetGroup) (*RulesetGroup, error) {
	reqURL := url.URL{
		Path: baseRulesetGroupPath,
	}

	cfg, err := json.Marshal(config)
	if err != nil {
		return nil, err
	}

	resp, err := a.Post(reqURL.String(), cfg)
	if err != nil {
		return nil, err
	}

	rulesetGroup := &RulesetGroup{}
	if err := json.Unmarshal(resp, rulesetGroup); err != nil {
		return nil, err
	}

	return rulesetGroup, nil
}

// DeleteRulesetGroup delete a rulesetGroup
func (a *API) DeleteRulesetGroup(bundle *RulesetGroup) (bool, error) {
	cid := CIDType(&bundle.CID)
	return a.DeleteRulesetGroupByCID(cid)
}

// DeleteRulesetGroupByCID delete a rulesetGroup by cid
func (a *API) DeleteRulesetGroupByCID(cid CIDType) (bool, error) {
	if matched, err := regexp.MatchString(rulesetGroupCIDRegex, string(*cid)); err != nil {
		return false, err
	} else if !matched {
		return false, fmt.Errorf("Invalid rule set group CID %v", cid)
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

// RulesetGroupSearch returns list of rulesetGroups matching a search query and/or filter
//    - a search query (see: https://login.circonus.com/resources/api#searching)
//    - a filter (see: https://login.circonus.com/resources/api#filtering)
func (a *API) RulesetGroupSearch(searchCriteria SearchQueryType, filterCriteria map[string]string) ([]RulesetGroup, error) {

	if searchCriteria == "" && len(filterCriteria) == 0 {
		return a.FetchRulesetGroups()
	}

	reqURL := url.URL{
		Path: baseRulesetGroupPath,
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

	var results []RulesetGroup
	if err := json.Unmarshal(resp, &results); err != nil {
		return nil, err
	}

	return results, nil
}
