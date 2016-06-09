package circonusgometrics

// abstracted in preparation of separate circonus-api-go package

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

// this is an arbitrary json structure
// we only would care about submission_url
type CheckDetails struct {
	SubmissionUrl string `json:"submission_url"`
}

type Check struct {
	Cid            string       `json:"_cid"`
	Active         bool         `json:"_active"`
	BrokerCid      string       `json:"_broker"`
	CheckBundleCid string       `json:"_check_bundle"`
	CheckUUID      string       `json:"_check_uuid"`
	Details        CheckDetails `json:"_details"`
}

// Use Circonus API to retrieve a check by ID
func (m *CirconusMetrics) fetchCheckById(id int) (*Check, error) {
	cid := fmt.Sprintf("/check/%d", id)
	return m.fetchCheckByCid(cid)
}

// Use Circonus API to retrieve a check by CID
func (m *CirconusMetrics) fetchCheckByCid(cid string) (*Check, error) {
	result, err := m.apiCall("GET", cid, nil)
	if err != nil {
		return nil, err
	}

	check := new(Check)
	json.Unmarshal(result, check)

	return check, nil
}

// Use Circonus API to retrieve a check by submission url
func (m *CirconusMetrics) fetchCheckBySubmissionUrl(submissionUrl string) (*Check, error) {

	u, err := url.Parse(submissionUrl)
	if err != nil {
		return nil, err
	}

	// valid trap url: scheme://host[:port]/module/httptrap/UUID/secret

	// does it smell like a valid trap url path
	if u.Path[0:17] != "/module/httptrap/" {
		return nil, fmt.Errorf("Invalid submission URL '%s', unrecognized path.", submissionUrl)
	}

	// extract uuid/secret
	pathParts := strings.Split(u.Path[17:], "/")
	if len(pathParts) != 2 {
		return nil, fmt.Errorf("Invalid submission URL '%s', UUID not where expected.", submissionUrl)
	}

	uuid := pathParts[0]

	query := fmt.Sprintf("/v2/check?f__check_uuid=%s", uuid)

	result, err := m.apiCall("GET", query, nil)
	if err != nil {
		return nil, err
	}

	var checks []Check
	json.Unmarshal(result, &checks)

	if len(checks) == 0 {
		return nil, fmt.Errorf("No checks found with UUID %s", uuid)
	}

	numActive := 0
	checkId := -1

	for idx, check := range checks {
		if check.Active {
			numActive++
			checkId = idx
		}
	}

	if numActive > 1 {
		return nil, fmt.Errorf("Multiple checks with same UUID %s", uuid)
	}

	return &checks[checkId], nil

}
