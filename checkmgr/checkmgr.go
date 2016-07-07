package checkmgr

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"sync"
	"time"

	"github.com/circonus-labs/circonus-gometrics/api"
)

// Check management offers:
//
// Create a check if one cannot be found matching specific criteria
// Manage metrics in the supplied check (enabling new metrics as they are submitted)
//
// To disable check management, leave Config.Api.Token.Key blank
//
// use cases:
// configure without api token - check management disabled
//  - configuration parameters other than Check.SubmissionUrl, Debug and Log are ignored
//  - note: SubmissionUrl is **required** in this case as there is no way to derive w/o api
// configure with api token - check management enabled
//  - all otehr configuration parameters affect how the trap url is obtained
//    1. provided (Check.SubmissionUrl)
//    2. via check lookup (CheckConfig.Id)
//    3. via a search using CheckConfig.InstanceId + CheckConfig.SearchTag
//    4. a new check is created

const (
	defaultCheckType = "httptrap"
)

type CheckConfig struct {
	// a specific submission url
	SubmissionUrl string
	// a specific check id (not check bundle id)
	Id int
	// unique instance id string
	// used to search for a check to use
	// used as check.target when creating a check
	InstanceId string
	// unique check searching tag
	// used to search for a check to use (combined with instanceid)
	// used as a regular tag when creating a check
	SearchTag string
	// httptrap check secret (for creating a check)
	Secret string
	// additional tags to add to a check (when creating a check)
	// these tags will not be added to an existing check
	Tags []string
}

type BrokerConfig struct {
	// a specific broker id (numeric portion of cid)
	Id int
	// a tag that can be used to select 1-n brokers from which to select
	// when creating a new check (e.g. datacenter:abc)
	SelectTag string
	// for a broker to be considered viable it must respond to a
	// connection attempt within this amount of time
	MaxResponseTime time.Duration
}

type Config struct {
	Api    api.Config
	Check  CheckConfig
	Broker BrokerConfig

	Log   *log.Logger
	Debug bool
}

type CheckManager struct {
	enabled               bool
	Log                   *log.Logger
	Debug                 bool
	apih                  *api.Api
	checkBundle           *api.CheckBundle
	activeMetrics         map[string]bool
	checkType             string
	ready                 bool
	trapUrl               string
	trapCN                string
	trapSSL               bool
	trapLastUpdate        time.Time
	trapmu                sync.Mutex
	submissionUrl         string
	checkId               int
	instanceId            string
	searchTag             string
	checkSecret           string
	tags                  []string
	brokerGroupId         int
	brokerSelectTag       string
	brokerMaxResponseTime time.Duration
}

func NewCheckManager(cmc *Config) (*CheckManager, error) {

	if cmc == nil {
		return nil, errors.New("Invalid Check Manager configuration (nil).")
	}

	cm := &CheckManager{
		enabled: false,
		ready:   false,
	}

	cm.Debug = cmc.Debug
	cm.Log = cmc.Log
	if cm.Log == nil {
		if cm.Debug {
			cm.Log = log.New(os.Stderr, "", log.LstdFlags)
		} else {
			cm.Log = log.New(ioutil.Discard, "", log.LstdFlags)
		}
	}

	if cmc.Check.SubmissionUrl != "" {
		cm.trapUrl = cmc.Check.SubmissionUrl
		cm.ready = true
	}

	if cmc.Api.Token.Key == "" {
		if cm.trapUrl != "" && cm.ready {
			return cm, nil
		} else {
			return nil, errors.New("Invalid check manager configuration (no API token AND no submission url).")
		}
	}

	// enable check manager (a blank api.Token.Key *disables* check management)

	cm.enabled = true
	cm.checkType = defaultCheckType

	cmc.Api.Debug = cm.Debug
	cmc.Api.Log = cm.Log

	apih, err := api.NewApi(&cmc.Api)
	if err != nil {
		return nil, err
	}
	cm.apih = apih

	return cm, nil
}

func (cm *CheckManager) GetTrapUrl() (string, error) {
	if cm.trapUrl != "" {
		return cm.trapUrl, nil
	}

	if !cm.enabled {
		return "", errors.New("No submission URL supplied and check manager disabled.")
	}

	if err := cm.initializeTrap(); err != nil {
		return "", err
	}

	return "", errors.New("Unable to initialze Circonus metrics trap.")
}
