# flagr
[![Go Reference](https://pkg.go.dev/badge/github.com/flga/flagr.svg)](https://pkg.go.dev/github.com/flga/flagr)
[![Go Report Card](https://goreportcard.com/badge/github.com/flga/flagr)](https://goreportcard.com/report/github.com/flga/flagr)

A wrapper around the standard `flag` package that provides a more ergonomic API.

## Why

The `flag` package works well, but the API is not very ergonomic due to the language limitations at the time it was designed.

Namely there's 2 different APIs for adding flags, we have the builtin types like `int` and `string`, and for everything else we have to resort to creating implementations of `flag.Value`.

This isn't bad, by any means, but it has a cumbersome consequence:
```go
type customInt [...]

n1 := flag.Int("n1", 1234, "usage")
n2 := flag.String("n2", "asd", "usage")
var myCustomInt customInt
flag.Var(&myCustomInt, "n3", "usage")
n4 := flag.Duration("n4", 5*time.Second, "usage")
n5 := flag.Bool("n5", true, "usage")

```
Due to the fact that we're using a custom `flag.Value` we have to explicitly break apart the declaration and initialization.

This isn't the end of the world, and it's pretty much the best we could do pre 1.18, as there was no way for `flag.Var` to return a meaningful type. Now there is.

Flagr aims to solve this by leveraging generics:
```go
n1 := flagr.Add(set, "n1", flagr.Int(1234), "usage")               // n1 is *int
n2 := flagr.Add(set, "n2", flagr.String("asd"), "usage")           // n2 is *string
n3 := flagr.Add(set, "n3", CustomInt(42), "usage")                 // n3 is *int64
n4 := flagr.Add(set, "n4", flagr.Duration(5*time.Second), "usage") // n4 is *time.Duration
n5 := flagr.Add(set, "n5", flagr.Bool(true), "usage")              // n5 is *bool

// essentially a re-implementation of flagr.Int64
func CustomInt(def int64) flagr.Getter[int64] {
    // flagr.Var is a helper function that can construct a flagr.Getter (akin to flag.Getter)
    return flagr.Var(defaultValue, func(value *int64, s string) error {
        // parse
		v, err := strconv.ParseInt(s, 0, 64)
		if err != nil {
			return err
		}

        // set
		*value = v
		return nil
	})
}

func UsingSetterFrom(def int) flagr.Getter[int] {
    // Setting a flag value means 2 things, parsing it and actually setting.
    // In cases where setting is just a simple assignment (as oposed to say, slices)
    // and we already have a function that conforms to func(string)(T, error)
    // we can leverage flagr.SetterFrom.
    return flagr.Var(42, flagr.SetterFrom(strconv.Atoi))
}
```

`flagr.Add` accepts a `flagr.Getter[X]`, so custom types are also supported like in the standard flag package:
```go
// counts the number of times a flag was used
type Counter struct {
	value *int
	set   bool
}

func NewCounter(defaultValue int) flagr.Getter[int] {
	return &Counter{value: &defaultValue}
}

// methods that satisfy the standard Getter interface
func (f Counter) IsBoolFlag() bool { return true }
func (f *Counter) Get() any        { return f.value }
func (f *Counter) String() string  { [...] }
func (f *Counter) Set(s string) error {
	if !f.set {
		*f.value = 0
		f.set = true
	}

	*f.value++
	return nil
}

// satisfies flagr.Getter, allowing us to decouple the flag type (a struct in this case)
// from the actual data we're interested on, the count.
func (f Counter) Val() *int { return f.value }

```
And using it is exactly the same as we've seen before:
```go
count := flagr.Add(set, "c", NewCounter(2), "c") // count is a *int
```

## Goals
Flagr aims to be just a thin wrapper around the standard flag package.

Our goal is to clean up the API, not to add new features or change existing behaviour.

## Supported types
We've taken this oportunity to add in some extra goodies like builtin slice support
and a few (non external) types that are useful for most go programs like `url.URL` and `netip.Addr`.

Sure this adds a little weight but shouldn't be too much of an issue in practice, we might revisit it later.
    
### Builtin types   
- int, []int
- int8, []int8
- int16, []int16
- int32, []int32
- int64, []int64
- uint, []uint
- uint16, []uint16
- uint32, []uint32
- uint64, []uint64
- uint8, []uint8
- float32, []float32
- float64, []float64
- complex128, []complex128
- complex64, []complex64
- bool, []bool
- string, []string

### Time
- time.Duration, []time.Duration
- time.Time, []time.Time

### Networking
- netip.Addr, []netip.Addr
- netip.AddrPort, []netip.AddrPort
- url.URL, []url.URL

## Full documentation
[![Go Reference](https://pkg.go.dev/badge/github.com/flga/flagr.svg)](https://pkg.go.dev/github.com/flga/flagr)
