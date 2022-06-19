package env_test

import (
	"bytes"
	"errors"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/flga/flagr"
	"github.com/flga/flagr/env"
	"github.com/google/go-cmp/cmp"
)

func TestDefaultMapper(t *testing.T) {
	tests := map[string]string{
		"":                        "",
		"asd":                     "ASD",
		"asdASD123.,$=/_@'\"\n\t": "ASDASD123___________",
	}

	normalizer := env.DefaultMapper(env.NoSplit)
	for name, want := range tests {
		normalized, separator := normalizer(name)
		if normalized != want {
			t.Errorf("DefaultMapper() normalized = %s, want %s", normalized, want)
		}
		if separator != env.NoSplit {
			t.Errorf("DefaultMapper() separator = %s, want %s", normalized, want)
		}
	}
}

func TestMatchesFlagsToEnv(t *testing.T) {
	var set flagr.Set
	a1 := flagr.Add(&set, "a1", flagr.String("a"), "")
	a2 := flagr.Add(&set, "a2", flagr.String("a"), "")
	a3 := flagr.Add(&set, "a3", flagr.String("a"), "")
	a4 := flagr.Add(&set, "a4", flagr.String("a"), "")
	b1 := flagr.Add(&set, "b1", flagr.Int(1), "")
	b2 := flagr.Add(&set, "b2", flagr.Int(1), "")
	b3 := flagr.Add(&set, "b3", flagr.Int(1), "")
	b4 := flagr.Add(&set, "b4", flagr.Int(1), "")
	c1 := flagr.Add(&set, "c1", flagr.Strings("c", "c"), "")
	c2 := flagr.Add(&set, "c2", flagr.Strings("c", "c"), "")
	c3 := flagr.Add(&set, "c3", flagr.Strings("c", "c"), "")
	c4 := flagr.Add(&set, "c4", flagr.Strings("c", "c"), "")
	d1 := flagr.Add(&set, "d1", flagr.Ints(1, 1), "")
	d2 := flagr.Add(&set, "d2", flagr.Ints(1, 1), "")
	d3 := flagr.Add(&set, "d3", flagr.Ints(1, 1), "")
	d4 := flagr.Add(&set, "d4", flagr.Ints(1, 1), "")

	args := []string{
		"-a1", "testdata/.env",
		"-b1", "10",
		"-c1", "flag1", "-c1", "flag2",
		"-d1", "10", "-d1", "11",
	}
	lookuper := testLookuper(
		"APP_A2", "env",
		"APP_B2", "100",
		"APP_C2", "env1,env2",
		"APP_D2", "100,101",
	)
	if err := set.Parse(
		args,
		env.Parse(
			env.WithPrefix("app"),
			env.WithLookupFunc(lookuper),
			env.WithDotEnv(a1, false),
			env.WithMapper(func(flagName string) (envName string, listSplitter env.Splitter) {
				var defaultMapperRegex = regexp.MustCompile("[^a-zA-Z0-9_]")
				splitter := env.Splitter("")

				// we only have lists after app_cx (inclusive)
				if flagName >= "app_c" {
					splitter = env.Splitter(',')
				}
				return strings.ToUpper(defaultMapperRegex.ReplaceAllLiteralString(flagName, "_")), splitter
			}),
		),
	); err != nil {
		t.Fatal(err)
	}

	if want := "testdata/.env"; *a1 != want {
		t.Errorf("a1 = %v, want %v", *a1, want)
	}
	if want := "env"; *a2 != want {
		t.Errorf("a2 = %v, want %v", *a2, want)
	}
	if want := "file"; *a3 != want {
		t.Errorf("a3 = %v, want %v", *a3, want)
	}
	if want := "a"; *a4 != want {
		t.Errorf("a4 = %v, want %v", *a4, want)
	}
	if want := 10; *b1 != want {
		t.Errorf("b1 = %v, want %v", *b1, want)
	}
	if want := 100; *b2 != want {
		t.Errorf("b2 = %v, want %v", *b2, want)
	}
	if want := 1000; *b3 != want {
		t.Errorf("b3 = %v, want %v", *b3, want)
	}
	if want := 1; *b4 != want {
		t.Errorf("b4 = %v, want %v", *b4, want)
	}
	if want := []string{"flag1", "flag2"}; !reflect.DeepEqual(*c1, want) {
		t.Errorf("c1 = %v, want %v", *c1, want)
	}
	if want := []string{"env1", "env2"}; !reflect.DeepEqual(*c2, want) {
		t.Errorf("c2 = %v, want %v", *c2, want)
	}
	if want := []string{"file1", "file2"}; !reflect.DeepEqual(*c3, want) {
		t.Errorf("c3 = %v, want %v", *c3, want)
	}
	if want := []string{"c", "c"}; !reflect.DeepEqual(*c4, want) {
		t.Errorf("c4 = %v, want %v", *c4, want)
	}
	if want := []int{10, 11}; !reflect.DeepEqual(*d1, want) {
		t.Errorf("d1 = %v, want %v", *d1, want)
	}
	if want := []int{100, 101}; !reflect.DeepEqual(*d2, want) {
		t.Errorf("d2 = %v, want %v", *d2, want)
	}
	if want := []int{1000, 1001}; !reflect.DeepEqual(*d3, want) {
		t.Errorf("d3 = %v, want %v", *d3, want)
	}
	if want := []int{1, 1}; !reflect.DeepEqual(*d4, want) {
		t.Errorf("d4 = %v, want %v", *d4, want)
	}

	var buf bytes.Buffer
	set.SetOutput(&buf)
	set.PrintValues()
	want := `Current configuration:
  -a1 testdata/.env  (flags)
  -a2 env            (env: APP_A2)
  -a3 file           (envfile[testdata/.env]: APP_A3)
  -a4 a              (default)
  -b1 10             (flags)
  -b2 100            (env: APP_B2)
  -b3 1000           (envfile[testdata/.env]: APP_B3)
  -b4 1              (default)
  -c1 [flag1, flag2] (flags)
  -c2 [env1, env2]   (env: APP_C2)
  -c3 [file1, file2] (envfile[testdata/.env]: APP_C3)
  -c4 [c, c]         (default)
  -d1 [10, 11]       (flags)
  -d2 [100, 101]     (env: APP_D2)
  -d3 [1000, 1001]   (envfile[testdata/.env]: APP_D3)
  -d4 [1, 1]         (default)
`
	if diff := cmp.Diff(want, buf.String()); diff != "" {
		t.Errorf("values mismatch (-want +got):\n%s", diff)
	}
}

func TestFailsOnInvalidVals(t *testing.T) {
	t.Run("singe vals", func(t *testing.T) {
		var set flagr.Set
		flagr.Add(&set, "a1", flagr.Int(1), "")
		err := set.Parse(
			nil,
			env.Parse(
				env.WithPrefix("app"),
				env.WithLookupFunc(testLookuper(
					"APP_A1", "not a number",
				)),
			),
		)

		if want := strconv.ErrSyntax; !errors.Is(err, want) {
			t.Fatalf("err = %v, want %v", err, want)
		}
	})

	t.Run("slice vals", func(t *testing.T) {
		var set flagr.Set
		flagr.Add(&set, "a1", flagr.Ints(1), "")
		err := set.Parse(
			nil,
			env.Parse(
				env.WithPrefix("app"),
				env.WithMapper(env.DefaultMapper(",")),
				env.WithLookupFunc(testLookuper(
					"APP_A1", "1,2,not a number",
				)),
			),
		)

		if want := strconv.ErrSyntax; !errors.Is(err, want) {
			t.Fatalf("err = %v, want %v", err, want)
		}
	})
}

func TestInvalidFile(t *testing.T) {
	t.Run("fails if not optional and file doesn't exist", func(t *testing.T) {
		var set flagr.Set
		flagr.Add(&set, "a1", flagr.Int(1), "")
		if err := set.Parse(
			nil,
			env.Parse(
				env.WithPrefix("app"),
				env.WithDotEnv(ptr("notarealdotenvfile"), false),
			),
		); err == nil {
			t.Fatal("err is nil")
		}
	})

	t.Run("does not fail if optional and file doesn't exist", func(t *testing.T) {
		var set flagr.Set
		flagr.Add(&set, "a1", flagr.Int(1), "")
		if err := set.Parse(
			nil,
			env.Parse(
				env.WithPrefix("app"),
				env.WithDotEnv(ptr("notarealdotenvfile"), true),
			),
		); err != nil {
			t.Fatal(err)
		}
	})
}

func testLookuper(kv ...string) env.LookupFunc {
	env := make(map[string]string)
	for i, kOrV := range kv {
		if i%2 == 1 {
			env[kv[i-1]] = kOrV
		}
	}
	return func(key string) (string, bool) {
		v, ok := env[key]
		return v, ok
	}
}

func ptr[T any](t T) *T { return &t }
