// Copyright 2016 Circonus, Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Annotation API support - Fetch, Create, Delete, and Update
// See: https://login.circonus.com/resources/api/calls/annotation

package api

import (
	"encoding/json"
	"fmt"
	"net/url"
	"regexp"
)

// Annotation defines a annotation
type Annotation struct {
	CID            string   `json:"_cid,omitempty"`
	Created        int      `json:"_created,omitempty"`
	LastModified   int      `json:"_last_modified,omitempty"`
	LastModifiedBy string   `json:"_last_modified_by,omitempty"`
	Category       string   `json:"category"`
	Description    string   `json:"description"`
	RelatedMetrics []string `json:"rel_metrics"`
	Start          int      `json:"start"`
	Stop           int      `json:"stop"`
	Title          string   `json:"title"`
}

const (
	baseAnnotationPath = "/annotation"
	annotationCIDRegex = "^" + baseAnnotationPath + "/[0-9]+$"
)

// FetchAnnotation retrieves a annotation definition
func (a *API) FetchAnnotation(cid CIDType) (*Annotation, error) {
	if matched, err := regexp.MatchString(annotationCIDRegex, string(cid)); err != nil {
		return nil, err
	} else if !matched {
		return nil, fmt.Errorf("Invalid annotation CID %v", cid)
	}

	result, err := a.Get(string(cid))
	if err != nil {
		return nil, err
	}

	annotation := new(Annotation)
	if err := json.Unmarshal(result, annotation); err != nil {
		return nil, err
	}

	return annotation, nil
}

// FetchAnnotations retrieves all annotations
func (a *API) FetchAnnotations() ([]Annotation, error) {
	result, err := a.Get(baseAnnotationPath)
	if err != nil {
		return nil, err
	}

	var annotations []Annotation
	if err := json.Unmarshal(result, &annotations); err != nil {
		return nil, err
	}

	return annotations, nil
}

// UpdateAnnotation update annotation definition
func (a *API) UpdateAnnotation(config *Annotation) (*Annotation, error) {
	if matched, err := regexp.MatchString(annotationCIDRegex, string(config.CID)); err != nil {
		return nil, err
	} else if !matched {
		return nil, fmt.Errorf("Invalid annotation CID %v", config.CID)
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

	annotation := &Annotation{}
	if err := json.Unmarshal(resp, annotation); err != nil {
		return nil, err
	}

	return annotation, nil
}

// CreateAnnotation create a new annotation
func (a *API) CreateAnnotation(config *Annotation) (*Annotation, error) {
	reqURL := url.URL{
		Path: baseAnnotationPath,
	}

	cfg, err := json.Marshal(config)
	if err != nil {
		return nil, err
	}

	resp, err := a.Post(reqURL.String(), cfg)
	if err != nil {
		return nil, err
	}

	annotation := &Annotation{}
	if err := json.Unmarshal(resp, annotation); err != nil {
		return nil, err
	}

	return annotation, nil
}

// DeleteAnnotation delete a annotation
func (a *API) DeleteAnnotation(bundle *Annotation) (bool, error) {
	cid := CIDType(bundle.CID)
	return a.DeleteAnnotationByCID(cid)
}

// DeleteAnnotationByCID delete a annotation by cid
func (a *API) DeleteAnnotationByCID(cid CIDType) (bool, error) {
	if matched, err := regexp.MatchString(annotationCIDRegex, string(cid)); err != nil {
		return false, err
	} else if !matched {
		return false, fmt.Errorf("Invalid annotation CID %v", cid)
	}

	reqURL := url.URL{
		Path: string(cid),
	}

	_, err := a.Delete(reqURL.String())
	if err != nil {
		return false, err
	}

	return true, nil
}

// AnnotationSearch returns list of annotations matching a search query and/or filter
//    - a search query (see: https://login.circonus.com/resources/api#searching)
//    - a filter (see: https://login.circonus.com/resources/api#filtering)
func (a *API) AnnotationSearch(searchCriteria SearchQueryType, filterCriteria map[string]string) ([]Annotation, error) {

	if searchCriteria == "" && len(filterCriteria) == 0 {
		return a.FetchAnnotations()
	}

	reqURL := url.URL{
		Path: baseAnnotationPath,
	}

	q := url.Values{}

	if searchCriteria != "" {
		q.Set("search", string(searchCriteria))
	}

	if len(filterCriteria) > 0 {
		for filter, criteria := range filterCriteria {
			q.Set(filter, criteria)
		}
	}

	reqURL.RawQuery = q.Encode()

	resp, err := a.Get(reqURL.String())
	if err != nil {
		return nil, fmt.Errorf("[ERROR] API call error %+v", err)
	}

	var results []Annotation
	if err := json.Unmarshal(resp, &results); err != nil {
		return nil, err
	}

	return results, nil
}
