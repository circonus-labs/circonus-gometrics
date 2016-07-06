package checkmgr

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/circonus-labs/circonus-gometrics/api"
)

// Initialize CirconusMetrics instance. Attempt to find a check otherwise create one.
// use cases:
//
// check [bundle] by submission url
// check [bundle] by *check* id (note, not check_bundle id)
// check [bundle] by search
// create check [bundle]
func (cm *CheckManager) initializeTrap() error {
	cm.trapmu.Lock()
	defer cm.trapmu.Unlock()

	if cm.ready {
		return nil
	}

	// short-circuit for non-ssl submission urls
	if cm.submissionUrl != "" && cm.submissionUrl[0:5] == "http:" {
		cm.trapUrl = cm.submissionUrl
		cm.trapSSL = false
		cm.ready = true
		cm.trapLastUpdate = time.Now()
		return nil
	}

	var err error
	var check *api.Check
	var checkBundle *api.CheckBundle
	var broker *api.Broker

	if cm.submissionUrl != "" {
		check, err = cm.apih.FetchCheckBySubmissionUrl(cm.submissionUrl)
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
			cm.checkId = id
			cm.submissionUrl = ""
		} else {
			cm.Log.Printf("[WARN] SubmissionUrl check to Check ID: unable to convert %s to int %q\n", check.Cid, err)
		}
	} else if cm.checkId != 0 {
		check, err = cm.apih.FetchCheckById(cm.checkId)
		if err != nil {
			return err
		}
	} else {
		searchCriteria := fmt.Sprintf("(active:1)(host:\"%s\")(type:\"%s\")(tags:%s)", cm.instanceId, cm.checkType, cm.searchTag)
		checkBundle, err = cm.checkBundleSearch(searchCriteria)
		if err != nil {
			return err
		}

		if checkBundle == nil {
			// err==nil && checkBundle==nil is "no check bundles matched"
			// an error *should* be returned for any other invalid scenario
			checkBundle, broker, err = cm.createNewCheck()
			if err != nil {
				return err
			}
		}
	}

	if checkBundle == nil {
		if check != nil {
			checkBundle, err = cm.apih.FetchCheckBundleByCid(check.CheckBundleCid)
			if err != nil {
				return err
			}
		} else {
			return fmt.Errorf("[ERROR] Unable to retrieve, find, or create check.")
		}
	}

	if broker == nil {
		broker, err = cm.apih.FetchBrokerByCid(checkBundle.Brokers[0])
		if err != nil {
			return err
		}
	}

	// retain to facilitate metric management (adding new metrics specifically)
	cm.checkBundle = checkBundle

	// url to which metrics should be PUT
	cm.trapUrl = checkBundle.Config.SubmissionUrl

	// mark for SSL
	if cm.trapUrl[0:6] == "https:" {
		cm.trapSSL = true
	}

	// FIX: this needs to move somewhere else
	// load the CA certificate for the broker hosting the submission url
	//CirconusMetrics.loadCACert()

	// used when sending as "ServerName" get around certs not having IP SANS
	// (cert created with server name as CN but IP used in trap url)
	cn, err := cm.GetBrokerCN(broker, cm.trapUrl)
	if err != nil {
		return err
	}
	cm.trapCN = cn

	cm.trapLastUpdate = time.Now()

	// all ready, flush can send metrics
	cm.ready = true

	// inventory active metrics
	cm.activeMetrics = make(map[string]bool)
	for _, metric := range checkBundle.Metrics {
		if metric.Status == "active" {
			cm.activeMetrics[metric.Name] = true
		}
	}

	return nil
}

// Search for a check bundle given a predetermined set of criteria
func (cm *CheckManager) checkBundleSearch(criteria string) (*api.CheckBundle, error) {
	checkBundles, err := cm.apih.SearchCheckBundles(criteria)
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
		return nil, fmt.Errorf("[ERROR] Multiple possibilities multiple check bundles match criteria %s\n", criteria)
	}

	return &checkBundles[checkId], nil
}

// Create a new check to receive metrics
func (cm *CheckManager) createNewCheck() (*api.CheckBundle, *api.Broker, error) {
	checkSecret := cm.checkSecret
	if checkSecret == "" {
		secret, err := makeSecret()
		if err != nil {
			secret = "myS3cr3t"
		}
		checkSecret = secret
	}

	broker, brokerErr := cm.GetBroker()
	if brokerErr != nil {
		return nil, nil, brokerErr
	}

	config := api.CheckBundle{
		Brokers:     []string{broker.Cid},
		Config:      api.CheckBundleConfig{AsyncMetrics: true, Secret: checkSecret},
		DisplayName: fmt.Sprintf("%s /%s", cm.instanceId, cm.checkType),
		Metrics:     []api.CheckBundleMetric{},
		/*
				api.CheckBundleMetric{
					Name:   "cgmplaceholder",
					Status: "active",
					Type:   "numeric",
				},
			},
		*/
		MetricLimit: 0,
		Notes:       "",
		Period:      60,
		Status:      "active",
		Tags:        append([]string{cm.searchTag}, cm.tags...),
		Target:      cm.instanceId,
		Timeout:     10,
		Type:        cm.checkType,
	}

	checkBundle, err := cm.apih.CreateCheckBundle(config)
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
func (cm *CheckManager) addNewCheckMetrics(newMetrics map[string]*api.CheckBundleMetric) {
	// only manage metrics if checkBundle has been populated
	if cm.checkBundle == nil {
		return
	}

	newCheckBundle := cm.checkBundle
	numCurrMetrics := len(newCheckBundle.Metrics)
	numNewMetrics := len(newMetrics)

	if numCurrMetrics+numNewMetrics >= cap(newCheckBundle.Metrics) {
		nm := make([]api.CheckBundleMetric, numCurrMetrics+numNewMetrics)
		copy(nm, newCheckBundle.Metrics)
		newCheckBundle.Metrics = nm
	}

	newCheckBundle.Metrics = newCheckBundle.Metrics[0 : numCurrMetrics+numNewMetrics]

	i := 0
	for _, metric := range newMetrics {
		newCheckBundle.Metrics[numCurrMetrics+i] = *metric
		i++
	}

	checkBundle, err := cm.apih.UpdateCheckBundle(newCheckBundle)
	if err != nil {
		cm.Log.Printf("[ERROR] updating check bundle with new metrics %v", err)
		return
	}

	for _, metric := range newMetrics {
		cm.activeMetrics[metric.Name] = true
	}

	cm.checkBundle = checkBundle
}
