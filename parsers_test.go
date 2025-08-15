//nolint:nilness
package configenv_test

import (
	"slices"
	"testing"

	"github.com/orsinium-labs/configenv"
)

func TestBool(t *testing.T) {
	check := func(t *testing.T, val string, exp bool) {
		env := []string{"BE_XYZ=" + val}
		var act bool
		vars := configenv.Vars{"XYZ": configenv.Bool(&act)}
		err := vars.Parse(configenv.Config{Environ: env, Prefix: "BE_"})
		if err != nil {
			t.Fatalf("parse error: %v", err)
		}
		if act != exp {
			t.Fatalf("want %v, got %v", act, exp)
		}
	}

	check(t, "true", true)
	check(t, "True", true)
	check(t, "TRUE", true)
	check(t, "1", true)

	check(t, "false", false)
	check(t, "False", false)
	check(t, "FALSE", false)
	check(t, "0", false)
	check(t, "", false)
}

func TestInt(t *testing.T) {
	check := func(t *testing.T, val string, exp int) {
		env := []string{"BE_XYZ=" + val}
		var act int
		vars := configenv.Vars{"XYZ": configenv.Int(&act)}
		err := vars.Parse(configenv.Config{Environ: env, Prefix: "BE_"})
		if err != nil {
			t.Fatalf("parse error: %v", err)
		}
		if act != exp {
			t.Fatalf("want %v, got %v", act, exp)
		}
	}

	check(t, "", 0)
	check(t, "1", 1)
	check(t, "13", 13)
	check(t, "-13", -13)
	check(t, "-0", 0)
}

func TestString(t *testing.T) {
	check := func(t *testing.T, val string, exp string) {
		env := []string{"BE_XYZ=" + val}
		var act string
		vars := configenv.Vars{"XYZ": configenv.String(&act)}
		err := vars.Parse(configenv.Config{Environ: env, Prefix: "BE_"})
		if err != nil {
			t.Fatalf("parse error: %v", err)
		}
		if act != exp {
			t.Fatalf("want %v, got %v", act, exp)
		}
	}

	check(t, "", "")
	check(t, "hi", "hi")
	check(t, "hello world!", "hello world!")
	check(t, "don't", "don't")
}

func TestStrings(t *testing.T) {
	check := func(t *testing.T, val string, exp []string) {
		env := []string{"BE_XYZ=" + val}
		var act []string
		vars := configenv.Vars{"XYZ": configenv.Strings(&act, ",")}
		err := vars.Parse(configenv.Config{Environ: env, Prefix: "BE_"})
		if err != nil {
			t.Fatalf("parse error: %v", err)
		}
		if !slices.Equal(act, exp) {
			t.Fatalf("want %v, got %v", act, exp)
		}
	}

	check(t, "", []string{})
	check(t, "h", []string{"h"})
	check(t, "hello", []string{"hello"})
	check(t, "hello world", []string{"hello world"})
	check(t, "hello,world", []string{"hello", "world"})
	check(t, "hello, world", []string{"hello", " world"})
}
