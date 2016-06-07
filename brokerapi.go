package circonusgometrics

// abstracted in preparation of separate circonus-api-go package

import (
	"encoding/json"
	"fmt"
)

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

type Broker struct {
	Cid       string         `json:"_cid"`
	Details   []BrokerDetail `json:"_details"`
	Latitude  string         `json:"_latitude"`
	Longitude string         `json:"_longitude"`
	Name      string         `json:"_name"`
	Tags      []string       `json:"_tags"`
	Type      string         `json:"_type"`
}

// Use Circonus API to retrieve  a specific broker by Broker Group ID
func (m *CirconusMetrics) fetchBrokerById(id int) (*Broker, error) {
	cid := fmt.Sprintf("/v2/broker/%d", id)
	return m.fetchBrokerByCid(cid)
}

// Use Circonus API to retreive a broker by CID
func (m *CirconusMetrics) fetchBrokerByCid(cid string) (*Broker, error) {
	result, err := m.apiCall("GET", cid, nil)
	if err != nil {
		return nil, err
	}

	response := new(Broker)
	if err := json.Unmarshal(result, &response); err != nil {
		return nil, err
	}

	return response, nil

}

// Use Circonus API to retreive a list of brokers
func (m *CirconusMetrics) fetchBrokerList() ([]Broker, error) {
	result, err := m.apiCall("GET", "/v2/broker", nil)
	if err != nil {
		return nil, err
	}

	var response []Broker
	json.Unmarshal(result, &response)

	return response, nil
}
