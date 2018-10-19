// Copyright 2016 Circonus, Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Tags helper functions

package circonusgometrics

import (
	"encoding/base64"
	"fmt"
	"sort"
	"strings"
	"unicode"
)

// Tag defines an individual tag
type Tag struct {
	Category string
	Value    string
}

// Tags defines a list of tags
type Tags []Tag

// SetMetricTags sets the tags for the named metric and flags a check update is needed
func (m *CirconusMetrics) SetMetricTags(name string, tags []string) bool {
	return m.check.AddMetricTags(name, tags, false)
}

// AddMetricTags appends tags to any existing tags for the named metric and flags a check update is needed
func (m *CirconusMetrics) AddMetricTags(name string, tags []string) bool {
	return m.check.AddMetricTags(name, tags, true)
}

// EncodeMetricStreamTags encodes Tags into a string suitable for use with
// stream tags. Tags directly embedded into metric names using the
// `metric_name|ST[<tags>]` syntax.
func EncodeMetricStreamTags(tags *Tags) string {
	if len(*tags) == 0 {
		return ""
	}

	tagEnc := make(map[string]bool)
	// additional deduplication step for stream tags since spaces are removed prior to encoding
	for _, tag := range prepTags(tags) {
		c := base64.StdEncoding.EncodeToString([]byte(strings.Map(removeSpaces, tag.Category)))
		v := base64.StdEncoding.EncodeToString([]byte(strings.Map(removeSpaces, tag.Value)))
		tagEnc[fmt.Sprintf(`b"%s":b"%s"`, c, v)] = true
	}
	i := 0
	tagList := make([]string, len(tagEnc))
	for t := range tagEnc {
		tagList[i] = t
		i++
	}
	return strings.Join(tagList, ",")
}

// EncodeMetricTags encodes Tags into an array of strings. The format check_bundle.metircs.metric.tags needs.
// This helper is intended to work with legacy check bundle metrics. Tags directly on named metrics are being
// deprecated in favor of stream tags.
func EncodeMetricTags(tags *Tags) []string {
	if len(*tags) == 0 {
		return []string{}
	}

	tagList := []string{}
	for _, tag := range prepTags(tags) {
		tagList = append(tagList, fmt.Sprintf(`%s:%s`, tag.Category, tag.Value))
	}
	return tagList
}

func removeSpaces(r rune) rune {
	if unicode.IsSpace(r) {
		return -1
	}
	return r
}

type compareFunc func(t1, t2 *Tag) int

type multiSorter struct {
	tags    Tags
	compare []compareFunc
}

func orderBy(compare ...compareFunc) *multiSorter {
	return &multiSorter{
		compare: compare,
	}
}

func (ms *multiSorter) Sort(tags Tags) {
	ms.tags = tags
	sort.Sort(ms)
}

func (ms *multiSorter) Len() int {
	return len(ms.tags)
}

func (ms *multiSorter) Swap(i, j int) {
	ms.tags[i], ms.tags[j] = ms.tags[j], ms.tags[i]
}

func (ms *multiSorter) Less(i, j int) bool {
	t1, t2 := &ms.tags[i], &ms.tags[j]
	var k int
	for k = 0; k < len(ms.compare)-1; k++ {
		if ms.compare[k](t1, t2) == -1 {
			return true
		}
	}
	return ms.compare[k](t1, t2) == -1
}

func prepTags(tags *Tags) Tags {
	if len(*tags) == 0 {
		return Tags{}
	}

	ltags := make(Tags, len(*tags))
	copy(ltags, *tags)

	// multisort
	category := func(t1, t2 *Tag) int {
		if t1.Category == t2.Category {
			return 0
		} else if t1.Category < t2.Category {
			return -1
		}
		return 1
	}
	value := func(t1, t2 *Tag) int {
		if t1.Value == t2.Value {
			return 0
		} else if t1.Value < t2.Value {
			return -1
		}
		return 1
	}

	orderBy(category, value).Sort(ltags)

	// deduplicate
	unique := make(map[string]map[string]bool)
	for _, tag := range ltags {
		if _, exists := unique[tag.Category]; !exists {
			unique[tag.Category] = make(map[string]bool)
		}
		unique[tag.Category][tag.Value] = true
	}

	if len(unique) == 0 {
		return Tags{}
	}

	tagList := Tags{}
	for tcat, tvals := range unique {
		for tval := range tvals {
			tagList = append(tagList, Tag{tcat, tval})
		}
	}

	return tagList
}
