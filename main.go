package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/tufin/oasdiff/build"
	"github.com/tufin/oasdiff/checker"
	"github.com/tufin/oasdiff/diff"
	"github.com/tufin/oasdiff/internal"
	"github.com/tufin/oasdiff/load"
	"github.com/tufin/oasdiff/utils"
)

var base, revision, filter, filterExtension, format, lang, warnIgnoreFile, errIgnoreFile string
var prefix_base, prefix_revision, strip_prefix_base, strip_prefix_revision, prefix string
var excludeExamples, excludeDescription, summary, breakingOnly, failOnDiff, failOnWarns, version, composed, checkBreaking, excludeEndpoints bool
var deprecationDays int
var includeChecks, excludeElements utils.StringList

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
	flag.StringVar(&format, "format", "", "output format: yaml, json, text or html")
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

func validateFormatFlag() *internal.ReturnError {
	var supportedFormats utils.StringSet

	if checkBreaking {
		if format == "" {
			format = "text"
		}
		supportedFormats = utils.StringList{"yaml", "json", "text"}.ToStringSet()
	} else {
		if format == "" {
			format = "yaml"
		}
		if format == "json" && !isExcludeEndpoints() {
			return internal.GetErrInvalidFlags(fmt.Errorf("json format requires \"-exclude-elements endpoints\""))
		}
		supportedFormats = utils.StringList{"yaml", "json", "text", "html"}.ToStringSet()
	}

	if !supportedFormats.Contains(format) {
		return internal.GetErrUnsupportedDiffFormat(format)
	}
	return nil
}

func validateFlags() *internal.ReturnError {
	if base == "" {
		return internal.GetErrInvalidFlags(fmt.Errorf("please specify the \"-base\" flag: the path of the original OpenAPI spec in YAML or JSON format"))
	}
	if revision == "" {
		return internal.GetErrInvalidFlags(fmt.Errorf("please specify the \"-revision\" flag: the path of the revised OpenAPI spec in YAML or JSON format"))
	}
	if err := validateFormatFlag(); err != nil {
		return err
	}
	if prefix != "" {
		if prefix_revision != "" {
			return internal.GetErrInvalidFlags(fmt.Errorf("\"-prefix\" and \"-prefix_revision\" can't be used simultaneously"))
		}
		prefix_revision = prefix
	}
	if failOnWarns {
		if !checkBreaking || !failOnDiff {
			return internal.GetErrInvalidFlags(fmt.Errorf("\"-fail-on-warns\" is relevant only with \"-check-breaking\" and \"-fail-on-diff\""))
		}
	}

	if invalidChecks := checker.ValidateIncludeChecks(includeChecks); len(invalidChecks) > 0 {
		return internal.GetErrInvalidFlags(fmt.Errorf("invalid include-checks: %s", invalidChecks.String()))
	}

	if invalidElements := diff.ValidateExcludeElements(excludeElements); len(invalidElements) > 0 {
		return internal.GetErrInvalidFlags(fmt.Errorf("invalid exclude-elements: %s", invalidElements.String()))
	}

	return nil
}

func main() {
	flag.Parse()

	if version {
		fmt.Printf("oasdiff version: %s\n", build.Version)
		os.Exit(0)
	}

	if err := validateFlags(); err != nil {
		exitWithError(err)
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
	var operationsSources *diff.OperationsSourcesMap

	if checkBreaking {
		config.IncludeExtensions.Add(checker.XStabilityLevelExtension)
		config.IncludeExtensions.Add(diff.SunsetExtension)
		config.IncludeExtensions.Add(checker.XExtensibleEnumExtension)
	}

	if composed {
		s1, err := load.FromGlob(loader, base)
		if err != nil {
			exitWithError(internal.GetErrFailedToLoadSpec("base", base, err))
		}

		s2, err := load.FromGlob(loader, revision)
		if err != nil {
			exitWithError(internal.GetErrFailedToLoadSpec("revision", revision, err))
		}
		diffReport, operationsSources, err = diff.GetPathsDiff(config, s1, s2)
		if err != nil {
			exitWithError(internal.GetErrDiffFailed(err))
		}
	} else {
		s1, err := checker.LoadOpenAPISpecInfo(base)
		if err != nil {
			exitWithError(internal.GetErrFailedToLoadSpec("base", base, err))
		}

		s2, err := checker.LoadOpenAPISpecInfo(revision)
		if err != nil {
			exitWithError(internal.GetErrFailedToLoadSpec("revision", revision, err))
		}
		diffReport, operationsSources, err = diff.GetWithOperationsSourcesMap(config, s1, s2)
		if err != nil {
			exitWithError(internal.GetErrDiffFailed(err))
		}
	}

	if checkBreaking {
		exit(internal.HandleBreakingChanges(diffReport, operationsSources, includeChecks, format, failOnWarns, lang, warnIgnoreFile, errIgnoreFile))
	}

	if summary {
		if err := internal.PrintYAML(diffReport.GetSummary()); err != nil {
			exitWithError(internal.GetErrFailedPrint("summary", err))
		}
		exitNormally(diffReport.Empty())
	}

	exit(diffReport.Empty(), internal.HandleDiff(diffReport, format))
}

func exit(diffEmpty bool, err *internal.ReturnError) {
	if err != nil {
		exitWithError(err)
	}
	exitNormally(diffEmpty)
}

func exitWithError(err *internal.ReturnError) {
	fmt.Fprintf(os.Stderr, "%v\n", err.Err)
	os.Exit(err.Code)
}

func exitNormally(diffEmpty bool) {
	if failOnDiff && !diffEmpty {
		os.Exit(1)
	}
	os.Exit(0)
}
