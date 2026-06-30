// Copyright (c) 2026, the go-ruby-faker/faker authors. All rights reserved.
// Use of this source code is governed by a BSD-3-Clause license that can be
// found in the LICENSE file.

package faker

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestConfigAccessors(t *testing.T) {
	r := NewRandom(7)
	f := New(r)
	if f.Random() != r {
		t.Error("Random accessor mismatch")
	}
	if f.Locale() != "en" {
		t.Errorf("default locale = %q", f.Locale())
	}
	f.SetLocale("fr")
	if f.Locale() != "fr" {
		t.Error("SetLocale failed")
	}
	r2 := NewRandom(8)
	f.SetRandom(r2)
	if f.Random() != r2 {
		t.Error("SetRandom failed")
	}
	if r.Seed() != 7 {
		t.Errorf("Seed = %d", r.Seed())
	}
}

func TestBaseHelpers(t *testing.T) {
	f := NewSeeded(1)
	if f.Sample(nil) != "" {
		t.Error("Sample(nil) must be empty")
	}
	if got := f.Sample([]string{"only"}); got != "only" {
		t.Errorf("Sample single = %q", got)
	}
	if got := f.SampleN([]string{"a", "b", "c"}, 2); len(got) != 2 {
		t.Errorf("SampleN len = %d", len(got))
	}
	if got := f.Shuffle([]string{"a", "b", "c", "d"}); len(got) != 4 {
		t.Errorf("Shuffle len = %d", len(got))
	}
	if got := f.Numerify("##", false); len(got) != 2 {
		t.Errorf("Numerify = %q", got)
	}
	if got := f.Numerify("##", true); len(got) != 2 {
		t.Errorf("Numerify leadingZero = %q", got)
	}
	if got := f.Letterify("??"); len(got) != 2 {
		t.Errorf("Letterify = %q", got)
	}
	if got := f.Bothify("#?"); len(got) != 2 {
		t.Errorf("Bothify = %q", got)
	}
	if got := f.FetchAll("color.name"); len(got) == 0 {
		t.Error("FetchAll returned empty")
	}
}

func TestSampleNFallbackPaths(t *testing.T) {
	f := NewSeeded(2)
	pool := make([]string, 12)
	for i := range pool {
		pool[i] = string(rune('a' + i))
	}
	// n in (6, length): partial Fisher-Yates branch.
	if got := f.SampleN(pool, 8); len(got) != 8 {
		t.Errorf("SampleN(8) len = %d", len(got))
	}
	// n >= length: full shuffle branch.
	if got := f.SampleN(pool, 20); len(got) != len(pool) {
		t.Errorf("SampleN(>=len) len = %d", len(got))
	}
	// empty / non-positive.
	if got := f.SampleN(nil, 3); len(got) != 0 {
		t.Errorf("SampleN(nil) = %v", got)
	}
	if got := f.SampleN(pool, 0); len(got) != 0 {
		t.Errorf("SampleN(0) = %v", got)
	}
}

func TestRandomCore(t *testing.T) {
	r := NewRandom(123)
	if v := r.Float(); v < 0 || v >= 1 {
		t.Errorf("Float out of [0,1): %v", v)
	}
	if v := r.Intn(10); v < 0 || v >= 10 {
		t.Errorf("Intn out of range: %v", v)
	}
	if v := r.IntInRange(5, 5); v != 5 {
		t.Errorf("IntInRange(5,5) = %d", v)
	}
	// Reversed range exercises the swap branch.
	if v := r.IntInRange(9, 1); v < 1 || v > 9 {
		t.Errorf("IntInRange reversed out of range: %d", v)
	}
	if v := r.FloatScaled(10); v < 0 || v >= 10 {
		t.Errorf("FloatScaled out of range: %v", v)
	}
	if v := r.FloatScaled(0); v < 0 || v >= 1 {
		t.Errorf("FloatScaled(0) fallback out of range: %v", v)
	}
	if v := r.FloatInRange(2, 5); v < 2 || v > 5 {
		t.Errorf("FloatInRange out of range: %v", v)
	}
	if v := r.FloatInRange(5, 2); v < 2 || v > 5 {
		t.Errorf("FloatInRange reversed out of range: %v", v)
	}
	b := r.Bytes(7)
	if len(b) != 7 {
		t.Errorf("Bytes len = %d", len(b))
	}
	// limit 0 short-circuit.
	if v := r.limitedRand(0); v != 0 {
		t.Errorf("limitedRand(0) = %d", v)
	}
}

func TestRandomIntnPanic(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Error("Intn(0) should panic")
		}
	}()
	NewRandom(1).Intn(0)
}

func TestLargeLimitTwoWord(t *testing.T) {
	// A limit beyond 2^32 forces the two-word assembly path in limitedRand.
	r := NewRandom(99)
	const big = uint64(1) << 40
	got := r.limitedRand(big)
	if got > big {
		t.Errorf("limitedRand(big) = %d exceeds limit", got)
	}
}

func TestSeedEdgeCases(t *testing.T) {
	// Zero seed: single-word init_genrand path.
	if k := seedKey(0); len(k) != 1 || k[0] != 0 {
		t.Errorf("seedKey(0) = %v", k)
	}
	// Negative seed: |seed| path.
	if k := seedKey(-5); len(k) != 1 || k[0] != 5 {
		t.Errorf("seedKey(-5) = %v", k)
	}
	// Large seed needing two 32-bit words triggers initByArray.
	big := int64(1)<<40 | 7
	if k := seedKey(big); len(k) != 2 {
		t.Errorf("seedKey(big) len = %d", len(k))
	}
	r := NewRandom(big)
	_ = r.Float() // exercise initByArray-seeded draw
}

func TestEntropyRandom(t *testing.T) {
	r := NewRandomEntropy()
	if r == nil {
		t.Fatal("NewRandomEntropy returned nil")
	}
	_ = r.Float()
}

func TestRoundN(t *testing.T) {
	if got := roundN(1.2345, 2); got != 1.23 {
		t.Errorf("roundN(1.2345,2) = %v", got)
	}
	if got := roundN(-1.235, 2); got != -1.24 {
		t.Errorf("roundN(-1.235,2) = %v", got)
	}
}

func TestUnique(t *testing.T) {
	u := NewUnique(0) // default budget
	seen := map[string]bool{}
	counter := 0
	gen := func() string {
		counter++
		return strings.Repeat("x", counter%5)
	}
	for i := 0; i < 4; i++ {
		v, err := u.Next(gen)
		if err != nil {
			t.Fatalf("Next error: %v", err)
		}
		if seen[v] {
			t.Errorf("duplicate value %q", v)
		}
		seen[v] = true
	}
	u.Clear()
	if len(u.seen) != 0 {
		t.Error("Clear did not reset")
	}

	// Exhaustion: a generator that always returns the same value.
	u2 := NewUnique(3)
	if _, err := u2.Next(func() string { return "const" }); err != nil {
		t.Fatalf("first Next: %v", err)
	}
	_, err := u2.Next(func() string { return "const" })
	if err == nil {
		t.Fatal("expected RetriesExhaustedError")
	}
	var re *RetriesExhaustedError
	if !errorsAs(err, &re) {
		t.Fatalf("wrong error type: %v", err)
	}
	if !strings.Contains(re.Error(), "retry limit") {
		t.Errorf("error message = %q", re.Error())
	}
}

// errorsAs is a tiny local shim to avoid importing errors solely for one As.
func errorsAs(err error, target **RetriesExhaustedError) bool {
	if e, ok := err.(*RetriesExhaustedError); ok {
		*target = e
		return true
	}
	return false
}

func TestRegexify(t *testing.T) {
	f := NewSeeded(3)
	cases := []string{
		`/[A-Z]{3}/`,
		`^\d{2}-\d{4}$`,
		`(cat|dog|bird)`,
		`[a-c]?x{2,4}`,
		`\w{5}`,
		`plain`,
		`[abc]{2}`,
		`a{0}b`, // {0} repeat
	}
	for _, c := range cases {
		got := f.Regexify(c)
		_ = got // value is RNG-driven; we exercise every branch
	}
	// Unterminated constructs hit the fallback byte-copy branches.
	for _, c := range []string{`[abc`, `(alt`, `x{2`, `\`} {
		_ = f.regexify(c)
	}
	// letterAt both halves.
	if letterAt(0) != 'A' || letterAt(51) != 'z' {
		t.Error("letterAt boundaries wrong")
	}
	// parseRange / atoi edges.
	if _, _, ok := parseRange("x"); ok {
		t.Error("parseRange non-numeric should fail")
	}
	if _, ok := atoi(""); ok {
		t.Error("atoi empty should fail")
	}
	if _, ok := atoi("1a"); ok {
		t.Error("atoi non-digit should fail")
	}
	if matchParen("(abc", 0) != -1 {
		t.Error("matchParen unbalanced should be -1")
	}
}

func TestRegexifyEmptyClassAndBadRange(t *testing.T) {
	f := NewSeeded(4)
	// Empty character class -> sampleClass returns '?'.
	if got := f.regexify("[]"); got == "" {
		t.Error("empty class produced empty output")
	}
	// Reversed range bounds in class are skipped (lo>hi).
	_ = f.regexify("[z-a]")
}

func TestDispatchFallbacks(t *testing.T) {
	f := NewSeeded(5)
	// Cross-class token with a method not in callMethod -> fetch fallback.
	if got := f.dispatch("address", "Color", "name"); got == "" {
		t.Error("cross-class fetch fallback empty")
	}
	// No-class token resolving via classFor self-method.
	if got := f.dispatch("name", "", "first_name"); got == "" {
		t.Error("self-method dispatch empty")
	}
	// classFor unknown key.
	if classFor("nope") != "" {
		t.Error("classFor unknown should be empty")
	}
	for _, k := range []string{"name", "address", "company", "phone_number", "cell_phone", "internet", "commerce"} {
		if classFor(k) == "" {
			t.Errorf("classFor(%q) empty", k)
		}
	}
	// snakeCase.
	if snakeCase("PhoneNumber") != "phone_number" {
		t.Errorf("snakeCase = %q", snakeCase("PhoneNumber"))
	}
}

func TestCallMethodCoverage(t *testing.T) {
	f := NewSeeded(6)
	pairs := [][2]string{
		{"Name", "name"}, {"Name", "name_with_middle"}, {"Name", "first_name"},
		{"Name", "male_first_name"}, {"Name", "female_first_name"},
		{"Name", "neutral_first_name"}, {"Name", "last_name"}, {"Name", "middle_name"},
		{"Name", "prefix"}, {"Name", "suffix"},
		{"Address", "city"}, {"Address", "city_prefix"}, {"Address", "city_suffix"},
		{"Address", "street_name"}, {"Address", "street_address"}, {"Address", "street_suffix"},
		{"Address", "secondary_address"}, {"Address", "building_number"},
		{"Address", "community"}, {"Address", "zip_code"}, {"Address", "state"},
		{"Address", "state_abbr"}, {"Address", "country"}, {"Address", "country_code"},
		{"Company", "name"}, {"Company", "suffix"},
		{"PhoneNumber", "area_code"}, {"PhoneNumber", "exchange_code"},
		{"Commerce", "product_name"},
	}
	for _, p := range pairs {
		if v, ok := f.callMethod(p[0], p[1]); !ok {
			t.Errorf("callMethod(%q,%q) not handled", p[0], p[1])
		} else if v == "" {
			t.Errorf("callMethod(%q,%q) empty", p[0], p[1])
		}
	}
	// Unknown class and unknown method -> ok=false.
	if _, ok := f.callMethod("Nope", "x"); ok {
		t.Error("unknown class should be ok=false")
	}
	if _, ok := f.callMethod("Name", "no_such"); ok {
		t.Error("unknown method should be ok=false")
	}
}

func TestParseTokensEdges(t *testing.T) {
	// splitClassMeth: leading non-alpha class prefix is treated as plain key.
	cls, meth := splitClassMeth("9foo.bar")
	if cls != "" || meth != "9foo.bar" {
		t.Errorf("splitClassMeth digit-prefix = (%q,%q)", cls, meth)
	}
	cls, meth = splitClassMeth("city_prefix")
	if cls != "" || meth != "city_prefix" {
		t.Errorf("splitClassMeth no-dot = (%q,%q)", cls, meth)
	}
	// parseTokens with no token returns nil; with trailing suffix retains it.
	if toks := parseTokens("no tokens here"); toks != nil {
		t.Errorf("parseTokens plain = %v", toks)
	}
	toks := parseTokens("#{Name.first_name} of #{city}")
	if len(toks) != 2 {
		t.Fatalf("parseTokens count = %d", len(toks))
	}
	// Unterminated #{ stops the scan.
	if toks := parseTokens("#{unterminated"); toks != nil {
		t.Errorf("parseTokens unterminated = %v", toks)
	}
}

func TestParseNumerifyFallback(t *testing.T) {
	// A format with no #{} tokens but with '#' digits gets numerified.
	f := NewSeeded(7)
	f.data = &dataset{
		raw:       map[string]json.RawMessage{"x.y": json.RawMessage(`["##-##"]`)},
		strCache:  map[string][]string{},
		listCache: map[string][][]string{},
	}
	got := f.Parse("x.y")
	if len(got) != 5 || got[2] != '-' {
		t.Errorf("Parse numerify fallback = %q", got)
	}
}

func TestDataAccessors(t *testing.T) {
	d := enData
	if len(d.strings("color.name")) == 0 {
		t.Error("strings missing color.name")
	}
	if d.strings("color.name") == nil {
		t.Error("strings cache miss")
	}
	// Missing key.
	if d.strings("does.not.exist") != nil {
		t.Error("missing key should be nil")
	}
	// Re-call hits cache.
	_ = d.strings("does.not.exist")
	// scalar from a string value.
	if d.scalar("separator") == "" {
		t.Error("scalar separator empty")
	}
	// scalar missing key.
	if d.scalar("nope") != "" {
		t.Error("scalar missing should be empty")
	}
	// scalar from an array value -> first element.
	if d.scalar("color.name") == "" {
		t.Error("scalar from array empty")
	}
	// lists.
	if len(d.lists("company.buzzwords")) == 0 {
		t.Error("lists empty")
	}
	_ = d.lists("company.buzzwords") // cache hit
	if d.lists("nope") != nil {
		t.Error("lists missing should be nil")
	}
	// object.
	if d.object("lorem.punctuation") == nil {
		t.Error("object missing lorem.punctuation")
	}
	if d.object("nope") != nil {
		t.Error("object missing should be nil")
	}
	// flatten.
	if len(d.flatten("company.buzzwords")) == 0 {
		t.Error("flatten empty")
	}
}

func TestDataMalformedDecodes(t *testing.T) {
	d := &dataset{
		raw: map[string]json.RawMessage{
			"arr":  json.RawMessage(`{"not":"array"}`), // not []string and not string
			"list": json.RawMessage(`["flat"]`),        // not [][]string
			"obj":  json.RawMessage(`["flat"]`),        // not object
		},
		strCache:  map[string][]string{},
		listCache: map[string][][]string{},
	}
	if d.strings("arr") != nil {
		t.Error("strings on object should be nil")
	}
	if d.lists("list") != nil {
		t.Error("lists on flat array should be nil")
	}
	if d.object("obj") != nil {
		t.Error("object on array should be nil")
	}
	if d.scalar("missing") != "" {
		t.Error("scalar missing empty")
	}
}

func TestMustLoadPanics(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Error("mustLoad should panic on invalid JSON")
		}
	}()
	mustLoad([]byte("not json"))
}
