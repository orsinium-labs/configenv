package configenv

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

type Vars map[string]parser

type parser func(string, *context) error

type context struct {
	env    map[string]string
	prefix string
}

func (ctx *context) pop(name string) string {
	val := ctx.env[name]
	delete(ctx.env, name)
	return val
}

type Config struct {
	Environ    []string
	AllowExtra bool
	RequireAll bool
}

func (vars Vars) Parse(cfg Config) error {
	if cfg.Environ == nil {
		cfg.Environ = os.Environ()
	}

	env := map[string]string{}
	for _, pair := range cfg.Environ {
		key, value, found := strings.Cut(pair, "=")
		if !found {
			return errors.New("got an env var pair without `=`")
		}
		key, hasPrefix := strings.CutPrefix(key, "OB_")
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
