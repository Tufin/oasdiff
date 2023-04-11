package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"reflect"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/tufin/oasdiff/build"
	"github.com/tufin/oasdiff/checker"
	"github.com/tufin/oasdiff/checker/localizations"
	"github.com/tufin/oasdiff/diff"
	"github.com/tufin/oasdiff/load"
	"github.com/tufin/oasdiff/report"
	"github.com/tufin/oasdiff/utils"
	"gopkg.in/yaml.v3"
)

var base, revision, filter, filterExtension, format, lang, warnIgnoreFile, errIgnoreFile string
var prefix_base, prefix_revision, strip_prefix_base, strip_prefix_revision, prefix string
var excludeExamples, excludeDescription, summary, breakingOnly, failOnDiff, failOnWarns, version, composed, checkBreaking, excludeEndpoints bool
var deprecationDays int
var includeChecks, excludeElements utils.StringList

const (
	formatYAML = "yaml"
	formatJSON = "json"
	formatText = "text"
	formatHTML = "html"
)

func init() {
	flag.StringVar(&base, "base", "", "path or URL (or a glob in Composed mode) of original OpenAPI spec in YAML or JSON format")
	flag.StringVar(&revision, "revision", "", "path or URL (or a glob in Composed mode) of revised OpenAPI spec in YAML or JSON format")
	flag.BoolVar(&composed, "composed", false, "work in 'composed' mode, compare paths in all specs matching base and revision globs")
	flag.StringVar(&prefix_base, "prefix-base", "", "if provided, paths in original (base) spec will be prefixed with the given prefix before comparison")
	flag.StringVar(&prefix_revision, "prefix-revision", "", "if provided, paths in revised (revision) spec will be prefixed with the given prefix before comparison")
	flag.StringVar(&strip_prefix_base, "strip-prefix-base", "", "if provided, this prefix will be stripped from paths in original (base) spec before comparison")
	flag.StringVar(&strip_prefix_revision, "strip-prefix-revision", "", "if provided, this prefix will be stripped from paths in revised (revision) spec before comparison")
	flag.StringVar(&prefix, "prefix", "", "deprecated. use '-prefix-revision' instead")
	flag.StringVar(&filter, "filter", "", "if provided, diff will include only paths that match this regular expression")
	flag.StringVar(&filterExtension, "filter-extension", "", "if provided, diff will exclude paths and operations with an OpenAPI Extension matching this regular expression")
	flag.BoolVar(&excludeExamples, "exclude-examples", false, "ignore changes to examples (deprecated, use '-exclude-elements examples' instead)")
	flag.BoolVar(&excludeDescription, "exclude-description", false, "ignore changes to descriptions (deprecated, use '-exclude-elements description' instead)")
	flag.BoolVar(&summary, "summary", false, "display a summary of the changes instead of the full diff")
	flag.BoolVar(&breakingOnly, "breaking-only", false, "display breaking changes only (old method)")
	flag.BoolVar(&checkBreaking, "check-breaking", false, "check for breaking changes (new method)")
	flag.StringVar(&warnIgnoreFile, "warn-ignore", "", "the configuration file for ignoring warnings with '-check-breaking'")
	flag.StringVar(&errIgnoreFile, "err-ignore", "", "the configuration file for ignoring errors with '-check-breaking'")
	flag.IntVar(&deprecationDays, "deprecation-days", 0, "minimal number of days required between deprecating a resource and removing it without being considered 'breaking'")
	flag.StringVar(&format, "format", formatYAML, "output format: yaml, json, text or html")
	flag.StringVar(&lang, "lang", "en", "language for localized breaking changes checks errors")
	flag.BoolVar(&failOnDiff, "fail-on-diff", false, "exit with return code 1 when any ERR-level breaking changes are found, used together with '-check-breaking'")
	flag.BoolVar(&failOnWarns, "fail-on-warns", false, "exit with return code 1 when any WARN-level breaking changes are found, used together with '-check-breaking' and '-fail-on-diff'")
	flag.BoolVar(&version, "version", false, "show version and quit")
	flag.IntVar(&openapi3.CircularReferenceCounter, "max-circular-dep", 5, "maximum allowed number of circular dependencies between objects in OpenAPI specs")
	flag.BoolVar(&excludeEndpoints, "exclude-endpoints", false, "exclude endpoints from output (deprecated, use '-exclude-elements endpoints' instead)")
	flag.Var(&includeChecks, "include-checks", "comma-separated list of optional breaking-changes checks")
	flag.Var(&excludeElements, "exclude-elements", "comma-separated list of elements to exclude from diff")
}

func isExcludeEndpoints() bool {
	return excludeEndpoints || excludeElements.Contains("endpoints")
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
	supportedFormats := map[string]bool{"yaml": true, "json": true, "text": true, "html": true}
	if !supportedFormats[format] {
		fmt.Fprintf(os.Stderr, "invalid format. Should be yaml, json text or html\n")
		return false
	}
	if format == "json" && !isExcludeEndpoints() {
		fmt.Fprintf(os.Stderr, "json format requires '-exclude-elements endpoints'\n")
		return false
	}
	if prefix != "" {
		if prefix_revision != "" {
			fmt.Fprintf(os.Stderr, "'-prefix' and '-prefix_revision' can't be used simultaneously\n")
			return false
		}
		prefix_revision = prefix
	}
	if failOnWarns {
		if !checkBreaking || !failOnDiff {
			fmt.Fprintf(os.Stderr, "'-fail-on-warns' is relevant only with '-check-breaking' and '-fail-on-diff'\n")
			return false
		}
	}

	if invalidChecks := checker.ValidateIncludeChecks(includeChecks); len(invalidChecks) > 0 {
		fmt.Fprintf(os.Stderr, "invalid include-checks: %s\n", invalidChecks.String())
		return false
	}

	if invalidElements := diff.ValidateExcludeElements(excludeElements); len(invalidElements) > 0 {
		fmt.Fprintf(os.Stderr, "invalid exclude-elements: %s\n", invalidElements.String())
		return false
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

	config := diff.NewConfig()
	config.PathFilter = filter
	config.FilterExtension = filterExtension
	config.PathPrefixBase = prefix_base
	config.PathPrefixRevision = prefix_revision
	config.PathStripPrefixBase = strip_prefix_base
	config.PathStripPrefixRevision = strip_prefix_revision
	config.BreakingOnly = breakingOnly
	config.DeprecationDays = deprecationDays
	config.SetExcludeElements(excludeElements.ToStringSet(), excludeExamples, excludeDescription, excludeEndpoints)

	var diffReport *diff.Diff
	var err error
	var operationsSources *diff.OperationsSourcesMap

	if checkBreaking {
		config.IncludeExtensions.Add(checker.XStabilityLevelExtension)
		config.IncludeExtensions.Add(diff.SunsetExtension)
		config.IncludeExtensions.Add(checker.XExtensibleEnumExtension)
	}

	if composed {
		s1, err := load.FromGlob(loader, base)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to load base spec from %q with %v\n", base, err)
			os.Exit(102)
		}

		s2, err := load.FromGlob(loader, revision)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to load revision spec from %q with %v\n", revision, err)
			os.Exit(103)
		}
		diffReport, operationsSources, err = diff.GetPathsDiff(config, s1, s2)
		if err != nil {
			fmt.Fprintf(os.Stderr, "diff failed with %v\n", err)
			os.Exit(104)
		}
	} else {
		s1, err := checker.LoadOpenAPISpecInfo(base)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to load base spec from %q with %v\n", base, err)
			os.Exit(102)
		}

		s2, err := checker.LoadOpenAPISpecInfo(revision)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to load revision spec from %q with %v\n", revision, err)
			os.Exit(103)
		}
		diffReport, operationsSources, err = diff.GetWithOperationsSourcesMap(config, s1, s2)
		if err != nil {
			fmt.Fprintf(os.Stderr, "diff failed with %v\n", err)
			os.Exit(104)
		}
	}

	if checkBreaking {
		c := checker.GetChecks(includeChecks)
		c.Localizer = *localizations.New(lang, "en")
		errs := checker.CheckBackwardCompatibility(c, diffReport, operationsSources)

		if warnIgnoreFile != "" {
			errs, err = checker.ProcessIgnoredBackwardCompatibilityErrors(checker.WARN, errs, warnIgnoreFile)
			if err != nil {
				fmt.Fprintf(os.Stderr, "can't process warn ignore file %v\n", err)
				os.Exit(121)
			}
		}

		if errIgnoreFile != "" {
			errs, err = checker.ProcessIgnoredBackwardCompatibilityErrors(checker.ERR, errs, errIgnoreFile)
			if err != nil {
				fmt.Fprintf(os.Stderr, "can't process err ignore file %v\n", err)
				os.Exit(122)
			}
		}
		countWarns := 0

		if format == formatJSON {
			if err = printJSON(errs); err != nil {
				fmt.Fprintf(os.Stderr, "failed to print diff JSON with %v\n", err)
				os.Exit(106)
			}
			for _, bcerr := range errs {
				if bcerr.Level == checker.WARN {
					countWarns++
				}
			}
		} else {
			// pretty output
			if len(errs) > 0 {
				fmt.Printf(c.Localizer.Get("messages.total-errors"), len(errs))
			}

			for _, bcerr := range errs {
				if bcerr.Level == checker.WARN {
					countWarns++
				}
				fmt.Printf("%s\n\n", bcerr.PrettyErrorText(c.Localizer))
			}
		}

		countErrs := len(errs) - countWarns

		diffEmpty := countErrs == 0
		if failOnWarns {
			diffEmpty = len(errs) == 0
		}
		exitNormally(diffEmpty)
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
	case format == formatJSON:
		if err = printJSON(diffReport); err != nil {
			fmt.Fprintf(os.Stderr, "failed to print diff JSON with %v\n", err)
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

func printJSON(output interface{}) error {
	if reflect.ValueOf(output).IsNil() {
		return nil
	}

	bytes, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		return err
	}
	fmt.Printf("%s\n", bytes)
	return nil
}
