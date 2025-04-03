package configenv_test

import (
	"fmt"

	"github.com/orsinium-labs/configenv"
)

type Config struct {
	Debug bool
	Env   string
}

func Example() {
	env := []string{
		"BE_DEBUG=true",
		"BE_ENV=prod",
		"ENV=fake",
	}
	config := Config{}
	vars := configenv.Vars{
		"DEBUG": configenv.Required(configenv.Bool(&config.Debug)),
		"ENV":   configenv.String(&config.Env),
	}
	err := vars.Parse(configenv.Config{
		Environ: env,
		Prefix:  "BE_",
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(config.Env)
	//Output: prod
}
