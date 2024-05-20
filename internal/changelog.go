package internal

import (
	"fmt"
	"io"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/spf13/cobra"
	"github.com/tufin/oasdiff/checker"
	"github.com/tufin/oasdiff/formatters"
	"github.com/tufin/oasdiff/load"
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

	addCommonDiffFlags(&cmd, &flags)
	addCommonBreakingFlags(&cmd, &flags)
	enumWithOptions(&cmd, newEnumValue(formatters.SupportedFormatsByContentType(formatters.OutputChangelog), string(formatters.FormatText), &flags.format), "format", "f", "output format")
	enumWithOptions(&cmd, newEnumValue([]string{LevelErr, LevelWarn, LevelInfo}, "", &flags.failOn), "fail-on", "o", "exit with return code 1 when output includes errors with this level or higher")
	enumWithOptions(&cmd, newEnumValue([]string{LevelErr, LevelWarn, LevelInfo}, LevelInfo, &flags.level), "level", "", "output errors with this level or higher")

	return &cmd
}

func enumWithOptions(cmd *cobra.Command, value enumVal, name, shorthand, usage string) {
	cmd.PersistentFlags().VarP(value, name, shorthand, usage+": "+value.listOf())
}

func runChangelog(flags Flags, stdout io.Writer) (bool, *ReturnError) {
	level, err := checker.NewLevel(flags.getLevel())
	if err != nil {
		return false, getErrInvalidFlags(fmt.Errorf("invalid level value %s", flags.getLevel()))
	}

	return getChangelog(flags, stdout, level)
}

func getChangelog(flags Flags, stdout io.Writer, level checker.Level) (bool, *ReturnError) {

	openapi3.CircularReferenceCounter = flags.getCircularReferenceCounter()

	diffResult, err := calcDiff(flags)
	if err != nil {
		return false, err
	}

	bcConfig := checker.GetAllChecks(flags.getIncludeChecks(), flags.getDeprecationDaysBeta(), flags.getDeprecationDaysStable())

	errs, returnErr := filterIgnored(
		checker.CheckBackwardCompatibilityUntilLevel(bcConfig, diffResult.diffReport, diffResult.operationsSources, level),
		flags.getWarnIgnoreFile(),
		flags.getErrIgnoreFile(),
		checker.NewLocalizer(flags.getLang()))

	if returnErr != nil {
		return false, returnErr
	}

	if returnErr := outputChangelog(flags.getFormat(), flags.getLang(), flags.getColor(), stdout, errs, diffResult.specInfoPair); returnErr != nil {
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

func filterIgnored(errs checker.Changes, warnIgnoreFile string, errIgnoreFile string, l checker.Localizer) (checker.Changes, *ReturnError) {

	if warnIgnoreFile != "" {
		var err error
		errs, err = checker.ProcessIgnoredBackwardCompatibilityErrors(checker.WARN, errs, warnIgnoreFile, l)
		if err != nil {
			return nil, getErrCantProcessIgnoreFile("warn", err)
		}
	}

	if errIgnoreFile != "" {
		var err error
		errs, err = checker.ProcessIgnoredBackwardCompatibilityErrors(checker.ERR, errs, errIgnoreFile, l)
		if err != nil {
			return nil, getErrCantProcessIgnoreFile("err", err)
		}
	}

	return errs, nil
}

func outputChangelog(format string, lang string, color string, stdout io.Writer, errs checker.Changes, specInfoPair *load.SpecInfoPair) *ReturnError {
	// formatter lookup
	formatter, err := formatters.Lookup(format, formatters.FormatterOpts{
		Language: lang,
	})
	if err != nil {
		return getErrUnsupportedChangelogFormat(format)
	}

	// render
	colorMode, err := checker.NewColorMode(color)
	if err != nil {
		return getErrInvalidColorMode(err)
	}

	bytes, err := formatter.RenderChangelog(errs, formatters.RenderOpts{ColorMode: colorMode}, specInfoPair)
	if err != nil {
		return getErrFailedPrint("changelog "+format, err)
	}

	// print output
	_, _ = fmt.Fprintf(stdout, "%s\n", bytes)

	return nil
}
