// Copyright (c) 2026, the go-ruby-faker/faker authors. All rights reserved.
// Use of this source code is governed by a BSD-3-Clause license that can be
// found in the LICENSE file.

package faker

import (
	crand "crypto/rand"
	"math"
)

// Random is Ruby's Random: an MT19937 generator seeded exactly as MRI does
// (init_by_array over the seed's 32-bit little-endian words), so seeded output
// matches MRI bit for bit. This is the keystone of Faker's deterministic-seed
// contract: NewRandom(seed) reproduces the same draw sequence MRI's
// Random.new(seed) produces.
type Random struct {
	mt   [624]uint32
	mti  int
	seed int64
}

// NewRandom returns a Random seeded exactly as MRI's Random.new(seed).
func NewRandom(seed int64) *Random {
	r := &Random{seed: seed}
	key := seedKey(seed)
	if len(key) == 1 {
		r.initGenrand(key[0]) // MRI seeds a single-word seed with init_genrand
	} else {
		r.initByArray(key)
	}
	return r
}

// NewRandomEntropy draws a non-deterministic seed (matching MRI's
// entropy-seeded Random.new with no argument) and returns a seeded Random.
func NewRandomEntropy() *Random {
	return NewRandom(randomSeed())
}

// Seed reports the seed the generator was created with.
func (r *Random) Seed() int64 { return r.seed }

// randomSeed draws a non-deterministic seed for Random.new with no argument,
// matching MRI's entropy-seeded default.
func randomSeed() int64 {
	var b [8]byte
	_, _ = crand.Read(b[:]) // on the rare error path b stays zero — a valid seed
	var s uint64
	for i := 0; i < 8; i++ {
		s |= uint64(b[i]) << (8 * uint(i))
	}
	return int64(s & 0x7fffffffffffffff)
}

// seedKey turns a seed into the 32-bit little-endian word array MRI feeds to
// init_by_array (|seed|; zero yields a single zero word).
func seedKey(seed int64) []uint32 {
	s := uint64(seed)
	if seed < 0 {
		s = uint64(-seed)
	}
	if s == 0 {
		return []uint32{0}
	}
	var key []uint32
	for s > 0 {
		key = append(key, uint32(s))
		s >>= 32
	}
	return key
}

func (r *Random) initGenrand(s uint32) {
	r.mt[0] = s
	for i := 1; i < 624; i++ {
		r.mt[i] = 1812433253*(r.mt[i-1]^(r.mt[i-1]>>30)) + uint32(i)
	}
	r.mti = 624
}

func (r *Random) initByArray(key []uint32) {
	r.initGenrand(19650218)
	i, j := 1, 0
	k := 624 // max(624, len(key)); seeds are int64, so len(key) <= 2 and 624 wins.
	for ; k > 0; k-- {
		r.mt[i] = (r.mt[i] ^ ((r.mt[i-1] ^ (r.mt[i-1] >> 30)) * 1664525)) + key[j] + uint32(j)
		i++
		j++
		if i >= 624 {
			r.mt[0] = r.mt[623]
			i = 1
		}
		if j >= len(key) {
			j = 0
		}
	}
	for k = 623; k > 0; k-- {
		r.mt[i] = (r.mt[i] ^ ((r.mt[i-1] ^ (r.mt[i-1] >> 30)) * 1566083941)) - uint32(i)
		i++
		if i >= 624 {
			r.mt[0] = r.mt[623]
			i = 1
		}
	}
	r.mt[0] = 0x80000000
}

func (r *Random) genrandInt32() uint32 {
	if r.mti >= 624 {
		for i := 0; i < 624; i++ {
			y := (r.mt[i] & 0x80000000) | (r.mt[(i+1)%624] & 0x7fffffff)
			r.mt[i] = r.mt[(i+397)%624] ^ (y >> 1)
			if y&1 != 0 {
				r.mt[i] ^= 0x9908b0df
			}
		}
		r.mti = 0
	}
	y := r.mt[r.mti]
	r.mti++
	y ^= y >> 11
	y ^= (y << 7) & 0x9d2c5680
	y ^= (y << 15) & 0xefc60000
	y ^= y >> 18
	return y
}

// res53 is MRI's genrand_real: a 53-bit float in [0, 1).
func (r *Random) res53() float64 {
	a := r.genrandInt32() >> 5
	b := r.genrandInt32() >> 6
	return (float64(a)*67108864.0 + float64(b)) / 9007199254740992.0
}

// limitedRand returns a uniform integer in [0, limit] (inclusive) using MRI's
// mask-and-reject scheme: the value is assembled high 32-bit word first, each
// word masked by the corresponding slice of the bit mask, retrying while the
// result exceeds limit. A 32-bit limit consumes one genrand_int32 per attempt.
func (r *Random) limitedRand(limit uint64) uint64 {
	if limit == 0 {
		return 0
	}
	mask := makeMask64(limit)
	for {
		val := uint64(0)
		for i := 1; i >= 0; i-- {
			if m := (mask >> (uint(i) * 32)) & 0xffffffff; m != 0 {
				val |= (uint64(r.genrandInt32()) & m) << (uint(i) * 32)
			}
		}
		if val <= limit {
			return val
		}
	}
}

func makeMask64(x uint64) uint64 {
	x |= x >> 1
	x |= x >> 2
	x |= x >> 4
	x |= x >> 8
	x |= x >> 16
	x |= x >> 32
	return x
}

// Float mirrors Random#rand with no argument: a float in [0, 1).
func (r *Random) Float() float64 { return r.res53() }

// Intn mirrors Random#rand(n) for a positive integer n: an integer in [0, n).
// It panics for n <= 0, matching MRI's ArgumentError.
func (r *Random) Intn(n int) int {
	if n <= 0 {
		panic("faker: Random.Intn requires n > 0")
	}
	return int(r.limitedRand(uint64(n) - 1))
}

// IntInRange mirrors Random#rand(lo..hi) (inclusive) over integers.
func (r *Random) IntInRange(lo, hi int) int {
	span := hi - lo
	if span < 0 {
		lo, hi = hi, lo
		span = -span
	}
	return lo + int(r.limitedRand(uint64(span)))
}

// FloatScaled mirrors Random#rand(f) for a positive float f: res53 * f.
func (r *Random) FloatScaled(f float64) float64 {
	if f <= 0 {
		return r.res53()
	}
	return r.res53() * f
}

// FloatInRange returns lo + res53*(hi-lo). NOTE: this is an approximation of
// MRI's Random#rand(lo..hi) for FLOAT ranges; MRI uses a higher-precision
// internal float-range algorithm that this does not bit-match. Integer-range
// draws (IntInRange) ARE bit-exact; see the package doc for the seed-parity
// scope.
func (r *Random) FloatInRange(lo, hi float64) float64 {
	if hi < lo {
		lo, hi = hi, lo
	}
	return lo + r.res53()*(hi-lo)
}

// Bytes returns n pseudo-random bytes, matching Random#bytes(n) (little-endian
// words from genrand_int32).
func (r *Random) Bytes(n int) []byte {
	b := make([]byte, n)
	for i := 0; i < n; i += 4 {
		w := r.genrandInt32()
		for k := 0; k < 4 && i+k < n; k++ {
			b[i+k] = byte(w >> (8 * uint(k)))
		}
	}
	return b
}

// ---- Ruby Array helpers reproduced over the same generator ----

// sampleOne mirrors Ruby's Array#sample(random:) for a single element: it draws
// one index via limitedRand(len-1), i.e. arr[rand(len)]. Returns ok=false for
// an empty slice.
func sampleIndex(r *Random, length int) (int, bool) {
	if length <= 0 {
		return 0, false
	}
	return int(r.limitedRand(uint64(length) - 1)), true
}

// sampleN mirrors Ruby's Array#sample(n, random:) for n <= 6 against a slice of
// the given length, returning the chosen indices in MRI's draw order. MRI's
// small-n path picks i=rand(len), then each subsequent index from the shrinking
// range and bumps it past any already-chosen index it meets or exceeds (the
// "insert into sorted, increment when >=" scheme). For n > 6 this falls back to
// a partial Fisher-Yates shuffle: correct and uniform, but NOT bit-matching
// MRI's hash-based large-n path (Faker's targeted generators use n <= 6).
func sampleIndices(r *Random, length, n int) []int {
	if length <= 0 || n <= 0 {
		return nil
	}
	if n >= length {
		// Full coverage: MRI returns a shuffle of all indices.
		return shuffleIndices(r, length)
	}
	if n <= 6 {
		idx := make([]int, 0, n)
		for picked := 0; picked < n; picked++ {
			g := int(r.limitedRand(uint64(length-picked) - 1))
			// Bump g past every already-chosen index it reaches, scanning the
			// chosen set in ascending order (MRI sorts then increments).
			sorted := append([]int(nil), idx...)
			insertionSort(sorted)
			for _, v := range sorted {
				if g >= v {
					g++
				}
			}
			idx = append(idx, g)
		}
		return idx
	}
	// n in (6, length): partial Fisher-Yates (distribution parity only).
	perm := make([]int, length)
	for i := range perm {
		perm[i] = i
	}
	for i := 0; i < n; i++ {
		j := i + int(r.limitedRand(uint64(length-i)-1))
		perm[i], perm[j] = perm[j], perm[i]
	}
	return perm[:n]
}

// shuffleIndices mirrors Ruby's Array#shuffle(random:): back-to-front
// Fisher-Yates with j = rand(i+1) for i from len-1 down to 1.
func shuffleIndices(r *Random, length int) []int {
	perm := make([]int, length)
	for i := range perm {
		perm[i] = i
	}
	for i := length - 1; i > 0; i-- {
		j := int(r.limitedRand(uint64(i)))
		perm[i], perm[j] = perm[j], perm[i]
	}
	return perm
}

func insertionSort(a []int) {
	for i := 1; i < len(a); i++ {
		v := a[i]
		j := i - 1
		for j >= 0 && a[j] > v {
			a[j+1] = a[j]
			j--
		}
		a[j+1] = v
	}
}

// roundN rounds x to n decimal places, matching Ruby's Float#round(n)
// (round-half-up away from zero).
func roundN(x float64, n int) float64 {
	p := math.Pow(10, float64(n))
	if x < 0 {
		return -math.Floor(-x*p+0.5) / p
	}
	return math.Floor(x*p+0.5) / p
}
