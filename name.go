// Copyright (c) 2026, the go-ruby-faker/faker authors. All rights reserved.
// Use of this source code is governed by a BSD-3-Clause license that can be
// found in the LICENSE file.

package faker

// NameGen is Faker::Name.
type NameGen struct{ f *Faker }

// Name returns the Faker::Name generator.
func (f *Faker) Name() *NameGen { return &NameGen{f} }

// Name produces a random full name (Faker::Name.name).
func (n *NameGen) Name() string { return n.f.Parse("name.name") }

// NameWithMiddle produces a name with a middle name.
func (n *NameGen) NameWithMiddle() string { return n.f.Parse("name.name_with_middle") }

// FirstName produces a random first name.
//
// This mirrors a subtle quirk of the gem: Faker::Name.first_name evaluates
// parse('name.first_name') TWICE — once for the `.empty?` guard and again to
// produce the return value — so it consumes two parse cycles' worth of draws.
// Reproducing that double evaluation is required for seed-exact parity.
func (n *NameGen) FirstName() string {
	if n.f.Parse("name.first_name") == "" {
		return n.f.Fetch("name.first_name")
	}
	return n.f.Parse("name.first_name")
}

// MaleFirstName produces a random male first name.
func (n *NameGen) MaleFirstName() string { return n.f.Fetch("name.male_first_name") }

// FemaleFirstName produces a random female first name.
func (n *NameGen) FemaleFirstName() string { return n.f.Fetch("name.female_first_name") }

// NeutralFirstName produces a random gender-neutral first name.
func (n *NameGen) NeutralFirstName() string { return n.f.Fetch("name.neutral_first_name") }

// LastName produces a random last name.
func (n *NameGen) LastName() string { return n.f.Parse("name.last_name") }

// Prefix produces a random name prefix (e.g. "Mr.").
func (n *NameGen) Prefix() string { return n.f.Fetch("name.prefix") }

// Suffix produces a random name suffix (e.g. "IV").
func (n *NameGen) Suffix() string { return n.f.Fetch("name.suffix") }

// Initials produces random uppercase initials (Faker::Name.initials).
func (n *NameGen) Initials(number int) string {
	b := make([]byte, number)
	for i := 0; i < number; i++ {
		b[i] = byte(n.f.rng.IntInRange(65, 90)) // 'A'..'Z'
	}
	return string(b)
}
