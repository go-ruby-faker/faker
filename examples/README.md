# Ruby examples

Pure-Ruby examples of the `faker` gem — the Ruby face of this library. They run
under [go-embedded-ruby](https://github.com/go-embedded-ruby/ruby) (rbgo) via
its `require "faker"` binding:

```sh
rbgo examples/faker_usage.rb
```

| File | Shows |
| --- | --- |
| [`faker_usage.rb`](faker_usage.rb) | Seeding for reproducible draws (`Faker::Config.random`), and generating names, emails, addresses, company data, commerce, colors, lorem text, and numbers. |
