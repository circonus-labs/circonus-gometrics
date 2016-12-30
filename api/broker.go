// Copyright 2016 Circonus, Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package api

import (
	"encoding/json"
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

// BrokerDetail instance attributes
type BrokerDetail struct {
	CN           string   `json:"cn"`
	ExternalHost string   `json:"external_host"`
	ExternalPort int      `json:"external_port"`
	IP           string   `json:"ipaddress"`
	MinVer       int      `json:"minimum_version_required"`
	Modules      []string `json:"modules"`
	Port         int      `json:"port"`
	Skew         string   `json:"skew"`
	Status       string   `json:"status"`
	Version      int      `json:"version"`
}

// Broker definition
type Broker struct {
	CID       string         `json:"_cid"`
	Details   []BrokerDetail `json:"_details"`
	Latitude  string         `json:"_latitude"`
	Longitude string         `json:"_longitude"`
	Name      string         `json:"_name"`
	Tags      []string       `json:"_tags"`
	Type      string         `json:"_type"`
}

const (
	baseBrokerPath = "/broker"
	brokerCIDRegex = "^" + baseBrokerPath + "/[0-9]+$"
)

// FetchBrokerByID fetch a broker configuration by [group]id
func (a *API) FetchBrokerByID(id IDType) (*Broker, error) {
	if id <= 0 {
		return nil, fmt.Errorf("Invalid broker ID [%d]", id)
	}
	cid := CIDType(fmt.Sprintf("%s/%d", baseBrokerPath, id))
	return a.FetchBroker(&cid)
}

// FetchBroker fetch a broker configuration by cid
func (a *API) FetchBroker(cid *CIDType) (*Broker, error) {
	if cid == nil || *cid == "" {
		return nil, fmt.Errorf("Invalid broker CID [none]")
	}

	brokerCID := string(*cid)

	matched, err := regexp.MatchString(brokerCIDRegex, brokerCID)
	if err != nil {
		return nil, err
	}
	if !matched {
		return nil, fmt.Errorf("Invalid broker CID [%s]", brokerCID)
	}

	reqURL := url.URL{
		Path: brokerCID,
	}

	result, err := a.Get(reqURL.String())
	if err != nil {
		return nil, err
	}

	response := new(Broker)
	if err := json.Unmarshal(result, &response); err != nil {
		return nil, err
	}

	return response, nil

}

// FetchBrokersByTag return list of brokers with a specific tag
func (a *API) FetchBrokersByTag(searchTags TagType) (*[]Broker, error) {
	if len(searchTags) == 0 {
		return a.FetchBrokers()
	}

	filter := map[string]string{
		"f__tags_has": strings.Replace(strings.Join(searchTags, ","), ",", "&f__tags_has=", -1),
	}

	return a.SearchBrokers(nil, &filter)
}

// // BrokerSearch return a list of brokers matching a query/filter
// func (a *API) BrokerSearch(query SearchQueryType) ([]Broker, error) {
// 	queryURL := fmt.Sprintf("/broker?%s", string(query))
//
// 	result, err := a.Get(queryURL)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	var brokers []Broker
// 	if err := json.Unmarshal(result, &brokers); err != nil {
// 		return nil, err
// 	}
//
// 	return brokers, nil
// }

// SearchBrokers returns list of annotations matching a search query and/or filter
//    - a search query (see: https://login.circonus.com/resources/api#searching)
//    - a filter (see: https://login.circonus.com/resources/api#filtering)
func (a *API) SearchBrokers(searchCriteria *SearchQueryType, filterCriteria *map[string]string) (*[]Broker, error) {

	if (searchCriteria == nil || *searchCriteria == "") && (filterCriteria == nil || len(*filterCriteria) == 0) {
		return a.FetchBrokers()
	}

	reqURL := url.URL{
		Path: baseBrokerPath,
	}

	q := url.Values{}

	if searchCriteria != nil && *searchCriteria != "" {
		q.Set("search", string(*searchCriteria))
	}

	if filterCriteria != nil && len(*filterCriteria) > 0 {
		for filter, criteria := range *filterCriteria {
			q.Set(filter, criteria)
		}
	}

	reqURL.RawQuery = q.Encode()

	resp, err := a.Get(reqURL.String())
	if err != nil {
		return nil, fmt.Errorf("[ERROR] API call error %+v", err)
	}

	var results []Broker
	if err := json.Unmarshal(resp, &results); err != nil {
		return nil, err
	}

	return &results, nil
}

// FetchBrokers return list of all brokers available to the api token/app
func (a *API) FetchBrokers() (*[]Broker, error) {
	result, err := a.Get("/broker")
	if err != nil {
		return nil, err
	}

	var response []Broker
	if err := json.Unmarshal(result, &response); err != nil {
		return nil, err
	}

	return &response, nil
}
