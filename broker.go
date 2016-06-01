package circonusgometrics

import (
	"fmt"
	"math/rand"
	"net"
	"net/url"
	"reflect"
	"strings"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func (m *CirconusMetrics) getBroker() (*Broker, error) {
	if m.BrokerGroupId != 0 {
		broker, err := m.fetchBrokerById(m.BrokerGroupId)
		if err != nil {
			return nil, fmt.Errorf("Error fetching designated broker %d\n", m.BrokerGroupId)
		}
		if !m.isValidBroker(broker) {
			return nil, fmt.Errorf("Error designated broker %d [%s] is invalid (not active or does not support required check type).\n", m.BrokerGroupId, broker.Name)
		}
		return broker, nil
	}
	broker, err := m.selectBroker()
	if err != nil {
		return nil, fmt.Errorf("Unable to fetch suitable broker %s", err)
	}
	return broker, nil
}

func (m *CirconusMetrics) getBrokerCN(broker *Broker, submissionUrl string) (string, error) {
	u, err := url.Parse(submissionUrl)
	if err != nil {
		return "", err
	}

	hostParts := strings.Split(u.Host, ":")
	host := hostParts[0]

	if net.ParseIP(host) == nil { // it's a non-ip string
		return u.Host, nil
	}

	cn := ""

	for _, detail := range broker.Details {
		if detail.IP == host {
			cn = detail.CN
			break
		}
	}

	if cn == "" {
		return "", fmt.Errorf("Unable to match URL host (%s) to Broker", u.Host)
	}

	return cn, nil

}

func (m *CirconusMetrics) selectBroker() (*Broker, error) {
	brokerList, err := m.fetchBrokerList()
	if err != nil {
		return nil, err
	}

	validBrokers := make(map[string]Broker)
	haveEnterprise := false

	for _, broker := range brokerList {
		if m.isValidBroker(&broker) {
			validBrokers[broker.Cid] = broker
			if broker.Type == "enterprise" {
				haveEnterprise = true
			}
		}
	}

	if haveEnterprise { // eliminate non-enterprise brokers from valid brokers
		for k, v := range validBrokers {
			if v.Type != "enterprise" {
				delete(validBrokers, k)
			}
		}
	}

	validBrokerKeys := reflect.ValueOf(validBrokers).MapKeys()
	selectedBroker := validBrokers[validBrokerKeys[rand.Intn(len(validBrokerKeys))].String()]

	return &selectedBroker, nil

}

func (m *CirconusMetrics) brokerSupportsCheckType(checkType string, details *BrokerDetail) bool {

	for _, module := range details.Modules {
		if module == checkType {
			return true
		}
	}

	return false

}

func (m *CirconusMetrics) isValidBroker(broker *Broker) bool {
	valid := false
	for _, detail := range broker.Details {
		if detail.Status != "active" {
			continue
		}

		if m.brokerSupportsCheckType(m.checkType, &detail) {
			valid = true
			break
		}
	}
	return valid
}
