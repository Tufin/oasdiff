package internal

import (
	"fmt"
	"io"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/spf13/cobra"
	"github.com/tufin/oasdiff/checker"
	"github.com/tufin/oasdiff/checker/localizations"
	"github.com/tufin/oasdiff/diff"
)

func getChangelogCmd() *cobra.Command {

	flags := ChangelogFlags{}

	cmd := cobra.Command{
		Use:   "changelog",
		Short: "Display changelog",
		PreRun: func(cmd *cobra.Command, args []string) {
			if returnErr := flags.validate(); returnErr != nil {
				exit(false, returnErr, cmd.ErrOrStderr())
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			failEmpty, returnErr := runChangelog(&flags, cmd.OutOrStdout())
			exit(failEmpty, returnErr, cmd.ErrOrStderr())
		},
	}

	// common
	cmd.PersistentFlags().BoolVarP(&flags.composed, "composed", "c", false, "work in 'composed' mode, compare paths in all specs matching base and revision globs")
	cmd.PersistentFlags().StringVarP(&flags.base, "base", "b", "", "path or URL (or a glob in Composed mode) of original OpenAPI spec in YAML or JSON format")
	cmd.PersistentFlags().StringVarP(&flags.revision, "revision", "r", "", "path or URL (or a glob in Composed mode) of revised OpenAPI spec in YAML or JSON format")
	cmd.PersistentFlags().StringVarP(&flags.format, "format", "f", "text", "output format: yaml, json, text")
	cmd.PersistentFlags().StringSliceVarP(&flags.excludeElements, "exclude-elements", "", nil, "comma-separated list of elements to exclude from diff")
	cmd.PersistentFlags().StringVarP(&flags.matchPath, "match-path", "", "", "include only paths that match this regular expression")
	cmd.PersistentFlags().StringVarP(&flags.filterExtension, "filter-extension", "", "", "exclude paths and operations with an OpenAPI Extension matching this regular expression")
	cmd.PersistentFlags().IntVarP(&flags.circularReferenceCounter, "max-circular-dep", "", 5, "maximum allowed number of circular dependencies between objects in OpenAPI specs")
	cmd.PersistentFlags().StringVarP(&flags.prefixBase, "prefix-base", "", "", "add this prefix to paths in 'base' spec before comparison")
	cmd.PersistentFlags().StringVarP(&flags.prefixRevision, "prefix-revision", "", "", "add this prefix to paths in 'revision' spec before comparison")
	cmd.PersistentFlags().StringVarP(&flags.stripPrefixBase, "strip-prefix-base", "", "", "strip this prefix from paths in 'base' spec before comparison")
	cmd.PersistentFlags().StringVarP(&flags.stripPrefixRevision, "strip-prefix-revision", "", "", "strip this prefix from paths in 'revision' spec before comparison")
	cmd.PersistentFlags().BoolVarP(&flags.matchPathParams, "match-path-params", "", false, "include path parameter names in endpoint matching")

	cmd.MarkPersistentFlagRequired("base")
	cmd.MarkPersistentFlagRequired("revision")

	// specific for breaking-changes
	cmd.PersistentFlags().VarP(&flags.failOn, "fail-on", "", "exit with return code 1 when output includes errors with this level or higher")
	// level
	// err-ignore
	// warn-ignore
	// info-ignore
	// deprecation-days
	// lang
	cmd.PersistentFlags().StringSliceVarP(&flags.includeChecks, "include-checks", "", nil, "comma-separated list of optional breaking-changes checks")
	return &cmd
}

func runChangelog(flags *ChangelogFlags, stdout io.Writer) (bool, *ReturnError) {
	return getChangelog(flags, stdout, checker.INFO)
}

func getChangelog(flags *ChangelogFlags, stdout io.Writer, level checker.Level) (bool, *ReturnError) {

	openapi3.CircularReferenceCounter = flags.circularReferenceCounter

	diffConfig := flags.toConfig()

	var diffReport *diff.Diff
	var operationsSources *diff.OperationsSourcesMap

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true

	if flags.composed {
		var err *ReturnError
		if diffReport, operationsSources, err = composedDiff(loader, flags.base, flags.revision, diffConfig); err != nil {
			return false, err
		}
	} else {
		var err *ReturnError
		if diffReport, operationsSources, err = normalDiff(loader, flags.base, flags.revision, diffConfig); err != nil {
			return false, err
		}
	}

	bcConfig := checker.GetChecks(flags.includeChecks)
	bcConfig.Localizer = *localizations.New(flags.lang, "en")

	errs, returnErr := filterIgnored(
		checker.CheckBackwardCompatibilityUntilLevel(bcConfig, diffReport, operationsSources, level),
		flags.warnIgnoreFile, flags.errIgnoreFile)

	if returnErr != nil {
		return false, returnErr
	}

	if returnErr := outputChangelog(bcConfig, flags.format, stdout, errs); returnErr != nil {
		return false, returnErr
	}

	if flags.failOn != "" {
		level, err := flags.failOn.ToLevel()
		if err != nil {
			return false, getErrInvalidFlags(fmt.Errorf("invalid fail-on value %s", flags.failOn))
		}
		return errs.HasLevelOrHigher(level), nil
	}

	return false, nil
}

func filterIgnored(errs checker.BackwardCompatibilityErrors, warnIgnoreFile string, errIgnoreFile string) (checker.BackwardCompatibilityErrors, *ReturnError) {

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

func outputChangelog(config checker.BackwardCompatibilityCheckConfig, format string, stdout io.Writer, errs checker.BackwardCompatibilityErrors) *ReturnError {
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
			fmt.Fprintf(stdout, config.Localizer.Get("messages.total-errors"), len(errs))
		}

		for _, bcerr := range errs {
			fmt.Fprintf(stdout, "%s\n\n", bcerr.PrettyErrorText(config.Localizer))
		}
	default:
		return getErrUnsupportedBreakingChangesFormat(format)
	}

	return nil
}
