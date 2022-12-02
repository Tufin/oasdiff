package main

import (
	"flag"
	"fmt"
	"os"
	"reflect"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/tufin/oasdiff/build"
	"github.com/tufin/oasdiff/diff"
	"github.com/tufin/oasdiff/load"
	"github.com/tufin/oasdiff/report"
	"gopkg.in/yaml.v3"
)

var base, revision, filter, filterExtension, format string
var prefix_base, prefix_revision, strip_prefix_base, strip_prefix_revision, prefix string
var excludeExamples, excludeDescription, summary, breakingOnly, failOnDiff, version bool
var deprecationDays int

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
	flag.StringVar(&strip_prefix_base, "strip-prefix-base", "", "if provided, this prefix will be stripped from paths in original (base) spec before comparison")
	flag.StringVar(&strip_prefix_revision, "strip-prefix-revision", "", "if provided, this prefix will be stripped from paths in revised (revision) spec before comparison")
	flag.StringVar(&prefix, "prefix", "", "deprecated. use -prefix-revision instead")
	flag.StringVar(&filter, "filter", "", "if provided, diff will include only paths that match this regular expression")
	flag.StringVar(&filterExtension, "filter-extension", "", "if provided, diff will exclude paths and operations with an OpenAPI Extension matching this regular expression")
	flag.BoolVar(&excludeExamples, "exclude-examples", false, "ignore changes to examples")
	flag.BoolVar(&excludeDescription, "exclude-description", false, "ignore changes to descriptions")
	flag.BoolVar(&summary, "summary", false, "display a summary of the changes instead of the full diff")
	flag.BoolVar(&breakingOnly, "breaking-only", false, "display breaking changes only")
	flag.IntVar(&deprecationDays, "deprecation-days", 0, "minimal number of days required between deprecating a resource and removing it without being considered 'breaking'")
	flag.StringVar(&format, "format", formatYAML, "output format: yaml, text or html")
	flag.BoolVar(&failOnDiff, "fail-on-diff", false, "fail with exit code 1 if a difference is found")
	flag.BoolVar(&version, "version", false, "show version and quit")
	flag.IntVar(&openapi3.CircularReferenceCounter, "max-circular-dep", 0, "Maximum number of circular dependencies between objects that allowed to process OpenAPI file")
}

func validateFlags() bool {
	if base == "" {
		fmt.Fprintf(os.Stderr, "please specify the '-base' flag: the path of the original OpenAPI spec in YAML or JSON format\n")
		return false
	}
	if revision == "" {
		fmt.Fprintf(os.Stderr, "please specify the '-revision' flag: the path of the revised OpenAPI spec in YAML or JSON format\n")
		return false
	}
	supportedFormats := map[string]bool{"": true, "yaml": true, "text": true, "html": true}
	if !supportedFormats[format] {
		fmt.Fprintf(os.Stderr, "invalid format. Should be yaml, text or html\n")
		return false
	}
	if prefix != "" {
		if prefix_revision != "" {
			fmt.Fprintf(os.Stderr, "'-prefix' and '-prefix_revision' can't be used simultaneously\n")
			return false
		}
		prefix_revision = prefix
	}
	return true
}

func main() {
	flag.Parse()

	if version {
		fmt.Printf("oasdiff version: %s\n", build.Version)
		os.Exit(0)
	}

	if !validateFlags() {
		os.Exit(101)
	}

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true

	s1, err := load.From(loader, base)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load base spec from %q with %v\n", base, err)
		os.Exit(102)
	}

	s2, err := load.From(loader, revision)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load revision spec from %q with %v\n", revision, err)
		os.Exit(103)
	}

	config := diff.NewConfig()
	config.ExcludeExamples = excludeExamples
	config.ExcludeDescription = excludeDescription
	config.PathFilter = filter
	config.FilterExtension = filterExtension
	config.PathPrefixBase = prefix_base
	config.PathPrefixRevision = prefix_revision
	config.PathStripPrefixBase = strip_prefix_base
	config.PathStripPrefixRevision = strip_prefix_revision
	config.BreakingOnly = breakingOnly
	config.DeprecationDays = deprecationDays

	diffReport, err := diff.Get(config, s1, s2)

	if err != nil {
		fmt.Fprintf(os.Stderr, "diff failed with %v\n", err)
		os.Exit(104)
	}

	if summary {
		if err = printYAML(diffReport.GetSummary()); err != nil {
			fmt.Fprintf(os.Stderr, "failed to print summary with %v\n", err)
			os.Exit(105)
		}
		exitNormally(diffReport.Empty())
	}

	switch {
	case format == formatYAML:
		if err = printYAML(diffReport); err != nil {
			fmt.Fprintf(os.Stderr, "failed to print diff YAML with %v\n", err)
			os.Exit(106)
		}
	case format == formatText:
		fmt.Printf("%s", report.GetTextReportAsString(diffReport))
	case format == formatHTML:
		html, err := report.GetHTMLReportAsString(diffReport)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to generate HTML diff report with %v\n", err)
			os.Exit(107)
		}
		fmt.Printf("%s", html)
	default:
		fmt.Fprintf(os.Stderr, "unknown output format %q\n", format)
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
