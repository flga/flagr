// Package file provides a flagr.Parser that is able to read config files and set the appropriate flags.
package file

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"

	"github.com/flga/flagr"
)

// KeyPathSeparator is the value used to separate sub paths in path expressions.
// A typical KeyPath will look something like "root.inner.prop".
const KeyPathSeparator = "."

// KeyPath represents a path expression which is used to walk the values produced
// by the decoding of the config file. It is akin to jsonpath but we only
// support dot notation.
type KeyPath string

// Split is a convenience method to split a [KeyPath] into sub paths.
func (k KeyPath) Split() []string {
	return strings.Split(string(k), KeyPathSeparator)
}

// Mapper converts a flag name into a [KeyPath].
//
// Sub paths must be separated by [KeyPathSeparator] and are not allowed to contain
// it within them.
//
// Given a flag name like "api_http_address" and a json config file structured as
//
//	{
//	    "api": {
//	        "http": {
//	            "address": "0.0.0.0"
//	        }
//	    }
//	}
//
// the corresponding [KeyPath] should be "api.http.address".
//
// If the json file were flat, no mapping would be necessary as long as
// the field name and the flag name are equal:
//
//	{
//		"api_http_address": "0.0.0.0"
//	}
//
// the corresponding [KeyPath] should be "api_http_address".
type Mapper func(flagName string) KeyPath

// NoopMapper returns the given flag name as is.
func NoopMapper(flagName string) KeyPath { return KeyPath(flagName) }

// DecoderFunc is a function that deserializes data into v. Functions like
// json.Unmarshal conform to this type.
type DecoderFunc func(data []byte, v interface{}) error

// Extension represents a file extension. Values must include the leading dot.
type Extension string

// Mux maps [Extension] to [DecoderFunc], it is valid for the same [DecoderFunc] to
// be used for multiple [Extension]. For example, if unmarshaling yaml it is
// encouraged to map both ".yaml" and ".yml" to yaml.Unmarshal.
type Mux map[Extension]DecoderFunc

func (m Mux) supportedExts() []Extension {
	ret := make([]Extension, len(m))
	i := 0
	for k := range m {
		ret[i] = k
		i++
	}
	return ret
}

// Options contains all the options used to parse a config file.
type Options struct {
	Mapper            Mapper // Maps flag names to property paths
	IgnoreMissingFile bool   // If true, we don't treat [fs.ErrNotExist] as an error.
	FS                fs.FS  // If provided, this will be used instead of the primary filesystem.
}

// Option is a function that mutates Options.
type Option func(*Options)

// WithMapper configures Parser to use the given [Mapper] when mapping
// flag names to [KeyPath].
func WithMapper(n Mapper) Option {
	return func(o *Options) {
		o.Mapper = n
	}
}

// IgnoreMissingFile makes it so that if the provided file doesn't exist, it is
// not considered an error.
func IgnoreMissingFile() Option {
	return func(o *Options) {
		o.IgnoreMissingFile = true
	}
}

// With FS configures the Parser such that the file is retrieved from the given
// fs instead of the primary filesystem.
func WithFS(fs fs.FS) Option {
	return func(o *Options) {
		o.FS = fs
	}
}

// Parse returns a [flagr.FlagParser] that parses the file stored in path and
// assigns the results to any flags that have not yet been set.
//
// Flags cannot have complex values, only primitive values are allowed: strings,
// bools, ints etc. An exception is made for slices as flags can be repeatable.
//
// This is a valid JSON declaration for a flag named "foo".
//
//	{
//		"foo": 42
//	}
//
// So is this (assuming that foo is repeatable, otherwise it will be set to the last value of the array):
//
//	{
//		"foo": [1, 2, 3]
//	}
//
// But this is not, as there is no way to map a json object to a flag value:
//
//	{
//		"foo": {"bar": 42}
//	}
//
// The file contents are read and decoded using the decoder mapped to the file's
// extension. If decoding fails it returns [ErrDecode], if no suitable decoder is found
// it returns [ErrUnsupported].
//
// File decoding is controlled by [Mux]. At least one mapping must be provided.
//
// If the file cannot be read an error is returned. If it cannot be found and
// [IgnoreMissingFile] has been set, the error is omitted and parsing stops.
//
// After decoding, we will go trough all the flags that have not yet been set,
// and try to find their values in the decoded result. Mapping flag names
// to the decoded properties is done using the given [Mapper]. If no [Mapper]
// is provided it'll use [NoopMapper].
//
// If no corresponding value can be found, the flag will remain unset.
// If a value is found it is converted back to a string and fed trough [flag.Value.Set],
// if this fails we will return the error.
// If we are unable to convert the value to a string (for example, if it's an object)
// [ErrVal] will be returned containing the key that failed and the error.
func Parse(path *string, mux Mux, options ...Option) flagr.Parser {
	if path == nil {
		panic("file: path cannot be nil")
	}

	opts := Options{
		Mapper:            NoopMapper,
		IgnoreMissingFile: false,
	}
	for _, opt := range options {
		opt(&opts)
	}

	if len(mux) == 0 {
		panic("file: len(mux) cannot be 0")
	}

	if opts.FS == nil {
		opts.FS = osFS{}
	}

	return func(set *flagr.Set) error {
		f, err := opts.FS.Open(*path)
		if err != nil {
			if errors.Is(err, fs.ErrNotExist) && opts.IgnoreMissingFile {
				return nil
			}
			return fmt.Errorf("file: %w", err)
		}
		defer f.Close()

		data, err := io.ReadAll(f)
		if err != nil {
			return fmt.Errorf("file: %w", err)
		}

		ext := Extension(filepath.Ext(*path))
		decoder, found := mux[ext]
		if !found {
			return ErrUnsupported{
				Ext:       ext,
				Available: mux.supportedExts(),
			}
		}

		var values map[string]any
		if err := decoder(data, &values); err != nil {
			return ErrDecode{err}
		}

		return set.VisitRemaining(func(f *flagr.Flag) error {
			key := opts.Mapper(f.Name)
			wrapper, ok := find(values, key)
			if !ok {
				return nil
			}

			var vals []string
			if err := stringify(wrapper, &vals); err != nil {
				return ErrVal{
					Key: key,
					Err: err,
				}
			}
			for _, val := range vals {
				if err := set.Set("file TODO", f.Name, val); err != nil {
					return err
				}
			}

			return nil
		})
	}
}

// Static is a helper for calling [Parse] with a static path.
func Static(path string) *string { return &path }

var _ fs.FS = osFS{}

type osFS struct{}

func (osFS) Open(name string) (fs.File, error) {
	return os.Open(name)
}

func find(root map[string]any, key KeyPath) (reflect.Value, bool) {
	rv := unwrap(reflect.ValueOf(root))
	for _, segment := range key.Split() {
		rv = unwrap(rv.MapIndex(reflect.ValueOf(segment)))
		if !rv.IsValid() {
			return reflect.Value{}, false
		}
	}
	return rv, true
}

func unwrap(rv reflect.Value) reflect.Value {
	switch rv.Kind() {
	case reflect.Interface, reflect.Pointer:
		return rv.Elem()
	default:
		return rv
	}
}

func stringify(v reflect.Value, values *[]string) error {
	switch v.Kind() {
	case reflect.Bool:
		*values = append(*values, strconv.FormatBool(v.Bool()))
		return nil

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		*values = append(*values, strconv.FormatInt(v.Int(), 10))
		return nil

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		*values = append(*values, strconv.FormatUint(v.Uint(), 10))
		return nil

	case reflect.Float32, reflect.Float64:
		*values = append(*values, strconv.FormatFloat(v.Float(), 'f', -1, 64))
		return nil

	case reflect.Interface, reflect.Pointer:
		return stringify(v.Elem(), values)

	case reflect.Slice:
		len := v.Len()
		for i := 0; i < len; i++ {
			if err := stringify(v.Index(i), values); err != nil {
				return err
			}
		}
		return nil

	case reflect.String:
		*values = append(*values, v.String())
		return nil

	default:
		return fmt.Errorf("unsupported type %q", v.Type().String())
	}
}

// ErrVal is returned when we're unable to convert a value to a string.
type ErrVal struct {
	Key KeyPath
	Err error
}

func (e ErrVal) Error() string {
	return fmt.Sprintf("file: invalid value for path %q: %s", e.Key, e.Err)
}

func (e ErrVal) Unwrap() error {
	return e.Err
}

// ErrUnsupported is returned when we could not find a [DecoderFunc] in the given
// [Mux] with the provided file's extension.
type ErrUnsupported struct {
	Ext       Extension
	Available []Extension
}

func (e ErrUnsupported) Error() string {
	var available strings.Builder
	for i, ext := range e.Available {
		if i > 0 {
			available.WriteString(", ")
		}
		available.WriteString(string(ext))
	}
	return fmt.Sprintf("file: unsupported extension %q, must be one of: %s", e.Ext, available.String())
}

// ErrDecode is returned if deserialization fails.
type ErrDecode struct {
	Err error
}

func (e ErrDecode) Error() string {
	return fmt.Sprintf("file: unable to decode: %s", e.Err.Error())
}

func (e ErrDecode) Unwrap() error {
	return e.Err
}
