package configenv

import (
	"cmp"
	"encoding/json"
	"fmt"
	"slices"
	"strconv"
	"strings"
)

// One-letter aliases for people living on the edge.
var (
	R = Required
	S = String[string]
	I = Int[int]
	U = Uint[uint]
	B = Bool[bool]
	F = Float64[float64]
)

type Parser[T any] func(*T) parser

var (
	_ Parser[bool]   = Bool
	_ Parser[int]    = Int
	_ Parser[string] = String
)

func Required(p parser) parser {
	return func(raw string, ctx *context) error {
		if raw == "" {
			return nil
		}
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
		if raw == "" {
			return nil
		}
		val, err := strconv.ParseInt(raw, 10, 0)
		*target = T(val)
		return err
	}
}

func Int8[T ~int8](target *T) parser {
	return func(raw string, ctx *context) error {
		if raw == "" {
			return nil
		}
		val, err := strconv.ParseInt(raw, 10, 8)
		*target = T(val)
		return err
	}
}

func Int16[T ~int16](target *T) parser {
	return func(raw string, ctx *context) error {
		if raw == "" {
			return nil
		}
		val, err := strconv.ParseInt(raw, 10, 16)
		*target = T(val)
		return err
	}
}

func Int32[T ~int32](target *T) parser {
	return func(raw string, ctx *context) error {
		if raw == "" {
			return nil
		}
		val, err := strconv.ParseInt(raw, 10, 32)
		*target = T(val)
		return err
	}
}

func Int64[T ~int64](target *T) parser {
	return func(raw string, ctx *context) error {
		if raw == "" {
			return nil
		}
		val, err := strconv.ParseInt(raw, 10, 64)
		*target = T(val)
		return err
	}
}

func Uint[T ~uint](target *T) parser {
	return func(raw string, ctx *context) error {
		if raw == "" {
			return nil
		}
		val, err := strconv.ParseUint(raw, 10, 0)
		*target = T(val)
		return err
	}
}

func Uint8[T ~uint8](target *T) parser {
	return func(raw string, ctx *context) error {
		if raw == "" {
			return nil
		}
		val, err := strconv.ParseUint(raw, 10, 8)
		*target = T(val)
		return err
	}
}

func Uint16[T ~uint16](target *T) parser {
	return func(raw string, ctx *context) error {
		if raw == "" {
			return nil
		}
		val, err := strconv.ParseUint(raw, 10, 16)
		*target = T(val)
		return err
	}
}

func Uint32[T ~uint32](target *T) parser {
	return func(raw string, ctx *context) error {
		if raw == "" {
			return nil
		}
		val, err := strconv.ParseUint(raw, 10, 32)
		*target = T(val)
		return err
	}
}

func Uint64[T ~uint64](target *T) parser {
	return func(raw string, ctx *context) error {
		if raw == "" {
			return nil
		}
		val, err := strconv.ParseUint(raw, 10, 64)
		*target = T(val)
		return err
	}
}

func Float32[T ~float32](target *T) parser {
	return func(raw string, ctx *context) error {
		if raw == "" {
			return nil
		}
		val, err := strconv.ParseFloat(raw, 32)
		*target = T(val)
		return err
	}
}

func Float64[T ~float64](target *T) parser {
	return func(raw string, ctx *context) error {
		if raw == "" {
			return nil
		}
		val, err := strconv.ParseFloat(raw, 64)
		*target = T(val)
		return err
	}
}

func Bool[T ~bool](target *T) parser {
	return func(raw string, ctx *context) error {
		if raw == "" {
			return nil
		}
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

func Strings[A ~[]V, V ~string](target *A, sep string) parser {
	return func(raw string, ctx *context) error {
		if raw == "" {
			return nil
		}
		for part := range strings.SplitSeq(raw, sep) {
			*target = append(*target, V(part))
		}
		return nil
	}
}

func Slice[A ~[]V, V any](target *A, sep string, p Parser[V]) parser {
	return func(raw string, ctx *context) error {
		if raw == "" {
			return nil
		}
		for part := range strings.SplitSeq(raw, sep) {
			var parsed V
			err := p(&parsed)(part, ctx)
			if err != nil {
				return err
			}
			*target = append(*target, parsed)
		}
		return nil
	}
}

func JSON[T any](target *T) parser {
	return func(raw string, ctx *context) error {
		if raw == "" {
			return nil
		}
		return json.Unmarshal([]byte(raw), target)
	}
}

func PrefixMap[M ~map[K]V, K ~string, V any](target *M, p Parser[V]) parser {
	return func(raw string, ctx *context) error {
		prefix := ctx.name
		res := make(M, 0)
		for name, val := range ctx.env {
			suffix, found := strings.CutPrefix(name, prefix)
			if !found {
				continue
			}
			ctx.name = name
			var parsed V
			err := p(&parsed)(val, ctx)
			if err != nil {
				return err
			}
			res[K(suffix)] = parsed
		}
		for key := range res {
			delete(ctx.env, prefix+string(key))
		}
		*target = res
		return nil
	}
}

func PrefixSlice[A ~[]V, V any](target *A, p Parser[V]) parser {
	return func(raw string, ctx *context) error {
		prefix := ctx.name
		type pair struct {
			i int64
			k string
			v V
		}
		pairs := []pair{}
		for name, val := range ctx.env {
			rawSuffix, found := strings.CutPrefix(name, prefix)
			if !found {
				continue
			}
			index, err := strconv.ParseInt(rawSuffix, 10, 0)
			if err != nil {
				return fmt.Errorf(
					"env var %s%s has invalid suffix: %v",
					ctx.prefix, name, err,
				)
			}
			ctx.name = name
			var parsed V
			err = p(&parsed)(val, ctx)
			if err != nil {
				return err
			}
			pairs = append(pairs, pair{
				i: index,
				k: name,
				v: parsed,
			})
		}

		slices.SortFunc(pairs, func(a, b pair) int {
			return cmp.Compare(a.i, b.i)
		})
		for _, p := range pairs {
			delete(ctx.env, p.k)
		}
		for _, p := range pairs {
			*target = append(*target, p.v)
		}
		return nil
	}
}
