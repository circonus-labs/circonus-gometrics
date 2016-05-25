package circonusgometrics

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
)

type CheckConfig struct {
	AsyncMetrics  bool   `json:"async_metrics"`
	Secret        string `json:"secret"`
	SubmissionUrl string `json:"submission_url"`
}

type CheckMetric struct {
	Name   string `json:"name"`
	Type   string `json:"type"`
	Units  string `json:"units"`
	Status string `json:"status"`
}

type Check struct {
	CheckUUIDs         []string `json:"_check_uuids,omitempty"`
	Checks             []string `json:"_checks,omitempty"`
	CID                string   `json:"_cid,omitempty"`
	Created            int      `json:"_created,omitempty"`
	LastModified       int      `json:"_last_modified,omitempty"`
	LastModifedBy      string   `json:"_last_modifed_by,omitempty"`
	ReverseConnectUrls []string `json:"_reverse_connection_urls,omitempty"`
	Brokers            []string `json:"brokers"`
	CheckConfig        `json:"config"`
	DisplayName        string        `json:"display_name"`
	Metrics            []CheckMetric `json:"metrics"`
	MetricLimit        int           `json:"metric_limit"`
	Notes              string        `json:"notes"`
	Period             int           `json:"period"`
	Status             string        `json:"status"`
	Tags               []string      `json:"tags"`
	Target             string        `json:"target"`
	Timeout            int           `json:"timeout"`
	Type               string        `json:"type"`
}

func (m *CirconusMetrics) getTrapUrl() (string, error) {
	if m.TrapUrl != "" {
		return m.TrapUrl, nil
	}

	trapUrl := ""

	check, err := m.apiCheckSearch()
	if err != nil {
		m.Log.Fatalf("%+v\n", err)
	}

	if check == nil {
		m.Log.Println("call create check")
		check, err = m.createCheck()
		if err != nil {
			m.Log.Fatalf("%+v\n", err)
		}
	}

	trapUrl = check.CheckConfig.SubmissionUrl

	return trapUrl, nil

}

func (m *CirconusMetrics) createCheck() (*Check, error) {
	checkSecret := m.CheckSecret
	if checkSecret == "" {
		secret, err := makeSecret()
		if err != nil {
			secret = "myS3cr3t"
		}
		checkSecret = secret
	}

	config := &Check{
		Brokers:     []string{fmt.Sprintf("/broker/%d", m.getBrokerGroupId())},
		CheckConfig: CheckConfig{AsyncMetrics: true, Secret: checkSecret},
		DisplayName: fmt.Sprintf("%s /%s", m.InstanceId, checkType),
		Metrics: []CheckMetric{
			CheckMetric{
				Name:   "placeholder",
				Status: "active",
				Type:   "numeric",
			},
		},
		MetricLimit: 0,
		Notes:       "",
		Period:      60,
		Status:      "active",
		Tags:        append([]string{circonusSearchTag}, m.Tags...),
		Target:      m.InstanceId,
		Timeout:     10,
		Type:        checkType,
	}

	cfgJson, err := json.Marshal(config)
	if err != nil {
		return nil, err
	}

	// m.Log.Printf(string(cfgJson))

	check, err := m.apiCreateCheck(cfgJson)
	if err != nil {
		return nil, err
	}

	return check, nil
}

func makeSecret() (string, error) {
	hash := sha256.New()
	x := make([]byte, 2048)
	if _, err := rand.Read(x); err != nil {
		return "", err
	}
	hash.Write(x)
	return hex.EncodeToString(hash.Sum(nil))[0:16], nil
}
