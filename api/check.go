package api

// abstracted in preparation of separate circonus-api-go package

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

// CheckDetails is an arbitrary json structure, we would only care about submission_url
type CheckDetails struct {
	SubmissionURL string `json:"submission_url"`
}

// Check definition
type Check struct {
	Cid            string       `json:"_cid"`
	Active         bool         `json:"_active"`
	BrokerCid      string       `json:"_broker"`
	CheckBundleCid string       `json:"_check_bundle"`
	CheckUUID      string       `json:"_check_uuid"`
	Details        CheckDetails `json:"_details"`
}

// FetchCheckByID Use Circonus API to retrieve a check by ID
func (a *API) FetchCheckByID(id int) (*Check, error) {
	cid := fmt.Sprintf("/check/%d", id)
	return a.FetchCheckByCID(cid)
}

// FetchCheckByCID Use Circonus API to retrieve a check by CID
func (a *API) FetchCheckByCID(cid string) (*Check, error) {
	result, err := a.Get(cid)
	if err != nil {
		return nil, err
	}

	check := new(Check)
	json.Unmarshal(result, check)

	return check, nil
}

// FetchCheckBySubmissionURL Use Circonus API to retrieve a check by submission url
func (a *API) FetchCheckBySubmissionURL(submissionURL string) (*Check, error) {

	u, err := url.Parse(submissionURL)
	if err != nil {
		return nil, err
	}

	// valid trap url: scheme://host[:port]/module/httptrap/UUID/secret

	// does it smell like a valid trap url path
	if u.Path[:17] != "/module/httptrap/" {
		return nil, fmt.Errorf("[ERROR] Invalid submission URL '%s', unrecognized path", submissionURL)
	}

	// extract uuid/secret
	pathParts := strings.Split(u.Path[17:len(u.Path)], "/")
	if len(pathParts) != 2 {
		return nil, fmt.Errorf("[ERROR] Invalid submission URL '%s', UUID not where expected", submissionURL)
	}

	uuid := pathParts[0]

	query := fmt.Sprintf("/check?f__check_uuid=%s", uuid)

	result, err := a.Get(query)
	if err != nil {
		return nil, err
	}

	var checks []Check
	json.Unmarshal(result, &checks)

	if len(checks) == 0 {
		return nil, fmt.Errorf("[ERROR] No checks found with UUID %s", uuid)
	}

	numActive := 0
	checkID := -1

	for idx, check := range checks {
		if check.Active {
			numActive++
			checkID = idx
		}
	}

	if numActive > 1 {
		return nil, fmt.Errorf("[ERROR] Multiple checks with same UUID %s", uuid)
	}

	return &checks[checkID], nil

}
