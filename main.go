package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/tufin/oasdiff/diff"
	"github.com/tufin/oasdiff/load"
	"github.com/tufin/oasdiff/text"
	"gopkg.in/yaml.v3"
)

var base, revision, prefix, filter, format string
var examples bool

const (
	formatJSON = "json"
	formatYAML = "yaml"
	formatText = "text"
)

func init() {
	flag.StringVar(&base, "base", "", "path of original OpenAPI spec")
	flag.StringVar(&revision, "revision", "", "path of revised OpenAPI spec")
	flag.StringVar(&prefix, "prefix", "", "path prefix that exists in base spec but not the revision (optional)")
	flag.StringVar(&filter, "filter", "", "regex to filter result paths (optional)")
	flag.BoolVar(&examples, "examples", false, "whether to include examples in the diff")
	flag.StringVar(&format, "format", formatYAML, "output format: yaml, json or text")
}

func main() {
	flag.Parse()

	swaggerLoader := openapi3.NewSwaggerLoader()
	swaggerLoader.IsExternalRefsAllowed = true

	s1, err := load.From(swaggerLoader, base)
	if err != nil {
		fmt.Printf("failed to load base spec from %q with %v", base, err)
		return
	}

	s2, err := load.From(swaggerLoader, revision)
	if err != nil {
		fmt.Printf("failed to load revision spec from %q with %v", revision, err)
		return
	}

	diffReport, err := diff.Get(&diff.Config{
		IncludeExamples: examples,
		Filter:          filter,
		Prefix:          prefix,
	}, s1, s2)

	if err != nil {
		fmt.Printf("diff failed with %v", err)
	}

	if format == formatJSON {
		bytes, err := json.MarshalIndent(diffReport, "", " ")
		if err != nil {
			fmt.Printf("failed to marshal result as %q with %v", format, err)
			return
		}
		fmt.Printf("%s\n", bytes)
	} else if format == formatYAML {
		bytes, err := yaml.Marshal(diffReport)
		if err != nil {
			fmt.Printf("failed to marshal result as %q with %v", format, err)
			return
		}
		fmt.Printf("%s\n", bytes)
	} else if format == formatText {
		text.Print(diffReport, os.Stdout)
	} else {
		fmt.Printf("unknown format %q\n", format)
	}
}
