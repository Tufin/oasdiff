package internal

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"
	"github.com/tufin/oasdiff/checker"
	"github.com/tufin/oasdiff/checker/localizations"
	"github.com/tufin/oasdiff/diff"
	"github.com/tufin/oasdiff/formatters"
)

func getBreakingChangesCmd() *cobra.Command {

	flags := ChangelogFlags{}

	cmd := cobra.Command{
		Use:   "breaking base revision [flags]",
		Short: "Display breaking changes",
		Long:  "Display breaking changes between base and revision specs." + specHelp,
		Args:  getParseArgs(&flags),
		RunE:  getRun(&flags, runBreakingChanges),
	}

	cmd.PersistentFlags().BoolVarP(&flags.composed, "composed", "c", false, "work in 'composed' mode, compare paths in all specs matching base and revision globs")
	enumWithOptions(&cmd, newEnumValue(formatters.SupportedFormatsByContentType(formatters.OutputBreaking), string(formatters.FormatText), &flags.format), "format", "f", "output format")
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
	enumWithOptions(&cmd, newEnumValue([]string{LevelErr, LevelWarn}, "", &flags.failOn), "fail-on", "o", "exit with return code 1 when output includes errors with this level or higher")
	enumWithOptions(&cmd, newEnumValue(localizations.GetSupportedLanguages(), localizations.LangDefault, &flags.lang), "lang", "l", "language for localized output")
	cmd.PersistentFlags().StringVarP(&flags.errIgnoreFile, "err-ignore", "", "", "configuration file for ignoring errors")
	cmd.PersistentFlags().StringVarP(&flags.warnIgnoreFile, "warn-ignore", "", "", "configuration file for ignoring warnings")
	cmd.PersistentFlags().VarP(newEnumSliceValue(checker.GetOptionalChecks(), nil, &flags.includeChecks), "include-checks", "i", "comma-separated list of optional checks (run 'oasdiff checks --required false' to see options)")
	cmd.PersistentFlags().IntVarP(&flags.deprecationDaysBeta, "deprecation-days-beta", "", checker.BetaDeprecationDays, "min days required between deprecating a beta resource and removing it")
	cmd.PersistentFlags().IntVarP(&flags.deprecationDaysStable, "deprecation-days-stable", "", checker.StableDeprecationDays, "min days required between deprecating a stable resource and removing it")
	enumWithOptions(&cmd, newEnumValue([]string{"auto", "always", "never"}, "auto", &flags.color), "color", "", "when to output colored escape sequences")

	return &cmd
}

func runBreakingChanges(flags Flags, stdout io.Writer) (bool, *ReturnError) {
	return getChangelog(flags, stdout, checker.WARN)
}

func outputBreakingChanges(format string, lang string, color string, stdout io.Writer, errs checker.Changes) *ReturnError {
	// formatter lookup
	formatter, err := formatters.Lookup(format, formatters.FormatterOpts{
		Language: lang,
	})
	if err != nil {
		return getErrUnsupportedBreakingChangesFormat(format)
	}

	// render
	colorMode, err := checker.NewColorMode(color)
	if err != nil {
		return getErrInvalidColorMode(err)
	}

	bytes, err := formatter.RenderBreakingChanges(errs, formatters.RenderOpts{ColorMode: colorMode})
	if err != nil {
		return getErrFailedPrint("breaking "+format, err)
	}

	// print output
	_, _ = fmt.Fprintf(stdout, "%s\n", bytes)

	return nil
}
