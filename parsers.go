package configenv

import (
	"fmt"
	"strconv"
	"strings"
)

func Required(p parser) parser {
	return func(raw string, ctx *context) error {
		if raw == "" {
			return fmt.Errorf("required but not found")
		}
		return p(raw, ctx)
	}
}

// Apply the function to the env var value before parsing it.
func Map(p parser, f func(string) string) parser {
	return func(raw string, ctx *context) error {
		raw = f(raw)
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

func Uint[T ~uint](target *T) parser {
	return func(raw string, ctx *context) error {
		val, err := strconv.ParseUint(raw, 10, 0)
		*target = T(val)
		return err
	}
}

func Uint8[T ~uint8](target *T) parser {
	return func(raw string, ctx *context) error {
		val, err := strconv.ParseUint(raw, 10, 8)
		*target = T(val)
		return err
	}
}

func Uint16[T ~uint16](target *T) parser {
	return func(raw string, ctx *context) error {
		val, err := strconv.ParseUint(raw, 10, 16)
		*target = T(val)
		return err
	}
}

func Uint32[T ~uint32](target *T) parser {
	return func(raw string, ctx *context) error {
		val, err := strconv.ParseUint(raw, 10, 32)
		*target = T(val)
		return err
	}
}

func Uint64[T ~uint64](target *T) parser {
	return func(raw string, ctx *context) error {
		val, err := strconv.ParseUint(raw, 10, 64)
		*target = T(val)
		return err
	}
}

func Float32[T ~uint64](target *T) parser {
	return func(raw string, ctx *context) error {
		val, err := strconv.ParseFloat(raw, 32)
		*target = T(val)
		return err
	}
}

func Float64[T ~uint64](target *T) parser {
	return func(raw string, ctx *context) error {
		val, err := strconv.ParseFloat(raw, 64)
		*target = T(val)
		return err
	}
}

func Bool[T ~bool](target *T) parser {
	return func(raw string, ctx *context) error {
		val, err := strconv.ParseBool(raw)
		*target = T(val)
		return err
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
