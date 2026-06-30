// Copyright (c) 2026, the go-ruby-faker/faker authors. All rights reserved.
// Use of this source code is governed by a BSD-3-Clause license that can be
// found in the LICENSE file.

package faker

import (
	"strconv"
	"strings"
)

// InternetGen is Faker::Internet.
type InternetGen struct{ f *Faker }

// Internet returns the Faker::Internet generator.
func (f *Faker) Internet() *InternetGen { return &InternetGen{f} }

var internetSeparators = []string{".", "_"}

// charPrepare reproduces Faker::Char.prepare for the en locale: fix umlauts,
// strip everything that is not a word char or '-', and downcase.
func charPrepare(s string) string {
	s = fixUmlauts(s)
	var b strings.Builder
	for _, r := range s {
		switch {
		case r >= 'a' && r <= 'z':
			b.WriteRune(r)
		case r >= 'A' && r <= 'Z':
			b.WriteRune(r + ('a' - 'A'))
		case r >= '0' && r <= '9':
			b.WriteRune(r)
		case r == '_' || r == '-':
			b.WriteRune(r)
		}
	}
	return b.String()
}

func fixUmlauts(s string) string {
	repl := strings.NewReplacer(
		"ä", "ae", "Ä", "ae",
		"ö", "oe", "Ö", "oe",
		"ü", "ue", "Ü", "ue",
		"ß", "ss",
	)
	return repl.Replace(s)
}

// Username produces a username (Faker::Internet.username with no specifier).
//
// Faithful to the gem's sample over two freshly-evaluated options:
//
//	sample([
//	  Char.prepare(Faker::Name.first_name),
//	  [Faker::Name.first_name, Faker::Name.last_name].map { Char.prepare }.join(sample(separators)),
//	])
//
// so the draw sequence is first_name (option 0), then first_name + last_name +
// sample(separators) (option 1), then the 2-element sample, then downcase. Note
// the two options each draw their OWN first_name — they are not shared.
func (in *InternetGen) Username() string {
	opt0 := charPrepare(in.f.Name().FirstName())
	first := charPrepare(in.f.Name().FirstName())
	last := charPrepare(in.f.Name().LastName())
	opt1 := first + in.f.Sample(internetSeparators) + last
	return strings.ToLower(in.f.Sample([]string{opt0, opt1}))
}

// DomainWord produces the leading word of a domain (Faker::Internet.domain_word):
// the prepared first token of a company name.
func (in *InternetGen) DomainWord() string {
	name := in.f.Company().Name()
	first := name
	if i := strings.IndexByte(name, ' '); i >= 0 {
		first = name[:i]
	}
	return charPrepare(first)
}

// DomainSuffix produces a domain suffix (Faker::Internet.domain_suffix).
func (in *InternetGen) DomainSuffix() string { return in.f.Fetch("internet.domain_suffix") }

// SafeDomainSuffix produces a safe domain suffix ("example"/"test").
func (in *InternetGen) SafeDomainSuffix() string { return in.f.Fetch("internet.safe_domain_suffix") }

// DomainName produces "word.suffix" using a safe suffix (Faker::Internet.domain_name).
func (in *InternetGen) DomainName() string {
	return in.DomainWord() + "." + in.SafeDomainSuffix()
}

// Email produces an email address (Faker::Internet.email): a username local
// part at a generated domain.
func (in *InternetGen) Email() string {
	local := sanitizeEmailLocal(in.Username())
	return local + "@" + in.DomainName()
}

// emailLocalAllowed is the gem's sanitize_email_local_part char_range:
// 0-9, A-Z, a-z and the RFC local-part specials "!#$%&'*+-/=?^_`{|}~.".
const emailLocalAllowed = "!#$%&'*+-/=?^_`{|}~."

// sanitizeEmailLocal mirrors Faker::Internet.sanitize_email_local_part: each
// character outside the allowed set is replaced by '#' (no collapsing, no
// trimming).
func sanitizeEmailLocal(s string) string {
	var b strings.Builder
	for _, r := range s {
		ok := (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') ||
			(r >= '0' && r <= '9') || strings.ContainsRune(emailLocalAllowed, r)
		if ok {
			b.WriteRune(r)
		} else {
			b.WriteByte('#')
		}
	}
	return b.String()
}

// Slug produces a URL slug of two words joined by '-' or '_'
// (Faker::Internet.slug).
func (in *InternetGen) Slug() string {
	glue := in.f.Sample([]string{"-", "_"})
	words := in.f.SampleN(in.f.data.strings("internet.slug"), 2)
	return strings.Join(words, glue)
}

// URL produces "http://domain/username" (Faker::Internet.url).
func (in *InternetGen) URL() string {
	return "http://" + in.DomainName() + "/" + in.Username()
}

// IPv4Address produces a dotted IPv4 address (Faker::Internet.ip_v4_address).
func (in *InternetGen) IPv4Address() string {
	parts := make([]string, 4)
	for i := range parts {
		parts[i] = strconv.Itoa(in.f.rng.IntInRange(0, 255))
	}
	return strings.Join(parts, ".")
}

// IPv6Address produces a colon-separated IPv6 address (Faker::Internet.ip_v6_address).
func (in *InternetGen) IPv6Address() string {
	parts := make([]string, 8)
	for i := range parts {
		parts[i] = strconv.FormatInt(int64(in.f.rng.Intn(65536)), 16)
	}
	return strings.Join(parts, ":")
}

// MacAddress produces a colon-separated MAC address (Faker::Internet.mac_address).
// prefix may be a partial "aa" or "aa:bb" hex prefix.
func (in *InternetGen) MacAddress(prefix string) string {
	var digits []int
	if prefix != "" {
		for _, h := range strings.Split(prefix, ":") {
			v, err := strconv.ParseInt(h, 16, 0)
			if err == nil {
				digits = append(digits, int(v))
			}
		}
	}
	for len(digits) < 6 {
		digits = append(digits, in.f.rng.Intn(256))
	}
	parts := make([]string, 6)
	const hexd = "0123456789abcdef"
	for i := 0; i < 6; i++ {
		v := digits[i] & 0xff
		parts[i] = string([]byte{hexd[v>>4], hexd[v&0xf]})
	}
	return strings.Join(parts, ":")
}

// Password produces a password between minLength and maxLength characters
// (Faker::Internet.password). With mixCase an uppercase letter is guaranteed;
// with specialChars a special character is guaranteed. The result is shuffled.
func (in *InternetGen) Password(minLength, maxLength int, mixCase, specialChars bool) string {
	if minLength < 1 {
		minLength = 1
	}
	if maxLength < minLength {
		maxLength = minLength
	}
	target := in.f.rng.IntInRange(minLength, maxLength)

	var password []string
	var bag []string

	lower := lettersLower()
	password = append(password, in.f.Sample(lower))
	bag = append(bag, lower...)

	digits := digitStrings()
	password = append(password, in.f.Sample(digits))
	bag = append(bag, digits...)

	if mixCase {
		upper := lettersUpper()
		password = append(password, in.f.Sample(upper))
		bag = append(bag, upper...)
	}
	if specialChars {
		special := []string{"!", "@", "#", "$", "%", "^", "&", "*"}
		password = append(password, in.f.Sample(special))
		bag = append(bag, special...)
	}
	for len(password) < target {
		password = append(password, in.f.Sample(bag))
	}
	return strings.Join(in.f.Shuffle(password), "")
}

func lettersLower() []string {
	out := make([]string, 26)
	for i := range out {
		out[i] = string(byte('a' + i))
	}
	return out
}

func lettersUpper() []string {
	out := make([]string, 26)
	for i := range out {
		out[i] = string(byte('A' + i))
	}
	return out
}

func digitStrings() []string {
	out := make([]string, 10)
	for i := range out {
		out[i] = strconv.Itoa(i)
	}
	return out
}
