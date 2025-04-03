package configenv

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

// One-letter aliases for people living on the edge.
type (
	C = Config
	V = Vars
)

type Vars map[string]parser

type parser func(string, *context) error

type context struct {
	env    map[string]string
	prefix string
}

func (ctx *context) pop(name string) string {
	val := ctx.env[ctx.prefix+name]
	delete(ctx.env, name)
	return val
}

type Config struct {
	// Pairs of name=value env vars. If not provided, [os.Environ] will be used.
	Environ []string
	// The prefix to add to all env var names.
	//
	// The prefix ensures that env vars between services don't conflict
	Prefix string
	// If true, will not return an error for unknown env vars with the given prefix.
	//
	// By default, extra env vars are forbidden which prevents typos in env var names.
	// For example, if you write `BE_DEUG=false` instead of `BE_DEBUG=false`.
	AllowExtra bool
	// Require all env vars to be present and non-empty, even if not wrapped in [Required].
	//
	// It's a good idea to set it on the production to ensure no env var is forgotten.
	RequireAll bool
}

func (vars Vars) Parse(cfg Config) error {
	if cfg.Environ == nil {
		cfg.Environ = os.Environ()
	}
	if cfg.Prefix == "" {
		cfg.AllowExtra = true
	}

	env := map[string]string{}
	for _, pair := range cfg.Environ {
		key, value, found := strings.Cut(pair, "=")
		if !found {
			return errors.New("got an env var pair without `=`")
		}
		key, hasPrefix := strings.CutPrefix(key, cfg.Prefix)
		if !hasPrefix {
			continue
		}
		env[key] = value
	}
	ctx := context{env: env}
	for name, parse := range vars {
		raw := ctx.pop(name)
		if raw == "" {
			if cfg.RequireAll {
				return fmt.Errorf(
					"parse %s%s: required but not found",
					ctx.prefix, name,
				)
			}
			continue
		}
		err := parse(raw, &ctx)
		if err != nil {
			return fmt.Errorf("parse %s%s: %v", ctx.prefix, name, err)
		}
	}
	if !cfg.AllowExtra {
		for key := range ctx.env {
			return fmt.Errorf("unsupported env var: %s", key)
		}
	}
	return nil
}
