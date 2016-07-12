package api

// abstracted in preparation of separate circonus-api-go package

import (
	"encoding/json"
	"fmt"
)

// BrokerIDType broker id - e.g. 1234 if _cid were '/broker/1234'
type BrokerIDType int

// BrokerCIDType full broker cid - e.g. '/broker/1234'
type BrokerCIDType string

// BrokerSearchTagType search/select tag type
type BrokerSearchTagType string

// BrokerDetail instance attributes
type BrokerDetail struct {
	CN      string   `json:"cn"`
	IP      string   `json:"ipaddress"`
	MinVer  int      `json:"minimum_version_required"`
	Modules []string `json:"modules"`
	Port    int      `json:"port"`
	Skew    string   `json:"skew"`
	Status  string   `json:"status"`
	Version int      `json:"version"`
}

// Broker definition
type Broker struct {
	Cid       string         `json:"_cid"`
	Details   []BrokerDetail `json:"_details"`
	Latitude  string         `json:"_latitude"`
	Longitude string         `json:"_longitude"`
	Name      string         `json:"_name"`
	Tags      []string       `json:"_tags"`
	Type      string         `json:"_type"`
}

// FetchBrokerByID Use Circonus API to retrieve  a specific broker by Broker Group ID
func (a *API) FetchBrokerByID(id BrokerIDType) (*Broker, error) {
	cid := BrokerCIDType(fmt.Sprintf("/broker/%d", id))
	return a.FetchBrokerByCID(cid)
}

// FetchBrokerByCID Use Circonus API to retreive a broker by CID
func (a *API) FetchBrokerByCID(cid BrokerCIDType) (*Broker, error) {
	result, err := a.Get(string(cid))
	if err != nil {
		return nil, err
	}

	response := new(Broker)
	if err := json.Unmarshal(result, &response); err != nil {
		return nil, err
	}

	return response, nil

}

// FetchBrokerListByTag Use Circonus API to retreive a list of brokers which have a specific tag
func (a *API) FetchBrokerListByTag(searchTag BrokerSearchTagType) ([]Broker, error) {
	result, err := a.Get(fmt.Sprintf("/broker?f__tags_has=%s", searchTag))
	if err != nil {
		return nil, err
	}

	var response []Broker
	json.Unmarshal(result, &response)

	return response, nil
}

// FetchBrokerList Use Circonus API to retreive a list of brokers
func (a *API) FetchBrokerList() ([]Broker, error) {
	result, err := a.Get("/broker")
	if err != nil {
		return nil, err
	}

	var response []Broker
	json.Unmarshal(result, &response)

	return response, nil
}
