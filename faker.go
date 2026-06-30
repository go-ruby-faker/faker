// Copyright (c) 2026, the go-ruby-faker/faker authors. All rights reserved.
// Use of this source code is governed by a BSD-3-Clause license that can be
// found in the LICENSE file.

// Package faker is a pure-Go (CGO=0) port of the Ruby `faker` gem: a seeded,
// MRI-faithful fake-data generator.
//
// # Deterministic seeding contract
//
// The gem's key correctness property is that, given
//
//	Faker::Config.random = Random.new(SEED)
//
// the sequence of generated values is reproducible. This package mirrors that:
// construct a [Faker] with [New] over a [Random] seeded by [NewRandom] and the
// draw sequence reproduces MRI's, because [Random] is a bit-exact port of
// Ruby's MT19937 (including init_by_array seeding and the mask-and-reject
// integer scheme) and the Array#sample / Array#shuffle helpers replay MRI's
// algorithms.
//
// # Seed-parity scope (honest)
//
// Exact-sequence parity with the gem holds for every generator whose randomness
// flows through INTEGER draws and single/sampled selections: Name, Address
// (city/street/zip/state/country/full_address), PhoneNumber, Lorem, Company,
// Commerce.product_name/department, Color.color_name, Boolean, Number's
// integer paths (number/digit/hexadecimal/binary/decimal), Internet's
// table-driven paths, and Date.between/forward/backward (Julian-day integer
// arithmetic).
//
// Distribution/format parity only (NOT bit-exact sequence) applies where the
// value derives from a FLOAT-range draw, because MRI's Random#rand(lo..hi) for
// float ranges uses a higher-precision internal algorithm this port does not
// reproduce: Address.latitude/longitude, Number.between/within/positive/
// negative/normal with non-integer bounds, Commerce.price, and Time.between/
// forward/backward. These match the gem's data tables, formats and statistical
// distribution; their exact per-seed value may differ in low-order digits.
package faker

import (
	"strconv"
	"strings"
)

// Faker holds the configuration (locale + RNG) and exposes the generators.
// It corresponds to Faker::Config plus the Faker::Base class methods.
type Faker struct {
	rng    *Random
	locale string
	data   *dataset
}

// New returns a Faker using the given seeded RNG and the default ("en") locale.
func New(rng *Random) *Faker {
	return &Faker{rng: rng, locale: "en", data: enData}
}

// NewSeeded is a convenience for New(NewRandom(seed)).
func NewSeeded(seed int64) *Faker { return New(NewRandom(seed)) }

// Random returns the underlying RNG (Faker::Config.random).
func (f *Faker) Random() *Random { return f.rng }

// SetRandom replaces the RNG (Faker::Config.random=).
func (f *Faker) SetRandom(r *Random) { f.rng = r }

// Locale reports the active locale (Faker::Config.locale).
func (f *Faker) Locale() string { return f.locale }

// SetLocale sets the locale. Only "en" ships embedded data; other locales fall
// back to "en" (mirroring the gem's en fallback for missing translations).
func (f *Faker) SetLocale(loc string) { f.locale = loc }

// ---- Faker::Base helpers ----

// Sample mirrors Faker::Base.sample(list): a single random element.
func (f *Faker) Sample(list []string) string {
	i, ok := sampleIndex(f.rng, len(list))
	if !ok {
		return ""
	}
	return list[i]
}

// SampleN mirrors Faker::Base.sample(list, n): n random elements.
func (f *Faker) SampleN(list []string, n int) []string {
	idx := sampleIndices(f.rng, len(list), n)
	out := make([]string, len(idx))
	for k, i := range idx {
		out[k] = list[i]
	}
	return out
}

// Shuffle mirrors Faker::Base.shuffle(list).
func (f *Faker) Shuffle(list []string) []string {
	idx := shuffleIndices(f.rng, len(list))
	out := make([]string, len(idx))
	for k, i := range idx {
		out[k] = list[i]
	}
	return out
}

// Numerify mirrors Faker::Base.numerify: each '#' becomes a digit. By default
// the FIRST '#' becomes 1..9 (no leading zero) and the rest 0..9; with
// leadingZero true every '#' becomes 0..9.
func (f *Faker) Numerify(s string, leadingZero bool) string {
	if leadingZero {
		return f.replaceHashes(s, true, false)
	}
	return f.replaceHashes(s, false, true)
}

// replaceHashes substitutes '#' tokens. When firstNonZero is true the first '#'
// draws 1..9 and the remainder 0..9; otherwise every '#' draws 0..9.
func (f *Faker) replaceHashes(s string, allZero, firstNonZero bool) string {
	var b strings.Builder
	first := true
	for _, ch := range s {
		if ch != '#' {
			b.WriteRune(ch)
			continue
		}
		if !allZero && firstNonZero && first {
			b.WriteString(strconv.Itoa(f.rng.IntInRange(1, 9)))
		} else {
			b.WriteString(strconv.Itoa(f.rng.Intn(10)))
		}
		first = false
	}
	return b.String()
}

// Letterify mirrors Faker::Base.letterify: each '?' becomes an uppercase A–Z.
func (f *Faker) Letterify(s string) string {
	var b strings.Builder
	for _, ch := range s {
		if ch == '?' {
			b.WriteByte(byte('A' + f.rng.Intn(26)))
		} else {
			b.WriteRune(ch)
		}
	}
	return b.String()
}

// Bothify mirrors Faker::Base.bothify: numerify then letterify.
func (f *Faker) Bothify(s string) string {
	return f.Letterify(f.Numerify(s, false))
}

// Fetch mirrors Faker::Base.fetch: sample one value from a locale array; if the
// chosen value is a /regex/ literal, regexify it.
func (f *Faker) Fetch(key string) string {
	v := f.Sample(f.data.strings(key))
	if len(v) >= 2 && v[0] == '/' && v[len(v)-1] == '/' {
		return f.Regexify(v)
	}
	return v
}

// FetchAll mirrors Faker::Base.fetch_all: the full locale array for a key.
func (f *Faker) FetchAll(key string) []string { return f.data.strings(key) }

// Parse mirrors Faker::Base.parse: fetch a format string from the locale and
// expand its #{token} placeholders by dispatching to generators (Name.,
// Address., …) or recursive locale lookups. If nothing parses, the string is
// numerified.
func (f *Faker) Parse(key string) string {
	fetched := f.Fetch(key)
	parts := parseTokens(fetched)
	if len(parts) == 0 {
		return f.Numerify(fetched, false)
	}
	keyPath := key
	if i := strings.IndexByte(key, '.'); i >= 0 {
		keyPath = key[:i]
	}
	var b strings.Builder
	for _, p := range parts {
		b.WriteString(p.prefix)
		b.WriteString(f.dispatch(keyPath, p.cls, p.meth))
		b.WriteString(p.suffix)
	}
	return b.String()
}

// Regexify mirrors Faker::Base.regexify: a deliberately simple regex generator
// matching the gem's supported subset (anchors, {n}/{m,n}, ?, [...] classes and
// ranges, (a|b) alternations, \d, \w). It draws with the same Array#sample over
// the configured RNG, so seeded output matches the gem for the patterns it
// supports.
func (f *Faker) Regexify(reg string) string {
	return f.regexify(reg)
}

// ---- parse token model ----

type token struct {
	prefix string // leading "(" if present
	cls    string // class prefix e.g. "Name", or "" for self
	meth   string // method/sub-key name
	suffix string // trailing literal text
}

// parseTokens reproduces Faker::Base#parse's scan:
// /(\(?)#\{([A-Za-z]+\.)?([^}]+)\}([^#]++)?/
func parseTokens(s string) []token {
	var out []token
	i := 0
	for i < len(s) {
		// Find next "#{"
		j := strings.Index(s[i:], "#{")
		if j < 0 {
			break
		}
		j += i
		prefix := ""
		if j > 0 && s[j-1] == '(' {
			prefix = "("
		}
		// content up to "}"
		k := strings.IndexByte(s[j+2:], '}')
		if k < 0 {
			break
		}
		inner := s[j+2 : j+2+k]
		cls, meth := splitClassMeth(inner)
		// suffix: text after "}" up to next "#"
		rest := s[j+2+k+1:]
		suffix := rest
		if h := strings.IndexByte(rest, '#'); h >= 0 {
			suffix = rest[:h]
			i = j + 2 + k + 1 + h
		} else {
			i = len(s)
		}
		out = append(out, token{prefix: prefix, cls: cls, meth: meth, suffix: suffix})
	}
	return out
}

// splitClassMeth splits "Name.first_name" into ("Name","first_name") and
// "city_prefix" into ("","city_prefix"). Only a leading [A-Za-z]+"." counts as
// a class prefix.
func splitClassMeth(inner string) (string, string) {
	dot := strings.IndexByte(inner, '.')
	if dot <= 0 {
		return "", inner
	}
	cand := inner[:dot]
	for _, c := range cand {
		if !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')) {
			return "", inner
		}
	}
	return cand, inner[dot+1:]
}

// snakeCase converts e.g. "PhoneNumber" to "phone_number" (the gem's just-enough
// snake-casing in #parse).
func snakeCase(s string) string {
	var b strings.Builder
	for i, c := range s {
		if c >= 'A' && c <= 'Z' {
			if i > 0 {
				b.WriteByte('_')
			}
			b.WriteByte(byte(c - 'A' + 'a'))
		} else {
			b.WriteRune(c)
		}
	}
	return b.String()
}
