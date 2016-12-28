// Copyright 2016 Circonus, Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// User API support - Fetch and Update
// See: https://login.circonus.com/resources/api/calls/user
// Note: Create and Delete are not supported directly via the User API
// endpoint. See the Account endpoint for inviting and removing users
// from specific accounts.

package api

import (
	"encoding/json"
	"fmt"
	"net/url"
	"regexp"
)

// UserContactInfo defines known contact details
type UserContactInfo struct {
	SMS  string `json:"sms,omitempty"`
	XMPP string `json:"xmpp,omitempty"`
}

// User definition
type User struct {
	CID         string          `json:"_cid,omitempty"`
	ContactInfo UserContactInfo `json:"contact_info,omitempty"`
	Email       string          `json:"email"`
	Firstname   string          `json:"firstname"`
	Lastname    string          `json:"lastname"`
}

const baseUserPath = "/user"

// FetchUser retrieves a user definition
func (a *API) FetchUser(cid CIDType) (*User, error) {
	if cid == "" {
		cid = CIDType(baseUserPath + "/current")
	}

	if matched, err := regexp.MatchString("^"+baseUserPath+"/([0-9]+|current)$", string(cid)); err != nil {
		return nil, err
	} else if !matched {
		return nil, fmt.Errorf("Invalid user CID %v", cid)
	}

	result, err := a.Get(string(cid))
	if err != nil {
		return nil, err
	}

	user := new(User)
	if err := json.Unmarshal(result, user); err != nil {
		return nil, err
	}

	return user, nil
}

// FetchUsers retrieves users for current account
func (a *API) FetchUsers() ([]User, error) {
	result, err := a.Get(baseUserPath)
	if err != nil {
		return nil, err
	}

	var users []User
	if err := json.Unmarshal(result, &users); err != nil {
		return nil, err
	}

	return users, nil
}

// UpdateUser update user information
func (a *API) UpdateUser(config *User) (*User, error) {
	if matched, err := regexp.MatchString("^"+baseUserPath+"/[0-9]+$", string(config.CID)); err != nil {
		return nil, err
	} else if !matched {
		return nil, fmt.Errorf("Invalid user CID %v", config.CID)
	}

	reqURL := url.URL{
		Path: config.CID,
	}

	cfg, err := json.Marshal(config)
	if err != nil {
		return nil, err
	}

	resp, err := a.Put(reqURL.String(), cfg)
	if err != nil {
		return nil, err
	}

	user := &User{}
	if err := json.Unmarshal(resp, user); err != nil {
		return nil, err
	}

	return user, nil
}
