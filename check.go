package circonusgometrics

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

// use cases:
//
// check [bundle] by submission url
// check [bundle] by *check* id (note, not check_bundle id)
// check [bundle] by search
// create check [bundle]

func (m *CirconusMetrics) initializeTrap() error {
	if m.ready {
		return nil
	}

	m.trapmu.Lock()
	defer m.trapmu.Unlock()

	var err error
	var check *Check
	var checkBundle *CheckBundle
	var broker *Broker

	if m.SubmissionUrl != "" {
		check, err = m.fetchCheckBySubmissionUrl(m.SubmissionUrl)
		if err != nil {
			return err
		}
	} else if m.CheckId != 0 {
		check, err = m.fetchCheckById(m.CheckId)
		if err != nil {
			return err
		}
	} else {
		searchCriteria := fmt.Sprintf("(active:1)(host:\"%s\")(type:\"%s\")(tags:%s)", m.InstanceId, m.checkType, m.SearchTag)
		checkBundle, err = m.checkBundleSearch(searchCriteria)
		if err != nil {
			return err
		}

		if checkBundle == nil {
			// err==nil && checkBundle==nil is "no check bundles matched"
			// an error *should* be returned for any other invalid scenario
			checkBundle, broker, err = m.createNewCheck()
			if err != nil {
				return err
			}
		}
	}

	if checkBundle == nil {
		if check != nil {
			checkBundle, err = m.fetchCheckBundleByCid(check.CheckBundleCid)
			if err != nil {
				return err
			}
		} else {
			return fmt.Errorf("Unable to determine check bundle")
		}
	}

	if broker == nil {
		broker, err = m.fetchBrokerByCid(checkBundle.Brokers[0])
		if err != nil {
			return err
		}
	}

	// retain to facilitate metric management (adding new metrics specifically)
	m.checkBundle = checkBundle
	// url to which metrics should be PUT
	m.trapUrl = checkBundle.Config.SubmissionUrl
	// used when sending as "ServerName" get around certs not having IP SANS
	// (cert created with server name as CN but IP used in trap url)
	cn, err := m.getBrokerCN(broker, m.trapUrl)
	if err != nil {
		return err
	}
	m.trapCN = cn
	// all ready, flush can send metrics
	m.ready = true

	// inventory actie metrics
	for _, metric := range checkBundle.Metrics {
		if metric.Status == "active" {
			m.activeMetrics[metric.Name] = true
		}
	}

	return nil
}

func (m *CirconusMetrics) checkBundleSearch(criteria string) (*CheckBundle, error) {
	checkBundles, err := m.searchCheckBundles(criteria)
	if err != nil {
		return nil, err
	}

	if len(checkBundles) == 0 {
		return nil, nil // trigger creation of a new check
		// return nil, fmt.Errorf("No checks found matching criteria %s", searchCriteria)
	}

	numActive := 0
	checkId := -1

	for idx, check := range checkBundles {
		if check.Status == "active" {
			numActive++
			checkId = idx
		}
	}

	if numActive > 1 {
		return nil, fmt.Errorf("Multiple possibilities multiple check bundles match criteria %s\n", criteria)
	}

	return &checkBundles[checkId], nil

}

func (m *CirconusMetrics) createNewCheck() (*CheckBundle, *Broker, error) {
	checkSecret := m.CheckSecret
	if checkSecret == "" {
		secret, err := makeSecret()
		if err != nil {
			secret = "myS3cr3t"
		}
		checkSecret = secret
	}

	broker, brokerErr := m.getBroker()
	if brokerErr != nil {
		return nil, nil, brokerErr
	}

	config := CheckBundle{
		Brokers:     []string{broker.Cid},
		Config:      CheckBundleConfig{AsyncMetrics: true, Secret: checkSecret},
		DisplayName: fmt.Sprintf("%s /%s", m.InstanceId, m.checkType),
		Metrics: []CheckBundleMetric{
			CheckBundleMetric{
				Name:   "cgmplaceholder",
				Status: "active",
				Type:   "numeric",
			},
		},
		MetricLimit: 0,
		Notes:       "",
		Period:      60,
		Status:      "active",
		Tags:        append([]string{m.SearchTag}, m.Tags...),
		Target:      m.InstanceId,
		Timeout:     10,
		Type:        m.checkType,
	}

	checkBundle, err := m.createCheckBundle(config)
	if err != nil {
		return nil, nil, err
	}

	return checkBundle, broker, nil
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

func (m *CirconusMetrics) addNewCheckMetrics(newMetrics map[string]*CheckBundleMetric) {
	newCheckBundle := m.checkBundle
	numCurrMetrics := len(newCheckBundle.Metrics)
	numNewMetrics := len(newMetrics)

	if numCurrMetrics+numNewMetrics >= cap(newCheckBundle.Metrics) {
		nm := make([]CheckBundleMetric, numCurrMetrics+numNewMetrics)
		copy(nm, newCheckBundle.Metrics)
		newCheckBundle.Metrics = nm
	}

	newCheckBundle.Metrics = newCheckBundle.Metrics[0 : numCurrMetrics+numNewMetrics]

	i := 0
	for _, metric := range newMetrics {
		newCheckBundle.Metrics[numCurrMetrics+i] = *metric
		i++
	}

	checkBundle, err := m.updateCheckBundle(newCheckBundle)
	if err != nil {
		m.Log.Printf("[ERROR] updating check bundle with new metrics %v", err)
	}

	for _, metric := range newMetrics {
		m.activeMetrics[metric.Name] = true
	}

	m.checkBundle = checkBundle
}
