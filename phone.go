// Copyright (c) 2026, the go-ruby-faker/faker authors. All rights reserved.
// Use of this source code is governed by a BSD-3-Clause license that can be
// found in the LICENSE file.

package faker

import "strconv"

// PhoneNumberGen is Faker::PhoneNumber.
type PhoneNumberGen struct{ f *Faker }

// PhoneNumber returns the Faker::PhoneNumber generator.
func (f *Faker) PhoneNumber() *PhoneNumberGen { return &PhoneNumberGen{f} }

// PhoneNumber produces a phone number in a random format (Faker::PhoneNumber.phone_number).
func (p *PhoneNumberGen) PhoneNumber() string { return p.f.Parse("phone_number.formats") }

// CellPhone produces a cell phone number (Faker::PhoneNumber.cell_phone).
func (p *PhoneNumberGen) CellPhone() string { return p.f.Parse("cell_phone.formats") }

// CountryCode produces "+NN" (Faker::PhoneNumber.country_code).
func (p *PhoneNumberGen) CountryCode() string { return "+" + p.f.Fetch("phone_number.country_code") }

// PhoneNumberWithCountryCode prefixes a country code (Faker::PhoneNumber.phone_number_with_country_code).
func (p *PhoneNumberGen) PhoneNumberWithCountryCode() string {
	return p.CountryCode() + " " + p.PhoneNumber()
}

// CellPhoneWithCountryCode prefixes a country code to a cell number.
func (p *PhoneNumberGen) CellPhoneWithCountryCode() string {
	return p.CountryCode() + " " + p.CellPhone()
}

// AreaCode produces an area code (Faker::PhoneNumber.area_code).
func (p *PhoneNumberGen) AreaCode() string { return p.f.Fetch("phone_number.area_code") }

// ExchangeCode produces an exchange code (Faker::PhoneNumber.exchange_code).
func (p *PhoneNumberGen) ExchangeCode() string { return p.f.Fetch("phone_number.exchange_code") }

// SubscriberNumber produces a numeric subscriber number of the given length
// (Faker::PhoneNumber.subscriber_number / extension). It mirrors the gem's
// PositionalGenerator int(length:): a single range [10^(n-1), 10^n - 1] is
// sampled from a one-element array (consuming one draw) and then a value in
// that inclusive range is drawn.
func (p *PhoneNumberGen) SubscriberNumber(length int) string {
	if length < 1 {
		return ""
	}
	lower := 1
	for i := 0; i < length-1; i++ {
		lower *= 10
	}
	upper := lower*10 - 1
	// The single-element ranges array is sampled first (one draw), matching MRI.
	_ = p.f.rng.limitedRand(0)
	v := p.f.rng.IntInRange(lower, upper)
	return zeroPad(v, length)
}

// Extension is an alias for SubscriberNumber.
func (p *PhoneNumberGen) Extension(length int) string { return p.SubscriberNumber(length) }

func zeroPad(v, width int) string {
	s := strconv.Itoa(v)
	for len(s) < width {
		s = "0" + s
	}
	return s
}
