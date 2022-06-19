package env_test

import (
	"errors"
	"os"

	"github.com/flga/flagr"
	"github.com/flga/flagr/env"
)

func Example() {
	var fs flagr.Set

	_ = flagr.Add(&fs, "flag", flagr.Int(1), "usage")
	envfile := flagr.Add(&fs, "envfile", flagr.String(".env"), "usage")

	if err := fs.Parse(
		os.Args[1:],
		env.Parse(
			env.WithPrefix("app"),                  // prefix every flag with "app" before mapping it
			env.WithMapper(env.DefaultMapper(",")), // use a custom mapper that treats "," as a separator for all values
			env.WithStaticDotEnv(".env", true),     // tries to read env values from a ".env" file if it exists
			env.WithDotEnv(envfile, true),          // same as above but allows the filename to be dynamic
			env.WithLookupFunc(os.LookupEnv),       // uses the given func to lookup env values (this is a noop, os.LookupEnv is the default)
		),
	); err != nil {
		if errors.Is(err, flagr.ErrHelp) {
			os.Exit(2)
		}
		os.Exit(1)
	}
}
