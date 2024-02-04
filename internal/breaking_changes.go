package internal

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"
	"github.com/tufin/oasdiff/checker"
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

	addCommonDiffFlags(&cmd, &flags)
	addCommonBreakingFlags(&cmd, &flags)
	enumWithOptions(&cmd, newEnumValue(formatters.SupportedFormatsByContentType(formatters.OutputBreaking), string(formatters.FormatText), &flags.format), "format", "f", "output format")
	enumWithOptions(&cmd, newEnumValue([]string{LevelErr, LevelWarn}, "", &flags.failOn), "fail-on", "o", "exit with return code 1 when output includes errors with this level or higher")

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
