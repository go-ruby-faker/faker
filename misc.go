// Copyright (c) 2026, the go-ruby-faker/faker authors. All rights reserved.
// Use of this source code is governed by a BSD-3-Clause license that can be
// found in the LICENSE file.

package faker

import "strings"

// BooleanGen is Faker::Boolean.
type BooleanGen struct{ f *Faker }

// Boolean returns the Faker::Boolean generator.
func (f *Faker) Boolean() *BooleanGen { return &BooleanGen{f} }

// Boolean returns a bool that is true with probability trueRatio
// (Faker::Boolean.boolean). trueRatio 0.5 gives an even coin.
func (b *BooleanGen) Boolean(trueRatio float64) bool {
	return b.f.rng.Float() < trueRatio
}

// ColorGen is Faker::Color.
type ColorGen struct{ f *Faker }

// Color returns the Faker::Color generator.
func (f *Faker) Color() *ColorGen { return &ColorGen{f} }

// ColorName produces the name of a color (Faker::Color.color_name).
func (c *ColorGen) ColorName() string { return c.f.Fetch("color.name") }

// HexColor produces a random hex color "#rrggbb" via HSL conversion
// (Faker::Color.hex_color with no fixed components). Distribution parity:
// saturation/lightness derive from float draws.
func (c *ColorGen) HexColor() string {
	h := float64(c.f.rng.IntInRange(0, 360))
	s := roundN(c.f.rng.Float(), 2)
	l := roundN(c.f.rng.Float(), 2)
	return hslToHex(h, s, l)
}

// RGBColor produces an [r,g,b] triple (Faker::Color.rgb_color).
func (c *ColorGen) RGBColor() [3]int {
	return [3]int{c.singleRGB(), c.singleRGB(), c.singleRGB()}
}

func (c *ColorGen) singleRGB() int { return c.f.rng.IntInRange(0, 255) }

func hslToHex(h, s, l float64) string {
	chroma := (1 - absF(2*l-1)) * s
	hPrime := h / 60
	x := chroma * (1 - absF(modF(hPrime, 2)-1))
	m := l - 0.5*chroma
	var r, g, b float64
	switch int(hPrime) {
	case 0:
		r, g, b = chroma, x, 0
	case 1:
		r, g, b = x, chroma, 0
	case 2:
		r, g, b = 0, chroma, x
	case 3:
		r, g, b = 0, x, chroma
	case 4:
		r, g, b = x, 0, chroma
	default:
		r, g, b = chroma, 0, x
	}
	ri := int(roundN((r+m)*255, 0))
	gi := int(roundN((g+m)*255, 0))
	bi := int(roundN((b+m)*255, 0))
	return "#" + hex2(ri) + hex2(gi) + hex2(bi)
}

func hex2(v int) string {
	const d = "0123456789abcdef"
	if v < 0 {
		v = 0
	}
	if v > 255 {
		v = 255
	}
	return string([]byte{d[v>>4], d[v&0xf]})
}

func absF(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

func modF(a, b float64) float64 {
	r := a - float64(int(a/b))*b
	if r < 0 {
		r += b
	}
	return r
}

// CompanyGen is Faker::Company.
type CompanyGen struct{ f *Faker }

// Company returns the Faker::Company generator.
func (f *Faker) Company() *CompanyGen { return &CompanyGen{f} }

// Name produces a company name (Faker::Company.name).
func (c *CompanyGen) Name() string { return c.f.Parse("company.name") }

// Suffix produces a company suffix (Faker::Company.suffix).
func (c *CompanyGen) Suffix() string { return c.f.Fetch("company.suffix") }

// Industry produces a company industry (Faker::Company.industry).
func (c *CompanyGen) Industry() string { return c.f.Fetch("company.industry") }

// Buzzword produces a single buzzword (Faker::Company.buzzword).
func (c *CompanyGen) Buzzword() string {
	return c.f.Sample(c.f.data.flatten("company.buzzwords"))
}

// CatchPhrase produces a company catch phrase by sampling one word from each
// buzzword column (Faker::Company.catch_phrase).
func (c *CompanyGen) CatchPhrase() string {
	return c.joinColumns("company.buzzwords")
}

// BS produces some company BS, one word from each bs column (Faker::Company.bs).
func (c *CompanyGen) BS() string {
	return c.joinColumns("company.bs")
}

func (c *CompanyGen) joinColumns(key string) string {
	cols := c.f.data.lists(key)
	parts := make([]string, 0, len(cols))
	for _, col := range cols {
		parts = append(parts, c.f.Sample(col))
	}
	return strings.Join(parts, " ")
}

// CommerceGen is Faker::Commerce.
type CommerceGen struct{ f *Faker }

// Commerce returns the Faker::Commerce generator.
func (f *Faker) Commerce() *CommerceGen { return &CommerceGen{f} }

// ProductName produces "Adjective Material Product" (Faker::Commerce.product_name).
func (c *CommerceGen) ProductName() string {
	return c.f.Fetch("commerce.product_name.adjective") + " " +
		c.f.Fetch("commerce.product_name.material") + " " +
		c.f.Fetch("commerce.product_name.product")
}

// Department produces a department name. With fixedAmount the number of
// categories is exactly max; otherwise it is 1 + rand(max)
// (Faker::Commerce.department). Multiple categories are joined "a, b & c".
func (c *CommerceGen) Department(max int, fixedAmount bool) string {
	num := max
	if !fixedAmount {
		num = 1 + c.f.rng.Intn(max)
	}
	cats := c.f.SampleN(c.f.data.strings("commerce.department"), num)
	switch {
	case len(cats) == 0:
		return ""
	case len(cats) == 1:
		return cats[0]
	default:
		sep := c.f.data.scalar("separator")
		comma := strings.Join(cats[:len(cats)-1], ", ")
		return comma + sep + cats[len(cats)-1]
	}
}

// Price produces a price in the half-open range [0, max), floored to cents
// (Faker::Commerce.price). Distribution parity only (float-range draw).
func (c *CommerceGen) Price(max float64) float64 {
	return float64(int64(c.f.rng.FloatScaled(max)*100)) / 100.0
}
