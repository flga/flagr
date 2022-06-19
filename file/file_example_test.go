package file_test

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"os"

	"github.com/flga/flagr"
	"github.com/flga/flagr/file"
)

func Example() {
	var fs flagr.Set

	_ = flagr.Add(&fs, "child.flagName", flagr.Int(1), "usage")
	configFile := flagr.Add(&fs, "config", flagr.String("config.json"), "usage")

	if err := fs.Parse(
		os.Args[1:],
		file.Parse(
			file.Static("config.json"), // uses the static file "config.json" as a source
			file.Mux{
				".json": json.Unmarshal, // tells the parser how to handler ".json" files
			},
			file.IgnoreMissingFile(),         // makes the file optional, if it does not exist it is not considered an error
			file.WithFS(nil),                 // uses the given fs.FS to lookup the file
			file.WithMapper(file.NoopMapper), // map flag names to jsonpath (root.child.leaf) using the given mapper, since the flags are already in dot notation no mapping is necessary
		),
		file.Parse(
			configFile, // this time, use a dynamic config file
			file.Mux{
				".json": json.Unmarshal,
				".xml":  xml.Unmarshal, // we're able to parse xml now too
			},
		),
	); err != nil {
		if errors.Is(err, flagr.ErrHelp) {
			os.Exit(2)
		}
		os.Exit(1)
	}
}
