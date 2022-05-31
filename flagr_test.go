package flagr_test

import (
	"flag"
	"net/netip"
	"net/url"
	"reflect"
	"testing"
	"time"

	"github.com/flga/flagr"
)

const timeLayout = "2006-01-02"

type flags struct {
	Int             *int
	Ints            *[]int
	Int8            *int8
	Int8s           *[]int8
	Int16           *int16
	Int16s          *[]int16
	Int32           *int32
	Int32s          *[]int32
	Int64           *int64
	Int64s          *[]int64
	Uint            *uint
	Uints           *[]uint
	Uint8           *uint8
	Uint8s          *[]uint8
	Uint16          *uint16
	Uint16s         *[]uint16
	Uint32          *uint32
	Uint32s         *[]uint32
	Uint64          *uint64
	Uint64s         *[]uint64
	Float32         *float32
	Float32s        *[]float32
	Float64         *float64
	Float64s        *[]float64
	Complex64       *complex64
	Complex64s      *[]complex64
	Complex128      *complex128
	Complex128s     *[]complex128
	Bool            *bool
	Bools           *[]bool
	String          *string
	Strings         *[]string
	Duration        *time.Duration
	Durations       *[]time.Duration
	Time            *time.Time
	MustTime        *time.Time
	Times           *[]time.Time
	MustTimes       *[]time.Time
	URL             **url.URL
	MustURL         **url.URL
	URLs            *[]*url.URL
	MustURLs        *[]*url.URL
	IPAddr          *netip.Addr
	MustIPAddr      *netip.Addr
	IPAddrs         *[]netip.Addr
	MustIPAddrs     *[]netip.Addr
	IPAddrPort      *netip.AddrPort
	MustIPAddrPort  *netip.AddrPort
	IPAddrPorts     *[]netip.AddrPort
	MustIPAddrPorts *[]netip.AddrPort
}

var defaults = struct {
	Int             int
	Ints            []int
	Int8            int8
	Int8s           []int8
	Int16           int16
	Int16s          []int16
	Int32           int32
	Int32s          []int32
	Int64           int64
	Int64s          []int64
	Uint            uint
	Uints           []uint
	Uint8           uint8
	Uint8s          []uint8
	Uint16          uint16
	Uint16s         []uint16
	Uint32          uint32
	Uint32s         []uint32
	Uint64          uint64
	Uint64s         []uint64
	Float32         float32
	Float32s        []float32
	Float64         float64
	Float64s        []float64
	Complex64       complex64
	Complex64s      []complex64
	Complex128      complex128
	Complex128s     []complex128
	Bool            bool
	Bools           []bool
	String          string
	Strings         []string
	Duration        time.Duration
	Durations       []time.Duration
	Time            time.Time
	MustTime        string
	Times           []time.Time
	MustTimes       []string
	URL             *url.URL
	MustURL         string
	URLs            []*url.URL
	MustURLs        []string
	IPAddr          netip.Addr
	MustIPAddr      string
	IPAddrs         []netip.Addr
	MustIPAddrs     []string
	IPAddrPort      netip.AddrPort
	MustIPAddrPort  string
	IPAddrPorts     []netip.AddrPort
	MustIPAddrPorts []string
}{
	Int:             42,
	Ints:            []int{42, 24},
	Int8:            42,
	Int8s:           []int8{42, 24},
	Int16:           42,
	Int16s:          []int16{42, 24},
	Int32:           42,
	Int32s:          []int32{42, 24},
	Int64:           42,
	Int64s:          []int64{42, 24},
	Uint:            42,
	Uints:           []uint{42, 24},
	Uint8:           42,
	Uint8s:          []uint8{42, 24},
	Uint16:          42,
	Uint16s:         []uint16{42, 24},
	Uint32:          42,
	Uint32s:         []uint32{42, 24},
	Uint64:          42,
	Uint64s:         []uint64{42, 24},
	Float32:         4.2,
	Float32s:        []float32{4.2, 2.4},
	Float64:         4.2,
	Float64s:        []float64{4.2, 2.4},
	Complex64:       42i,
	Complex64s:      []complex64{42i, 24i},
	Complex128:      42i,
	Complex128s:     []complex128{42i, 24i},
	Bool:            true,
	Bools:           []bool{true, false},
	String:          "asd",
	Strings:         []string{"asd", "dsa"},
	Duration:        42 * time.Second,
	Durations:       []time.Duration{42 * time.Second, 24 * time.Second},
	Time:            musttime("4242-02-24"),
	MustTime:        "4242-02-24",
	Times:           []time.Time{musttime("4242-02-24"), musttime("2000-02-24")},
	MustTimes:       []string{"4242-02-24", "2000-02-24"},
	URL:             musturl("https://go.dev"),
	MustURL:         "https://go.dev",
	URLs:            []*url.URL{musturl("https://go.dev"), musturl("https://go.dev/tour/")},
	MustURLs:        []string{"https://go.dev", "https://go.dev/tour/"},
	IPAddr:          netip.MustParseAddr("127.0.0.1"),
	MustIPAddr:      "127.0.0.1",
	IPAddrs:         []netip.Addr{netip.MustParseAddr("127.0.0.1"), netip.MustParseAddr("127.0.0.2")},
	MustIPAddrs:     []string{"127.0.0.1", "127.0.0.2"},
	IPAddrPort:      netip.MustParseAddrPort("127.0.0.1:80"),
	MustIPAddrPort:  "127.0.0.1:80",
	IPAddrPorts:     []netip.AddrPort{netip.MustParseAddrPort("127.0.0.1:80"), netip.MustParseAddrPort("127.0.0.1:81")},
	MustIPAddrPorts: []string{"127.0.0.1:80", "127.0.0.1:81"},
}

func makeFlags(s *flagr.Set) flags {
	var vals flags
	vals.Int = flagr.Add(s, "a00", flagr.Int(defaults.Int), "usage for a00")
	vals.Ints = flagr.Add(s, "a01", flagr.Ints(defaults.Ints...), "usage for a01")
	vals.Int8 = flagr.Add(s, "a02", flagr.Int8(defaults.Int8), "usage for a02")
	vals.Int8s = flagr.Add(s, "a03", flagr.Int8s(defaults.Int8s...), "usage for a03")
	vals.Int16 = flagr.Add(s, "a04", flagr.Int16(defaults.Int16), "usage for a04")
	vals.Int16s = flagr.Add(s, "a05", flagr.Int16s(defaults.Int16s...), "usage for a05")
	vals.Int32 = flagr.Add(s, "a06", flagr.Int32(defaults.Int32), "usage for a06")
	vals.Int32s = flagr.Add(s, "a07", flagr.Int32s(defaults.Int32s...), "usage for a07")
	vals.Int64 = flagr.Add(s, "a08", flagr.Int64(defaults.Int64), "usage for a08")
	vals.Int64s = flagr.Add(s, "a09", flagr.Int64s(defaults.Int64s...), "usage for a09")
	vals.Uint = flagr.Add(s, "a10", flagr.Uint(defaults.Uint), "usage for a10")
	vals.Uints = flagr.Add(s, "a11", flagr.Uints(defaults.Uints...), "usage for a11")
	vals.Uint8 = flagr.Add(s, "a12", flagr.Uint8(defaults.Uint8), "usage for a12")
	vals.Uint8s = flagr.Add(s, "a13", flagr.Uint8s(defaults.Uint8s...), "usage for a13")
	vals.Uint16 = flagr.Add(s, "a14", flagr.Uint16(defaults.Uint16), "usage for a14")
	vals.Uint16s = flagr.Add(s, "a15", flagr.Uint16s(defaults.Uint16s...), "usage for a15")
	vals.Uint32 = flagr.Add(s, "a16", flagr.Uint32(defaults.Uint32), "usage for a16")
	vals.Uint32s = flagr.Add(s, "a17", flagr.Uint32s(defaults.Uint32s...), "usage for a17")
	vals.Uint64 = flagr.Add(s, "a18", flagr.Uint64(defaults.Uint64), "usage for a18")
	vals.Uint64s = flagr.Add(s, "a19", flagr.Uint64s(defaults.Uint64s...), "usage for a19")
	vals.Float32 = flagr.Add(s, "a20", flagr.Float32(defaults.Float32), "usage for a20")
	vals.Float32s = flagr.Add(s, "a21", flagr.Float32s(defaults.Float32s...), "usage for a21")
	vals.Float64 = flagr.Add(s, "a22", flagr.Float64(defaults.Float64), "usage for a22")
	vals.Float64s = flagr.Add(s, "a23", flagr.Float64s(defaults.Float64s...), "usage for a23")
	vals.Complex64 = flagr.Add(s, "a24", flagr.Complex64(defaults.Complex64), "usage for a24")
	vals.Complex64s = flagr.Add(s, "a25", flagr.Complex64s(defaults.Complex64s...), "usage for a25")
	vals.Complex128 = flagr.Add(s, "a26", flagr.Complex128(defaults.Complex128), "usage for a26")
	vals.Complex128s = flagr.Add(s, "a27", flagr.Complex128s(defaults.Complex128s...), "usage for a27")
	vals.Bool = flagr.Add(s, "a28", flagr.Bool(defaults.Bool), "usage for a28")
	vals.Bools = flagr.Add(s, "a29", flagr.Bools(defaults.Bools...), "usage for a29")
	vals.String = flagr.Add(s, "a30", flagr.String(defaults.String), "usage for a30")
	vals.Strings = flagr.Add(s, "a31", flagr.Strings(defaults.Strings...), "usage for a31")
	vals.Duration = flagr.Add(s, "a32", flagr.Duration(defaults.Duration), "usage for a32")
	vals.Durations = flagr.Add(s, "a33", flagr.Durations(defaults.Durations...), "usage for a33")
	vals.Time = flagr.Add(s, "a34", flagr.Time(timeLayout, defaults.Time), "usage for a34")
	vals.MustTime = flagr.Add(s, "a35", flagr.MustTime(timeLayout, defaults.MustTime), "usage for a35")
	vals.Times = flagr.Add(s, "a36", flagr.Times(timeLayout, defaults.Times...), "usage for a36")
	vals.MustTimes = flagr.Add(s, "a37", flagr.MustTimes(timeLayout, defaults.MustTimes...), "usage for a37")
	vals.URL = flagr.Add(s, "a38", flagr.URL(defaults.URL), "usage for a38")
	vals.MustURL = flagr.Add(s, "a39", flagr.MustURL(defaults.MustURL), "usage for a39")
	vals.URLs = flagr.Add(s, "a40", flagr.URLs(defaults.URLs...), "usage for a40")
	vals.MustURLs = flagr.Add(s, "a41", flagr.MustURLs(defaults.MustURLs...), "usage for a41")
	vals.IPAddr = flagr.Add(s, "a42", flagr.IPAddr(defaults.IPAddr), "usage for a42")
	vals.MustIPAddr = flagr.Add(s, "a43", flagr.MustIPAddr(defaults.MustIPAddr), "usage for a43")
	vals.IPAddrs = flagr.Add(s, "a44", flagr.IPAddrs(defaults.IPAddrs...), "usage for a44")
	vals.MustIPAddrs = flagr.Add(s, "a45", flagr.MustIPAddrs(defaults.MustIPAddrs...), "usage for a45")
	vals.IPAddrPort = flagr.Add(s, "a46", flagr.IPAddrPort(defaults.IPAddrPort), "usage for a46")
	vals.MustIPAddrPort = flagr.Add(s, "a47", flagr.MustIPAddrPort(defaults.MustIPAddrPort), "usage for a47")
	vals.IPAddrPorts = flagr.Add(s, "a48", flagr.IPAddrPorts(defaults.IPAddrPorts...), "usage for a48")
	vals.MustIPAddrPorts = flagr.Add(s, "a49", flagr.MustIPAddrPorts(defaults.MustIPAddrPorts...), "usage for a49")
	return vals
}

func TestDefaults(t *testing.T) {
	var s flagr.Set
	vals := makeFlags(&s)
	if err := s.Parse(nil); err != nil {
		t.Fatal(err)
	}

	want := flags{
		Int:             ptrTo(defaults.Int),
		Ints:            ptrTo(defaults.Ints),
		Int8:            ptrTo(defaults.Int8),
		Int8s:           ptrTo(defaults.Int8s),
		Int16:           ptrTo(defaults.Int16),
		Int16s:          ptrTo(defaults.Int16s),
		Int32:           ptrTo(defaults.Int32),
		Int32s:          ptrTo(defaults.Int32s),
		Int64:           ptrTo(defaults.Int64),
		Int64s:          ptrTo(defaults.Int64s),
		Uint:            ptrTo(defaults.Uint),
		Uints:           ptrTo(defaults.Uints),
		Uint8:           ptrTo(defaults.Uint8),
		Uint8s:          ptrTo(defaults.Uint8s),
		Uint16:          ptrTo(defaults.Uint16),
		Uint16s:         ptrTo(defaults.Uint16s),
		Uint32:          ptrTo(defaults.Uint32),
		Uint32s:         ptrTo(defaults.Uint32s),
		Uint64:          ptrTo(defaults.Uint64),
		Uint64s:         ptrTo(defaults.Uint64s),
		Float32:         ptrTo(defaults.Float32),
		Float32s:        ptrTo(defaults.Float32s),
		Float64:         ptrTo(defaults.Float64),
		Float64s:        ptrTo(defaults.Float64s),
		Complex64:       ptrTo(defaults.Complex64),
		Complex64s:      ptrTo(defaults.Complex64s),
		Complex128:      ptrTo(defaults.Complex128),
		Complex128s:     ptrTo(defaults.Complex128s),
		Bool:            ptrTo(defaults.Bool),
		Bools:           ptrTo(defaults.Bools),
		String:          ptrTo(defaults.String),
		Strings:         ptrTo(defaults.Strings),
		Duration:        ptrTo(defaults.Duration),
		Durations:       ptrTo(defaults.Durations),
		Time:            ptrTo(defaults.Time),
		MustTime:        ptrTo(defaults.Time),
		Times:           ptrTo(defaults.Times),
		MustTimes:       ptrTo(defaults.Times),
		URL:             ptrTo(defaults.URL),
		MustURL:         ptrTo(defaults.URL),
		URLs:            ptrTo(defaults.URLs),
		MustURLs:        ptrTo(defaults.URLs),
		IPAddr:          ptrTo(defaults.IPAddr),
		MustIPAddr:      ptrTo(defaults.IPAddr),
		IPAddrs:         ptrTo(defaults.IPAddrs),
		MustIPAddrs:     ptrTo(defaults.IPAddrs),
		IPAddrPort:      ptrTo(defaults.IPAddrPort),
		MustIPAddrPort:  ptrTo(defaults.IPAddrPort),
		IPAddrPorts:     ptrTo(defaults.IPAddrPorts),
		MustIPAddrPorts: ptrTo(defaults.IPAddrPorts),
	}
	if !reflect.DeepEqual(want, vals) {
		t.Errorf("Parse() vals don't match")
	}
}

func TestSetters(t *testing.T) {
	var s flagr.Set
	vals := makeFlags(&s)
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
		"-a25", "1i", "-a25", "2i", "-a25", "3i",
		"-a26", "1i",
		"-a27", "1i", "-a27", "2i", "-a27", "3i",
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

	want := flags{
		Int:             ptrTo(int(1)),
		Ints:            ptrTo([]int{1, 2, 3}),
		Int8:            ptrTo(int8(1)),
		Int8s:           ptrTo([]int8{1, 2, 3}),
		Int16:           ptrTo(int16(1)),
		Int16s:          ptrTo([]int16{1, 2, 3}),
		Int32:           ptrTo(int32(1)),
		Int32s:          ptrTo([]int32{1, 2, 3}),
		Int64:           ptrTo(int64(1)),
		Int64s:          ptrTo([]int64{1, 2, 3}),
		Uint:            ptrTo(uint(1)),
		Uints:           ptrTo([]uint{1, 2, 3}),
		Uint8:           ptrTo(uint8(1)),
		Uint8s:          ptrTo([]uint8{1, 2, 3}),
		Uint16:          ptrTo(uint16(1)),
		Uint16s:         ptrTo([]uint16{1, 2, 3}),
		Uint32:          ptrTo(uint32(1)),
		Uint32s:         ptrTo([]uint32{1, 2, 3}),
		Uint64:          ptrTo(uint64(1)),
		Uint64s:         ptrTo([]uint64{1, 2, 3}),
		Float32:         ptrTo(float32(1)),
		Float32s:        ptrTo([]float32{1, 2, 3}),
		Float64:         ptrTo(float64(1)),
		Float64s:        ptrTo([]float64{1, 2, 3}),
		Complex64:       ptrTo(complex64(1i)),
		Complex64s:      ptrTo([]complex64{1i, 2i, 3i}),
		Complex128:      ptrTo(complex128(1i)),
		Complex128s:     ptrTo([]complex128{1i, 2i, 3i}),
		Bool:            ptrTo(false),
		Bools:           ptrTo([]bool{false, true, false}),
		String:          ptrTo("qwe"),
		Strings:         ptrTo([]string{"qwe", "rty", "uio"}),
		Duration:        ptrTo(1 * time.Second),
		Durations:       ptrTo([]time.Duration{1 * time.Second, 2 * time.Second, 3 * time.Second}),
		Time:            ptrTo(musttime("0000-01-01")),
		MustTime:        ptrTo(musttime("0000-01-01")),
		Times:           ptrTo([]time.Time{musttime("0000-01-01"), musttime("0000-01-02"), musttime("0000-01-03")}),
		MustTimes:       ptrTo([]time.Time{musttime("0000-01-01"), musttime("0000-01-02"), musttime("0000-01-03")}),
		URL:             ptrTo(musturl("https://a.com")),
		MustURL:         ptrTo(musturl("https://a.com")),
		URLs:            ptrTo([]*url.URL{musturl("https://a.com"), musturl("https://b.com"), musturl("https://c.com")}),
		MustURLs:        ptrTo([]*url.URL{musturl("https://a.com"), musturl("https://b.com"), musturl("https://c.com")}),
		IPAddr:          ptrTo(netip.MustParseAddr("0.0.0.0")),
		MustIPAddr:      ptrTo(netip.MustParseAddr("0.0.0.0")),
		IPAddrs:         ptrTo([]netip.Addr{netip.MustParseAddr("0.0.0.0"), netip.MustParseAddr("0.0.0.1"), netip.MustParseAddr("0.0.0.2")}),
		MustIPAddrs:     ptrTo([]netip.Addr{netip.MustParseAddr("0.0.0.0"), netip.MustParseAddr("0.0.0.1"), netip.MustParseAddr("0.0.0.2")}),
		IPAddrPort:      ptrTo(netip.MustParseAddrPort("0.0.0.0:80")),
		MustIPAddrPort:  ptrTo(netip.MustParseAddrPort("0.0.0.0:80")),
		IPAddrPorts:     ptrTo([]netip.AddrPort{netip.MustParseAddrPort("0.0.0.0:80"), netip.MustParseAddrPort("0.0.0.0:81"), netip.MustParseAddrPort("0.0.0.0:82")}),
		MustIPAddrPorts: ptrTo([]netip.AddrPort{netip.MustParseAddrPort("0.0.0.0:80"), netip.MustParseAddrPort("0.0.0.0:81"), netip.MustParseAddrPort("0.0.0.0:82")}),
	}

	if !reflect.DeepEqual(want, vals) {
		t.Errorf("Parse() vals don't match")
	}
}

func TestDoesntClobberDefault(t *testing.T) {
	strings := []string{"a", "b", "c"}
	urlv := &url.URL{Host: "foo.com"}
	time := musttime("2000-01-01")

	fs := flagr.NewSet("", flagr.ContinueOnError)
	flagr.Add(&fs, "a", flagr.Strings(strings...), "")
	flagr.Add(&fs, "b", flagr.URL(urlv), "")
	flagr.Add(&fs, "c", flagr.Time(timeLayout, musttime("0001-01-01")), "")
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
	if want := musttime("2000-01-01"); !reflect.DeepEqual(time, want) {
		t.Fatalf("want %v, got %v", want, urlv)
	}
}

func TestZeroVal(t *testing.T) {
	var set flagr.Set

	if got := set.Name(); got != "" {
		t.Errorf("Name() = %v, want %v", got, "")
	}
	if got := set.ErrorHandling(); got != flagr.ContinueOnError {
		t.Errorf("ErrorHandling() = %v, want %v", got, flagr.ContinueOnError)
	}
	if err := set.Parse(nil); err != nil {
		t.Errorf("Parse() = %v", err)
	}

	set.Init("foo", flag.ExitOnError)
	if got := set.Name(); got != "foo" {
		t.Errorf("Name() = %v, want %v", got, "foo")
	}
	if got := set.ErrorHandling(); got != flagr.ExitOnError {
		t.Errorf("ErrorHandling() = %v, want %v", got, flagr.ExitOnError)
	}
}

func TestUsage(t *testing.T) {
	var set flagr.Set
	a := 0
	set.SetUsage(func() {
		a++
	})
	set.Usage()

	if a != 1 {
		t.Fatal()
	}
}

func musttime(t string) time.Time {
	v, err := time.Parse(timeLayout, t)
	if err != nil {
		panic(err)
	}
	return v
}
func musturl(t string) *url.URL {
	v, err := url.Parse(t)
	if err != nil {
		panic(err)
	}
	return v
}

func ptrTo[T any](t T) *T {
	return &t
}
