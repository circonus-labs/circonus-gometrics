package checkmgr

import (
	"io/ioutil"
	"log"
	"os"
	"sync"
	"time"

	"github.com/circonus-labs/circonus-gometrics/api"
)

const (
	defaultCheckType = "httptrap"
)

type Check struct {
	SubmissionUrl string
	CheckId       int
	InstanceId    string
	SearchTag     string
	Secret        string
	Tags          []string
}

type Broker struct {
	Id int
	// a tag that can be used to select 1-n brokers from which to select
	// when creating a new check (e.g. datacenter:abc)
	SelectTag string
	// for a broker to be considered viable it must respond to a
	// connection attempt within this amount of time
	MaxResponseTime time.Duration
}

type Config struct {
	Api *api.Config
	Check
	Broker

	Log   *log.Logger
	Debug bool
}

type CheckManager struct {
	Log   *log.Logger
	Debug bool

	apih           *api.Api
	checkBundle    *api.CheckBundle
	activeMetrics  map[string]bool
	checkType      string
	ready          bool
	trapUrl        string
	trapCN         string
	trapSSL        bool
	trapLastUpdate time.Time
	trapmu         sync.Mutex

	submissionUrl string
	checkId       int
	instanceId    string
	searchTag     string
	checkSecret   string
	tags          []string

	brokerGroupId         int
	brokerSelectTag       string
	brokerMaxResponseTime time.Duration
}

func NewCheckManager(cmc *Config) (*CheckManager, error) {

	cm := &CheckManager{}

	cm.Debug = cmc.Debug
	cm.Log = cmc.Log
	if cm.Log == nil {
		if cm.Debug {
			cm.Log = log.New(os.Stderr, "", log.LstdFlags)
		} else {
			cm.Log = log.New(ioutil.Discard, "", log.LstdFlags)
		}
	}

	cmc.Api.Debug = cm.Debug
	cmc.Api.Log = cm.Log

	apih, err := api.NewApi(cmc.Api)
	if err != nil {
		return nil, err
	}
	cm.apih = apih

	cm.checkType = defaultCheckType

	return cm, nil
}

func (cm *CheckManager) GetTrapUrl() (string, error) {
	if cm.ready {
		return cm.trapUrl, nil
	}

	if cm.manage {
		if err := cm.initializeTrap(); err != nil {
			return "", err
		}
	} else if cm.submissionUrl != "" {
		cm.trapUrl = cm.submissionUrl
		cm.ready = true
		return cm.trapUrl, nil
	}

	return "", errors.New("Unable to initialze Circonus metrics trap.")
}
