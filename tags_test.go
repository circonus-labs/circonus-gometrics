// Copyright 2016 Circonus, Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Tags helper functions

package circonusgometrics

import (
	"encoding/base64"
	"fmt"
	"regexp"
	"strings"
	"testing"
)

func TestEncodeMetricTags(t *testing.T) {
	inputTags := Tags{
		{"cat1", "val1"},
		{"cat2", "val2"},
		{"cat2", "val1"},
		{"cat3", "compound:val"},
	}
	expectTags := Tags{
		{"cat1", "val1"},
		{"cat2", "val1"},
		{"cat2", "val2"},
		{"cat3", "compound:val"},
	}

	tl := EncodeMetricTags(inputTags)
	if len(tl) != len(expectTags) {
		t.Fatalf("expected %d tags, got %d", len(expectTags), len(tl))
	}
	for idx, tag := range tl {
		otag := fmt.Sprintf("%s:%s", expectTags[idx].Category, expectTags[idx].Value)
		if tag != otag {
			t.Fatalf("expected '%s' got '%s'", otag, tag)
		}
	}
}

func TestEncodeMetricStreamTags(t *testing.T) {
	inputTags := Tags{
		{"cat1", "val1"},
		{"cat1", "val2"},
		{"cat 1", "val2"}, // should have space removed and then be deduplicated
		{"cat2", "val2"},
		{"cat2", "val1"}, // should be sorted above previous one
		{"cat2", "val1"}, // duplicate should be omitted
		{"cat3", "compound:val"},
		{"cat3", fmt.Sprintf(`b"%s"`, base64.StdEncoding.EncodeToString([]byte("bar")))}, // manually base64 encoded and formatted (e.g. `b"base64encodedstr"`), do not double encode
	}
	expectTags := Tags{
		{"cat1", "val1"},
		{"cat1", "val2"},
		{"cat2", "val1"},
		{"cat2", "val2"},
		{"cat3", "compound:val"},
		{"cat3", "foo"},
	}

	t.Logf("tags: %v\n", inputTags)
	// expect ts to be in format b"b64cat":b"b64val",...
	ts := EncodeMetricStreamTags(inputTags)
	tl := strings.Split(ts, ",")
	if len(tl) != len(expectTags) {
		t.Fatalf("expected %d tags, got %d", len(expectTags), len(tl))
	}
	rx := regexp.MustCompile(`^b"(?P<cat>[^"]+)":b"(?P<val>[^"]+)"$`)
	for id, tag := range tl {
		matches := rx.FindStringSubmatch(string(tag))
		if len(matches) < 2 {
			t.Fatalf("tag did not match (%s)", tag)
		}
		result := make(map[string]string)
		for i, name := range rx.SubexpNames() {
			if i != 0 && name != "" {
				result[name] = matches[i]
			}
		}
		if cat, found := result["cat"]; !found {
			t.Fatalf("category: named match not found '%s'", cat)
		} else if cat == "" {
			t.Fatalf("category: invalid (empty) '%s'", cat)
			if dcat, err := base64.StdEncoding.DecodeString(cat); err != nil {
				t.Fatalf("category: error decoding base64 '%s' (%s)", cat, err)
			} else if string(dcat) != expectTags[id].Category {
				t.Fatalf("category: expected '%s' got '%s'->'%s'", expectTags[id].Category, cat, string(dcat))
			}
		}

		if val, found := result["val"]; !found {
			t.Fatalf("value: named match not found '%s'", val)
		} else if val == "" {
			t.Fatalf("value: invalid (empty) '%s'", val)
			if dval, err := base64.StdEncoding.DecodeString(val); err != nil {
				t.Fatalf("value: error decoding base64 '%s' (%s)", val, err)
			} else if string(dval) != expectTags[id].Value {
				t.Fatalf("value: expected '%s' got '%s'->'%s'", expectTags[id].Value, val, string(dval))
			}
		}
	}
}
