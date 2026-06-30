// Copyright (c) 2026, the go-ruby-faker/faker authors. All rights reserved.
// Use of this source code is governed by a BSD-3-Clause license that can be
// found in the LICENSE file.

package faker

import (
	_ "embed"
	"encoding/json"
	"strconv"
	"strings"
	"testing"
	"time"
)

// Oracle vectors captured from the Ruby `faker` gem (3.8.0) with
// Faker::Config.random = Random.new(seed). These pin the deterministic-seed
// contract: every key below is in the EXACT-SEQUENCE parity scope. The gem is
// not available in CI, so this test is data-driven from the committed JSON and
// holds 100% on its own.
//
//go:embed testdata/oracle.json
var oracleJSON []byte

// gen returns the n-th (1-indexed in the oracle, 0-indexed here) value for a key
// from a fresh Faker; the vectors call each generator n times in sequence, so we
// drive a single Faker and collect n outputs.
func collect(f *Faker, n int, fn func(*Faker) string) []string {
	out := make([]string, n)
	for i := range out {
		out[i] = fn(f)
	}
	return out
}

func generatorFor(key string) func(*Faker) string {
	switch key {
	case "name.name":
		return func(f *Faker) string { return f.Name().Name() }
	case "name.first_name":
		return func(f *Faker) string { return f.Name().FirstName() }
	case "name.last_name":
		return func(f *Faker) string { return f.Name().LastName() }
	case "name.name_with_middle":
		return func(f *Faker) string { return f.Name().NameWithMiddle() }
	case "name.prefix":
		return func(f *Faker) string { return f.Name().Prefix() }
	case "name.suffix":
		return func(f *Faker) string { return f.Name().Suffix() }
	case "name.initials":
		return func(f *Faker) string { return f.Name().Initials(3) }
	case "address.city":
		return func(f *Faker) string { return f.Address().City() }
	case "address.street_name":
		return func(f *Faker) string { return f.Address().StreetName() }
	case "address.street_address":
		return func(f *Faker) string { return f.Address().StreetAddress() }
	case "address.zip_code":
		return func(f *Faker) string { return f.Address().ZipCode() }
	case "address.state":
		return func(f *Faker) string { return f.Address().State() }
	case "address.state_abbr":
		return func(f *Faker) string { return f.Address().StateAbbr() }
	case "address.country":
		return func(f *Faker) string { return f.Address().Country() }
	case "address.country_code":
		return func(f *Faker) string { return f.Address().CountryCode() }
	case "address.full_address":
		return func(f *Faker) string { return f.Address().FullAddress() }
	case "address.building_number":
		return func(f *Faker) string { return f.Address().BuildingNumber() }
	case "address.secondary_address":
		return func(f *Faker) string { return f.Address().SecondaryAddress() }
	case "address.community":
		return func(f *Faker) string { return f.Address().Community() }
	case "phone.phone_number":
		return func(f *Faker) string { return f.PhoneNumber().PhoneNumber() }
	case "phone.cell_phone":
		return func(f *Faker) string { return f.PhoneNumber().CellPhone() }
	case "phone.country_code":
		return func(f *Faker) string { return f.PhoneNumber().CountryCode() }
	case "phone.area_code":
		return func(f *Faker) string { return f.PhoneNumber().AreaCode() }
	case "phone.exchange_code":
		return func(f *Faker) string { return f.PhoneNumber().ExchangeCode() }
	case "phone.subscriber_number":
		return func(f *Faker) string { return f.PhoneNumber().SubscriberNumber(4) }
	case "lorem.word":
		return func(f *Faker) string { return f.Lorem().Word() }
	case "lorem.words3":
		return func(f *Faker) string { return strings.Join(f.Lorem().Words(3, false), ",") }
	case "lorem.words5":
		return func(f *Faker) string { return strings.Join(f.Lorem().Words(5, false), ",") }
	case "lorem.sentence":
		return func(f *Faker) string { return f.Lorem().Sentence(4, false) }
	case "lorem.sentences":
		return func(f *Faker) string { return strings.Join(f.Lorem().Sentences(2, false), "|") }
	case "lorem.paragraph":
		return func(f *Faker) string { return f.Lorem().Paragraph(3, false) }
	case "lorem.characters":
		return func(f *Faker) string { return f.Lorem().Characters(12) }
	case "number.number":
		return func(f *Faker) string { return strconv.FormatInt(f.Number().Number(10), 10) }
	case "number.digit":
		return func(f *Faker) string { return strconv.Itoa(f.Number().Digit()) }
	case "number.hexadecimal":
		return func(f *Faker) string { return f.Number().Hexadecimal(6) }
	case "number.binary":
		return func(f *Faker) string { return f.Number().Binary(8) }
	case "number.decimal":
		return func(f *Faker) string { return rubyFloat(f.Number().Decimal(4, 2)) }
	case "number.between_int":
		return func(f *Faker) string { return strconv.Itoa(f.Number().BetweenInt(1, 100)) }
	case "company.name":
		return func(f *Faker) string { return f.Company().Name() }
	case "company.suffix":
		return func(f *Faker) string { return f.Company().Suffix() }
	case "company.industry":
		return func(f *Faker) string { return f.Company().Industry() }
	case "company.catch_phrase":
		return func(f *Faker) string { return f.Company().CatchPhrase() }
	case "company.bs":
		return func(f *Faker) string { return f.Company().BS() }
	case "company.buzzword":
		return func(f *Faker) string { return f.Company().Buzzword() }
	case "commerce.product_name":
		return func(f *Faker) string { return f.Commerce().ProductName() }
	case "commerce.department":
		return func(f *Faker) string { return f.Commerce().Department(3, false) }
	case "color.name":
		return func(f *Faker) string { return f.Color().ColorName() }
	case "internet.username":
		return func(f *Faker) string { return f.Internet().Username() }
	case "internet.domain_word":
		return func(f *Faker) string { return f.Internet().DomainWord() }
	case "internet.domain_name":
		return func(f *Faker) string { return f.Internet().DomainName() }
	case "internet.email":
		return func(f *Faker) string { return f.Internet().Email() }
	case "internet.slug":
		return func(f *Faker) string { return f.Internet().Slug() }
	case "internet.url":
		return func(f *Faker) string { return f.Internet().URL() }
	case "internet.ip_v4":
		return func(f *Faker) string { return f.Internet().IPv4Address() }
	case "internet.ip_v6":
		return func(f *Faker) string { return f.Internet().IPv6Address() }
	case "internet.mac":
		return func(f *Faker) string { return f.Internet().MacAddress("") }
	case "date.between":
		from := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
		to := time.Date(2025, 12, 31, 0, 0, 0, 0, time.UTC)
		return func(f *Faker) string { return f.Date().Between(from, to).Format("2006-01-02") }
	default:
		return nil
	}
}

// rubyFloat formats a float the way Ruby's Float#to_s does for the small decimals
// the oracle uses (trailing-zero-trimmed, always a fractional part).
func rubyFloat(x float64) string {
	s := strconv.FormatFloat(x, 'f', -1, 64)
	if !strings.Contains(s, ".") {
		s += ".0"
	}
	return s
}

func TestOracleParity(t *testing.T) {
	var vectors map[string]map[string][]string
	if err := json.Unmarshal(oracleJSON, &vectors); err != nil {
		t.Fatalf("parse oracle: %v", err)
	}
	if len(vectors) == 0 {
		t.Fatal("no oracle vectors")
	}
	checked := 0
	for seedStr, byKey := range vectors {
		seed, err := strconv.ParseInt(seedStr, 10, 64)
		if err != nil {
			t.Fatalf("bad seed %q: %v", seedStr, err)
		}
		for key, want := range byKey {
			fn := generatorFor(key)
			if fn == nil {
				t.Fatalf("no generator mapped for oracle key %q", key)
			}
			f := NewSeeded(seed)
			got := collect(f, len(want), fn)
			for i := range want {
				if got[i] != want[i] {
					t.Errorf("seed=%d key=%s [%d]: got %q want %q", seed, key, i, got[i], want[i])
				}
				checked++
			}
		}
	}
	t.Logf("verified %d oracle values across %d seeds", checked, len(vectors))
}
