// Copyright (c) 2026, the go-ruby-faker/faker authors. All rights reserved.
// Use of this source code is governed by a BSD-3-Clause license that can be
// found in the LICENSE file.

package faker

import (
	"math"
	"strconv"
	"strings"
)

// NumberGen is Faker::Number.
type NumberGen struct{ f *Faker }

// Number returns the Faker::Number generator.
func (f *Faker) Number() *NumberGen { return &NumberGen{f} }

// Number produces a random integer with the given number of digits. The first
// digit is non-zero (1..9). Returns 0 for digits < 1 (the gem returns nil).
func (n *NumberGen) Number(digits int) int64 {
	if digits < 1 {
		return 0
	}
	if digits == 1 {
		return int64(n.f.rng.IntInRange(0, 9))
	}
	var b strings.Builder
	b.WriteString(strconv.Itoa(n.nonZeroDigit()))
	for i := 0; i < digits-1; i++ {
		b.WriteString(strconv.Itoa(n.Digit()))
	}
	v, _ := strconv.ParseInt(b.String(), 10, 64)
	return v
}

// Decimal produces a float with lDigits integer digits and rDigits fractional
// digits (Faker::Number.decimal). The last fractional digit is non-zero.
func (n *NumberGen) Decimal(lDigits, rDigits int) float64 {
	ld := n.Number(lDigits)
	var b strings.Builder
	for i := 0; i < rDigits-1; i++ {
		b.WriteString(strconv.Itoa(n.Digit()))
	}
	b.WriteString(strconv.Itoa(n.nonZeroDigit()))
	v, _ := strconv.ParseFloat(strconv.FormatInt(ld, 10)+"."+b.String(), 64)
	return v
}

// Digit produces a single-digit integer 0..9 (Faker::Number.digit).
func (n *NumberGen) Digit() int { return n.f.rng.Intn(10) }

func (n *NumberGen) nonZeroDigit() int { return n.f.rng.IntInRange(1, 9) }

// NonZeroDigit produces a non-zero single-digit integer 1..9.
func (n *NumberGen) NonZeroDigit() int { return n.nonZeroDigit() }

// Hexadecimal produces a lowercase hex string of the given length
// (Faker::Number.hexadecimal).
func (n *NumberGen) Hexadecimal(digits int) string {
	const hexDigits = "0123456789abcdef"
	b := make([]byte, digits)
	for i := 0; i < digits; i++ {
		b[i] = hexDigits[n.f.rng.Intn(16)]
	}
	return string(b)
}

// Binary produces a binary string of the given length (Faker::Number.binary).
func (n *NumberGen) Binary(digits int) string {
	b := make([]byte, digits)
	for i := 0; i < digits; i++ {
		b[i] = byte('0' + n.f.rng.Intn(2))
	}
	return string(b)
}

// Between returns a float in [from, to] (inclusive). NOTE: float-range
// distribution parity only, not bit-exact sequence (see package doc).
func (n *NumberGen) Between(from, to float64) float64 {
	return n.f.rng.FloatInRange(from, to)
}

// BetweenInt returns an integer in [from, to] (inclusive) — bit-exact with the
// gem when both bounds are integers.
func (n *NumberGen) BetweenInt(from, to int) int { return n.f.rng.IntInRange(from, to) }

// Positive returns a positive float in the range (Faker::Number.positive).
func (n *NumberGen) Positive(from, to float64) float64 {
	v := n.Between(from, to)
	if v < 0 {
		return -v
	}
	return v
}

// Negative returns a negative float in the range (Faker::Number.negative).
func (n *NumberGen) Negative(from, to float64) float64 {
	v := n.Between(from, to)
	if v > 0 {
		return -v
	}
	return v
}

// Normal returns a Gaussian float given mean and standard deviation
// (Faker::Number.normal, Box–Muller). Distribution parity only.
func (n *NumberGen) Normal(mean, stddev float64) float64 {
	theta := 2 * math.Pi * n.f.rng.Float()
	rho := math.Sqrt(-2 * math.Log(1-n.f.rng.Float()))
	return mean + stddev*rho*math.Cos(theta)
}
