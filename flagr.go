// Pakckage flagr is a thin wrapper over the standard flag package creating a generalized API for defining flags of any type.
package flagr

import (
	stdflag "flag"
	"fmt"
	"io"
	"net/netip"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"
)

const intSize = 32 << (^uint(0) >> 63)

// ErrorHandling defines how Set.Parse behaves if the parse fails.
type ErrorHandling = stdflag.ErrorHandling

// A Flag represents the state of a flag.
type Flag = stdflag.Flag

// These constants cause Set.Parse to behave as described if the parse fails.
const (
	ContinueOnError = stdflag.ContinueOnError // Return a descriptive error.
	ExitOnError     = stdflag.ExitOnError     // Call os.Exit(2) or for -h/-help Exit(0).
	PanicOnError    = stdflag.PanicOnError    // Call panic with a descriptive error.
)

// ErrHelp is the error returned if the -help or -h flag is invoked but no such flag is defined.
var ErrHelp = stdflag.ErrHelp

// A Set represents a set of defined flags. The zero value of a Set
// has no name and has ContinueOnError error handling.
//
// Flag names must be unique within a Set. An attempt to define a flag whose
// name is already in use will cause a panic.
type Set struct {
	fs *stdflag.FlagSet
}

// NewSet returns a new, empty flag set with the specified name and
// error handling property. If the name is not empty, it will be printed
// in the default usage message and in error messages.
func NewSet(name string, errorHandling ErrorHandling) Set {
	fs := stdflag.NewFlagSet(name, errorHandling)
	return Set{fs: fs}
}

func (set *Set) init() {
	if set.fs == nil {
		set.fs = &stdflag.FlagSet{}
	}

	if set.fs.Usage == nil {
		set.fs.Usage = func() {
			if set.fs.Name() == "" {
				fmt.Fprintf(set.fs.Output(), "Usage:\n")
			} else {
				fmt.Fprintf(set.fs.Output(), "Usage of %s:\n", set.fs.Name())
			}
			set.fs.PrintDefaults()
		}
	}
}

// SetUsage overrides the Set's usage func.
func (set *Set) SetUsage(usage func()) {
	set.init()
	set.fs.Usage = usage
}

// Usage invokes the usage function provided with SetUsage.
func (set *Set) Usage() {
	set.init()
	set.fs.Usage()
}

// Output returns the destination for usage and error messages. os.Stderr is returned if
// output was not set or was set to nil.
func (set *Set) Output() io.Writer { set.init(); return set.fs.Output() }

// Name returns the name of the flag set.
func (set *Set) Name() string { set.init(); return set.fs.Name() }

// ErrorHandling returns the error handling behavior of the flag set.
func (set *Set) ErrorHandling() ErrorHandling { set.init(); return set.fs.ErrorHandling() }

// SetOutput sets the destination for usage and error messages.
// If output is nil, os.Stderr is used.
func (set *Set) SetOutput(output io.Writer) { set.init(); set.fs.SetOutput(output) }

// VisitAll visits the flags in lexicographical order, calling fn for each.
// It visits all flags, even those not set.
func (set *Set) VisitAll(fn func(*Flag)) { set.init(); set.fs.VisitAll(fn) }

// Visit visits the flags in lexicographical order, calling fn for each.
// It visits only those flags that have been set.
func (set *Set) Visit(fn func(*Flag)) { set.init(); set.fs.Visit(fn) }

// Lookup returns the Flag structure of the named flag, returning nil if none exists.
func (set *Set) Lookup(name string) *Flag { set.init(); return set.fs.Lookup(name) }

// Set sets the value of the named flag.
func (set *Set) Set(name, value string) error { set.init(); return set.fs.Set(name, value) }

// UnquoteUsage extracts a back-quoted name from the usage
// string for a flag and returns it and the un-quoted usage.
// Given "a `name` to show" it returns ("name", "a name to show").
// If there are no back quotes, the name is an educated guess of the
// type of the flag's value, or the empty string if the flag is boolean.
func UnquoteUsage(flag *Flag) (name string, usage string) {
	return stdflag.UnquoteUsage(flag)
}

// PrintDefaults prints, to standard error unless configured otherwise, the
// default values of all defined command-line flags in the set. See the
// documentation for the global function PrintDefaults for more information.
func (set *Set) PrintDefaults() { set.init(); set.fs.PrintDefaults() }

// NFlag returns the number of flags that have been set.
func (set *Set) NFlag() int { set.init(); return set.fs.NFlag() }

// Arg returns the i'th argument. Arg(0) is the first remaining argument
// after flags have been processed. Arg returns an empty string if the
// requested element does not exist.
func (set *Set) Arg(i int) string { set.init(); return set.fs.Arg(i) }

// NArg is the number of arguments remaining after flags have been processed.
func (set *Set) NArg() int { set.init(); return set.fs.NArg() }

// Args returns the non-flag arguments.
func (set *Set) Args() []string { set.init(); return set.fs.Args() }

// Parse parses flag definitions from the argument list, which should not
// include the command name. Must be called after all flags in the FlagSet
// are defined and before flags are accessed by the program.
// The return value will be ErrHelp if -help or -h were set but not defined.
func (set *Set) Parse(arguments []string) error { set.init(); return set.fs.Parse(arguments) }

// Parsed reports whether f.Parse has been called.
func (set *Set) Parsed() bool { set.init(); return set.fs.Parsed() }

// Init sets the name and error handling property for a flag set.
// By default, the zero FlagSet uses an empty name and the
// ContinueOnError error handling policy.
func (set *Set) Init(name string, errorHandling ErrorHandling) {
	set.init()
	set.fs.Init(name, errorHandling)
}

// Add creates a new flag on the given Set, returning the underlying value of the provided Getter.
func Add[T any](set *Set, name string, value Getter[T], usage string) *T {
	set.init()
	set.fs.Var(value, name, usage)
	return value.Val()
}

// Int returns a Getter that can parse values of type int.
func Int(defaultValue int) Getter[int] {
	return Var(defaultValue, set(parseInt[int]))
}

// Ints returns a Getter that can parse and accumulate values of type int.
func Ints(defaults ...int) Getter[[]int] {
	return Slice(defaults, parseInt[int])
}

// Int8 returns a Getter that can parse values of type int8.
func Int8(defaultValue int8) Getter[int8] {
	return Var(defaultValue, set(parseInt[int8]))
}

// Int8s returns a Getter that can parse and accumulate values of type int8.
func Int8s(defaults ...int8) Getter[[]int8] {
	return Slice(defaults, parseInt[int8])
}

// Int16 returns a Getter that can parse values of type int16.
func Int16(defaultValue int16) Getter[int16] {
	return Var(defaultValue, set(parseInt[int16]))
}

// Int16s returns a Getter that can parse and accumulate values of type int16.
func Int16s(defaults ...int16) Getter[[]int16] {
	return Slice(defaults, parseInt[int16])
}

// Int32 returns a Getter that can parse values of type int32.
func Int32(defaultValue int32) Getter[int32] {
	return Var(defaultValue, set(parseInt[int32]))
}

// Int32s returns a Getter that can parse and accumulate values of type int32.
func Int32s(defaults ...int32) Getter[[]int32] {
	return Slice(defaults, parseInt[int32])
}

// Int64 returns a Getter that can parse values of type int64.
func Int64(defaultValue int64) Getter[int64] {
	return Var(defaultValue, set(parseInt[int64]))
}

// Int64s returns a Getter that can parse and accumulate values of type int64.
func Int64s(defaults ...int64) Getter[[]int64] {
	return Slice(defaults, parseInt[int64])
}

// Uint returns a Getter that can parse values of type uint.
func Uint(defaultValue uint) Getter[uint] {
	return Var(defaultValue, set(parseUint[uint]))
}

// Uints returns a Getter that can parse and accumulate values of type uint.
func Uints(defaults ...uint) Getter[[]uint] {
	return Slice(defaults, parseUint[uint])
}

// Uint8 returns a Getter that can parse values of type uint8.
func Uint8(defaultValue uint8) Getter[uint8] {
	return Var(defaultValue, set(parseUint[uint8]))
}

// Uint8s returns a Getter that can parse and accumulate values of type uint8.
func Uint8s(defaults ...uint8) Getter[[]uint8] {
	return Slice(defaults, parseUint[uint8])
}

// Uint16 returns a Getter that can parse values of type uint16.
func Uint16(defaultValue uint16) Getter[uint16] {
	return Var(defaultValue, set(parseUint[uint16]))
}

// Uint16s returns a Getter that can parse and accumulate values of type uint16.
func Uint16s(defaults ...uint16) Getter[[]uint16] {
	return Slice(defaults, parseUint[uint16])
}

// Uint32 returns a Getter that can parse values of type uint32.
func Uint32(defaultValue uint32) Getter[uint32] {
	return Var(defaultValue, set(parseUint[uint32]))
}

// Uint32s returns a Getter that can parse and accumulate values of type uint32.
func Uint32s(defaults ...uint32) Getter[[]uint32] {
	return Slice(defaults, parseUint[uint32])
}

// Uint64 returns a Getter that can parse values of type uint64.
func Uint64(defaultValue uint64) Getter[uint64] {
	return Var(defaultValue, set(parseUint[uint64]))
}

// Uint64s returns a Getter that can parse and accumulate values of type uint64.
func Uint64s(defaults ...uint64) Getter[[]uint64] {
	return Slice(defaults, parseUint[uint64])
}

// Float32 returns a Getter that can parse values of type float32.
func Float32(defaultValue float32) Getter[float32] {
	return Var(defaultValue, set(parseFloat[float32]))
}

// Float32s returns a Getter that can parse and accumulate values of type float32.
func Float32s(defaults ...float32) Getter[[]float32] {
	return Slice(defaults, parseFloat[float32])
}

// Float64 returns a Getter that can parse values of type float64.
func Float64(defaultValue float64) Getter[float64] {
	return Var(defaultValue, set(parseFloat[float64]))
}

// Float64s returns a Getter that can parse and accumulate values of type float64.
func Float64s(defaults ...float64) Getter[[]float64] {
	return Slice(defaults, parseFloat[float64])
}

// Complex64 returns a Getter that can parse values of type complex64.
func Complex64(defaultValue complex64) Getter[complex64] {
	return Var(defaultValue, set(parseComplex64))
}

// Complex64s returns a Getter that can parse and accumulate values of type complex64.
func Complex64s(defaults ...complex64) Getter[[]complex64] {
	return Slice(defaults, parseComplex64)
}

// Complex128 returns a Getter that can parse values of type complex128.
func Complex128(defaultValue complex128) Getter[complex128] {
	return Var(defaultValue, set(parseComplex128))
}

// Complex128s returns a Getter that can parse and accumulate values of type complex128.
func Complex128s(defaults ...complex128) Getter[[]complex128] {
	return Slice(defaults, parseComplex128)
}

// Bool returns a Getter that can parse values of type bool.
func Bool(defaultValue bool) Getter[bool] {
	return Var(defaultValue, set(strconv.ParseBool))
}

// Bools returns a Getter that can parse and accumulate values of type bool.
func Bools(defaults ...bool) Getter[[]bool] {
	return Slice(defaults, strconv.ParseBool)
}

// String returns a Getter that can parse values of type string.
func String(defaultValue string) Getter[string] {
	return Var(defaultValue, set(parseString))
}

// Strings returns a Getter that can parse and accumulate values of type string.
func Strings(defaults ...string) Getter[[]string] {
	return Slice(defaults, parseString)
}

// Duration returns a Getter that can parse values of type time.Duration.
func Duration(defaultValue time.Duration) Getter[time.Duration] {
	return Var(defaultValue, set(time.ParseDuration))
}

// Durations returns a Getter that can parse and accumulate values of type time.Duration.
func Durations(defaults ...time.Duration) Getter[[]time.Duration] {
	return Slice(defaults, time.ParseDuration)
}

// Time returns a Getter that can parse values of type time.Time.
func Time(layout string, defaultValue time.Time) Getter[time.Time] {
	return Var(defaultValue, set(ptime(layout)))
}

// MustTime, like Time, returns a Getter that can parse values of type time.Time, but
// allowing the default value to be provided as a string. It panics if the given string cannot be parsed
// as time.Time.
func MustTime(layout string, defaultValue string) Getter[time.Time] {
	return MustVar(defaultValue, set(ptime(layout)))
}

// Times returns a Getter that can parse and accumulate values of type time.Time.
func Times(layout string, defaults ...time.Time) Getter[[]time.Time] {
	return Slice(defaults, ptime(layout))
}

// MustTimes, like Times, returns a Getter that can parse values of type time.Time and accumulate them, but
// allowing the default values to be provided as strings. It panics if any given string cannot be parsed
// as time.Time.
func MustTimes(layout string, defaults ...string) Getter[[]time.Time] {
	return MustSlice(defaults, ptime(layout))
}

// URL returns a Getter that can parse values of type *url.URL.
func URL(defaultValue *url.URL) Getter[*url.URL] {
	return Var(defaultValue, set(url.Parse))
}

// MustURL, like URL, returns a Getter that can parse values of type *url.URL, but
// allowing the default value to be provided as a string. It panics if the given string cannot be parsed
// as *url.URL.
func MustURL(defaultValue string) Getter[*url.URL] {
	return MustVar(defaultValue, set(url.Parse))
}

// URLs returns a Getter that can parse and accumulate values of type *url.URL.
func URLs(defaults ...*url.URL) Getter[[]*url.URL] {
	return Slice(defaults, url.Parse)
}

// MustURLs, like URLs, returns a Getter that can parse values of type *url.URL and accumulate them, but
// allowing the default values to be provided as strings. It panics if any given string cannot be parsed
// as *url.URL.
func MustURLs(defaults ...string) Getter[[]*url.URL] {
	return MustSlice(defaults, url.Parse)
}

// IPAddr returns a Getter that can parse values of type netip.Addr.
func IPAddr(defaultValue netip.Addr) Getter[netip.Addr] {
	return Var(defaultValue, set(netip.ParseAddr))
}

// MustIPAddr, like IPAddr, returns a Getter that can parse values of type netip.Addr, but
// allowing the default value to be provided as a string. It panics if the given string cannot be parsed
// as netip.Addr.
func MustIPAddr(defaultValue string) Getter[netip.Addr] {
	return MustVar(defaultValue, set(netip.ParseAddr))
}

// IPAddrs returns a Getter that can parse and accumulate values of type netip.Addr.
func IPAddrs(defaults ...netip.Addr) Getter[[]netip.Addr] {
	return Slice(defaults, netip.ParseAddr)
}

// MustIPAddrs, like IPAddrs, returns a Getter that can parse values of type netip.Addr and accumulate them, but
// allowing the default values to be provided as strings. It panics if any given string cannot be parsed
// as netip.Addr.
func MustIPAddrs(defaults ...string) Getter[[]netip.Addr] {
	return MustSlice(defaults, netip.ParseAddr)
}

// IPAddrPort returns a Getter that can parse values of type netip.AddrPort.
func IPAddrPort(defaultValue netip.AddrPort) Getter[netip.AddrPort] {
	return Var(defaultValue, set(netip.ParseAddrPort))
}

// MustIPAddrPort, like IPAddrPort, returns a Getter that can parse values of type netip.AddrPort, but
// allowing the default value to be provided as a string. It panics if the given string cannot be parsed
// as netip.AddrPort.
func MustIPAddrPort(defaultValue string) Getter[netip.AddrPort] {
	return MustVar(defaultValue, set(netip.ParseAddrPort))
}

// IPAddrPorts returns a Getter that can parse and accumulate values of type netip.AddrPort.
func IPAddrPorts(defaults ...netip.AddrPort) Getter[[]netip.AddrPort] {
	return Slice(defaults, netip.ParseAddrPort)
}

// MustIPAddrPorts, like IPAddrPorts, returns a Getter that can parse values of type netip.AddrPort and accumulate them, but
// allowing the default values to be provided as strings. It panics if any given string cannot be parsed
// as netip.AddrPort.
func MustIPAddrPorts(defaults ...string) Getter[[]netip.AddrPort] {
	return MustSlice(defaults, netip.ParseAddrPort)
}

// Getter is any type that satisfies flag.Getter and provides a new method Val()
// that returns a pointer to the actual value of type.
//
// This allows constructing a flag.Getter such that its type and the type of its
// value can (but don't need too) be diferent.
//
// In scenarios where the implementation is simple both types should be the same,
// such as custom builtin type, for example.
//
// In scenarios where the implementation is complex, we might want to use different
// types.
//
// Let's consider a Getter implementation that counts the times a flag has been provided.
//
// It could be described as a simple integer, but this is insufficient for a correct
// implementation given that we would need to know, somehow, if the current value
// is the default value and we should reset it, or if we can just increment the count.
//
// Since we need to keep track of that extra state, we'd need to use a struct to
// implement Getter, but because the type that implements a Getter and the type
// that a Getter returns as its value need not be the same, we can define a
// struct MyCounter that implements a Getter[int], which is the actual value the
// caller is interested in. The fact that we had to use a struct is an implementation
// detail.
//
// You can look at the CustomTypes example for a concrete implementation.
type Getter[T any] interface {
	stdflag.Getter
	IsBoolFlag() bool
	Val() *T
}

// Parser is a func that parses a string into T.
type Parser[T any] func(string) (T, error)

// Setter is a function that can, given a string, construct a meaningful value
// and assign it to *T.
type Setter[T any] func(*T, string) error

// Setter from returns a Setter that assigns into *T the result of Parser[T].
func SetterFrom[T any](parser Parser[T]) Setter[T] {
	return set(parser)
}

func set[T any](parse Parser[T]) Setter[T] {
	return func(t *T, s string) error {
		v, err := parse(s)
		if err != nil {
			return err
		}
		*t = T(v)
		return nil
	}
}

var _ Getter[any] = value[any]{}

type value[T any] struct {
	Value  *T
	Setter Setter[T]
}

// Var returns a Getter[T] with the given default value and Setter.
//
// Setter is called to parse the provided value and assign it to the underlying value.
func Var[T any](val T, setter Setter[T]) Getter[T] {
	return value[T]{
		Value:  &val,
		Setter: setter,
	}
}

// MustVar, like Var, returns a Getter[T] with the given default value and Setter,
// but allows the default value to be provided as a string.
//
// When calling MustVar the given setter will be used to convert the default value,
// into T. If it returns an error, MustVar panics.
//
// Setter is called to parse the provided value and assign it to the underlying value.
func MustVar[T any](defaultValue string, setter Setter[T]) value[T] {
	v := new(T)
	if err := setter(v, defaultValue); err != nil {
		panic(fmt.Errorf("flag: invalid default value %q: %w", defaultValue, err))
	}
	return value[T]{
		Value:  v,
		Setter: setter,
	}
}

func (v value[T]) Get() any {
	return v.Value
}

func (v value[T]) Val() *T {
	return v.Value
}

func (v value[T]) Set(s string) error {
	return v.Setter(v.Value, s)
}

func (v value[T]) String() string {
	if v.Value == nil {
		return "<nil>"
	}
	return fmt.Sprint(*v.Value)
}

func (v value[T]) IsBoolFlag() bool {
	return reflect.TypeOf(*v.Value).Kind() == reflect.Bool // seems insufficiently general, we'll see
}

var _ Getter[[]any] = &slice[any, []any]{}

type slice[T any, S ~[]T] struct {
	Value   *S
	Default S
	Parse   Parser[T]
	written bool
}

// Slice returns a Getter[S] with the given default value.
//
// If the same flag is provided multiple times, the result will be
// accumulated in S.
//
// The value will be initialized with a shallow copy of defaultValue.
func Slice[T any, S ~[]T](defaultValue S, parse Parser[T]) *slice[T, S] {
	vcopy := make(S, len(defaultValue))
	copy(vcopy, defaultValue)
	return &slice[T, S]{
		Value: &vcopy,
		Parse: parse,
	}
}

// MustSlice, returns a Getter[[]T] with the given default value,
// but allows the default values to be provided as a strings. Unlike Slice
// custom slice implementations are not supported. This could change in the future.
//
// When calling MustSlice each value will be parsed with the given Parser,
// any error will cause MustSlice to panic.
//
// If the same flag is provided multiple times, the result will be
// accumulated in a []T.
func MustSlice[T any](defaults []string, parse Parser[T]) *slice[T, []T] {
	vcopy := make([]T, len(defaults))
	for i, def := range defaults {
		v, err := parse(def)
		if err != nil {
			panic(fmt.Errorf("flag: invalid default value %q: %w", def, err))
		}
		vcopy[i] = v
	}
	return &slice[T, []T]{
		Value: &vcopy,
		Parse: parse,
	}
}

func (s *slice[T, S]) Get() any {
	return s.Value
}

func (s *slice[T, S]) Val() *S {
	return s.Value
}

func (f *slice[T, S]) Set(s string) error {
	if !f.written {
		*f.Value = (*f.Value)[:0]
		f.written = true
	}
	v, err := f.Parse(s)
	if err != nil {
		return err
	}
	*f.Value = append(*f.Value, v)
	return nil
}

func (s *slice[T, S]) String() string {
	if s.Value == nil {
		return "<nil>"
	}

	var buf strings.Builder
	buf.WriteByte('[')
	for i, v := range *s.Value {
		if i > 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(fmt.Sprint(v))
	}
	buf.WriteByte(']')
	return buf.String()
}

func (s *slice[T, S]) IsBoolFlag() bool {
	return reflect.TypeOf(S{}).Elem().Kind() == reflect.Bool
}

func parseInt[T ~int8 | ~int16 | ~int32 | ~int64 | ~int](s string) (T, error) {
	v, err := strconv.ParseInt(s, 0, intSize)
	if err != nil {
		return 0, err
	}
	return T(v), nil
}

func parseUint[T ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uint](s string) (T, error) {
	v, err := strconv.ParseUint(s, 0, intSize)
	if err != nil {
		return 0, err
	}
	return T(v), nil
}

func parseFloat[T ~float32 | ~float64](s string) (T, error) {
	v, err := strconv.ParseFloat(s, intSize)
	if err != nil {
		return 0, err
	}
	return T(v), nil
}

func parseComplex64(s string) (complex64, error) {
	v, err := strconv.ParseComplex(s, 64)
	if err != nil {
		return 0, err
	}
	return complex64(v), nil
}
func parseComplex128(s string) (complex128, error) {
	v, err := strconv.ParseComplex(s, 128)
	if err != nil {
		return 0, err
	}
	return complex128(v), nil
}

func parseString(s string) (string, error) { return s, nil }

func ptime(layout string) func(string) (time.Time, error) {
	return func(s string) (time.Time, error) { return time.Parse(layout, s) }
}
