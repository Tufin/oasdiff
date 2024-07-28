package internal

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"
	"github.com/tufin/oasdiff/diff"
	"github.com/tufin/oasdiff/formatters"
)

const summaryCmd = "summary"

func getSummaryCmd() *cobra.Command {

	flags := NewDiffFlags()

	cmd := cobra.Command{
		Use:   "summary base revision [flags]",
		Short: "Generate a diff summary",
		Long:  "Display a summary of changes between base and revision specs." + specHelp,
		Args:  getParseArgs(flags),
		RunE:  getRun(flags, runSummary),
	}

	addCommonDiffFlags(&cmd, flags)
	enumWithOptions(&cmd, newEnumSliceValue(diff.ExcludeDiffOptions, nil), "exclude-elements", "e", "elements to exclude")
	enumWithOptions(&cmd, newEnumValue(formatters.SupportedFormatsByContentType(formatters.OutputSummary), string(formatters.FormatYAML)), "format", "f", "output format")
	cmd.PersistentFlags().BoolP("fail-on-diff", "", false, "exit with return code 1 when any change is found")

	bindViperFlags(&cmd, flags.getViper())

	return &cmd
}

func runSummary(flags Flags, stdout io.Writer) (bool, *ReturnError) {

	diffResult, err := calcDiff(flags)
	if err != nil {
		return false, err
	}

	if err := outputSummary(stdout, diffResult.diffReport, flags.getFormat()); err != nil {
		return false, err
	}

	return flags.getFailOnDiff() && !diffResult.diffReport.Empty(), nil
}

func outputSummary(stdout io.Writer, diffReport *diff.Diff, format string) *ReturnError {
	// formatter lookup
	formatter, err := formatters.Lookup(format, formatters.DefaultFormatterOpts())
	if err != nil {
		return getErrUnsupportedFormat(format, summaryCmd)
	}

	// render
	bytes, err := formatter.RenderSummary(diffReport, formatters.NewRenderOpts())
	if err != nil {
		return getErrFailedPrint(summaryCmd+" "+format, err)
	}

	// print output
	_, _ = fmt.Fprintf(stdout, "%s\n", bytes)

	return nil
}
