package main

import (
	"flag"
	"fmt"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/tufin/oasdiff/diff"
	"github.com/tufin/oasdiff/load"
	"github.com/tufin/oasdiff/report"
	"gopkg.in/yaml.v3"
)

var base, revision, prefix, filter, format string
var excludeExamples, excludeDescription, summary, breaking bool

const (
	formatYAML = "yaml"
	formatText = "text"
	formatHTML = "html"
)

func init() {
	flag.StringVar(&base, "base", "", "path of original OpenAPI spec in YAML or JSON format")
	flag.StringVar(&revision, "revision", "", "path of revised OpenAPI spec in YAML or JSON format")
	flag.StringVar(&prefix, "prefix", "", "path prefix that exists in base spec but not the revision (optional)")
	flag.StringVar(&filter, "filter", "", "regex to filter result paths (optional)")
	flag.BoolVar(&excludeExamples, "exclude-examples", false, "exclude changes to examples")
	flag.BoolVar(&excludeDescription, "exclude-description", false, "exclude changes to descriptions")
	flag.BoolVar(&summary, "summary", false, "display a summary of the changes instead of the full diff")
	flag.BoolVar(&breaking, "breaking", false, "display breaking changes only")
	flag.StringVar(&format, "format", formatYAML, "output format: yaml, text or html")
}

func main() {
	flag.Parse()

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true

	s1, err := load.From(loader, base)
	if err != nil {
		fmt.Printf("failed to load base spec from %q with %v", base, err)
		return
	}

	s2, err := load.From(loader, revision)
	if err != nil {
		fmt.Printf("failed to load revision spec from %q with %v", revision, err)
		return
	}

	diffReport, err := diff.Get(&diff.Config{
		ExcludeExamples:    excludeExamples,
		ExcludeDescription: excludeDescription,
		PathFilter:         filter,
		PathPrefix:         prefix,
		BreakingOnly:       breaking,
	}, s1, s2)

	if err != nil {
		fmt.Printf("diff failed with %v", err)
		return
	}

	if summary {
		summary := diffReport.GetSummary()
		printYAML(summary)
		return
	}

	if format == formatYAML {
		printYAML(diffReport)
	} else if format == formatText {
		fmt.Printf("%s\n", report.GetTextReportAsString(diffReport))
	} else if format == formatHTML {
		html, err := report.GetHTMLReportAsString(diffReport)
		if err != nil {
			fmt.Printf("failed to generate HTML with %v", err)
			return
		}
		fmt.Printf("%s\n", html)
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
