package main

import (
	"encoding/json"
	"flag"
	"fmt"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/tufin/oasdiff/diff"
	"github.com/tufin/oasdiff/load"
	"gopkg.in/yaml.v3"
)

var base, revision, prefix, filter, format string
var examples bool

const (
	formatJSON = "json"
	formatYAML = "yaml"
)

func init() {
	flag.StringVar(&base, "base", "", "original OpenAPI spec")
	flag.StringVar(&revision, "revision", "", "revised OpenAPI spec")
	flag.StringVar(&prefix, "prefix", "", "path prefix that exists in base spec but not the revision")
	flag.StringVar(&filter, "filter", "", "regex to filter result paths")
	flag.BoolVar(&examples, "examples", false, "whether to include examples in the diff")
	flag.StringVar(&format, "format", formatYAML, "output format: yaml or json")
}

func main() {
	flag.Parse()

	swaggerLoader := openapi3.NewSwaggerLoader()
	swaggerLoader.IsExternalRefsAllowed = true

	s1, err := load.From(swaggerLoader, base)
	if err != nil {
		fmt.Printf("Failed to load base spec from '%s' with '%v'", base, err)
		return
	}

	s2, err := load.From(swaggerLoader, revision)
	if err != nil {
		fmt.Printf("Failed to load revision spec from '%s' with '%v'", revision, err)
		return
	}

	diffReport, err := diff.Get(&diff.Config{
		IncludeExamples: examples,
		Filter:          filter,
		Prefix:          prefix,
	}, s1, s2)

	if err != nil {
		fmt.Printf("Diff failed with '%v'", err)
	}

	if format == formatJSON {
		bytes, err := json.MarshalIndent(diffReport, "", " ")
		if err != nil {
			fmt.Printf("Failed to marshal result as %s with '%v'", format, err)
			return
		}
		fmt.Printf("%s\n", bytes)
	} else if format == formatYAML {
		bytes, err := yaml.Marshal(diffReport)
		if err != nil {
			fmt.Printf("Failed to marshal result as %s with '%v'", format, err)
			return
		}
		fmt.Printf("%s\n", bytes)
	} else {
		fmt.Printf("Unknown format %s\n", format)
	}
}
