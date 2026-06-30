// Copyright (c) 2026, the go-ruby-faker/faker authors. All rights reserved.
// Use of this source code is governed by a BSD-3-Clause license that can be
// found in the LICENSE file.

package faker

import "strings"

// LoremGen is Faker::Lorem.
type LoremGen struct{ f *Faker }

// Lorem returns the Faker::Lorem generator.
func (f *Faker) Lorem() *LoremGen { return &LoremGen{f} }

// alphaNums is Ruby's ALPHANUMS = LLetters + Numbers ('a'..'z' then '0'..'9').
var alphaNums = func() []string {
	out := make([]string, 0, 36)
	for c := byte('a'); c <= 'z'; c++ {
		out = append(out, string(c))
	}
	for c := byte('0'); c <= '9'; c++ {
		out = append(out, string(c))
	}
	return out
}()

// Word returns a single random lorem word (Faker::Lorem.word).
func (l *LoremGen) Word() string {
	w := l.Words(1, false)
	if len(w) == 0 {
		return ""
	}
	return w[0]
}

// Words returns number random lorem words. With supplemental true the
// supplemental word list is appended to the pool first. This mirrors the gem's
// word_list *= ((num/len)+1) expansion before sampling, so seeded output
// matches even when number exceeds the pool size.
func (l *LoremGen) Words(number int, supplemental bool) []string {
	pool := append([]string(nil), l.f.data.strings("lorem.words")...)
	if supplemental {
		pool = append(pool, l.f.data.strings("lorem.supplemental")...)
	}
	if len(pool) == 0 {
		return nil
	}
	mult := (number / len(pool)) + 1
	if mult > 1 {
		expanded := make([]string, 0, len(pool)*mult)
		for i := 0; i < mult; i++ {
			expanded = append(expanded, pool...)
		}
		pool = expanded
	}
	return l.f.SampleN(pool, number)
}

// Characters returns number random alphanumeric characters
// (Faker::Lorem.characters), drawing from a-z and 0-9.
func (l *LoremGen) Characters(number int) string {
	if number < 1 {
		return ""
	}
	var b strings.Builder
	for i := 0; i < number; i++ {
		b.WriteString(l.f.Sample(alphaNums))
	}
	return b.String()
}

// Sentence returns a capitalized sentence of wordCount words ended with the
// locale period.
func (l *LoremGen) Sentence(wordCount int, supplemental bool) string {
	words := l.Words(wordCount, supplemental)
	joined := strings.Join(words, l.localeSpace())
	return capitalize(joined) + l.localePeriod()
}

// Sentences returns number sentences (each of 3 words, matching the gem).
func (l *LoremGen) Sentences(number int, supplemental bool) []string {
	out := make([]string, 0, number)
	for i := 0; i < number; i++ {
		out = append(out, l.Sentence(3, supplemental))
	}
	return out
}

// Paragraph returns a paragraph of sentenceCount sentences joined by spaces.
func (l *LoremGen) Paragraph(sentenceCount int, supplemental bool) string {
	return strings.Join(l.Sentences(sentenceCount, supplemental), l.localeSpace())
}

// Paragraphs returns number paragraphs (each of 3 sentences).
func (l *LoremGen) Paragraphs(number int, supplemental bool) []string {
	out := make([]string, 0, number)
	for i := 0; i < number; i++ {
		out = append(out, l.Paragraph(3, supplemental))
	}
	return out
}

func (l *LoremGen) localeSpace() string {
	if m := l.f.data.object("lorem.punctuation"); m != nil {
		if v, ok := m["space"]; ok {
			return v
		}
	}
	return " "
}

func (l *LoremGen) localePeriod() string {
	if m := l.f.data.object("lorem.punctuation"); m != nil {
		if v, ok := m["period"]; ok {
			return v
		}
	}
	return "."
}

// capitalize mirrors Ruby String#capitalize for ASCII: upcase the first letter,
// downcase the rest.
func capitalize(s string) string {
	if s == "" {
		return ""
	}
	b := []byte(strings.ToLower(s))
	if b[0] >= 'a' && b[0] <= 'z' {
		b[0] -= 'a' - 'A'
	}
	return string(b)
}
