package checkmgr

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path"
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
	enabled bool
	Log     *log.Logger
	Debug   bool
	apih    *api.Api

	// check
	checkType          string
	checkId            int
	checkInstanceId    string
	checkSearchTag     string
	checkSecret        string
	checkTags          []string
	checkSubmissionUrl string

	// broker
	brokerId              int
	brokerSelectTag       string
	brokerMaxResponseTime time.Duration

	// state
	checkBundle    *api.CheckBundle
	activeMetrics  map[string]bool
	trapUrl        string
	trapCN         string
	trapLastUpdate time.Time
	trapmu         sync.Mutex
	certPool       *x509.CertPool
}

type Trap struct {
	Url *url.URL
	Tls *tls.Config
}

/*
return &Trap{
    url.Parse(cm.trapUrl),
    &tls.Config{
       RootCAs:    cm.certPool,
       ServerName: cm.trapCN,
    },
}
*/

func NewCheckManager(cfg *Config) (*CheckManager, error) {

	if cfg == nil {
		return nil, errors.New("Invalid Check Manager configuration (nil).")
	}

	cm := &CheckManager{
		enabled: false,
	}

	cm.Debug = cfg.Debug
	cm.Log = cfg.Log
	if cm.Log == nil {
		if cm.Debug {
			cm.Log = log.New(os.Stderr, "", log.LstdFlags)
		} else {
			cm.Log = log.New(ioutil.Discard, "", log.LstdFlags)
		}
	}

	if cfg.Check.SubmissionUrl != "" {
		cm.checkSubmissionUrl = cfg.Check.SubmissionUrl
	}

	if cfg.Api.Token.Key == "" {
		if cm.checkSubmissionUrl == "" {
			return nil, errors.New("Invalid check manager configuration (no API token AND no submission url).")
		}
		cm.trapUrl = cm.checkSubmissionUrl
		return cm, nil
	}

	// enable check manager (a blank api.Token.Key *disables* check management)

	cm.enabled = true

	// initialize api handle

	cfg.Api.Debug = cm.Debug
	cfg.Api.Log = cm.Log

	apih, err := api.NewApi(&cfg.Api)
	if err != nil {
		return nil, err
	}
	cm.apih = apih

	// initialize check related data

	cm.checkType = defaultCheckType
	cm.checkId = cfg.Check.Id
	cm.checkInstanceId = cfg.Check.InstanceId
	cm.checkSearchTag = cfg.Check.SearchTag
	cm.checkSecret = cfg.Check.Secret
	cm.checkTags = cfg.Check.Tags

	_, an := path.Split(os.Args[0])

	if cm.checkInstanceId == "" {
		hn, err := os.Hostname()
		if err != nil {
			hn = "unknown"
		}
		cm.checkInstanceId = fmt.Sprintf("%s:%s", hn, an)
	}

	if cm.checkSearchTag == "" {
		cm.checkSearchTag = fmt.Sprintf("service:%s", an)
	}

	return cm, nil
}

func (cm *CheckManager) GetTrap() (*Trap, error) {
	if cm.trapUrl == "" {
		if err := cm.initializeTrapUrl(); err != nil {
			return nil, err
		}
	}

	trap := &Trap{}

	u, err := url.Parse(cm.trapUrl)
	if err != nil {
		return nil, err
	}

	trap.Url = u

	if u.Scheme == "https" {
		if cm.certPool == nil {
			cm.loadCACert()
		}
		t := &tls.Config{
			RootCAs: cm.certPool,
		}
		if cm.trapCN != "" {
			t.ServerName = cm.trapCN
		}
		// trap.Tls := &tls.Config{
		//     RootCAs:    cm.certPool,
		//     ServerName: cm.trapCN,
		// }
		trap.Tls = t
	}

	return trap, nil
}

// func (cm *CheckManager) GetTrapUrl() (string, error) {
// 	if cm.trapUrl != "" {
// 		return cm.trapUrl, nil
// 	}
//
// 	if err := cm.initializeTrapUrl(); err != nil {
// 		return "", err
// 	}
//
// 	return cm.trapUrl, nil
// }

func (cm *CheckManager) ResetTrap() error {
	if cm.trapUrl == "" {
		return nil
	}

	cm.trapUrl = ""
	err := cm.initializeTrapUrl()
	return err
}
