package configenv

import (
	"fmt"
	"strconv"
	"strings"
)

func Required(p parser) parser {
	return func(raw string, ctx *context) error {
		if ctx.env[raw] == "" {
			return fmt.Errorf("required but not found")
		}
		return p(raw, ctx)
	}
}

func Int[T ~int](target *T) parser {
	return func(raw string, ctx *context) error {
		val, err := strconv.ParseInt(raw, 10, 0)
		*target = T(val)
		return err
	}
}

func Int8[T ~int8](target *T) parser {
	return func(raw string, ctx *context) error {
		val, err := strconv.ParseInt(raw, 10, 8)
		*target = T(val)
		return err
	}
}

func Int16[T ~int16](target *T) parser {
	return func(raw string, ctx *context) error {
		val, err := strconv.ParseInt(raw, 10, 16)
		*target = T(val)
		return err
	}
}

func Int32[T ~int32](target *T) parser {
	return func(raw string, ctx *context) error {
		val, err := strconv.ParseInt(raw, 10, 32)
		*target = T(val)
		return err
	}
}

func Int64[T ~int64](target *T) parser {
	return func(raw string, ctx *context) error {
		val, err := strconv.ParseInt(raw, 10, 64)
		*target = T(val)
		return err
	}
}

func Bool[T ~bool](target *T) parser {
	return func(raw string, ctx *context) error {
		switch raw {
		case "true", "yes":
			*target = true
		case "false", "", "no":
			*target = false
		default:
			return fmt.Errorf("must be true or false but got '%s'", raw)
		}
		return nil
	}
}

func String[T ~string](target *T) parser {
	return func(raw string, ctx *context) error {
		*target = T(raw)
		return nil
	}
}

func Strings[T ~[]string](target *T, sep string) parser {
	return func(raw string, ctx *context) error {
		parts := strings.Split(raw, sep)
		*target = T(parts)
		return nil
	}
}
