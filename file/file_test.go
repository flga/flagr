package file_test

import (
	"encoding/json"
	"errors"
	"io/fs"
	"net/netip"
	"net/url"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/flga/flagr"
	"github.com/flga/flagr/file"
	"github.com/flga/flagr/internal/testflags"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestParserInvariants(t *testing.T) {
	t.Run("panics if path is nil", func(t *testing.T) {
		defer func() {
			got := recover()
			want := "file: path cannot be nil"
			if !reflect.DeepEqual(got, want) {
				t.Fatalf("panic = %v, want %v", got, want)
			}
		}()
		file.Parse(nil, nil)
	})

	t.Run("panics if no mappings provided", func(t *testing.T) {
		defer func() {
			got := recover()
			want := "file: len(mux) cannot be 0"
			if !reflect.DeepEqual(got, want) {
				t.Fatalf("panic = %v, want %v", got, want)
			}
		}()
		file.Parse(new(string), nil)
	})

	t.Run("uses system's filesystem if no FS provided", func(t *testing.T) {
		var set flagr.Set
		flagr.Add(&set, "my-flag-name", flagr.String(""), "")

		err := set.Parse(
			nil,
			file.Parse(
				file.Static("testdata/barebones.json"),
				file.Mux{".json": json.Unmarshal},
			),
		)
		if err != nil {
			t.Fatalf("err = %v", err)
		}
	})

	t.Run("uses identity normalizer if none provided", func(t *testing.T) {
		var set flagr.Set
		val := flagr.Add(&set, "my-flag-name", flagr.String(""), "")

		err := set.Parse(
			nil,
			file.Parse(
				file.Static("testdata/barebones.json"),
				file.Mux{".json": json.Unmarshal},
			),
		)
		if err != nil {
			t.Fatal(err)
		}

		if want := "asd"; *val != want {
			t.Errorf("val = %q, want %q", *val, want)
		}
	})

	t.Run("uses the given normalizer", func(t *testing.T) {
		var set flagr.Set
		val := flagr.Add(&set, "my/flag/name", flagr.String(""), "")

		err := set.Parse(
			nil,
			file.Parse(
				file.Static("testdata/barebones.json"),
				file.Mux{".json": json.Unmarshal},
				file.WithMapper(func(flagName string) file.KeyPath {
					return file.KeyPath(strings.ReplaceAll(flagName, "/", "-"))
				}),
			),
		)
		if err != nil {
			t.Fatal(err)
		}

		if want := "asd"; *val != want {
			t.Errorf("val = %q, want %q", *val, want)
		}
	})

	t.Run("fails if it can't read the given file", func(t *testing.T) {
		var set flagr.Set
		err := set.Parse(
			nil,
			file.Parse(
				file.Static("testdata/______not a file.json"),
				file.Mux{".json": json.Unmarshal},
			),
		)
		if want := fs.ErrNotExist; !errors.Is(err, want) {
			t.Fatalf("err = %v, want %v", err, want)
		}
	})
	t.Run("does not fail if file doesn't exist but ignore missing is true", func(t *testing.T) {
		var set flagr.Set
		err := set.Parse(
			nil,
			file.Parse(
				file.Static("testdata/______not a file.json"),
				file.Mux{".json": json.Unmarshal},
				file.IgnoreMissingFile(),
			),
		)
		if err != nil {
			t.Fatalf("err = %v", err)
		}
	})

	t.Run("fails if no decoder found", func(t *testing.T) {
		var set flagr.Set
		err := set.Parse(
			nil,
			file.Parse(
				file.Static("testdata/barebones.json"),
				file.Mux{".notjson": json.Unmarshal},
			),
		)
		if want := (file.ErrUnsupported{}); !errors.As(err, &want) {
			t.Fatalf("err = %v, want %v", err, want)
		}
	})

	t.Run("fails if it can't decode", func(t *testing.T) {
		var set flagr.Set
		err := set.Parse(
			nil,
			file.Parse(
				file.Static("testdata/invalid_json.json"),
				file.Mux{".json": json.Unmarshal},
			),
		)
		if want := (file.ErrDecode{}); !errors.As(err, &want) {
			t.Fatalf("err = %v, want %v", err, want)
		}
	})

	t.Run("fails if it can't convert a value", func(t *testing.T) {
		var set flagr.Set
		flagr.Add(&set, "my-flag-name", flagr.String(""), "")
		err := set.Parse(
			nil,
			file.Parse(
				file.Static("testdata/invalid_flag_value.json"),
				file.Mux{".json": json.Unmarshal},
			),
		)
		if want := (file.ErrVal{}); !errors.As(err, &want) {
			t.Fatalf("err = %v, want %v", err, want)
		}
	})

	t.Run("fails if it can't set a value", func(t *testing.T) {
		var set flagr.Set
		flagr.Add(&set, "my-flag-name", faultyFlag(), "")
		err := set.Parse(
			nil,
			file.Parse(
				file.Static("testdata/barebones.json"),
				file.Mux{".json": json.Unmarshal},
			),
		)
		if want := errSentinel; !errors.Is(err, want) {
			t.Fatalf("err = %v, want %v", err, want)
		}
	})

	t.Run("only sets a value if not provided before", func(t *testing.T) {
		var set flagr.Set
		val := flagr.Add(&set, "my-flag-name", flagr.String("asd"), "")
		err := set.Parse(
			[]string{
				"-my-flag-name", "123",
			},
			file.Parse(
				file.Static("testdata/barebones.json"),
				file.Mux{".json": json.Unmarshal},
			),
		)
		if err != nil {
			t.Fatalf("err = %v", err)
		}

		if want := "123"; *val != want {
			t.Errorf("val = %q, want %q", *val, want)
		}
	})

	t.Run("ignores unknown properties", func(t *testing.T) {
		var set flagr.Set
		flagr.Add(&set, "this-is-not-in-the-json-file", flagr.String("asd"), "")
		err := set.Parse(
			nil,
			file.Parse(
				file.Static("testdata/barebones.json"),
				file.Mux{".json": json.Unmarshal},
			),
		)
		if err != nil {
			t.Fatalf("err = %v", err)
		}
	})

	t.Run("uses the given FS", func(t *testing.T) {
		var set flagr.Set
		mockFS := &mockFS{}
		err := set.Parse(
			nil,
			file.Parse(
				file.Static("bananas.json"),
				file.Mux{".json": json.Unmarshal},
				file.WithFS(mockFS),
			),
		)
		if want := errMockFS; !errors.Is(err, want) {
			t.Fatalf("err = %v, want %v", err, want)
		}

		if !reflect.DeepEqual(mockFS.calls, []string{"bananas.json"}) {
			t.Fatal("fs wasn't used")
		}
	})

	t.Run("propagates read errors", func(t *testing.T) {
		var set flagr.Set
		err := set.Parse(
			nil,
			file.Parse(
				file.Static("bananas.json"),
				file.Mux{".json": json.Unmarshal},
				file.WithFS(faultyFS{}),
			),
		)
		if want := errFaultyFile; !errors.Is(err, want) {
			t.Fatalf("err = %v, want %v", err, want)
		}
	})
}

func TestFlatJson(t *testing.T) {
	var set flagr.Set
	flags, _ := testflags.Make(&set, "")
	err := set.Parse(
		nil,
		file.Parse(
			file.Static("testdata/flat.json"),
			file.Mux{".json": json.Unmarshal},
		),
	)
	if err != nil {
		t.Fatal(err)
	}
	want := testflags.Flags{
		Int:             ptr(int(10)),
		Ints:            ptr([]int{10, 20}),
		Int8:            ptr(int8(10)),
		Int8s:           ptr([]int8{10, 20}),
		Int16:           ptr(int16(10)),
		Int16s:          ptr([]int16{10, 20}),
		Int32:           ptr(int32(10)),
		Int32s:          ptr([]int32{10, 20}),
		Int64:           ptr(int64(10)),
		Int64s:          ptr([]int64{10, 20}),
		Uint:            ptr(uint(10)),
		Uints:           ptr([]uint{10, 20}),
		Uint8:           ptr(uint8(10)),
		Uint8s:          ptr([]uint8{10, 20}),
		Uint16:          ptr(uint16(10)),
		Uint16s:         ptr([]uint16{10, 20}),
		Uint32:          ptr(uint32(10)),
		Uint32s:         ptr([]uint32{10, 20}),
		Uint64:          ptr(uint64(10)),
		Uint64s:         ptr([]uint64{10, 20}),
		Float32:         ptr(float32(1.0)),
		Float32s:        ptr([]float32{1.0, 2.0}),
		Float64:         ptr(float64(1.0)),
		Float64s:        ptr([]float64{1.0, 2.0}),
		Complex64:       ptr(complex64(1i)),
		Complex64s:      ptr([]complex64{1i, 2i}),
		Complex128:      ptr(complex128(1i)),
		Complex128s:     ptr([]complex128{1i, 2i}),
		Bool:            ptr(false),
		Bools:           ptr([]bool{false, true}),
		String:          ptr("qwe"),
		Strings:         ptr([]string{"qwe", "zxc"}),
		Duration:        ptr(1 * time.Second),
		Durations:       ptr([]time.Duration{1 * time.Second, 2 * time.Second}),
		Time:            ptr(testflags.MustTime("4242-02-25")),
		MustTime:        ptr(testflags.MustTime("4242-02-25")),
		Times:           ptr([]time.Time{testflags.MustTime("4242-02-25"), testflags.MustTime("2000-02-25")}),
		MustTimes:       ptr([]time.Time{testflags.MustTime("4242-02-25"), testflags.MustTime("2000-02-25")}),
		URL:             ptr(testflags.MustURL("https://go.devs")),
		MustURL:         ptr(testflags.MustURL("https://go.devs")),
		URLs:            ptr([]*url.URL{testflags.MustURL("https://go.devs"), testflags.MustURL("https://go.devs/tour/")}),
		MustURLs:        ptr([]*url.URL{testflags.MustURL("https://go.devs"), testflags.MustURL("https://go.devs/tour/")}),
		IPAddr:          ptr(netip.MustParseAddr("127.0.0.2")),
		MustIPAddr:      ptr(netip.MustParseAddr("127.0.0.2")),
		IPAddrs:         ptr([]netip.Addr{netip.MustParseAddr("127.0.0.2"), netip.MustParseAddr("127.0.0.3")}),
		MustIPAddrs:     ptr([]netip.Addr{netip.MustParseAddr("127.0.0.2"), netip.MustParseAddr("127.0.0.3")}),
		IPAddrPort:      ptr(netip.MustParseAddrPort("127.0.0.1:81")),
		MustIPAddrPort:  ptr(netip.MustParseAddrPort("127.0.0.1:81")),
		IPAddrPorts:     ptr([]netip.AddrPort{netip.MustParseAddrPort("127.0.0.1:81"), netip.MustParseAddrPort("127.0.0.1:82")}),
		MustIPAddrPorts: ptr([]netip.AddrPort{netip.MustParseAddrPort("127.0.0.1:81"), netip.MustParseAddrPort("127.0.0.1:82")}),
	}
	if diff := cmp.Diff(want, flags, cmpopts.IgnoreUnexported(netip.Addr{}, netip.AddrPort{})); diff != "" {
		t.Errorf("mismatch (-want +got):\n%s", diff)
	}
}

func TestNestedJson(t *testing.T) {
	var set flagr.Set
	flags, _ := testflags.Make(&set, "foo.bar.baz.")
	err := set.Parse(
		nil,
		file.Parse(
			file.Static("testdata/nested.json"),
			file.Mux{".json": json.Unmarshal},
		),
	)
	if err != nil {
		t.Fatal(err)
	}
	want := testflags.Flags{
		Int:             ptr(int(10)),
		Ints:            ptr([]int{10, 20}),
		Int8:            ptr(int8(10)),
		Int8s:           ptr([]int8{10, 20}),
		Int16:           ptr(int16(10)),
		Int16s:          ptr([]int16{10, 20}),
		Int32:           ptr(int32(10)),
		Int32s:          ptr([]int32{10, 20}),
		Int64:           ptr(int64(10)),
		Int64s:          ptr([]int64{10, 20}),
		Uint:            ptr(uint(10)),
		Uints:           ptr([]uint{10, 20}),
		Uint8:           ptr(uint8(10)),
		Uint8s:          ptr([]uint8{10, 20}),
		Uint16:          ptr(uint16(10)),
		Uint16s:         ptr([]uint16{10, 20}),
		Uint32:          ptr(uint32(10)),
		Uint32s:         ptr([]uint32{10, 20}),
		Uint64:          ptr(uint64(10)),
		Uint64s:         ptr([]uint64{10, 20}),
		Float32:         ptr(float32(1.0)),
		Float32s:        ptr([]float32{1.0, 2.0}),
		Float64:         ptr(float64(1.0)),
		Float64s:        ptr([]float64{1.0, 2.0}),
		Complex64:       ptr(complex64(1i)),
		Complex64s:      ptr([]complex64{1i, 2i}),
		Complex128:      ptr(complex128(1i)),
		Complex128s:     ptr([]complex128{1i, 2i}),
		Bool:            ptr(false),
		Bools:           ptr([]bool{false, true}),
		String:          ptr("qwe"),
		Strings:         ptr([]string{"qwe", "zxc"}),
		Duration:        ptr(1 * time.Second),
		Durations:       ptr([]time.Duration{1 * time.Second, 2 * time.Second}),
		Time:            ptr(testflags.MustTime("4242-02-25")),
		MustTime:        ptr(testflags.MustTime("4242-02-25")),
		Times:           ptr([]time.Time{testflags.MustTime("4242-02-25"), testflags.MustTime("2000-02-25")}),
		MustTimes:       ptr([]time.Time{testflags.MustTime("4242-02-25"), testflags.MustTime("2000-02-25")}),
		URL:             ptr(testflags.MustURL("https://go.devs")),
		MustURL:         ptr(testflags.MustURL("https://go.devs")),
		URLs:            ptr([]*url.URL{testflags.MustURL("https://go.devs"), testflags.MustURL("https://go.devs/tour/")}),
		MustURLs:        ptr([]*url.URL{testflags.MustURL("https://go.devs"), testflags.MustURL("https://go.devs/tour/")}),
		IPAddr:          ptr(netip.MustParseAddr("127.0.0.2")),
		MustIPAddr:      ptr(netip.MustParseAddr("127.0.0.2")),
		IPAddrs:         ptr([]netip.Addr{netip.MustParseAddr("127.0.0.2"), netip.MustParseAddr("127.0.0.3")}),
		MustIPAddrs:     ptr([]netip.Addr{netip.MustParseAddr("127.0.0.2"), netip.MustParseAddr("127.0.0.3")}),
		IPAddrPort:      ptr(netip.MustParseAddrPort("127.0.0.1:81")),
		MustIPAddrPort:  ptr(netip.MustParseAddrPort("127.0.0.1:81")),
		IPAddrPorts:     ptr([]netip.AddrPort{netip.MustParseAddrPort("127.0.0.1:81"), netip.MustParseAddrPort("127.0.0.1:82")}),
		MustIPAddrPorts: ptr([]netip.AddrPort{netip.MustParseAddrPort("127.0.0.1:81"), netip.MustParseAddrPort("127.0.0.1:82")}),
	}
	if diff := cmp.Diff(want, flags, cmpopts.IgnoreUnexported(netip.Addr{}, netip.AddrPort{})); diff != "" {
		t.Errorf("mismatch (-want +got):\n%s", diff)
	}
}

func ptr[T any](t T) *T {
	return &t
}

var errSentinel = errors.New("sentinel")

func faultyFlag() flagr.Getter[string] {
	return flagr.Var("", func(t *string, s string) error {
		return errSentinel
	})
}

var _ fs.FS = &mockFS{}

var errMockFS = errors.New("errMockFS")

type mockFS struct {
	calls []string
}

func (m *mockFS) Open(path string) (fs.File, error) {
	m.calls = append(m.calls, path)
	return nil, errMockFS
}

type faultyFS struct{}

func (faultyFS) Open(path string) (fs.File, error) {
	return faultyFile{}, nil
}

var errFaultyFile = errors.New("errFaultyFile")

type faultyFile struct{}

func (faultyFile) Stat() (fs.FileInfo, error) {
	panic("unsupported")
}
func (faultyFile) Read([]byte) (int, error) {
	return 0, errFaultyFile
}
func (faultyFile) Close() error {
	return nil
}
