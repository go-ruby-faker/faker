<p align="center"><img src="https://raw.githubusercontent.com/go-ruby-faker/brand/main/social/go-ruby-faker-faker.png" alt="go-ruby-faker/faker" width="720"></p>

# faker — go-ruby-faker

[![Docs](https://img.shields.io/badge/docs-mkdocs--material-DC2626)](https://go-ruby-faker.github.io/docs/)
[![License](https://img.shields.io/badge/license-BSD--3--Clause-blue)](LICENSE)
[![Go](https://img.shields.io/badge/go-1.26.4%2B-00ADD8)](https://go.dev/dl/)
[![Coverage](https://img.shields.io/badge/coverage-100%25-1a7f37)](#tests--coverage)

**A pure-Go (no cgo) reimplementation of Ruby's [`faker`](https://github.com/faker-ruby/faker)
gem (3.8.0)** — a seeded, MRI-faithful fake-data generator. Given the same seed,
it reproduces the gem's draw sequence **bit-for-bit**, because its random core is
an exact port of Ruby's `Random` (MT19937 with `init_by_array` seeding and the
mask-and-reject integer scheme) and of `Array#sample` / `Array#shuffle`.

It ships the gem's `en` i18n data tables **embedded**, so it is complete offline —
no network, no Ruby runtime, **CGO=0** on all six 64-bit Go targets (amd64,
arm64, riscv64, loong64, ppc64le, s390x). It is the fake-data generator for
[go-embedded-ruby](https://github.com/go-embedded-ruby/ruby), but is a
**standalone, reusable** module — a sibling of
[go-ruby-regexp](https://github.com/go-ruby-regexp/regexp) and
[go-ruby-erb](https://github.com/go-ruby-erb/erb).

## The deterministic-seed contract

The gem's key correctness property is that

```ruby
Faker::Config.random = Random.new(SEED)
```

makes the sequence of generated values reproducible. This package mirrors it:

```go
f := faker.NewSeeded(42)             // Faker::Config.random = Random.new(42)
f.Name().Name()                      // Faker::Name.name
f.Internet().Email()                 // Faker::Internet.email
```

`faker.NewRandom(SEED)` is a bit-exact port of `Random.new(SEED)`: its raw
`rand`, `rand(n)`, `rand(a..b)`, `sample(random:)` and `shuffle(random:)` outputs
match MRI. The suite pins ~2000 values across seven seeds to live gem output.

## Generators

`Name` · `Internet` · `Address` · `PhoneNumber` · `Lorem` · `Number` · `Company`
· `Commerce` · `Color` · `Boolean` · `Date` · `Time`, plus the `Faker::Base`
helpers `Sample` / `SampleN` / `Shuffle` / `Fetch` / `Numerify` / `Letterify` /
`Bothify` / `Regexify` / `Parse`, the `Faker::Config` locale + RNG accessors, and
the `unique` modifier (`NewUnique` — retry + exhaustion error).

| Class | Methods |
| --- | --- |
| `Name` | name, name_with_middle, first_name, male/female/neutral first_name, last_name, prefix, suffix, initials |
| `Internet` | username, email, domain_word, domain_name, domain_suffix, slug, url, ip_v4/ip_v6/mac_address, password |
| `Address` | city, street_name/address, secondary_address, building_number, community, zip_code (zip/postcode), state, state_abbr, country, country_code(_long), time_zone, full_address, latitude, longitude |
| `PhoneNumber` | phone_number, cell_phone, country_code, area_code, exchange_code, subscriber_number/extension, with_country_code variants |
| `Lorem` | word, words, characters, sentence, sentences, paragraph, paragraphs |
| `Number` | number, decimal, digit, non_zero_digit, hexadecimal, binary, between, positive, negative, normal |
| `Company` | name, suffix, industry, buzzword, catch_phrase, bs |
| `Commerce` | product_name, department, price |
| `Color` | color_name, hex_color, rgb_color |
| `Boolean` | boolean(true_ratio) |
| `Date` / `Time` | between, forward, backward |

## Seed-parity scope (honest)

**Exact-sequence parity** with the gem holds for every generator whose randomness
flows through INTEGER draws and `sample`/`shuffle` selections: `Name`, `Address`
(city/street/zip/state/country/full_address), `PhoneNumber`, `Lorem`, `Company`,
`Commerce.product_name`/`department`, `Color.color_name`, `Boolean`, `Number`'s
integer paths, `Internet`'s table-driven paths (username/email/domain/url/slug),
and `Date.between`/`forward`/`backward` (whole-day integer arithmetic).

**Distribution / data-table / format parity only** (not bit-exact sequence)
applies where the value derives from a FLOAT-range draw, because MRI's
`Random#rand(lo..hi)` for float ranges uses a higher-precision internal algorithm
this port does not reproduce: `Address.latitude`/`longitude`,
`Number.between`/`positive`/`negative`/`normal` with non-integer bounds,
`Commerce.price`, `Color.hex_color`, and `Time.between`/`forward`/`backward`.
These match the gem's data tables, formats and statistical distribution; their
exact per-seed value may differ in low-order digits.

## Tests & coverage

```sh
GOWORK=off go test -race -cover ./...
```

The suite is **100% statement coverage**. The deterministic, Ruby-free tests
cover every branch on their own (so the Windows and qemu cross-arch CI lanes hold
the gate); where the `ruby` binary and the `faker` gem are present, the oracle
vectors in `oracle_test.go` pin the seeded output to live gem 3.8.0 values.

## License

BSD-3-Clause — see [LICENSE](LICENSE). Copyright (c) 2026, the
go-ruby-faker/faker authors.
