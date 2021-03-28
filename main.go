package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/tufin/oasdiff/diff"
	"github.com/tufin/oasdiff/load"
	"github.com/tufin/oasdiff/report/text"
	"gopkg.in/yaml.v3"
)

var base, revision, prefix, filter, format string
var examples, summary bool

const (
	formatYAML = "yaml"
	formatText = "text"
)

func init() {
	flag.StringVar(&base, "base", "", "path of original OpenAPI spec")
	flag.StringVar(&revision, "revision", "", "path of revised OpenAPI spec")
	flag.StringVar(&prefix, "prefix", "", "path prefix that exists in base spec but not the revision (optional)")
	flag.StringVar(&filter, "filter", "", "regex to filter result paths (optional)")
	flag.BoolVar(&examples, "examples", false, "whether to include examples in the diff")
	flag.BoolVar(&summary, "summary", false, "whether to output full diff (default) or just summary")
	flag.StringVar(&format, "format", formatYAML, "output format: yaml or text")
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

	if summary {
		summary := diffReport.GetSummary()
		printYAML(summary)
		return
	}

	if format == formatYAML {
		printYAML(diffReport)
	} else if format == formatText {
		report := text.Report{
			Writer: os.Stdout,
		}
		report.Print(diffReport)
	} else {
		fmt.Printf("unknown format %q\n", format)
	}
}

func printYAML(output interface{}) {
	bytes, err := yaml.Marshal(output)
	if err != nil {
		fmt.Printf("failed to marshal result as %q with %v", format, err)
		return
	}
	fmt.Printf("%s\n", bytes)
}
