# frozen_string_literal: true
#
# Basic usage of the `faker` gem — a seeded, MRI-faithful fake-data generator.
#
# Runs under go-embedded-ruby (rbgo); see examples/README.md.

require "faker"

# Seeding the generator makes every draw reproducible: given the same seed, the
# sequence of generated values is bit-for-bit identical to the Ruby gem's.
Faker::Config.random = Random.new(42)

# People and contact details.
puts Faker::Name.name              # => "Brittany Klocko"
puts Faker::Internet.email         # => "lyndon@bayer-leuschke.test"
puts Faker::Internet.ip_v4_address # => "50.107.54.243"
puts Faker::PhoneNumber.phone_number

# Places, businesses, and commerce.
puts Faker::Address.full_address
puts Faker::Company.name
puts Faker::Company.catch_phrase
puts Faker::Commerce.product_name
puts Faker::Color.hex_color        # => "#d6efa3"

# Filler text and numbers (methods that take arguments accept counts/sizes).
puts Faker::Lorem.sentence         # => "Et officiis possimus dolores."
p    Faker::Lorem.words(3)         # => ["hic", "neque", "harum"]
puts Faker::Number.number(5)       # => 18687
p    Faker::Boolean.boolean        # => false

# Re-seeding with the same value replays the exact same sequence.
Faker::Config.random = Random.new(42)
puts Faker::Name.name              # => "Brittany Klocko" (identical to above)
