// Copyright (c) 2026, the go-ruby-faker/faker authors. All rights reserved.
// Use of this source code is governed by a BSD-3-Clause license that can be
// found in the LICENSE file.

package faker

import "strings"

// regexify reproduces Faker::Base.regexify's gsub pipeline for the subset of
// regex syntax the gem supports: anchors (^ $ and surrounding slashes), {n}
// and {m,n} repetition, ? (==> {0,1}), [..] character classes including a-z
// ranges, (a|b|c) alternations, \d and \w. It draws each random choice with the
// same Array#sample over the configured RNG, so seeded output matches the gem
// for these patterns.
//
// It intentionally does NOT handle ., *, unbounded {n,}, lookahead, nested
// groups — exactly the gem's documented limitations.
func (f *Faker) regexify(reg string) string {
	s := reg
	// Ditch leading "/", "^" and trailing "$", "/".
	s = strings.TrimPrefix(s, "/")
	s = strings.TrimPrefix(s, "^")
	s = strings.TrimSuffix(s, "/")
	s = strings.TrimSuffix(s, "$")

	var b strings.Builder
	i := 0
	for i < len(s) {
		c := s[i]
		switch {
		case c == '\\' && i+1 < len(s):
			next := s[i+1]
			i += 2
			rep := f.parseRepeat(s, &i)
			tok := func() string {
				switch next {
				case 'd':
					return string(rune('0' + f.rng.Intn(10)))
				case 'w':
					return string(letterAt(f.rng.Intn(52)))
				default:
					return string(next)
				}
			}
			for r := 0; r < rep; r++ {
				b.WriteString(tok())
			}
		case c == '[':
			end := strings.IndexByte(s[i:], ']')
			if end < 0 {
				b.WriteByte(c)
				i++
				continue
			}
			class := s[i+1 : i+end]
			i += end + 1
			rep := f.parseRepeat(s, &i)
			for r := 0; r < rep; r++ {
				b.WriteByte(f.sampleClass(class))
			}
		case c == '(':
			end := matchParen(s, i)
			if end < 0 {
				b.WriteByte(c)
				i++
				continue
			}
			body := s[i+1 : end]
			i = end + 1
			rep := f.parseRepeat(s, &i)
			alts := strings.Split(body, "|")
			for r := 0; r < rep; r++ {
				b.WriteString(f.Sample(alts))
			}
		default:
			i++
			rep := f.parseRepeat(s, &i)
			for r := 0; r < rep; r++ {
				b.WriteByte(c)
			}
		}
	}
	return b.String()
}

// parseRepeat reads an optional quantifier at *i and returns how many copies to
// emit: "?" -> 0 or 1; "{n}" -> n; "{m,n}" -> a sample of [m, n]. With no
// quantifier it returns 1 and leaves *i unchanged.
func (f *Faker) parseRepeat(s string, i *int) int {
	if *i >= len(s) {
		return 1
	}
	if s[*i] == '?' {
		*i++
		return f.rng.Intn(2) // sample([0,1]) over a 2-element range
	}
	if s[*i] == '{' {
		end := strings.IndexByte(s[*i:], '}')
		if end < 0 {
			return 1
		}
		inner := s[*i+1 : *i+end]
		*i += end + 1
		lo, hi, ok := parseRange(inner)
		if !ok {
			return 1
		}
		if lo == hi {
			return lo
		}
		// sample(Array(lo..hi)) -> uniform in [lo, hi]
		return f.rng.IntInRange(lo, hi)
	}
	return 1
}

// sampleClass picks one character from a character class body, expanding a-z
// style ranges first (the gem expands "A-Z" to a single sampled char, then
// samples the resulting chars).
func (f *Faker) sampleClass(class string) byte {
	var chars []byte
	for j := 0; j < len(class); j++ {
		if j+2 < len(class) && class[j+1] == '-' {
			lo, hi := class[j], class[j+2]
			if lo <= hi {
				for c := lo; c <= hi; c++ {
					chars = append(chars, c)
				}
			}
			j += 2
		} else {
			chars = append(chars, class[j])
		}
	}
	if len(chars) == 0 {
		return '?'
	}
	idx := f.rng.Intn(len(chars))
	return chars[idx]
}

func parseRange(inner string) (lo, hi int, ok bool) {
	if comma := strings.IndexByte(inner, ','); comma >= 0 {
		a, ok1 := atoi(inner[:comma])
		b, ok2 := atoi(inner[comma+1:])
		return a, b, ok1 && ok2
	}
	a, ok1 := atoi(inner)
	return a, a, ok1
}

func atoi(s string) (int, bool) {
	n := 0
	if s == "" {
		return 0, false
	}
	for _, c := range s {
		if c < '0' || c > '9' {
			return 0, false
		}
		n = n*10 + int(c-'0')
	}
	return n, true
}

func matchParen(s string, open int) int {
	depth := 0
	for i := open; i < len(s); i++ {
		switch s[i] {
		case '(':
			depth++
		case ')':
			depth--
			if depth == 0 {
				return i
			}
		}
	}
	return -1
}

// letterAt maps 0..51 to A..Z then a..z, matching Faker's Letters constant
// (ULetters + LLetters).
func letterAt(n int) byte {
	if n < 26 {
		return byte('A' + n)
	}
	return byte('a' + n - 26)
}
