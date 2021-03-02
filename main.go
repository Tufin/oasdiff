package main

import (
	"encoding/json"
	"flag"
	"fmt"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/tufin/oasdiff/diff"
	"github.com/tufin/oasdiff/load"
)

var base, revision, prefix, filter string
var examples bool

func init() {
	flag.StringVar(&base, "base", "", "original OpenAPI spec")
	flag.StringVar(&revision, "revision", "", "revised OpenAPI spec")
	flag.StringVar(&prefix, "prefix", "", "path prefix that exists in base spec but not the revision")
	flag.StringVar(&filter, "filter", "", "regex to filter result paths")
	flag.BoolVar(&examples, "examples", false, "whether to include examples in the diff")
}

func main() {
	flag.Parse()

	swaggerLoader := openapi3.NewSwaggerLoader()
	swaggerLoader.IsExternalRefsAllowed = true

	loader := load.NewOASLoader(swaggerLoader)

	s1, err := loader.From(base)
	if err != nil {
		fmt.Printf("Failed to load base spec from '%s' with '%v'", base, err)
		return
	}

	s2, err := loader.From(revision)
	if err != nil {
		fmt.Printf("Failed to load revision spec from '%s' with '%v'", revision, err)
		return
	}

	config := diff.Config{
		Examples: examples,
		Filter:   filter,
		Prefix:   prefix,
	}

	bytes, err := json.MarshalIndent(diff.Get(&config, s1, s2), "", " ")
	if err != nil {
		fmt.Printf("Failed to marshal result with '%v'", err)
		return
	}

	fmt.Printf("%s\n", bytes)
}
