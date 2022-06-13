package flagr_test

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"strconv"

	"github.com/flga/flagr"
)

func Example_slices() {
	var set flagr.Set

	// We can use flagr.Slice to create repeatable flags. Custom slice types are supported.
	type bools []bool
	a := flagr.Add[bools](&set, "a", flagr.Slice(bools{true, false, true}, strconv.ParseBool), "usage")

	// Making use of MustSlice so that we can pass default values as strings and delegate the parsing to flagr.
	// Makes it much more comfortable to use things that have to be parsed, like urls.
	b := flagr.Add[[]*url.URL](&set, "b", flagr.MustSlice([]string{"https://a.com", "https://b.com"}, url.Parse), "usage")

	// Using helper functions allows us to have cleaner signatures (type inference is still a little wonky)
	c := flagr.Add(&set, "c", Urls("https://a.com", "https://b.com"), "usage")

	args := []string{
		"-a", "-a", // std/flag bool semantics apply
		"-b", "https://c.com", "-b", "https://d.com",
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
	// Output:
	// a = [true true] and is a *flagr_test.bools
	// b = [https://c.com https://d.com] and is a *[]*url.URL
	// c = [https://a.com https://b.com] and is a *[]*url.URL
}

func Urls(defaults ...string) flagr.Getter[[]*url.URL] {
	return flagr.MustSlice(defaults, url.Parse)
}
