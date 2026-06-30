// Copyright (c) 2026, the go-ruby-faker/faker authors. All rights reserved.
// Use of this source code is governed by a BSD-3-Clause license that can be
// found in the LICENSE file.

package faker

import "strings"

// dispatch resolves a parse() token: if cls names another generator class with
// the requested method, call it; otherwise fall back to a locale lookup. This
// mirrors Faker::Base#parse's "cls.respond_to?(meth) ? cls.send(meth) : fetch".
//
// keyPath is the snake-cased name of the class the format string came from
// (e.g. "address"), used for the same-class fetch fallback.
func (f *Faker) dispatch(keyPath, cls, meth string) string {
	if cls != "" {
		if v, ok := f.callMethod(cls, meth); ok {
			return v
		}
		// Cross-class fetch fallback: faker.<snake(cls)>.<meth>
		return f.Fetch(snakeCase(cls) + "." + strings.ToLower(meth))
	}
	// No class prefix: try the current class's own methods, else same-class fetch.
	if v, ok := f.callMethod(classFor(keyPath), meth); ok {
		return v
	}
	return f.Fetch(keyPath + "." + strings.ToLower(meth))
}

// classFor maps a snake key path back to the generator class name used for
// self-method dispatch (e.g. "address" -> "Address").
func classFor(keyPath string) string {
	switch keyPath {
	case "name":
		return "Name"
	case "address":
		return "Address"
	case "company":
		return "Company"
	case "phone_number", "cell_phone":
		return "PhoneNumber"
	case "internet":
		return "Internet"
	case "commerce":
		return "Commerce"
	default:
		return ""
	}
}

// callMethod invokes the named generator method on the named class, returning
// ok=false when the class/method pair is not one parse() needs to dispatch.
func (f *Faker) callMethod(cls, meth string) (string, bool) {
	switch cls {
	case "Name":
		switch meth {
		case "name":
			return f.Name().Name(), true
		case "name_with_middle":
			return f.Name().NameWithMiddle(), true
		case "first_name":
			return f.Name().FirstName(), true
		case "male_first_name":
			return f.Name().MaleFirstName(), true
		case "female_first_name":
			return f.Name().FemaleFirstName(), true
		case "neutral_first_name":
			return f.Name().NeutralFirstName(), true
		case "last_name", "middle_name":
			return f.Name().LastName(), true
		case "prefix":
			return f.Name().Prefix(), true
		case "suffix":
			return f.Name().Suffix(), true
		}
	case "Address":
		switch meth {
		case "city":
			return f.Address().City(), true
		case "city_prefix":
			return f.Address().CityPrefix(), true
		case "city_suffix":
			return f.Address().CitySuffix(), true
		case "street_name":
			return f.Address().StreetName(), true
		case "street_address":
			return f.Address().StreetAddress(), true
		case "street_suffix":
			return f.Address().StreetSuffix(), true
		case "secondary_address":
			return f.Address().SecondaryAddress(), true
		case "building_number":
			return f.Address().BuildingNumber(), true
		case "community":
			return f.Address().Community(), true
		case "zip_code", "zip", "postcode":
			return f.Address().ZipCode(), true
		case "state":
			return f.Address().State(), true
		case "state_abbr":
			return f.Address().StateAbbr(), true
		case "country":
			return f.Address().Country(), true
		case "country_code":
			return f.Address().CountryCode(), true
		}
	case "Company":
		switch meth {
		case "name":
			return f.Company().Name(), true
		case "suffix":
			return f.Company().Suffix(), true
		}
	case "PhoneNumber":
		switch meth {
		case "area_code":
			return f.PhoneNumber().AreaCode(), true
		case "exchange_code":
			return f.PhoneNumber().ExchangeCode(), true
		}
	case "Commerce":
		switch meth {
		case "product_name":
			return f.Commerce().ProductName(), true
		}
	}
	return "", false
}
