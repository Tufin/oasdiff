package internal

import (
	"fmt"
	"io"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/spf13/cobra"
	"github.com/tufin/oasdiff/checker"
	"github.com/tufin/oasdiff/diff"
)

func getChangelogCmd() *cobra.Command {

	flags := ChangelogFlags{}

	cmd := cobra.Command{
		Use:   "changelog base revision [flags]",
		Short: "Display changelog",
		Long:  "Display a changelog between base and revision specs." + specHelp,
		Args:  getParseArgs(&flags),
		RunE:  getRun(&flags, runChangelog),
	}

	cmd.PersistentFlags().BoolVarP(&flags.composed, "composed", "c", false, "work in 'composed' mode, compare paths in all specs matching base and revision globs")
	enumWithOptions(&cmd, newEnumValue([]string{FormatYAML, FormatJSON, FormatText}, FormatText, &flags.format), "format", "f", "output format")
	enumWithOptions(&cmd, newEnumSliceValue(diff.ExcludeDiffOptions, nil, &flags.excludeElements), "exclude-elements", "e", "comma-separated list of elements to exclude")
	cmd.PersistentFlags().StringVarP(&flags.matchPath, "match-path", "p", "", "include only paths that match this regular expression")
	cmd.PersistentFlags().StringVarP(&flags.filterExtension, "filter-extension", "", "", "exclude paths and operations with an OpenAPI Extension matching this regular expression")
	cmd.PersistentFlags().IntVarP(&flags.circularReferenceCounter, "max-circular-dep", "", 5, "maximum allowed number of circular dependencies between objects in OpenAPI specs")
	cmd.PersistentFlags().StringVarP(&flags.prefixBase, "prefix-base", "", "", "add this prefix to paths in base-spec before comparison")
	cmd.PersistentFlags().StringVarP(&flags.prefixRevision, "prefix-revision", "", "", "add this prefix to paths in revised-spec before comparison")
	cmd.PersistentFlags().StringVarP(&flags.stripPrefixBase, "strip-prefix-base", "", "", "strip this prefix from paths in base-spec before comparison")
	cmd.PersistentFlags().StringVarP(&flags.stripPrefixRevision, "strip-prefix-revision", "", "", "strip this prefix from paths in revised-spec before comparison")
	cmd.PersistentFlags().BoolVarP(&flags.includePathParams, "include-path-params", "", false, "include path parameter names in endpoint matching")
	cmd.PersistentFlags().BoolVarP(&flags.flatten, "flatten", "", false, "merge subschemas under allOf before diff")
	enumWithOptions(&cmd, newEnumValue([]string{LangEn, LangRu}, LangDefault, &flags.lang), "lang", "l", "language for localized output")
	cmd.PersistentFlags().StringVarP(&flags.errIgnoreFile, "err-ignore", "", "", "configuration file for ignoring errors")
	cmd.PersistentFlags().StringVarP(&flags.warnIgnoreFile, "warn-ignore", "", "", "configuration file for ignoring warnings")
	cmd.PersistentFlags().VarP(newEnumSliceValue(checker.GetOptionalChecks(), nil, &flags.includeChecks), "include-checks", "i", "comma-separated list of optional checks (run 'oasdiff checks' to see options)")
	cmd.PersistentFlags().IntVarP(&flags.deprecationDaysBeta, "deprecation-days-beta", "", checker.BetaDeprecationDays, "min days required between deprecating a beta resource and removing it")
	cmd.PersistentFlags().IntVarP(&flags.deprecationDaysStable, "deprecation-days-stable", "", checker.StableDeprecationDays, "min days required between deprecating a stable resource and removing it")

	return &cmd
}

func enumWithOptions(cmd *cobra.Command, value enumVal, name, shorthand, usage string) {
	cmd.PersistentFlags().VarP(value, name, shorthand, usage+": "+value.listOf())
}

func runChangelog(flags Flags, stdout io.Writer) (bool, *ReturnError) {
	return getChangelog(flags, stdout, checker.INFO, getChangelogTitle)
}

func getChangelog(flags Flags, stdout io.Writer, level checker.Level, getOutputTitle GetOutputTitle) (bool, *ReturnError) {

	openapi3.CircularReferenceCounter = flags.getCircularReferenceCounter()

	diffReport, operationsSources, err := calcDiff(flags)
	if err != nil {
		return false, err
	}

	bcConfig := checker.GetAllChecks(flags.getIncludeChecks(), flags.getDeprecationDaysBeta(), flags.getDeprecationDaysStable())
	bcConfig.Localize = checker.NewLocalizer(flags.getLang(), LangDefault)

	errs, returnErr := filterIgnored(
		checker.CheckBackwardCompatibilityUntilLevel(bcConfig, diffReport, operationsSources, level),
		flags.getWarnIgnoreFile(), flags.getErrIgnoreFile())

	if returnErr != nil {
		return false, returnErr
	}

	if returnErr := outputChangelog(bcConfig, flags.getFormat(), stdout, errs, getOutputTitle); returnErr != nil {
		return false, returnErr
	}

	if flags.getFailOn() != "" {
		level, err := checker.NewLevel(flags.getFailOn())
		if err != nil {
			return false, getErrInvalidFlags(fmt.Errorf("invalid fail-on value %s", flags.getFailOn()))
		}
		return errs.HasLevelOrHigher(level), nil
	}

	return false, nil
}

func filterIgnored(errs checker.Changes, warnIgnoreFile string, errIgnoreFile string) (checker.Changes, *ReturnError) {

	if warnIgnoreFile != "" {
		var err error
		errs, err = checker.ProcessIgnoredBackwardCompatibilityErrors(checker.WARN, errs, warnIgnoreFile)
		if err != nil {
			return nil, getErrCantProcessIgnoreFile("warn", err)
		}
	}

	if errIgnoreFile != "" {
		var err error
		errs, err = checker.ProcessIgnoredBackwardCompatibilityErrors(checker.ERR, errs, errIgnoreFile)
		if err != nil {
			return nil, getErrCantProcessIgnoreFile("err", err)
		}
	}

	return errs, nil
}

func getChangelogTitle(config checker.Config, errs checker.Changes) string {
	count := errs.GetLevelCount()

	return config.Localize(
		"total-changes",
		len(errs),
		count[checker.ERR],
		checker.ERR.PrettyString(),
		count[checker.WARN],
		checker.WARN.PrettyString(),
		count[checker.INFO],
		checker.INFO.PrettyString(),
	)
}

type GetOutputTitle func(config checker.Config, errs checker.Changes) string

func outputChangelog(config checker.Config, format string, stdout io.Writer, errs checker.Changes, getOutputTitle GetOutputTitle) *ReturnError {
	switch format {
	case FormatYAML:
		if err := printYAML(stdout, errs); err != nil {
			return getErrFailedPrint("breaking changes YAML", err)
		}
	case FormatJSON:
		if err := printJSON(stdout, errs); err != nil {
			return getErrFailedPrint("breaking changes JSON", err)
		}
	case FormatText:
		if len(errs) > 0 {
			fmt.Fprint(stdout, getOutputTitle(config, errs))
		}

		for _, bcerr := range errs {
			fmt.Fprintf(stdout, "%s\n\n", bcerr.PrettyErrorText(config.Localize))
		}
	default:
		return getErrUnsupportedFormat(format)
	}

	return nil
}
