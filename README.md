# configenv

Go package for parsing env vars.

The main idea is to enforce (on the package API design level) the best practice of parsing env vars: prevent duplicate names, wrong types, typos, etc. See [Writing safe-to-use Go libraries](https://blog.orsinium.dev/posts/go/safe-api/) on the design principles behind my Go packages, including this one.

Also, try out [cliff](https://github.com/orsinium-labs/cliff), a package for safely parsing CLI args.

## Installation

```bash
go get github.com/orsinium-labs/configenv
```

## Usage

Given the following parser:

```go
type Config struct {
    Debug bool
    Env   string
}

config := Config{}
vars := configenv.Vars{
    "DEBUG": configenv.Required(configenv.Bool(&config.Debug)),
    "ENV":   configenv.String(&config.Env),
}
err := vars.Parse(configenv.Config{Prefix: "BE_"})
```

And the following env vars:

```bash
export BE_DEBUG=true
export BE_ENV=prod
```

You will get `Config{Debug: true, Env: "prod"}`

If you instead of `BE_DEBUG=true` provide `DEBUG=true`, `BE_DEUG=true`, or `BE_DEBUG=maybe`, the parser will detect it and return an error.
