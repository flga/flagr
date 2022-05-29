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

// ErrorHandling is an alias for flag.ErrorHandling
type ErrorHandling = stdflag.ErrorHandling

// Flag is an alias for flag.ErrorHandling
type Flag = stdflag.Flag

// Aliases for flag.ErrorHandling values
const (
	ContinueOnError = stdflag.ContinueOnError
	ExitOnError     = stdflag.ExitOnError
	PanicOnError    = stdflag.PanicOnError
)

// ErrHelp is an alias for flag.ErrHelp
var ErrHelp = stdflag.ErrHelp

// Set is akin to flag.FlagSet, with value methods (.Int, .String, etc) omitted
type Set struct {
	fs *stdflag.FlagSet
}

// NewSet behaves the same as flag.NewFlagSet
func NewSet(name string, errorHandling ErrorHandling) Set {
	fs := stdflag.NewFlagSet(name, errorHandling)
	return Set{fs}
}

// Output is an alias for flag.FlagSet.Output
func (set Set) Output() io.Writer { return set.fs.Output() }

// Name is an alias for flag.FlagSet.Name
func (set Set) Name() string { return set.fs.Name() }

// ErrorHandling is an alias for flag.FlagSet.ErrorHandling
func (set Set) ErrorHandling() ErrorHandling { return set.fs.ErrorHandling() }

// SetOutput is an alias for flag.FlagSet.SetOutput
func (set Set) SetOutput(output io.Writer) { set.fs.SetOutput(output) }

// VisitAll is an alias for flag.FlagSet.VisitAll
func (set Set) VisitAll(fn func(*Flag)) { set.fs.VisitAll(fn) }

// Visit is an alias for flag.FlagSet.Visit
func (set Set) Visit(fn func(*Flag)) { set.fs.Visit(fn) }

// Lookup is an alias for flag.FlagSet.Lookup
func (set Set) Lookup(name string) *Flag { return set.fs.Lookup(name) }

// Set is an alias for flag.FlagSet.Set
func (set Set) Set(name, value string) error { return set.fs.Set(name, value) }

// PrintDefaults is an alias for flag.FlagSet.PrintDefaults
func (set Set) PrintDefaults() { set.fs.PrintDefaults() }

// NFlag is an alias for flag.FlagSet.NFlag
func (set Set) NFlag() int { return set.fs.NFlag() }

// Arg is an alias for flag.FlagSet.Arg
func (set Set) Arg(i int) string { return set.fs.Arg(i) }

// NArg is an alias for flag.FlagSet.NArg
func (set Set) NArg() int { return set.fs.NArg() }

// Args is an alias for flag.FlagSet.Args
func (set Set) Args() []string { return set.fs.Args() }

// Parse is an alias for flag.FlagSet.Parse
func (set Set) Parse(arguments []string) error { return set.fs.Parse(arguments) }

// Parsed is an alias for flag.FlagSet.Parsed
func (set Set) Parsed() bool { return set.fs.Parsed() }

// Init is an alias for flag.FlagSet.Init
func (set Set) Init(name string, errorHandling ErrorHandling) { set.fs.Init(name, errorHandling) }

// Add creates a new flag on the given Set, returning the underlying value of the provided Getter.
func Add[T any](set Set, name string, value Getter[T], usage string) *T {
	set.fs.Var(value, name, usage)
	return value.Val()
}

func Int(defaultValue int) Getter[int]            { return Var(defaultValue, set(parseInt[int])) }
func Ints(defaults ...int) Getter[[]int]          { return Slice(defaults, parseInt[int]) }
func Int8(defaultValue int8) Getter[int8]         { return Var(defaultValue, set(parseInt[int8])) }
func Int8s(defaults ...int8) Getter[[]int8]       { return Slice(defaults, parseInt[int8]) }
func Int16(defaultValue int16) Getter[int16]      { return Var(defaultValue, set(parseInt[int16])) }
func Int16s(defaults ...int16) Getter[[]int16]    { return Slice(defaults, parseInt[int16]) }
func Int32(defaultValue int32) Getter[int32]      { return Var(defaultValue, set(parseInt[int32])) }
func Int32s(defaults ...int32) Getter[[]int32]    { return Slice(defaults, parseInt[int32]) }
func Int64(defaultValue int64) Getter[int64]      { return Var(defaultValue, set(parseInt[int64])) }
func Int64s(defaults ...int64) Getter[[]int64]    { return Slice(defaults, parseInt[int64]) }
func Uint(defaultValue uint) Getter[uint]         { return Var(defaultValue, set(parseUint[uint])) }
func Uints(defaults ...uint) Getter[[]uint]       { return Slice(defaults, parseUint[uint]) }
func Uint8(defaultValue uint8) Getter[uint8]      { return Var(defaultValue, set(parseUint[uint8])) }
func Uint8s(defaults ...uint8) Getter[[]uint8]    { return Slice(defaults, parseUint[uint8]) }
func Uint16(defaultValue uint16) Getter[uint16]   { return Var(defaultValue, set(parseUint[uint16])) }
func Uint16s(defaults ...uint16) Getter[[]uint16] { return Slice(defaults, parseUint[uint16]) }
func Uint32(defaultValue uint32) Getter[uint32]   { return Var(defaultValue, set(parseUint[uint32])) }
func Uint32s(defaults ...uint32) Getter[[]uint32] { return Slice(defaults, parseUint[uint32]) }
func Uint64(defaultValue uint64) Getter[uint64]   { return Var(defaultValue, set(parseUint[uint64])) }
func Uint64s(defaults ...uint64) Getter[[]uint64] { return Slice(defaults, parseUint[uint64]) }

func Float32(defaultValue float32) Getter[float32] {
	return Var(defaultValue, set(parseFloat[float32]))
}
func Float32s(defaults ...float32) Getter[[]float32] {
	return Slice(defaults, parseFloat[float32])
}
func Float64(defaultValue float64) Getter[float64] {
	return Var(defaultValue, set(parseFloat[float64]))
}
func Float64s(defaults ...float64) Getter[[]float64] {
	return Slice(defaults, parseFloat[float64])
}

func Complex64(defaultValue complex64) Getter[complex64] {
	return Var(defaultValue, set(parseComplex64))
}
func Complex64s(defaults ...complex64) Getter[[]complex64] {
	return Slice(defaults, parseComplex64)
}
func Complex128(defaultValue complex128) Getter[complex128] {
	return Var(defaultValue, set(parseComplex128))
}
func Complex128s(defaults ...complex128) Getter[[]complex128] {
	return Slice(defaults, parseComplex128)
}

func Bool(defaultValue bool) Getter[bool]   { return Var(defaultValue, set(strconv.ParseBool)) }
func Bools(defaults ...bool) Getter[[]bool] { return Slice(defaults, strconv.ParseBool) }

func String(defaultValue string) Getter[string]   { return Var(defaultValue, set(parseString)) }
func Strings(defaults ...string) Getter[[]string] { return Slice(defaults, parseString) }

func Duration(defaultValue time.Duration) Getter[time.Duration] {
	return Var(defaultValue, set(time.ParseDuration))
}
func Durations(defaults ...time.Duration) Getter[[]time.Duration] {
	return Slice(defaults, time.ParseDuration)
}
func Time(layout string, defaultValue time.Time) Getter[time.Time] {
	return Var(defaultValue, set(ptime(layout)))
}
func MustTime(layout string, defaultValue string) Getter[time.Time] {
	return MustVar(defaultValue, set(ptime(layout)))
}
func Times(layout string, defaults ...time.Time) Getter[[]time.Time] {
	return Slice(defaults, ptime(layout))
}
func MustTimes(layout string, defaults ...string) Getter[[]time.Time] {
	return MustSlice(defaults, ptime(layout))
}

func URL(defaultValue *url.URL) Getter[*url.URL]     { return Var(defaultValue, set(url.Parse)) }
func MustURL(defaultValue string) Getter[*url.URL]   { return MustVar(defaultValue, set(url.Parse)) }
func URLs(defaults ...*url.URL) Getter[[]*url.URL]   { return Slice(defaults, url.Parse) }
func MustURLs(defaults ...string) Getter[[]*url.URL] { return MustSlice(defaults, url.Parse) }

func IPAddr(defaultValue netip.Addr) Getter[netip.Addr] {
	return Var(defaultValue, set(netip.ParseAddr))
}
func MustIPAddr(defaultValue string) Getter[netip.Addr] {
	return MustVar(defaultValue, set(netip.ParseAddr))
}
func IPAddrs(defaults ...netip.Addr) Getter[[]netip.Addr] {
	return Slice(defaults, netip.ParseAddr)
}
func MustIPAddrs(defaults ...string) Getter[[]netip.Addr] {
	return MustSlice(defaults, netip.ParseAddr)
}
func IPAddrPort(defaultValue netip.AddrPort) Getter[netip.AddrPort] {
	return Var(defaultValue, set(netip.ParseAddrPort))
}
func MustIPAddrPort(defaultValue string) Getter[netip.AddrPort] {
	return MustVar(defaultValue, set(netip.ParseAddrPort))
}
func IPAddrPorts(defaults ...netip.AddrPort) Getter[[]netip.AddrPort] {
	return Slice(defaults, netip.ParseAddrPort)
}
func MustIPAddrPorts(defaults ...string) Getter[[]netip.AddrPort] {
	return MustSlice(defaults, netip.ParseAddrPort)
}

type Getter[T any] interface {
	stdflag.Getter
	Val() *T
	IsBoolFlag() bool
}

var _ Getter[any] = value[any]{}

type value[T any] struct {
	Value  *T
	Setter Setter[T]
}

func Var[T any](val T, setter Setter[T]) value[T] {
	return value[T]{
		Value:  &val,
		Setter: setter,
	}
}
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

type Parser[T any] func(string) (T, error)

var _ Getter[[]any] = &slice[any, []any]{}

type slice[T any, S ~[]T] struct {
	Value   *S
	Default S
	Parse   Parser[T]
	written bool
}

func Slice[T any, S ~[]T](defaultValue S, parse Parser[T]) *slice[T, S] {
	vcopy := make(S, len(defaultValue))
	copy(vcopy, defaultValue)
	return &slice[T, S]{
		Value: &vcopy,
		Parse: parse,
	}
}

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

type Setter[T any] func(*T, string) error

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
