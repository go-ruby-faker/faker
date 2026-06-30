// Copyright (c) 2026, the go-ruby-faker/faker authors. All rights reserved.
// Use of this source code is governed by a BSD-3-Clause license that can be
// found in the LICENSE file.

package faker

import (
	"encoding/json"
	"strings"
	"testing"
)

// withData returns a Faker whose dataset is built from the given key→raw-JSON map.
func withData(seed int64, kv map[string]string) *Faker {
	raw := make(map[string]json.RawMessage, len(kv))
	for k, v := range kv {
		raw[k] = json.RawMessage(v)
	}
	f := NewSeeded(seed)
	f.data = &dataset{raw: raw, strCache: map[string][]string{}, listCache: map[string][][]string{}}
	return f
}

func TestFetchRegexLiteral(t *testing.T) {
	// A fetched value that is a /regex/ literal is regexified.
	f := withData(1, map[string]string{"k.regex": `["/[A-Z]{3}/"]`})
	got := f.Fetch("k.regex")
	if len(got) != 3 {
		t.Errorf("Fetch regex literal = %q, want 3 letters", got)
	}
	for _, r := range got {
		if r < 'A' || r > 'Z' {
			t.Errorf("Fetch regex produced non A-Z: %q", got)
		}
	}
}

func TestCharPrepareDigits(t *testing.T) {
	// The 0-9 branch of charPrepare.
	if got := charPrepare("Bob42"); got != "bob42" {
		t.Errorf("charPrepare digits = %q", got)
	}
}

func TestLocalePunctuationFallback(t *testing.T) {
	// Missing lorem.punctuation object forces the default " " / "." fallbacks.
	f := withData(2, map[string]string{"lorem.words": `["foo","bar","baz"]`})
	l := f.Lorem()
	if l.localeSpace() != " " {
		t.Errorf("localeSpace fallback = %q", l.localeSpace())
	}
	if l.localePeriod() != "." {
		t.Errorf("localePeriod fallback = %q", l.localePeriod())
	}
	// Object present but lacking the keys also falls back.
	f2 := withData(3, map[string]string{"lorem.punctuation": `{"other":"x"}`})
	if f2.Lorem().localeSpace() != " " {
		t.Error("localeSpace missing-key fallback wrong")
	}
	if f2.Lorem().localePeriod() != "." {
		t.Error("localePeriod missing-key fallback wrong")
	}
}

func TestNumberSignBranches(t *testing.T) {
	f := NewSeeded(4)
	n := f.Number()
	// No-flip branches: already-positive / already-negative ranges.
	if v := n.Positive(1, 5); v < 0 {
		t.Errorf("Positive no-flip negative: %v", v)
	}
	if v := n.Negative(-5, -1); v > 0 {
		t.Errorf("Negative no-flip positive: %v", v)
	}
	// Flip branches: an all-negative range for Positive (v<0 -> -v) and an
	// all-positive range for Negative (v>0 -> -v).
	if v := n.Positive(-5, -1); v < 0 {
		t.Errorf("Positive flip stayed negative: %v", v)
	}
	if v := n.Negative(1, 5); v > 0 {
		t.Errorf("Negative flip stayed positive: %v", v)
	}
}

func TestZeroPadNoPad(t *testing.T) {
	// A value already at the requested width takes the no-pad branch.
	if got := zeroPad(12345, 3); got != "12345" {
		t.Errorf("zeroPad no-pad = %q", got)
	}
	if got := zeroPad(5, 3); got != "005" {
		t.Errorf("zeroPad pad = %q", got)
	}
}

func TestDataStringsScalarBranches(t *testing.T) {
	// strings: a bare scalar string value is wrapped in a single-element slice.
	f := withData(5, map[string]string{"k.scalar": `"hello"`})
	if got := f.data.strings("k.scalar"); len(got) != 1 || got[0] != "hello" {
		t.Errorf("strings scalar wrap = %v", got)
	}
	// scalar from a true scalar value.
	if got := f.data.scalar("k.scalar"); got != "hello" {
		t.Errorf("scalar = %q", got)
	}
}

func TestParseTokensTrailingSuffixNoHash(t *testing.T) {
	// A token whose suffix runs to end-of-string (no following '#') covers the
	// "i = len(s)" branch of parseTokens.
	toks := parseTokens("#{Name.first_name} suffix-to-end")
	if len(toks) != 1 {
		t.Fatalf("parseTokens count = %d", len(toks))
	}
	if !strings.Contains(toks[0].suffix, "suffix-to-end") {
		t.Errorf("suffix not captured: %q", toks[0].suffix)
	}
}

func TestRegexifyParenAndUnboundedRepeat(t *testing.T) {
	f := NewSeeded(6)
	// "?" quantifier on a parenthesised alternation, plus {m,n} on a class.
	_ = f.regexify("(a|b)?[x-z]{2,3}")
	// A bare quantifier with malformed {n} (no closing brace) leaves *i intact.
	_ = f.regexify("a{2")
	// parseRepeat with a non-numeric {} returns 1.
	_ = f.regexify("a{x}b")
	// \w token -> letterAt; \d already covered, plain escape (\.) default branch.
	got := f.regexify(`\w{4}`)
	if len(got) != 4 {
		t.Errorf(`\w{4} = %q`, got)
	}
	if dot := f.regexify(`\.`); dot != "." {
		t.Errorf(`\. = %q`, dot)
	}
}

func TestParseTokensParenPrefix(t *testing.T) {
	// A token immediately preceded by '(' captures the "(" prefix.
	toks := parseTokens("(#{Name.suffix})")
	if len(toks) != 1 {
		t.Fatalf("count = %d", len(toks))
	}
	if toks[0].prefix != "(" {
		t.Errorf("prefix = %q, want (", toks[0].prefix)
	}
}

func TestScalarEmptyArrayBranch(t *testing.T) {
	// A key whose value is an empty array: scalar() falls through both the
	// string decode and the array-first-element path to return "".
	f := withData(7, map[string]string{"k.empty": `[]`})
	if got := f.data.scalar("k.empty"); got != "" {
		t.Errorf("scalar empty array = %q", got)
	}
}
