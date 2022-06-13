package flagr_test

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"strconv"

	"github.com/flga/flagr"
)

func Example_var() {
	var set flagr.Set

	// Using flagr.Var to reimplement flagr.Int64.
	a := flagr.Add(&set, "a", flagr.Var(42, func(value *int64, s string) error {
		v, err := strconv.ParseInt(s, 0, 64)
		if err != nil {
			return err
		}
		*value = v
		return nil
	}), "usage")

	// Using helper functions allows us to have cleaner signatures
	b := flagr.Add(&set, "b", CustomInt64(42), "usage")

	// Making use of flagr.SetterFrom when setting a value is just a simple assignment (like above).
	// Also, if the underlying type is a ~bool std/flag bool semantics apply.
	c := flagr.Add(&set, "c", flagr.Var(false, flagr.SetterFrom(strconv.ParseBool)), "usage")

	// Making use of MustVar so that we can pass default values as strings and delegate the parsing to flagr.
	// Makes it much more comfortable to use things that have to be parsed, like urls.
	d := flagr.Add[*url.URL](&set, "d", flagr.MustVar("http://a.com", flagr.SetterFrom(url.Parse)), "usage")

	args := []string{
		"-a", "1",
		"-c", // std/flag bool semantics apply
		"-d", "http://b.com",
	}
	if err := set.Parse(args); err != nil {
		if errors.Is(err, flagr.ErrHelp) {
			os.Exit(0)
		}
		os.Exit(2)
	}

	fmt.Printf("a = %v and is a %T\n", *a, a)
	fmt.Printf("b = %v and is a %T\n", *b, b)
	fmt.Printf("c = %v and is a %T\n", *c, c)
	fmt.Printf("d = %v and is a %T\n", *d, d)
	// Output:
	// a = 1 and is a *int64
	// b = 42 and is a *int64
	// c = true and is a *bool
	// d = http://b.com and is a **url.URL
}

// an implementation equivalent to flagr.Int64
func CustomInt64(defaultValue int64) flagr.Getter[int64] {
	return flagr.Var(defaultValue, func(value *int64, s string) error {
		v, err := strconv.ParseInt(s, 0, 64)
		if err != nil {
			return err
		}
		*value = v
		return nil
	})
}
