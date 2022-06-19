// Package env provides a flagr.Parser that is able to read environment variables (and .env files) and set the appropriate flags.
package env

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"regexp"
	"strings"

	"github.com/flga/flagr"
	"github.com/hashicorp/go-envparse"
)

type Splitter string

// NoSplit disables value splitting.
const NoSplit Splitter = ""

// Mapper maps a flag to the corresponding env var.
type Mapper func(flagName string) (envName string, listSplitter Splitter)

var defaultMapperRegex = regexp.MustCompile("[^a-zA-Z0-9_]")

// DefaultMapper converts flags to env vars by replacing any non alphanumeric
// character (in the ascii sense) with an underscore.
//
// If splitter is not NoSplit (an empty string), it will be used to
// try to split env values to lists. If you need to control the splitter per
// flag, use a custom Mapper.
func DefaultMapper(splitter Splitter) Mapper {
	return func(flagName string) (envName string, listSplitter Splitter) {
		return strings.ToUpper(defaultMapperRegex.ReplaceAllLiteralString(flagName, "_")), splitter
	}
}

// LookupFunc returns the value for the given variable and whether it was found.
type LookupFunc func(varName string) (string, bool)

type options struct {
	prefix          string
	mapper          Mapper
	lookupFunc      LookupFunc
	envFile         *string
	envFileOptional bool
}

type Option func(*options)

// WithLookupFunc replaces the default lookup method (os.LookupEnv) with the given func.
func WithLookupFunc(fn LookupFunc) Option {
	return func(o *options) {
		o.lookupFunc = fn
	}
}

// WithDotEnv tells the parser to also parse the given .env file.
// Env vars take precedence over anything defined in it.
//
// Parsing will fail if the file does not exist, unless ignoreMissing is true.
//
// Path will be resolved just in time, so it may be set by other parsers up in the chain.
func WithDotEnv(path *string, optional bool) Option {
	return func(o *options) {
		o.envFile = path
		o.envFileOptional = optional
	}
}

// WithStaticDotEnv is an alias for [WithDotEnv] but with a static path.
func WithStaticDotEnv(path string, optional bool) Option {
	return WithDotEnv(&path, optional)
}

// WithMapper tells the parser how to map flags to env vars and, if the value
// is expected to be a list, how to split it.
func WithMapper(fn Mapper) Option {
	return func(o *options) {
		o.mapper = fn
	}
}

// WithPrefix prefixes every flag with s before mapping it to the corresponding env var.
// The prefix need not end in an underscore as one will be added automatically.
func WithPrefix(s string) Option {
	return func(o *options) {
		o.prefix = s + "_"
	}
}

func Parse(opts ...Option) flagr.Parser {
	options := options{
		prefix:     "",
		mapper:     DefaultMapper(""),
		lookupFunc: os.LookupEnv,
		envFile:    nil,
	}
	for _, opt := range opts {
		opt(&options)
	}

	return func(fs *flagr.Set) error {
		var fileData map[string]string
		if options.envFile != nil {
			fd, err := maybeParseEnvFile(*options.envFile, options.envFileOptional)
			if err != nil {
				return fmt.Errorf("env: unable to parse env file %q: %w", *options.envFile, err)
			}
			fileData = fd
		}

		return fs.VisitRemaining(func(flag *flagr.Flag) error {
			name, splitValBy := options.mapper(options.prefix + flag.Name)
			src := flagr.Source("env: " + name)
			val, ok := options.lookupFunc(name)
			if !ok {
				val, ok = fileData[name]
				if options.envFile != nil {
					src = flagr.Source(fmt.Sprintf("envfile[%s]: %s", *options.envFile, name))
				}
			}
			if !ok {
				return nil
			}

			switch {
			case splitValBy != "":
				for _, val := range strings.Split(val, string(splitValBy)) {
					if err := fs.Set(src, flag.Name, val); err != nil {
						return fmt.Errorf("env: %w", err)
					}
				}
				return nil

			default:
				if err := fs.Set(src, flag.Name, val); err != nil {
					return fmt.Errorf("env: %w", err)
				}
			}

			return nil
		})
	}
}

func maybeParseEnvFile(path string, ignoreMissing bool) (map[string]string, error) {
	f, err := os.Open(path)
	if err != nil {
		if ignoreMissing && errors.Is(err, fs.ErrNotExist) {
			return nil, nil
		}
		return nil, err
	}
	defer f.Close()

	return envparse.Parse(f)
}
