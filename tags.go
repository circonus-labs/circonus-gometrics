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
// Note: does not work with checks using metric_filters (the default) use metric
// `*WithTags` helper methods or manual manage stream tags in metric names.
func (m *CirconusMetrics) SetMetricTags(name string, tags []string) bool {
	return m.check.AddMetricTags(name, tags, false)
}

// AddMetricTags appends tags to any existing tags for the named metric and flags a check update is needed
// Note: does not work with checks using metric_filters (the default) use metric
// `*WithTags` helper methods or manual manage stream tags in metric names.
func (m *CirconusMetrics) AddMetricTags(name string, tags []string) bool {
	return m.check.AddMetricTags(name, tags, true)
}

// MetricNameWithStreamTags will encode tags as stream tags into supplied metric name.
// Note: if metric name already has stream tags it is assumed the metric name and
// embedded stream tags are being managed manually and calling this method will nave no effect.
func (m *CirconusMetrics) MetricNameWithStreamTags(metric string, tags Tags) string {
	if len(tags) == 0 {
		return metric
	}

	if strings.Contains(metric, "|ST[") {
		return metric
	}

	taglist := m.EncodeMetricStreamTags(metric, tags)
	if taglist != "" {
		return metric + "|ST[" + taglist + "]"
	}

	return metric
}

// EncodeMetricStreamTags encodes Tags into a string suitable for use with
// stream tags. Tags directly embedded into metric names using the
// `metric_name|ST[<tags>]` syntax.
func (m *CirconusMetrics) EncodeMetricStreamTags(metricName string, tags Tags) string {
	if len(tags) == 0 {
		return ""
	}

	tmpTags := m.EncodeMetricTags(metricName, tags)
	if len(tmpTags) == 0 {
		return ""
	}

	tagList := make([]string, len(tmpTags))
	for i, tag := range tmpTags {
		tagParts := strings.SplitN(tag, ":", 2)
		if len(tagParts) != 2 {
			m.Log.Printf("%s has invalid tag (%s)", metricName, tag)
			continue // invalid tag, skip it
		}
		encodeFmt := `b"%s"`
		encodedSig := `b"` // has cat or val been previously (or manually) base64 encoded and formatted
		tc := tagParts[0]
		tv := tagParts[1]
		if !strings.HasPrefix(tc, encodedSig) {
			tc = fmt.Sprintf(encodeFmt, base64.StdEncoding.EncodeToString([]byte(tc)))
		}
		if !strings.HasPrefix(tv, encodedSig) && tv != "" {
			tv = fmt.Sprintf(encodeFmt, base64.StdEncoding.EncodeToString([]byte(tv)))
		}
		tagList[i] = tc + ":" + tv
	}

	return strings.Join(tagList, ",")
}

// EncodeMetricTags encodes Tags into an array of strings. The format
// check_bundle.metircs.metric.tags needs. This helper is intended to work
// with legacy check bundle metrics. Tags directly on named metrics are being
// deprecated in favor of stream tags.
func (m *CirconusMetrics) EncodeMetricTags(metricName string, tags Tags) []string {
	if len(tags) == 0 {
		return []string{}
	}

	uniqueTags := make(map[string]bool)
	for _, t := range tags {
		tc := strings.Map(removeSpaces, strings.ToLower(t.Category))
		tv := strings.TrimSpace(t.Value)
		if tc == "" {
			m.Log.Printf("%s has invalid tag (%#v)", metricName, t)
			continue
		}
		tag := tc + ":"
		if tv != "" {
			tag += tv
		}
		uniqueTags[tag] = true
	}
	tagList := make([]string, len(uniqueTags))
	idx := 0
	for t := range uniqueTags {
		tagList[idx] = t
		idx++
	}
	sort.Strings(tagList)
	return tagList
}

func removeSpaces(r rune) rune {
	if unicode.IsSpace(r) {
		return -1
	}
	return r
}
