package internal

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"
	"github.com/tufin/oasdiff/checker"
	"github.com/tufin/oasdiff/diff"
	"github.com/tufin/oasdiff/formatters"
)

func getBreakingChangesCmd() *cobra.Command {

	flags := ChangelogFlags{}

	cmd := cobra.Command{
		Use:   "breaking base revision [flags]",
		Short: "Display breaking changes",
		Long: `Display breaking changes between base and revision specs.
Base and revision can be a path to a file or a URL.
In 'composed' mode, base and revision can be a glob and oasdiff will compare matching endpoints between the two sets of files.
`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {

			flags.base = args[0]
			flags.revision = args[1]

			// by now flags have been parsed successfully so we don't need to show usage on any errors
			cmd.Root().SilenceUsage = true

			failEmpty, err := runBreakingChanges(&flags, cmd.OutOrStdout())
			if err != nil {
				setReturnValue(cmd, err.Code)
				return err
			}

			if failEmpty {
				setReturnValue(cmd, 1)
			}

			return nil
		},
	}

	cmd.PersistentFlags().BoolVarP(&flags.composed, "composed", "c", false, "work in 'composed' mode, compare paths in all specs matching base and revision globs")
	enumWithOptions(&cmd, newEnumValue(formatters.SupportedFormatsByContentType("breaking-changes"), string(formatters.FormatText), &flags.format), "format", "f", "output format")
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
	enumWithOptions(&cmd, newEnumValue([]string{LangEn, LangRu}, LangDefault, &flags.lang), "lang", "l", "language for localized output")
	cmd.PersistentFlags().StringVarP(&flags.errIgnoreFile, "err-ignore", "", "", "configuration file for ignoring errors")
	cmd.PersistentFlags().StringVarP(&flags.warnIgnoreFile, "warn-ignore", "", "", "configuration file for ignoring warnings")
	enumWithOptions(&cmd, newEnumSliceValue(checker.GetOptionalChecks(), nil, &flags.includeChecks), "include-checks", "i", "comma-separated list of optional checks (run 'oasdiff checks' to see options)")
	cmd.PersistentFlags().IntVarP(&flags.deprecationDaysBeta, "deprecation-days-beta", "", checker.BetaDeprecationDays, "min days required between deprecating a beta resource and removing it")
	cmd.PersistentFlags().IntVarP(&flags.deprecationDaysStable, "deprecation-days-stable", "", checker.StableDeprecationDays, "min days required between deprecating a stable resource and removing it")

	return &cmd
}

func runBreakingChanges(flags *ChangelogFlags, stdout io.Writer) (bool, *ReturnError) {
	return getChangelog(flags, stdout, checker.WARN)
}

func outputBreakingChanges(config checker.Config, format string, lang string, stdout io.Writer, errs checker.Changes, level checker.Level) *ReturnError {
	// formatter lookup
	formatter, err := formatters.Lookup(format, formatters.FormatterOpts{
		Language: lang,
	})
	if err != nil {
		return getErrUnsupportedBreakingChangesFormat(format)
	}

	// render
	bytes, err := formatter.RenderBreakingChanges(nil, nil, errs, formatters.RenderOpts{})
	if err != nil {
		return getErrFailedPrint("diff "+format, err)
	}

	// print output
	_, _ = fmt.Fprintf(stdout, "%s\n", bytes)

	return nil
}
