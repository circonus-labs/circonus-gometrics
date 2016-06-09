package circonusgometrics

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Initialize CirconusMetrics instance. Attempt to find a check otherwise create one.
// use cases:
//
// check [bundle] by submission url
// check [bundle] by *check* id (note, not check_bundle id)
// check [bundle] by search
// create check [bundle]
func (m *CirconusMetrics) initializeTrap() error {
	m.trapmu.Lock()
	defer m.trapmu.Unlock()

	if m.ready {
		return nil
	}

	// short-circuit for non-ssl submission urls
	if m.SubmissionUrl != "" && m.SubmissionUrl[0:5] == "http:" {
		m.trapUrl = m.SubmissionUrl
		m.trapSSL = false
		m.ready = true
		m.trapLastUpdate = time.Now()
		return nil
	}

	var err error
	var check *Check
	var checkBundle *CheckBundle
	var broker *Broker

	if m.SubmissionUrl != "" {
		check, err = m.fetchCheckBySubmissionUrl(m.SubmissionUrl)
		if err != nil {
			return err
		}
		// extract check id from check object returned from looking up using submission url
		// set m.CheckId to the id
		// set m.SubmissionUrl to "" to prevent trying to search on it going forward
		// use case: if the broker is changed in the UI metrics would stop flowing
		// unless the new submission url can be fetched with the API (which is no
		// longer possible using the original submission url)
		id, err := strconv.Atoi(strings.Replace(check.Cid, "/check/", "", -1))
		if err == nil {
			m.CheckId = id
			m.SubmissionUrl = ""
		} else {
			m.Log.Printf("SubmissionUrl check to Check ID: unable to convert %s to int %q\n", check.Cid, err)
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

	// mark for SSL
	if m.trapUrl[0:6] == "https:" {
		m.trapSSL = true
	}

	// load the CA certificate for the broker hosting the submission url
	m.loadCACert()

	// used when sending as "ServerName" get around certs not having IP SANS
	// (cert created with server name as CN but IP used in trap url)
	cn, err := m.getBrokerCN(broker, m.trapUrl)
	if err != nil {
		return err
	}
	m.trapCN = cn

	m.trapLastUpdate = time.Now()

	// all ready, flush can send metrics
	m.ready = true

	// inventory active metrics
	m.activeMetrics = make(map[string]bool)
	for _, metric := range checkBundle.Metrics {
		if metric.Status == "active" {
			m.activeMetrics[metric.Name] = true
		}
	}

	return nil
}

// Search for a check bundle given a predetermined set of criteria
func (m *CirconusMetrics) checkBundleSearch(criteria string) (*CheckBundle, error) {
	checkBundles, err := m.searchCheckBundles(criteria)
	if err != nil {
		return nil, err
	}

	if len(checkBundles) == 0 {
		return nil, nil // trigger creation of a new check
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

// Create a new check to receive metrics
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

// Create a dynamic secret to use with a new check
func makeSecret() (string, error) {
	hash := sha256.New()
	x := make([]byte, 2048)
	if _, err := rand.Read(x); err != nil {
		return "", err
	}
	hash.Write(x)
	return hex.EncodeToString(hash.Sum(nil))[0:16], nil
}

// Add new metrics to an existing check
func (m *CirconusMetrics) addNewCheckMetrics(newMetrics map[string]*CheckBundleMetric) {
	// only manage metrics if checkBundle has been populated
	if m.checkBundle == nil {
		return
	}

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
		return
	}

	for _, metric := range newMetrics {
		m.activeMetrics[metric.Name] = true
	}

	m.checkBundle = checkBundle
}
