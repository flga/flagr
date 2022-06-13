package flagr_test

import (
	"errors"
	"fmt"
	"os"

	"github.com/flga/flagr"
)

func Example_customTypes() {
	var set flagr.Set

	// A fully custom flagr.Getter implementation.
	a := flagr.Add(&set, "a", NewFileMode(0), "usage")

	// Another fully custom implementation.
	// Notice that while the flag type is Counter, f is a *int.
	b := flagr.Add(&set, "b", NewCounter(2), "usage")

	args := []string{
		"-a", "rw",
		"-b", "-b", "-b",
	}
	if err := set.Parse(args); err != nil {
		if errors.Is(err, flagr.ErrHelp) {
			os.Exit(0)
		}
		os.Exit(2)
	}

	fmt.Printf("a = %s and is a %T\n", a.String(), a)
	fmt.Printf("b = %v and is a %T\n", *b, b)
	// Output:
	// a = rw and is a *flagr_test.FileMode
	// b = 3 and is a *int
}

var _ flagr.Getter[FileMode] = new(FileMode)

// a fully custom implementation of flagr.Getter
type FileMode byte

const (
	Read FileMode = 1 << iota
	Write
)

func NewFileMode(defaultValue FileMode) flagr.Getter[FileMode] {
	return &defaultValue
}

func (f *FileMode) Get() any {
	return f
}

func (f *FileMode) Set(s string) error {
	for _, rune := range s {
		switch rune {
		case 'r':
			*f |= Read
		case 'w':
			*f |= Write
		default:
			return fmt.Errorf("invalid file mode %q", s)
		}
	}
	return nil
}

func (f *FileMode) String() string {
	if f == nil {
		return ""
	}

	var ret string
	if *f&Read > 0 {
		ret += "r"
	}
	if *f&Write > 0 {
		ret += "w"
	}

	return ret
}

func (f *FileMode) IsBoolFlag() bool {
	return false
}

func (f *FileMode) Val() *FileMode {
	return f
}

// a fully custom implementation of flagr.Getter
type Counter struct {
	value *int
	set   bool
}

func NewCounter(defaultValue int) flagr.Getter[int] {
	return &Counter{value: &defaultValue}
}

func (f *Counter) Get() any {
	return f.value
}

func (f *Counter) Set(s string) error {
	if !f.set {
		*f.value = 0
		f.set = true
	}

	*f.value++
	return nil
}

func (f *Counter) String() string {
	if f == nil || f.value == nil {
		return ""
	}

	return fmt.Sprint(f.value)
}

func (f Counter) Val() *int {
	return f.value // this is what we care about, the wrapper struct is just an implementation constraint
}

func (f Counter) IsBoolFlag() bool {
	return true
}
