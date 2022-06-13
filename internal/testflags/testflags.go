package testflags

import (
	"net/netip"
	"net/url"
	"time"

	"github.com/flga/flagr"
)

const TimeLayout = "2006-01-02"

type Flags struct {
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

type Defaults struct {
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
}

func Make(s *flagr.Set) (Flags, Defaults) {
	defaults := Defaults{
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
		Time:            MustTime("4242-02-24"),
		MustTime:        "4242-02-24",
		Times:           []time.Time{MustTime("4242-02-24"), MustTime("2000-02-24")},
		MustTimes:       []string{"4242-02-24", "2000-02-24"},
		URL:             MustURL("https://go.dev"),
		MustURL:         "https://go.dev",
		URLs:            []*url.URL{MustURL("https://go.dev"), MustURL("https://go.dev/tour/")},
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

	var vals Flags
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
	vals.Time = flagr.Add(s, "a34", flagr.Time(TimeLayout, defaults.Time), "usage for a34")
	vals.MustTime = flagr.Add(s, "a35", flagr.MustTime(TimeLayout, defaults.MustTime), "usage for a35")
	vals.Times = flagr.Add(s, "a36", flagr.Times(TimeLayout, defaults.Times...), "usage for a36")
	vals.MustTimes = flagr.Add(s, "a37", flagr.MustTimes(TimeLayout, defaults.MustTimes...), "usage for a37")
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
	return vals, defaults
}

func MustTime(t string) time.Time {
	v, err := time.Parse(TimeLayout, t)
	if err != nil {
		panic(err)
	}
	return v
}

func MustURL(t string) *url.URL {
	v, err := url.Parse(t)
	if err != nil {
		panic(err)
	}
	return v
}
