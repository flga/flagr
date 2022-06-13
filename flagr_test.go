package flagr_test

import (
	"bytes"
	"errors"
	"flag"
	"io/ioutil"
	"net/netip"
	"net/url"
	"reflect"
	"strconv"
	"testing"
	"time"

	"github.com/flga/flagr"
	"github.com/flga/flagr/internal/testflags"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestNewSet(t *testing.T) {
	t.Run("new", func(t *testing.T) {
		name := "a"
		strat := flag.PanicOnError
		s := flagr.NewSet(name, strat)
		if got := s.Name(); name != got {
			t.Errorf("name = %q, want %q", got, name)
		}
		if got := s.ErrorHandling(); strat != got {
			t.Errorf("strat = %q, want %q", got, strat)
		}
	})

	t.Run("init", func(t *testing.T) {
		var s flagr.Set
		name := "a"
		strat := flag.PanicOnError
		s.Init(name, strat)
		if got := s.Name(); name != got {
			t.Errorf("name = %q, want %q", got, name)
		}
		if got := s.ErrorHandling(); strat != got {
			t.Errorf("strat = %q, want %q", got, strat)
		}
	})
}

func TestUsage(t *testing.T) {
	var set flagr.Set
	a := 0
	set.SetUsage(func() { a++ })
	set.Usage()

	if a != 1 {
		t.Fatal()
	}
}

func TestStdGetterSetters(t *testing.T) {
	t.Run("output", func(t *testing.T) {
		var set flagr.Set
		want := ioutil.Discard
		set.SetOutput(want)
		if got := set.Output(); got != want {
			t.Errorf("Output = %v, want %v", got, want)
		}
	})

	t.Run("unquote", func(t *testing.T) {
		var set flagr.Set
		flagr.Add(&set, "a", flagr.Int(0), "a `name` to show")
		if err := set.Parse(nil); err != nil {
			t.Fatal(err)
		}

		name, usage := flagr.UnquoteUsage(set.Lookup("a"))
		if want := "name"; name != want {
			t.Errorf("UnquoteUsage() name = %q, want %q", name, want)
		}
		if want := "a name to show"; usage != want {
			t.Errorf("UnquoteUsage() usage = %q, want %q", usage, want)
		}
	})

	t.Run("usage/prints", func(t *testing.T) {
		t.Run("with name", func(t *testing.T) {
			var set flagr.Set
			set.Init("mycmd", flag.ContinueOnError) // skips stdflag.defaultUsage
			flagr.Add(&set, "a", flagr.Int(0), "a `name` to show")
			flagr.Add(&set, "b", flagr.Int(0), "")
			flagr.Add(&set, "c", flagr.Int(0), "")
			flagr.Add(&set, "d", flagr.Int(0), "")
			flagr.Add(&set, "e", flagr.Int(0), "")
			if err := set.Parse([]string{"-e", "1"}); err != nil {
				t.Fatal(err)
			}

			if err := set.Set("env: APP_B_VAL", "b", "1"); err != nil {
				t.Fatal(err)
			}
			if err := set.Set("file: config.yml", "c", "2"); err != nil {
				t.Fatal(err)
			}
			if err := set.Set("", "d", "not a number"); !errors.As(err, new(*strconv.NumError)) {
				t.Fatal(err)
			}

			var output bytes.Buffer
			set.SetOutput(&output)

			// usage
			set.Usage()
			wantUsage := `Usage of mycmd:
  -a name
    	a name to show (default 0)
  -b	 (default 0)
  -c	 (default 0)
  -d	 (default 0)
  -e	 (default 0)
`
			if diff := cmp.Diff(wantUsage, output.String()); diff != "" {
				t.Errorf("usage mismatch (-want +got):\n%s", diff)
			}

			// defaults
			output.Reset()
			set.PrintDefaults()
			wantDefaults := `  -a name
    	a name to show (default 0)
  -b	 (default 0)
  -c	 (default 0)
  -d	 (default 0)
  -e	 (default 0)
`
			if diff := cmp.Diff(wantDefaults, output.String()); diff != "" {
				t.Errorf("defaults mismatch (-want +got):\n%s", diff)
			}

			// values
			output.Reset()
			set.PrintValues()
			wantValues := `Current configuration of mycmd:
  -a 0 (default)
  -b 1 (env: APP_B_VAL)
  -c 2 (file: config.yml)
  -d 0 (default)
  -e 1 (flags)
`
			if diff := cmp.Diff(wantValues, output.String()); diff != "" {
				t.Errorf("values mismatch (-want +got):\n%s", diff)
			}
		})

		t.Run("without name", func(t *testing.T) {
			var set flagr.Set
			flagr.Add(&set, "a", flagr.Int(0), "a `name` to show")
			flagr.Add(&set, "b", flagr.Int(0), "")
			flagr.Add(&set, "c", flagr.Int(0), "")
			flagr.Add(&set, "d", flagr.Int(0), "")
			flagr.Add(&set, "e", flagr.Int(0), "")
			if err := set.Parse([]string{"-e", "1"}); err != nil {
				t.Fatal(err)
			}

			if err := set.Set("env: APP_B_VAL", "b", "1"); err != nil {
				t.Fatal(err)
			}
			if err := set.Set("file: config.yml", "c", "2"); err != nil {
				t.Fatal(err)
			}
			if err := set.Set("", "d", "not a number"); !errors.As(err, new(*strconv.NumError)) {
				t.Fatal(err)
			}

			var output bytes.Buffer
			set.SetOutput(&output)

			// usage
			set.Usage()
			wantUsage := `Usage:
  -a name
    	a name to show (default 0)
  -b	 (default 0)
  -c	 (default 0)
  -d	 (default 0)
  -e	 (default 0)
`
			if diff := cmp.Diff(wantUsage, output.String()); diff != "" {
				t.Errorf("usage mismatch (-want +got):\n%s", diff)
			}

			// defaults
			output.Reset()
			set.PrintDefaults()
			wantDefaults := `  -a name
    	a name to show (default 0)
  -b	 (default 0)
  -c	 (default 0)
  -d	 (default 0)
  -e	 (default 0)
`
			if diff := cmp.Diff(wantDefaults, output.String()); diff != "" {
				t.Errorf("defaults mismatch (-want +got):\n%s", diff)
			}

			// values
			output.Reset()
			set.PrintValues()
			wantValues := `Current configuration:
  -a 0 (default)
  -b 1 (env: APP_B_VAL)
  -c 2 (file: config.yml)
  -d 0 (default)
  -e 1 (flags)
`
			if diff := cmp.Diff(wantValues, output.String()); diff != "" {
				t.Errorf("values mismatch (-want +got):\n%s", diff)
			}
		})
	})

	t.Run("args", func(t *testing.T) {
		var set flagr.Set
		flagr.Add(&set, "a", flagr.Int(0), "")
		flagr.Add(&set, "b", flagr.Int(0), "")
		flagr.Add(&set, "c", flagr.Int(0), "")
		flagr.Add(&set, "d", flagr.Int(0), "")
		if err := set.Parse([]string{
			"-a", "2",
			"foo",
			"bar",
		}); err != nil {
			t.Fatal(err)
		}

		if err := set.Set("", "b", "42"); err != nil {
			t.Fatal(err)
		}
		if err := set.Set("", "d", "not a number"); !errors.As(err, new(*strconv.NumError)) {
			t.Fatal(err)
		}

		if want, got := "42", set.Lookup("b").Value.String(); got != want {
			t.Errorf("Lookup() = %v, want %v", got, want)
		}
		if want, got := 2, set.NFlag(); got != want {
			t.Errorf("NFlag() = %v, want %v", got, want)
		}
		if want, got := "bar", set.Arg(1); got != want {
			t.Errorf("Arg() = %v, want %v", got, want)
		}
		if want, got := 2, set.NArg(); got != want {
			t.Errorf("NArg() = %v, want %v", got, want)
		}
		if want, got := []string{"foo", "bar"}, set.Args(); !reflect.DeepEqual(got, want) {
			t.Errorf("Args() = %v, want %v", got, want)
		}
	})
}

func TestWalks(t *testing.T) {
	var set flagr.Set
	flagr.Add(&set, "a", flagr.Int(0), "default")
	flagr.Add(&set, "b", flagr.Int(0), "set from cli")
	flagr.Add(&set, "c", flagr.Int(0), "set manually")
	flagr.Add(&set, "d", flagr.Int(0), "default")
	if err := set.Parse([]string{
		"-b", "1",
	}); err != nil {
		t.Fatal(err)
	}

	if err := set.Set("", "c", "2"); err != nil {
		t.Fatal(err)
	}

	t.Run("visits the right things", func(t *testing.T) {
		failif := func(err error) {
			if err != nil {
				t.Fatal(err)
			}
		}

		var allFlags, setFlags, remainingFlags []string
		failif(set.VisitAll(func(f *flagr.Flag) error {
			allFlags = append(allFlags, f.Name)
			return nil
		}))
		failif(set.Visit(func(f *flagr.Flag) error {
			setFlags = append(setFlags, f.Name)
			return nil
		}))
		failif(set.VisitRemaining(func(f *flagr.Flag) error {
			remainingFlags = append(remainingFlags, f.Name)
			return nil
		}))

		wantAll := []string{"a", "b", "c", "d"}
		wantSet := []string{"b", "c"}
		wantRemaining := []string{"a", "d"}
		if !reflect.DeepEqual(wantAll, allFlags) {
			t.Errorf("VisitAll() = %v, want %v", allFlags, wantAll)
		}
		if !reflect.DeepEqual(wantSet, setFlags) {
			t.Errorf("Visit() = %v, want %v", setFlags, wantSet)
		}
		if !reflect.DeepEqual(wantRemaining, remainingFlags) {
			t.Errorf("VisitRemaining() = %v, want %v", remainingFlags, wantRemaining)
		}
	})

	t.Run("stops walking if there's an error", func(t *testing.T) {
		sentinel := errors.New("sentinel")

		var allFlags, setFlags, remainingFlags []string
		allErr := set.VisitAll(func(f *flagr.Flag) error {
			allFlags = append(allFlags, f.Name)
			return sentinel
		})
		setErr := set.Visit(func(f *flagr.Flag) error {
			setFlags = append(setFlags, f.Name)
			return sentinel
		})
		remainingErr := set.VisitRemaining(func(f *flagr.Flag) error {
			remainingFlags = append(remainingFlags, f.Name)
			return sentinel
		})

		if allErr != sentinel {
			t.Errorf("VisitAll() err = %v, want %v", allErr, sentinel)
		}
		if setErr != sentinel {
			t.Errorf("Visit() err = %v, want %v", setErr, sentinel)
		}
		if remainingErr != sentinel {
			t.Errorf("VisitRemaining() err = %v, want %v", remainingErr, sentinel)
		}

		wantAll := []string{"a"}
		wantSet := []string{"b"}
		wantRemaining := []string{"a"}
		if !reflect.DeepEqual(wantAll, allFlags) {
			t.Errorf("VisitAll() = %v, want %v", allFlags, wantAll)
		}
		if !reflect.DeepEqual(wantSet, setFlags) {
			t.Errorf("Visit() = %v, want %v", setFlags, wantSet)
		}
		if !reflect.DeepEqual(wantRemaining, remainingFlags) {
			t.Errorf("VisitRemaining() = %v, want %v", remainingFlags, wantRemaining)
		}
	})
}

func TestParse(t *testing.T) {
	t.Run("calls extra parsers in the order they're defined", func(t *testing.T) {
		type pass struct {
			id   string
			seen []string
		}

		testParser := func(parserID string, flagToSet string, out *[]pass) flagr.Parser {
			return func(set *flagr.Set) error {
				var seen []string
				err := set.VisitRemaining(func(f *flagr.Flag) error {
					seen = append(seen, f.Name)
					if f.Name == flagToSet {
						set.Set(flagr.Source(parserID), f.Name, "42")
					}
					return nil
				})
				if err != nil {
					return err
				}

				*out = append(*out, pass{
					id:   parserID,
					seen: seen,
				})

				return nil
			}
		}

		var set flagr.Set
		flagr.Add(&set, "a", flagr.Int(0), "default")
		flagr.Add(&set, "b", flagr.Int(0), "cli args")
		flagr.Add(&set, "c", flagr.Int(0), "parser 1")
		flagr.Add(&set, "d", flagr.Int(0), "parser 2")

		var passes []pass
		if err := set.Parse(
			[]string{"-b", "1"},
			testParser("parser 1", "c", &passes),
			testParser("parser 2", "d", &passes),
		); err != nil {
			t.Fatal(err)
		}

		if !set.Parsed() {
			t.Errorf("Parsed() should be true")
		}

		want := []pass{
			{id: "parser 1", seen: []string{"a", "c", "d"}},
			{id: "parser 2", seen: []string{"a", "d"}},
		}
		if diff := cmp.Diff(want, passes, cmp.AllowUnexported(pass{})); diff != "" {
			t.Errorf("mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("if flag parsing fails, parsing is interrupted", func(t *testing.T) {
		var goodCalls int
		goodParser := func(set *flagr.Set) error {
			goodCalls++
			return set.VisitRemaining(func(f *flagr.Flag) error {
				return nil
			})
		}

		var set flagr.Set
		flagr.Add(&set, "a", flagr.Int(0), "")

		err := set.Parse(
			[]string{"-imnotaflag"},
			goodParser,
		)
		// std/flag returns a errors.New
		if want := "flag provided but not defined: -imnotaflag"; err.Error() != want {
			t.Fatalf("err = %v, want %q", err, want)
		}

		if !set.Parsed() {
			t.Errorf("Parsed() should be true")
		}

		if want := 0; goodCalls != want {
			t.Errorf("goodCalls = %d, want %d", goodCalls, want)
		}
	})

	t.Run("if a parser fails, parsing is interrupted", func(t *testing.T) {
		sentinel := errors.New("sentinel")
		faultyParser := func(set *flagr.Set) error {
			return set.VisitRemaining(func(f *flagr.Flag) error {
				return sentinel
			})
		}

		var goodCalls int
		goodParser := func(set *flagr.Set) error {
			goodCalls++
			return set.VisitRemaining(func(f *flagr.Flag) error {
				return nil
			})
		}

		var set flagr.Set
		flagr.Add(&set, "a", flagr.Int(0), "")

		if err := set.Parse(
			nil,
			goodParser,
			faultyParser,
			goodParser,
		); err != sentinel {
			t.Fatalf("err = %v, want %v", err, sentinel)
		}

		if !set.Parsed() {
			t.Errorf("Parsed() should be true")
		}

		if want := 1; goodCalls != want {
			t.Errorf("goodCalls = %d, want %d", goodCalls, want)
		}
	})
}

func TestDefaults(t *testing.T) {
	var s flagr.Set
	vals, defaults := testflags.Make(&s)
	if err := s.Parse(nil); err != nil {
		t.Fatal(err)
	}

	want := testflags.Flags{
		Int:             ptr(defaults.Int),
		Ints:            ptr(defaults.Ints),
		Int8:            ptr(defaults.Int8),
		Int8s:           ptr(defaults.Int8s),
		Int16:           ptr(defaults.Int16),
		Int16s:          ptr(defaults.Int16s),
		Int32:           ptr(defaults.Int32),
		Int32s:          ptr(defaults.Int32s),
		Int64:           ptr(defaults.Int64),
		Int64s:          ptr(defaults.Int64s),
		Uint:            ptr(defaults.Uint),
		Uints:           ptr(defaults.Uints),
		Uint8:           ptr(defaults.Uint8),
		Uint8s:          ptr(defaults.Uint8s),
		Uint16:          ptr(defaults.Uint16),
		Uint16s:         ptr(defaults.Uint16s),
		Uint32:          ptr(defaults.Uint32),
		Uint32s:         ptr(defaults.Uint32s),
		Uint64:          ptr(defaults.Uint64),
		Uint64s:         ptr(defaults.Uint64s),
		Float32:         ptr(defaults.Float32),
		Float32s:        ptr(defaults.Float32s),
		Float64:         ptr(defaults.Float64),
		Float64s:        ptr(defaults.Float64s),
		Complex64:       ptr(defaults.Complex64),
		Complex64s:      ptr(defaults.Complex64s),
		Complex128:      ptr(defaults.Complex128),
		Complex128s:     ptr(defaults.Complex128s),
		Bool:            ptr(defaults.Bool),
		Bools:           ptr(defaults.Bools),
		String:          ptr(defaults.String),
		Strings:         ptr(defaults.Strings),
		Duration:        ptr(defaults.Duration),
		Durations:       ptr(defaults.Durations),
		Time:            ptr(defaults.Time),
		MustTime:        ptr(defaults.Time),
		Times:           ptr(defaults.Times),
		MustTimes:       ptr(defaults.Times),
		URL:             ptr(defaults.URL),
		MustURL:         ptr(defaults.URL),
		URLs:            ptr(defaults.URLs),
		MustURLs:        ptr(defaults.URLs),
		IPAddr:          ptr(defaults.IPAddr),
		MustIPAddr:      ptr(defaults.IPAddr),
		IPAddrs:         ptr(defaults.IPAddrs),
		MustIPAddrs:     ptr(defaults.IPAddrs),
		IPAddrPort:      ptr(defaults.IPAddrPort),
		MustIPAddrPort:  ptr(defaults.IPAddrPort),
		IPAddrPorts:     ptr(defaults.IPAddrPorts),
		MustIPAddrPorts: ptr(defaults.IPAddrPorts),
	}
	if diff := cmp.Diff(want, vals, cmpopts.IgnoreUnexported(netip.Addr{}, netip.AddrPort{})); diff != "" {
		t.Errorf("mismatch (-want +got):\n%s", diff)
	}
}

func TestSetters(t *testing.T) {
	var s flagr.Set
	vals, _ := testflags.Make(&s)
	args := []string{
		"-a00", "1",
		"-a01", "1", "-a01", "2", "-a01", "3",
		"-a02", "1",
		"-a03", "1", "-a03", "2", "-a03", "3",
		"-a04", "1",
		"-a05", "1", "-a05", "2", "-a05", "3",
		"-a06", "1",
		"-a07", "1", "-a07", "2", "-a07", "3",
		"-a08", "1",
		"-a09", "1", "-a09", "2", "-a09", "3",
		"-a10", "1",
		"-a11", "1", "-a11", "2", "-a11", "3",
		"-a12", "1",
		"-a13", "1", "-a13", "2", "-a13", "3",
		"-a14", "1",
		"-a15", "1", "-a15", "2", "-a15", "3",
		"-a16", "1",
		"-a17", "1", "-a17", "2", "-a17", "3",
		"-a18", "1",
		"-a19", "1", "-a19", "2", "-a19", "3",
		"-a20", "1",
		"-a21", "1", "-a21", "2", "-a21", "3",
		"-a22", "1",
		"-a23", "1", "-a23", "2", "-a23", "3",
		"-a24", "1i",
		"-a25", "1i", "-a25", "2i",
		"-a26", "1i",
		"-a27", "1i", "-a27", "2i",
		"-a28=false",
		"-a29=false", "-a29=true", "-a29=false",
		"-a30", "qwe",
		"-a31", "qwe", "-a31", "rty", "-a31", "uio",
		"-a32", "1s",
		"-a33", "1s", "-a33", "2s", "-a33", "3s",
		"-a34", "0000-01-01",
		"-a35", "0000-01-01",
		"-a36", "0000-01-01", "-a36", "0000-01-02", "-a36", "0000-01-03",
		"-a37", "0000-01-01", "-a37", "0000-01-02", "-a37", "0000-01-03",
		"-a38", "https://a.com",
		"-a39", "https://a.com",
		"-a40", "https://a.com", "-a40", "https://b.com", "-a40", "https://c.com",
		"-a41", "https://a.com", "-a41", "https://b.com", "-a41", "https://c.com",
		"-a42", "0.0.0.0",
		"-a43", "0.0.0.0",
		"-a44", "0.0.0.0", "-a44", "0.0.0.1", "-a44", "0.0.0.2",
		"-a45", "0.0.0.0", "-a45", "0.0.0.1", "-a45", "0.0.0.2",
		"-a46", "0.0.0.0:80",
		"-a47", "0.0.0.0:80",
		"-a48", "0.0.0.0:80", "-a48", "0.0.0.0:81", "-a48", "0.0.0.0:82",
		"-a49", "0.0.0.0:80", "-a49", "0.0.0.0:81", "-a49", "0.0.0.0:82",
	}
	if err := s.Parse(args); err != nil {
		t.Fatal(err)
	}

	want := testflags.Flags{
		Int:             ptr(int(1)),
		Ints:            ptr([]int{1, 2, 3}),
		Int8:            ptr(int8(1)),
		Int8s:           ptr([]int8{1, 2, 3}),
		Int16:           ptr(int16(1)),
		Int16s:          ptr([]int16{1, 2, 3}),
		Int32:           ptr(int32(1)),
		Int32s:          ptr([]int32{1, 2, 3}),
		Int64:           ptr(int64(1)),
		Int64s:          ptr([]int64{1, 2, 3}),
		Uint:            ptr(uint(1)),
		Uints:           ptr([]uint{1, 2, 3}),
		Uint8:           ptr(uint8(1)),
		Uint8s:          ptr([]uint8{1, 2, 3}),
		Uint16:          ptr(uint16(1)),
		Uint16s:         ptr([]uint16{1, 2, 3}),
		Uint32:          ptr(uint32(1)),
		Uint32s:         ptr([]uint32{1, 2, 3}),
		Uint64:          ptr(uint64(1)),
		Uint64s:         ptr([]uint64{1, 2, 3}),
		Float32:         ptr(float32(1)),
		Float32s:        ptr([]float32{1, 2, 3}),
		Float64:         ptr(float64(1)),
		Float64s:        ptr([]float64{1, 2, 3}),
		Complex64:       ptr(complex64(1i)),
		Complex64s:      ptr([]complex64{1i, 2i}),
		Complex128:      ptr(complex128(1i)),
		Complex128s:     ptr([]complex128{1i, 2i}),
		Bool:            ptr(false),
		Bools:           ptr([]bool{false, true, false}),
		String:          ptr("qwe"),
		Strings:         ptr([]string{"qwe", "rty", "uio"}),
		Duration:        ptr(1 * time.Second),
		Durations:       ptr([]time.Duration{1 * time.Second, 2 * time.Second, 3 * time.Second}),
		Time:            ptr(testflags.MustTime("0000-01-01")),
		MustTime:        ptr(testflags.MustTime("0000-01-01")),
		Times:           ptr([]time.Time{testflags.MustTime("0000-01-01"), testflags.MustTime("0000-01-02"), testflags.MustTime("0000-01-03")}),
		MustTimes:       ptr([]time.Time{testflags.MustTime("0000-01-01"), testflags.MustTime("0000-01-02"), testflags.MustTime("0000-01-03")}),
		URL:             ptr(testflags.MustURL("https://a.com")),
		MustURL:         ptr(testflags.MustURL("https://a.com")),
		URLs:            ptr([]*url.URL{testflags.MustURL("https://a.com"), testflags.MustURL("https://b.com"), testflags.MustURL("https://c.com")}),
		MustURLs:        ptr([]*url.URL{testflags.MustURL("https://a.com"), testflags.MustURL("https://b.com"), testflags.MustURL("https://c.com")}),
		IPAddr:          ptr(netip.MustParseAddr("0.0.0.0")),
		MustIPAddr:      ptr(netip.MustParseAddr("0.0.0.0")),
		IPAddrs:         ptr([]netip.Addr{netip.MustParseAddr("0.0.0.0"), netip.MustParseAddr("0.0.0.1"), netip.MustParseAddr("0.0.0.2")}),
		MustIPAddrs:     ptr([]netip.Addr{netip.MustParseAddr("0.0.0.0"), netip.MustParseAddr("0.0.0.1"), netip.MustParseAddr("0.0.0.2")}),
		IPAddrPort:      ptr(netip.MustParseAddrPort("0.0.0.0:80")),
		MustIPAddrPort:  ptr(netip.MustParseAddrPort("0.0.0.0:80")),
		IPAddrPorts:     ptr([]netip.AddrPort{netip.MustParseAddrPort("0.0.0.0:80"), netip.MustParseAddrPort("0.0.0.0:81"), netip.MustParseAddrPort("0.0.0.0:82")}),
		MustIPAddrPorts: ptr([]netip.AddrPort{netip.MustParseAddrPort("0.0.0.0:80"), netip.MustParseAddrPort("0.0.0.0:81"), netip.MustParseAddrPort("0.0.0.0:82")}),
	}
	if diff := cmp.Diff(want, vals, cmpopts.IgnoreUnexported(netip.Addr{}, netip.AddrPort{})); diff != "" {
		t.Errorf("mismatch (-want +got):\n%s", diff)
	}
}

func TestDoesntClobberDefault(t *testing.T) {
	strings := []string{"a", "b", "c"}
	urlv := &url.URL{Host: "foo.com"}
	time := testflags.MustTime("2000-01-01")

	fs := flagr.NewSet("", flagr.ContinueOnError)
	flagr.Add(fs, "a", flagr.Strings(strings...), "")
	flagr.Add(fs, "b", flagr.URL(urlv), "")
	flagr.Add(fs, "c", flagr.Time(testflags.TimeLayout, testflags.MustTime("0001-01-01")), "")
	err := fs.Parse([]string{"-a", "1", "-b", "https://username:password@bar.com"})
	if err != nil {
		t.Fatal(err)
	}

	if want := []string{"a", "b", "c"}; !reflect.DeepEqual(strings, want) {
		t.Fatalf("want %v, got %v", want, strings)
	}
	if want := (&url.URL{Host: "foo.com"}); !reflect.DeepEqual(urlv, want) {
		t.Fatalf("want %v, got %v", want, urlv)
	}
	if want := testflags.MustTime("2000-01-01"); !reflect.DeepEqual(time, want) {
		t.Fatalf("want %v, got %v", want, urlv)
	}
}

func ptr[T any](t T) *T {
	return &t
}
