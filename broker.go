package circonusgometrics

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
	CID       string         `json:"_cid"`
	Details   []BrokerDetail `json:"_details"`
	Latitude  string         `json:"_latitude"`
	Longitude string         `json:"_longitude"`
	Name      string         `json:"name"`
	Tags      []string       `json:"tags"`
	Type      string         `json:"type"`
}

func (m *CirconusMetrics) getBrokerGroupId() int {
	bgi := m.BrokerGroupId
	if bgi == 0 {
		bgi = circonusTrapBroker
		// add lookup and leverage selection method unless we can add a method for the account admin to designate a "default" broker (maybe an overall then potentially by check type so that
		// if POST check_bundle {..., brokers: [], ...} is sent, the default broker for the checktype or the generic dfeault broker would be used)
	}
	return bgi
}
