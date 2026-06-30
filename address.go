// Copyright (c) 2026, the go-ruby-faker/faker authors. All rights reserved.
// Use of this source code is governed by a BSD-3-Clause license that can be
// found in the LICENSE file.

package faker

// AddressGen is Faker::Address.
type AddressGen struct{ f *Faker }

// Address returns the Faker::Address generator.
func (f *Faker) Address() *AddressGen { return &AddressGen{f} }

// City produces a city name (Faker::Address.city).
func (a *AddressGen) City() string { return a.f.Parse("address.city") }

// CityPrefix produces a city prefix (Faker::Address.city_prefix).
func (a *AddressGen) CityPrefix() string { return a.f.Fetch("address.city_prefix") }

// CitySuffix produces a city suffix (Faker::Address.city_suffix).
func (a *AddressGen) CitySuffix() string { return a.f.Fetch("address.city_suffix") }

// StreetName produces a street name (Faker::Address.street_name).
func (a *AddressGen) StreetName() string { return a.f.Parse("address.street_name") }

// StreetSuffix produces a street suffix (Faker::Address.street_suffix).
func (a *AddressGen) StreetSuffix() string { return a.f.Fetch("address.street_suffix") }

// StreetAddress produces a street address (Faker::Address.street_address).
func (a *AddressGen) StreetAddress() string {
	return a.f.Numerify(a.f.Parse("address.street_address"), false)
}

// SecondaryAddress produces a secondary address (Faker::Address.secondary_address).
func (a *AddressGen) SecondaryAddress() string {
	return a.f.Bothify(a.f.Fetch("address.secondary_address"))
}

// BuildingNumber produces a building number (Faker::Address.building_number).
func (a *AddressGen) BuildingNumber() string {
	return a.f.Bothify(a.f.Fetch("address.building_number"))
}

// Community produces a community name (Faker::Address.community).
func (a *AddressGen) Community() string { return a.f.Parse("address.community") }

// ZipCode produces a zip/postal code (Faker::Address.zip_code). The fetched
// postcode pattern is letterified then numerified with leading zeros allowed.
func (a *AddressGen) ZipCode() string {
	lettered := a.f.Letterify(a.f.Fetch("address.postcode"))
	return a.f.Numerify(lettered, true)
}

// Zip is an alias for ZipCode.
func (a *AddressGen) Zip() string { return a.ZipCode() }

// Postcode is an alias for ZipCode.
func (a *AddressGen) Postcode() string { return a.ZipCode() }

// State produces a state name (Faker::Address.state).
func (a *AddressGen) State() string { return a.f.Fetch("address.state") }

// StateAbbr produces a state abbreviation (Faker::Address.state_abbr).
func (a *AddressGen) StateAbbr() string { return a.f.Fetch("address.state_abbr") }

// Country produces a country name (Faker::Address.country).
func (a *AddressGen) Country() string { return a.f.Fetch("address.country") }

// CountryCode produces an ISO 3166 alpha-2 country code (Faker::Address.country_code).
func (a *AddressGen) CountryCode() string { return a.f.Fetch("address.country_code") }

// CountryCodeLong produces an ISO 3166 alpha-3 country code.
func (a *AddressGen) CountryCodeLong() string { return a.f.Fetch("address.country_code_long") }

// TimeZone produces a time-zone name (Faker::Address.time_zone).
func (a *AddressGen) TimeZone() string { return a.f.Fetch("address.time_zone") }

// FullAddress produces a full address (Faker::Address.full_address).
func (a *AddressGen) FullAddress() string { return a.f.Parse("address.full_address") }

// Latitude produces a latitude in [-90, 90). Distribution parity only (float draw).
func (a *AddressGen) Latitude() float64 { return a.f.rng.Float()*180 - 90 }

// Longitude produces a longitude in [-180, 180). Distribution parity only.
func (a *AddressGen) Longitude() float64 { return a.f.rng.Float()*360 - 180 }
