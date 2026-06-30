// Copyright (c) 2026, the go-ruby-faker/faker authors. All rights reserved.
// Use of this source code is governed by a BSD-3-Clause license that can be
// found in the LICENSE file.

package faker

import (
	_ "embed"
	"encoding/json"
)

// en.json is the committed en-locale dataset extracted verbatim from the faker
// gem's i18n tables, so the library is complete offline (no network, no gem).
//
//go:embed data/en.json
var enJSON []byte

// dataset is a parsed locale table. Most keys hold string arrays; a few hold
// nested structures (company.buzzwords / company.bs are arrays-of-arrays,
// lorem.punctuation is an object). Those are kept as json.RawMessage and decoded
// on demand by the generators that need them.
type dataset struct {
	raw map[string]json.RawMessage
	// cached decoded forms
	strCache  map[string][]string
	listCache map[string][][]string
}

var enData = mustLoad(enJSON)

func mustLoad(b []byte) *dataset {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(b, &raw); err != nil {
		panic("faker: cannot parse embedded locale data: " + err.Error())
	}
	return &dataset{
		raw:       raw,
		strCache:  make(map[string][]string),
		listCache: make(map[string][][]string),
	}
}

// strings returns the locale array for key as []string. A missing key yields an
// empty slice (the generators that call it always have their keys present in
// the embedded en data). A scalar-string value is wrapped in a single-element
// slice (matching the gem treating e.g. "separator" uniformly).
func (d *dataset) strings(key string) []string {
	if v, ok := d.strCache[key]; ok {
		return v
	}
	raw, ok := d.raw[key]
	if !ok {
		d.strCache[key] = nil
		return nil
	}
	var arr []string
	if err := json.Unmarshal(raw, &arr); err == nil {
		d.strCache[key] = arr
		return arr
	}
	var s string
	if err := json.Unmarshal(raw, &s); err == nil {
		arr = []string{s}
		d.strCache[key] = arr
		return arr
	}
	d.strCache[key] = nil
	return nil
}

// scalar returns a scalar string value (e.g. "separator"); for an array value
// it returns the first element.
func (d *dataset) scalar(key string) string {
	raw, ok := d.raw[key]
	if !ok {
		return ""
	}
	var s string
	if err := json.Unmarshal(raw, &s); err == nil {
		return s
	}
	if arr := d.strings(key); len(arr) > 0 {
		return arr[0]
	}
	return ""
}

// lists returns an array-of-arrays value (company.buzzwords / company.bs).
func (d *dataset) lists(key string) [][]string {
	if v, ok := d.listCache[key]; ok {
		return v
	}
	raw, ok := d.raw[key]
	if !ok {
		d.listCache[key] = nil
		return nil
	}
	var ll [][]string
	if err := json.Unmarshal(raw, &ll); err != nil {
		d.listCache[key] = nil
		return nil
	}
	d.listCache[key] = ll
	return ll
}

// object returns an object value as a map (lorem.punctuation).
func (d *dataset) object(key string) map[string]string {
	raw, ok := d.raw[key]
	if !ok {
		return nil
	}
	var m map[string]string
	if err := json.Unmarshal(raw, &m); err != nil {
		return nil
	}
	return m
}

// flatten returns all elements of an array-of-arrays as a single slice
// (Company.buzzword flattens buzzwords).
func (d *dataset) flatten(key string) []string {
	var out []string
	for _, l := range d.lists(key) {
		out = append(out, l...)
	}
	return out
}
