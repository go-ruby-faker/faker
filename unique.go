// Copyright (c) 2026, the go-ruby-faker/faker authors. All rights reserved.
// Use of this source code is governed by a BSD-3-Clause license that can be
// found in the LICENSE file.

package faker

import "fmt"

// RetriesExhaustedError is returned (as the error value) when Unique cannot find
// a fresh value within the retry budget, mirroring the gem's
// Faker::UniqueGenerator::RetryLimitExceeded.
type RetriesExhaustedError struct {
	Retries int
}

func (e *RetriesExhaustedError) Error() string {
	return fmt.Sprintf("faker: retry limit %d exceeded for unique generator", e.Retries)
}

// Unique wraps a generator function so each returned value is distinct, mirroring
// the gem's `Faker::Name.unique.first_name` modifier. It calls gen until it
// produces a value not seen before, giving up after maxRetries attempts.
//
// Pass maxRetries <= 0 to use the gem's default of 10_000.
type Unique struct {
	maxRetries int
	seen       map[string]struct{}
}

// NewUnique returns a Unique tracker. maxRetries <= 0 uses the gem default 10000.
func NewUnique(maxRetries int) *Unique {
	if maxRetries <= 0 {
		maxRetries = 10000
	}
	return &Unique{maxRetries: maxRetries, seen: make(map[string]struct{})}
}

// Next runs gen until it yields a value not previously returned by this tracker,
// or returns a *RetriesExhaustedError after maxRetries attempts.
func (u *Unique) Next(gen func() string) (string, error) {
	for tries := 0; tries < u.maxRetries; tries++ {
		v := gen()
		if _, dup := u.seen[v]; !dup {
			u.seen[v] = struct{}{}
			return v, nil
		}
	}
	return "", &RetriesExhaustedError{Retries: u.maxRetries}
}

// Clear forgets all previously seen values (Faker::UniqueGenerator#clear).
func (u *Unique) Clear() { u.seen = make(map[string]struct{}) }
