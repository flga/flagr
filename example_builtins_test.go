package flagr_test

import (
	"errors"
	"fmt"
	"net/netip"
	"os"

	"github.com/flga/flagr"
)

func Example_builtins() {
	set := flagr.NewSet("mycmd", flagr.ContinueOnError)
	// Builtin types like int, float, etc.
	a := flagr.Add(&set, "a", flagr.Int(42), "usage")

	// Slices of things, every flag provided by flagr has a slice counterpart.
	b := flagr.Add(&set, "b", flagr.Ints(1, 2, 3), "usage")

	// A bool flag has the same semantics as std/flag bools.
	c := flagr.Add(&set, "c", flagr.Bool(false), "usage")

	// One of the extra types flagr provides, along with urls, times, etc.
	d := flagr.Add(&set, "d", flagr.IPAddrPort(netip.MustParseAddrPort("127.0.0.1:8080")), "usage")

	// Same as above but allows defaults to be provided as strings, delegating parsing to flagr.
	e := flagr.Add(&set, "e", flagr.MustIPAddrPort("127.0.0.1:8080"), "usage")

	args := []string{
		"-b", "5", "-b", "6",
		"-c",
		"-d", "0.0.0.0:80",
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
	fmt.Printf("e = %v and is a %T\n", *e, e)
	// Output:
	// a = 42 and is a *int
	// b = [5 6] and is a *[]int
	// c = true and is a *bool
	// d = 0.0.0.0:80 and is a *netip.AddrPort
	// e = 127.0.0.1:8080 and is a *netip.AddrPort
}
