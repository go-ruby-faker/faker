// Copyright (c) 2026, the go-ruby-faker/faker authors. All rights reserved.
// Use of this source code is governed by a BSD-3-Clause license that can be
// found in the LICENSE file.

package faker

import (
	"encoding/json"
	"strings"
	"testing"
	"time"
)

// newF returns a deterministically-seeded Faker for ruby-free coverage tests.
func newF() *Faker { return NewSeeded(42) }

func TestNameGenerators(t *testing.T) {
	f := newF()
	n := f.Name()
	checks := []struct {
		name string
		got  string
	}{
		{"name", n.Name()},
		{"name_with_middle", n.NameWithMiddle()},
		{"first_name", n.FirstName()},
		{"male_first_name", n.MaleFirstName()},
		{"female_first_name", n.FemaleFirstName()},
		{"neutral_first_name", n.NeutralFirstName()},
		{"last_name", n.LastName()},
		{"prefix", n.Prefix()},
		{"suffix", n.Suffix()},
	}
	for _, c := range checks {
		if c.got == "" {
			t.Errorf("Name.%s produced empty string", c.name)
		}
	}
	if got := n.Initials(3); len(got) != 3 {
		t.Errorf("Initials(3) = %q, want length 3", got)
	}
	for _, r := range n.Initials(2) {
		if r < 'A' || r > 'Z' {
			t.Errorf("Initials contains non A-Z rune %q", r)
		}
	}
}

func TestNameFirstNameEmptyFallback(t *testing.T) {
	// With an empty first_name list, Parse yields "" and the Fetch fallback
	// branch runs. Use a custom dataset to drive the empty path.
	f := newF()
	f.data = &dataset{
		raw:       map[string]json.RawMessage{"name.first_name": json.RawMessage(`[""]`)},
		strCache:  map[string][]string{},
		listCache: map[string][][]string{},
	}
	if got := f.Name().FirstName(); got != "" {
		t.Errorf("FirstName with empty list = %q, want \"\"", got)
	}
}

func TestAddressGenerators(t *testing.T) {
	f := newF()
	a := f.Address()
	if a.City() == "" {
		t.Error("City empty")
	}
	if a.CityPrefix() == "" {
		t.Error("CityPrefix empty")
	}
	if a.CitySuffix() == "" {
		t.Error("CitySuffix empty")
	}
	if a.StreetName() == "" {
		t.Error("StreetName empty")
	}
	if a.StreetSuffix() == "" {
		t.Error("StreetSuffix empty")
	}
	if a.StreetAddress() == "" {
		t.Error("StreetAddress empty")
	}
	if a.SecondaryAddress() == "" {
		t.Error("SecondaryAddress empty")
	}
	if a.BuildingNumber() == "" {
		t.Error("BuildingNumber empty")
	}
	if a.Community() == "" {
		t.Error("Community empty")
	}
	if a.ZipCode() == "" {
		t.Error("ZipCode empty")
	}
	if a.Zip() == "" {
		t.Error("Zip empty")
	}
	if a.Postcode() == "" {
		t.Error("Postcode empty")
	}
	if a.State() == "" {
		t.Error("State empty")
	}
	if a.StateAbbr() == "" {
		t.Error("StateAbbr empty")
	}
	if a.Country() == "" {
		t.Error("Country empty")
	}
	if a.CountryCode() == "" {
		t.Error("CountryCode empty")
	}
	if a.CountryCodeLong() == "" {
		t.Error("CountryCodeLong empty")
	}
	if a.TimeZone() == "" {
		t.Error("TimeZone empty")
	}
	if a.FullAddress() == "" {
		t.Error("FullAddress empty")
	}
	if lat := a.Latitude(); lat < -90 || lat >= 90 {
		t.Errorf("Latitude %v out of range", lat)
	}
	if lon := a.Longitude(); lon < -180 || lon >= 180 {
		t.Errorf("Longitude %v out of range", lon)
	}
}

func TestPhoneGenerators(t *testing.T) {
	f := newF()
	p := f.PhoneNumber()
	if p.PhoneNumber() == "" {
		t.Error("PhoneNumber empty")
	}
	if p.CellPhone() == "" {
		t.Error("CellPhone empty")
	}
	if !strings.HasPrefix(p.CountryCode(), "+") {
		t.Error("CountryCode must start with +")
	}
	if p.AreaCode() == "" {
		t.Error("AreaCode empty")
	}
	if p.ExchangeCode() == "" {
		t.Error("ExchangeCode empty")
	}
	if got := p.SubscriberNumber(4); len(got) != 4 {
		t.Errorf("SubscriberNumber(4) = %q want len 4", got)
	}
	if got := p.Extension(3); len(got) != 3 {
		t.Errorf("Extension(3) = %q want len 3", got)
	}
	if p.SubscriberNumber(0) != "" {
		t.Error("SubscriberNumber(0) must be empty")
	}
	if !strings.HasPrefix(p.PhoneNumberWithCountryCode(), "+") {
		t.Error("PhoneNumberWithCountryCode must start with +")
	}
	if !strings.HasPrefix(p.CellPhoneWithCountryCode(), "+") {
		t.Error("CellPhoneWithCountryCode must start with +")
	}
}

func TestLoremGenerators(t *testing.T) {
	f := newF()
	l := f.Lorem()
	if l.Word() == "" {
		t.Error("Word empty")
	}
	if got := l.Words(3, false); len(got) != 3 {
		t.Errorf("Words(3) len = %d", len(got))
	}
	// supplemental + count exceeding pool exercises the multiplier path.
	if got := l.Words(2000, true); len(got) != 2000 {
		t.Errorf("Words(2000, supplemental) len = %d", len(got))
	}
	if got := l.Characters(12); len(got) != 12 {
		t.Errorf("Characters(12) len = %d", len(got))
	}
	if l.Characters(0) != "" {
		t.Error("Characters(0) must be empty")
	}
	s := l.Sentence(4, false)
	if !strings.HasSuffix(s, ".") || s[0] < 'A' || s[0] > 'Z' {
		t.Errorf("Sentence not capitalized/terminated: %q", s)
	}
	if got := l.Sentences(2, false); len(got) != 2 {
		t.Errorf("Sentences(2) len = %d", len(got))
	}
	if l.Paragraph(3, false) == "" {
		t.Error("Paragraph empty")
	}
	if got := l.Paragraphs(2, false); len(got) != 2 {
		t.Errorf("Paragraphs(2) len = %d", len(got))
	}
}

func TestLoremEmptyPoolAndCapitalize(t *testing.T) {
	f := newF()
	f.data = &dataset{
		raw:       map[string]json.RawMessage{"lorem.words": json.RawMessage(`[]`)},
		strCache:  map[string][]string{},
		listCache: map[string][][]string{},
	}
	l := f.Lorem()
	if got := l.Words(3, false); got != nil {
		t.Errorf("Words with empty pool = %v, want nil", got)
	}
	if l.Word() != "" {
		t.Error("Word with empty pool must be empty")
	}
	if capitalize("") != "" {
		t.Error("capitalize(\"\") must be empty")
	}
	if capitalize("9abc") != "9abc" {
		t.Errorf("capitalize leading digit unchanged: %q", capitalize("9abc"))
	}
}

func TestNumberGenerators(t *testing.T) {
	f := newF()
	n := f.Number()
	if n.Number(0) != 0 {
		t.Error("Number(0) must be 0")
	}
	if v := n.Number(1); v < 0 || v > 9 {
		t.Errorf("Number(1) = %d out of [0,9]", v)
	}
	if v := n.Number(5); v < 10000 {
		t.Errorf("Number(5) = %d should have 5 digits", v)
	}
	if d := n.Digit(); d < 0 || d > 9 {
		t.Errorf("Digit = %d", d)
	}
	if d := n.NonZeroDigit(); d < 1 || d > 9 {
		t.Errorf("NonZeroDigit = %d", d)
	}
	if got := n.Hexadecimal(6); len(got) != 6 {
		t.Errorf("Hexadecimal(6) = %q", got)
	}
	if got := n.Binary(8); len(got) != 8 {
		t.Errorf("Binary(8) = %q", got)
	}
	if v := n.Decimal(4, 2); v <= 0 {
		t.Errorf("Decimal positive expected, got %v", v)
	}
	if v := n.BetweenInt(1, 100); v < 1 || v > 100 {
		t.Errorf("BetweenInt out of range: %d", v)
	}
	if v := n.Between(1.0, 2.0); v < 1 || v > 2 {
		t.Errorf("Between out of range: %v", v)
	}
	if v := n.Positive(-5, 5); v < 0 {
		t.Errorf("Positive returned negative: %v", v)
	}
	if v := n.Negative(-5, 5); v > 0 {
		t.Errorf("Negative returned positive: %v", v)
	}
	_ = n.Normal(0, 1) // distribution only; just exercise the path
}

func TestCompanyAndCommerce(t *testing.T) {
	f := newF()
	c := f.Company()
	if c.Name() == "" {
		t.Error("Company.Name empty")
	}
	if c.Suffix() == "" {
		t.Error("Company.Suffix empty")
	}
	if c.Industry() == "" {
		t.Error("Company.Industry empty")
	}
	if c.Buzzword() == "" {
		t.Error("Company.Buzzword empty")
	}
	if c.CatchPhrase() == "" {
		t.Error("Company.CatchPhrase empty")
	}
	if c.BS() == "" {
		t.Error("Company.BS empty")
	}
	cm := f.Commerce()
	if cm.ProductName() == "" {
		t.Error("ProductName empty")
	}
	if cm.Department(3, false) == "" {
		t.Error("Department empty")
	}
	if cm.Department(1, true) == "" {
		t.Error("Department fixed single empty")
	}
	// Multi-category department exercises the join-with-separator branch.
	found := false
	for i := 0; i < 20 && !found; i++ {
		if strings.Contains(cm.Department(5, true), " & ") {
			found = true
		}
	}
	if !found {
		t.Error("expected a multi-category department with separator")
	}
	if p := cm.Price(100); p < 0 || p >= 100 {
		t.Errorf("Price out of range: %v", p)
	}
}

func TestDepartmentEmpty(t *testing.T) {
	f := newF()
	f.data = &dataset{
		raw:       map[string]json.RawMessage{"commerce.department": json.RawMessage(`[]`)},
		strCache:  map[string][]string{},
		listCache: map[string][][]string{},
	}
	if got := f.Commerce().Department(3, true); got != "" {
		t.Errorf("Department with empty data = %q", got)
	}
}

func TestMiscColorBoolean(t *testing.T) {
	f := newF()
	if f.Color().ColorName() == "" {
		t.Error("ColorName empty")
	}
	hex := f.Color().HexColor()
	if len(hex) != 7 || hex[0] != '#' {
		t.Errorf("HexColor malformed: %q", hex)
	}
	rgb := f.Color().RGBColor()
	for _, v := range rgb {
		if v < 0 || v > 255 {
			t.Errorf("RGB component out of range: %d", v)
		}
	}
	b := f.Boolean()
	if !b.Boolean(1.0) {
		t.Error("Boolean(1.0) must be true")
	}
	if b.Boolean(0.0) {
		t.Error("Boolean(0.0) must be false")
	}
}

func TestHexColorAllHueSectors(t *testing.T) {
	// Drive hslToHex across all six hue sectors and the hex2 clamp paths.
	for h := 0.0; h < 360; h += 30 {
		got := hslToHex(h, 0.5, 0.5)
		if len(got) != 7 || got[0] != '#' {
			t.Errorf("hslToHex(%v) malformed: %q", h, got)
		}
	}
	if hex2(-1) != "00" {
		t.Errorf("hex2(-1) = %q", hex2(-1))
	}
	if hex2(999) != "ff" {
		t.Errorf("hex2(999) = %q", hex2(999))
	}
	if absF(-3) != 3 || absF(3) != 3 {
		t.Error("absF wrong")
	}
	if modF(-1, 3) != 2 {
		t.Errorf("modF(-1,3) = %v", modF(-1, 3))
	}
	if modF(7, 3) != 1 {
		t.Errorf("modF(7,3) = %v", modF(7, 3))
	}
}

func TestDateGenerators(t *testing.T) {
	f := newF()
	d := f.Date()
	from := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC)
	got := d.Between(from, to)
	if got.Before(from) || got.After(to) {
		t.Errorf("Between out of range: %v", got)
	}
	// Reversed bounds exercise the lo/hi swap branch; the offset is added from
	// the first argument, so the span is non-negative either way.
	got = d.Between(to, from)
	if got.Before(to) {
		t.Errorf("Between reversed before base: %v", got)
	}
	base := time.Date(2010, 6, 15, 12, 0, 0, 0, time.UTC)
	fwd := d.Forward(base, 30)
	if !fwd.After(truncDay(base)) {
		t.Errorf("Forward not in future: %v", fwd)
	}
	bwd := d.Backward(base, 30)
	if !bwd.Before(truncDay(base)) {
		t.Errorf("Backward not in past: %v", bwd)
	}
}

func TestTimeGenerators(t *testing.T) {
	f := newF()
	tg := f.Time()
	from := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2000, 1, 2, 0, 0, 0, 0, time.UTC)
	got := tg.Between(from, to)
	if got.Before(from) || got.After(to) {
		t.Errorf("Time.Between out of range: %v", got)
	}
	base := time.Date(2010, 6, 15, 12, 0, 0, 0, time.UTC)
	if !tg.Forward(base, 5).After(base) {
		t.Error("Time.Forward not in future")
	}
	if !tg.Backward(base, 5).Before(base) {
		t.Error("Time.Backward not in past")
	}
}

func TestInternetGenerators(t *testing.T) {
	f := newF()
	in := f.Internet()
	if in.Username() == "" {
		t.Error("Username empty")
	}
	if in.DomainWord() == "" {
		t.Error("DomainWord empty")
	}
	if !strings.Contains(in.DomainName(), ".") {
		t.Error("DomainName missing dot")
	}
	if in.DomainSuffix() == "" {
		t.Error("DomainSuffix empty")
	}
	if in.SafeDomainSuffix() == "" {
		t.Error("SafeDomainSuffix empty")
	}
	if !strings.Contains(in.Email(), "@") {
		t.Error("Email missing @")
	}
	if in.Slug() == "" {
		t.Error("Slug empty")
	}
	if !strings.HasPrefix(in.URL(), "http://") {
		t.Error("URL missing scheme")
	}
	ip4 := in.IPv4Address()
	if strings.Count(ip4, ".") != 3 {
		t.Errorf("IPv4 malformed: %q", ip4)
	}
	ip6 := in.IPv6Address()
	if strings.Count(ip6, ":") != 7 {
		t.Errorf("IPv6 malformed: %q", ip6)
	}
	mac := in.MacAddress("")
	if strings.Count(mac, ":") != 5 {
		t.Errorf("MAC malformed: %q", mac)
	}
	macp := in.MacAddress("aa:bb")
	if !strings.HasPrefix(macp, "aa:bb:") {
		t.Errorf("MAC prefix not honored: %q", macp)
	}
	// Invalid hex group in prefix is skipped.
	if in.MacAddress("zz") == "" {
		t.Error("MacAddress with bad prefix should still produce a value")
	}
}

func TestInternetPassword(t *testing.T) {
	f := newF()
	in := f.Internet()
	p := in.Password(10, 20, true, true)
	if len(p) < 10 || len(p) > 20 {
		t.Errorf("Password length %d out of [10,20]", len(p))
	}
	if !strings.ContainsAny(p, "ABCDEFGHIJKLMNOPQRSTUVWXYZ") {
		t.Errorf("mixCase password lacks uppercase: %q", p)
	}
	if !strings.ContainsAny(p, "!@#$%^&*") {
		t.Errorf("specialChars password lacks special: %q", p)
	}
	// Clamp branches: min<1, max<min.
	if got := in.Password(0, 0, false, false); len(got) < 1 {
		t.Error("Password(0,0) should clamp to >=1")
	}
}

func TestCharPrepareAndSanitize(t *testing.T) {
	if got := charPrepare("Jöhn-Döe!"); got != "joehn-doee" {
		t.Errorf("charPrepare = %q", got)
	}
	if got := sanitizeEmailLocal("a b"); got != "a#b" {
		t.Errorf("sanitizeEmailLocal = %q", got)
	}
	if got := sanitizeEmailLocal("a.b_c"); got != "a.b_c" {
		t.Errorf("sanitizeEmailLocal allowed punctuation altered: %q", got)
	}
	if got := fixUmlauts("ÄßöÜ"); got != "aessoeue" {
		t.Errorf("fixUmlauts = %q", got)
	}
}
