package main

import (
	"flag"
	"fmt"
	"os"
	"reflect"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/tufin/oasdiff/diff"
	"github.com/tufin/oasdiff/load"
	"github.com/tufin/oasdiff/report"
	"gopkg.in/yaml.v3"
)

var base, revision, prefix_base, prefix_revision, filter, filterExtension, format string
var excludeExamples, excludeDescription, summary, breakingOnly, failOnDiff bool

const (
	formatYAML = "yaml"
	formatText = "text"
	formatHTML = "html"
)

func init() {
	flag.StringVar(&base, "base", "", "path of original OpenAPI spec in YAML or JSON format")
	flag.StringVar(&revision, "revision", "", "path of revised OpenAPI spec in YAML or JSON format")
	flag.StringVar(&prefix_base, "prefix-base", "", "if provided, paths in original (base) spec will be prefixed with the given prefix before comparison")
	flag.StringVar(&prefix_revision, "prefix-revision", "", "if provided, paths in revised (revision) spec will be prefixed with the given prefix before comparison")
	flag.StringVar(&filter, "filter", "", "if provided, diff will include only paths that match this regular expression")
	flag.StringVar(&filterExtension, "filter-extension", "", "if provided, diff will exclude paths and operations with an OpenAPI Extension matching this regular expression")
	flag.BoolVar(&excludeExamples, "exclude-examples", false, "ignore changes to examples")
	flag.BoolVar(&excludeDescription, "exclude-description", false, "ignore changes to descriptions")
	flag.BoolVar(&summary, "summary", false, "display a summary of the changes instead of the full diff")
	flag.BoolVar(&breakingOnly, "breaking-only", false, "display breaking changes only")
	flag.StringVar(&format, "format", formatYAML, "output format: yaml, text or html")
	flag.BoolVar(&failOnDiff, "fail-on-diff", false, "fail with exit code 1 if a difference is found")
}

func validateFlags() bool {
	if base == "" {
		fmt.Printf("please specify the '-base' flag: the path of the original OpenAPI spec in YAML or JSON format\n")
		return false
	}
	if revision == "" {
		fmt.Printf("please specify the '-revision' flag: the path of the revised OpenAPI spec in YAML or JSON format\n")
		return false
	}
	supportedFormats := map[string]bool{"": true, "yaml": true, "text": true, "html": true}
	if !supportedFormats[format] {
		fmt.Printf("invalid format. Should be yaml, text or html\n")
		return false
	}
	return true
}

func main() {
	flag.Parse()

	if !validateFlags() {
		os.Exit(101)
	}

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true

	s1, err := load.From(loader, base)
	if err != nil {
		fmt.Printf("failed to load base spec from %q with %v\n", base, err)
		os.Exit(102)
	}

	s2, err := load.From(loader, revision)
	if err != nil {
		fmt.Printf("failed to load revision spec from %q with %v\n", revision, err)
		os.Exit(103)
	}

	config := diff.NewConfig()
	config.ExcludeExamples = excludeExamples
	config.ExcludeDescription = excludeDescription
	config.PathFilter = filter
	config.FilterExtension = filterExtension
	config.PathPrefixBase = prefix_base
	config.PathPrefixRevision = prefix_revision
	config.BreakingOnly = breakingOnly

	diffReport, err := diff.Get(config, s1, s2)

	if err != nil {
		fmt.Printf("diff failed with %v\n", err)
		os.Exit(104)
	}

	if summary {
		if err = printYAML(diffReport.GetSummary()); err != nil {
			fmt.Printf("failed to print summary with %v\n", err)
			os.Exit(105)
		}
		exitNormally(diffReport.Empty())
	}

	if format == formatYAML {
		if err = printYAML(diffReport); err != nil {
			fmt.Printf("failed to print diff YAML with %v\n", err)
			os.Exit(106)
		}
	} else if format == formatText {
		fmt.Printf("%s", report.GetTextReportAsString(diffReport))
	} else if format == formatHTML {
		html, err := report.GetHTMLReportAsString(diffReport)
		if err != nil {
			fmt.Printf("failed to generate HTML diff report with %v\n", err)
			os.Exit(107)
		}
		fmt.Printf("%s", html)
	} else {
		fmt.Printf("unknown output format %q\n", format)
		os.Exit(108)
	}

	exitNormally(diffReport.Empty())
}

func exitNormally(diffEmpty bool) {
	if failOnDiff && !diffEmpty {
		os.Exit(1)
	}
	os.Exit(0)
}

func printYAML(output interface{}) error {
	if reflect.ValueOf(output).IsNil() {
		return nil
	}

	bytes, err := yaml.Marshal(output)
	if err != nil {
		return err
	}
	fmt.Printf("%s", bytes)
	return nil
}
