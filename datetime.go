// Copyright (c) 2026, the go-ruby-faker/faker authors. All rights reserved.
// Use of this source code is governed by a BSD-3-Clause license that can be
// found in the LICENSE file.

package faker

import "time"

// DateGen is Faker::Date. Date arithmetic uses whole-day (Julian-day) integer
// draws, so DateGen is bit-exact with the gem's seeded sequence.
type DateGen struct{ f *Faker }

// Date returns the Faker::Date generator.
func (f *Faker) Date() *DateGen { return &DateGen{f} }

// Between returns a date in [from, to] inclusive (Faker::Date.between). Both
// bounds are truncated to whole days; the draw is an exact integer day offset.
func (d *DateGen) Between(from, to time.Time) time.Time {
	from = truncDay(from)
	to = truncDay(to)
	lo := julianDay(from)
	hi := julianDay(to)
	if hi < lo {
		lo, hi = hi, lo
	}
	span := hi - lo
	off := int(d.f.rng.limitedRand(uint64(span)))
	return from.AddDate(0, 0, off)
}

// Forward returns a date 1..days into the future from base (Faker::Date.forward).
func (d *DateGen) Forward(base time.Time, days int) time.Time {
	base = truncDay(base)
	return d.Between(base.AddDate(0, 0, 1), base.AddDate(0, 0, days))
}

// Backward returns a date 1..days in the past from base (Faker::Date.backward).
func (d *DateGen) Backward(base time.Time, days int) time.Time {
	base = truncDay(base)
	return d.Between(base.AddDate(0, 0, -days), base.AddDate(0, 0, -1))
}

func truncDay(t time.Time) time.Time {
	y, m, dd := t.Date()
	return time.Date(y, m, dd, 0, 0, 0, 0, t.Location())
}

// julianDay returns the count of whole days since the Unix epoch for a
// day-truncated time. Only the difference between two julianDay values is used,
// so the absolute origin is irrelevant to correctness.
func julianDay(t time.Time) int {
	return int(truncDay(t).Unix() / 86400)
}

// TimeGen is Faker::Time. Time draws use float ranges, so TimeGen matches the
// gem's data/format/distribution but NOT its exact per-seed value (see package
// doc on seed-parity scope).
type TimeGen struct{ f *Faker }

// Time returns the Faker::Time generator.
func (f *Faker) Time() *TimeGen { return &TimeGen{f} }

// Between returns a time in [from, to] (Faker::Time.between). Distribution
// parity only.
func (tg *TimeGen) Between(from, to time.Time) time.Time {
	lo := float64(from.UnixNano()) / 1e9
	hi := float64(to.UnixNano()) / 1e9
	sec := tg.f.rng.FloatInRange(lo, hi)
	whole := int64(sec)
	nsec := int64((sec - float64(whole)) * 1e9)
	return time.Unix(whole, nsec).In(from.Location())
}

// Forward returns a time 1..days into the future (Faker::Time.forward).
func (tg *TimeGen) Forward(base time.Time, days int) time.Time {
	return tg.Between(base.Add(24*time.Hour), base.Add(time.Duration(days)*24*time.Hour))
}

// Backward returns a time 1..days in the past (Faker::Time.backward).
func (tg *TimeGen) Backward(base time.Time, days int) time.Time {
	return tg.Between(base.Add(-time.Duration(days)*24*time.Hour), base.Add(-24*time.Hour))
}
