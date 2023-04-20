package internal

import (
	"flag"
	"fmt"

	"github.com/tufin/oasdiff/checker"
	"github.com/tufin/oasdiff/diff"
	"github.com/tufin/oasdiff/utils"
)

type InputFlags struct {
	base                     string
	revision                 string
	composed                 bool
	prefixBase               string
	prefixRevision           string
	stripPrefixBase          string
	strip_prefix_revision    string
	prefix                   string
	filter                   string
	filterExtension          string
	excludeExamples          bool
	excludeDescription       bool
	summary                  bool
	breakingOnly             bool
	checkBreaking            bool
	warnIgnoreFile           string
	errIgnoreFile            string
	deprecationDays          int
	format                   string
	lang                     string
	failOnDiff               bool
	failOnWarns              bool
	version                  bool
	circularReferenceCounter int
	excludeEndpoints         bool
	includeChecks            utils.StringList
	excludeElements          utils.StringList
}

func parseFlags(args []string) (*InputFlags, *ReturnError) {

	if len(args) < 1 {
		return nil, getErrInvalidFlags(fmt.Errorf("empty argument list"))
	}

	flags := flag.NewFlagSet(args[0], flag.ContinueOnError)

	inputFlags := InputFlags{}
	flags.StringVar(&inputFlags.base, "base", "", "path or URL (or a glob in Composed mode) of original OpenAPI spec in YAML or JSON format")
	flags.StringVar(&inputFlags.revision, "revision", "", "path or URL (or a glob in Composed mode) of revised OpenAPI spec in YAML or JSON format")
	flags.BoolVar(&inputFlags.composed, "composed", false, "work in 'composed' mode, compare paths in all specs matching base and revision globs")
	flags.StringVar(&inputFlags.prefixBase, "prefix-base", "", "if provided, paths in original (base) spec will be prefixed with the given prefix before comparison")
	flags.StringVar(&inputFlags.prefixRevision, "prefix-revision", "", "if provided, paths in revised (revision) spec will be prefixed with the given prefix before comparison")
	flags.StringVar(&inputFlags.stripPrefixBase, "strip-prefix-base", "", "if provided, this prefix will be stripped from paths in original (base) spec before comparison")
	flags.StringVar(&inputFlags.strip_prefix_revision, "strip-prefix-revision", "", "if provided, this prefix will be stripped from paths in revised (revision) spec before comparison")
	flags.StringVar(&inputFlags.prefix, "prefix", "", "deprecated. use '-prefix-revision' instead")
	flags.StringVar(&inputFlags.filter, "filter", "", "if provided, diff will include only paths that match this regular expression")
	flags.StringVar(&inputFlags.filterExtension, "filter-extension", "", "if provided, diff will exclude paths and operations with an OpenAPI Extension matching this regular expression")
	flags.BoolVar(&inputFlags.excludeExamples, "exclude-examples", false, "ignore changes to examples (deprecated, use '-exclude-elements examples' instead)")
	flags.BoolVar(&inputFlags.excludeDescription, "exclude-description", false, "ignore changes to descriptions (deprecated, use '-exclude-elements description' instead)")
	flags.BoolVar(&inputFlags.summary, "summary", false, "display a summary of the changes instead of the full diff")
	flags.BoolVar(&inputFlags.breakingOnly, "breaking-only", false, "display breaking changes only (old method)")
	flags.BoolVar(&inputFlags.checkBreaking, "check-breaking", false, "check for breaking changes (new method)")
	flags.StringVar(&inputFlags.warnIgnoreFile, "warn-ignore", "", "the configuration file for ignoring warnings with '-check-breaking'")
	flags.StringVar(&inputFlags.errIgnoreFile, "err-ignore", "", "the configuration file for ignoring errors with '-check-breaking'")
	flags.IntVar(&inputFlags.deprecationDays, "deprecation-days", 0, "minimal number of days required between deprecating a resource and removing it without being considered 'breaking'")
	flags.StringVar(&inputFlags.format, "format", "", "output format=yaml, json, text or html")
	flags.StringVar(&inputFlags.lang, "lang", "en", "language for localized breaking changes checks errors")
	flags.BoolVar(&inputFlags.failOnDiff, "fail-on-diff", false, "exit with return code 1 when any ERR-level breaking changes are found, used together with '-check-breaking'")
	flags.BoolVar(&inputFlags.failOnWarns, "fail-on-warns", false, "exit with return code 1 when any WARN-level breaking changes are found, used together with '-check-breaking' and '-fail-on-diff'")
	flags.BoolVar(&inputFlags.version, "version", false, "show version and quit")
	flags.IntVar(&inputFlags.circularReferenceCounter, "max-circular-dep", 5, "maximum allowed number of circular dependencies between objects in OpenAPI specs")
	flags.BoolVar(&inputFlags.excludeEndpoints, "exclude-endpoints", false, "exclude endpoints from output (deprecated, use '-exclude-elements endpoints' instead)")
	flags.Var(&inputFlags.includeChecks, "include-checks", "comma-separated list of optional breaking-changes checks")
	flags.Var(&inputFlags.excludeElements, "exclude-elements", "comma-separated list of elements to exclude from diff")

	if err := flags.Parse(args[1:]); err != nil {
		return nil, getErrInvalidFlags(err)
	}

	return &inputFlags, nil
}

func isExcludeEndpoints(inputFlags *InputFlags) bool {
	return inputFlags.excludeEndpoints || inputFlags.excludeElements.Contains("endpoints")
}

func validateFormatFlag(inputFlags *InputFlags) *ReturnError {
	var supportedFormats utils.StringSet

	if inputFlags.checkBreaking {
		if inputFlags.format == "" {
			inputFlags.format = "text"
		}
		supportedFormats = utils.StringList{"yaml", "json", "text"}.ToStringSet()
	} else {
		if inputFlags.format == "" {
			inputFlags.format = "yaml"
		}
		if inputFlags.format == "json" && !isExcludeEndpoints(inputFlags) {
			return getErrInvalidFlags(fmt.Errorf("json format requires \"-exclude-elements endpoints\""))
		}
		supportedFormats = utils.StringList{"yaml", "json", "text", "html"}.ToStringSet()
	}

	if !supportedFormats.Contains(inputFlags.format) {
		return getErrUnsupportedDiffFormat(inputFlags.format)
	}
	return nil
}

func validateFlags(inputFlags *InputFlags) *ReturnError {
	if inputFlags.base == "" {
		return getErrInvalidFlags(fmt.Errorf("please specify the \"-base\" flag=the path of the original OpenAPI spec in YAML or JSON format"))
	}
	if inputFlags.revision == "" {
		return getErrInvalidFlags(fmt.Errorf("please specify the \"-revision\" flag=the path of the revised OpenAPI spec in YAML or JSON format"))
	}
	if returnErr := validateFormatFlag(inputFlags); returnErr != nil {
		return returnErr
	}
	if inputFlags.prefix != "" {
		if inputFlags.prefixRevision != "" {
			return getErrInvalidFlags(fmt.Errorf("\"-prefix\" and \"-prefix_revision\" can't be used simultaneously"))
		}
		inputFlags.prefixRevision = inputFlags.prefix
	}
	if inputFlags.failOnWarns {
		if !inputFlags.checkBreaking || !inputFlags.failOnDiff {
			return getErrInvalidFlags(fmt.Errorf("\"-fail-on-warns\" is relevant only with \"-check-breaking\" and \"-fail-on-diff\""))
		}
	}

	if invalidChecks := checker.ValidateIncludeChecks(inputFlags.includeChecks); len(invalidChecks) > 0 {
		return getErrInvalidFlags(fmt.Errorf("invalid include-checks=%s", inputFlags.includeChecks))
	}

	if invalidElements := diff.ValidateExcludeElements(inputFlags.excludeElements); len(invalidElements) > 0 {
		return getErrInvalidFlags(fmt.Errorf("invalid exclude-elements=%s", inputFlags.excludeElements))
	}

	return nil
}
