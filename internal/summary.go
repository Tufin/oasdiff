package internal

import (
	"fmt"
	"io"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/spf13/cobra"
	"github.com/tufin/oasdiff/diff"
	"github.com/tufin/oasdiff/formatters"
)

func getSummaryCmd() *cobra.Command {
	flags := DiffFlags{}

	cmd := cobra.Command{
		Use:   "summary base revision [flags]",
		Short: "Generate a diff summary",
		Long: `Display a summary of changes between base and revision specs.
Base and revision can be a path to a file or a URL.
In 'composed' mode, base and revision can be a glob and oasdiff will compare matching endpoints between the two sets of files.
`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {

			flags.base = args[0]
			flags.revision = args[1]

			// by now flags have been parsed successfully, so we don't need to show usage on any errors
			cmd.Root().SilenceUsage = true

			failEmpty, err := runSummary(&flags, cmd.OutOrStdout())
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
	enumWithOptions(&cmd, newEnumValue(formatters.SupportedFormatsByContentType("summary"), string(formatters.FormatYAML), &flags.format), "format", "f", "output format")
	cmd.PersistentFlags().VarP(newEnumSliceValue(diff.ExcludeDiffOptions, nil, &flags.excludeElements), "exclude-elements", "e", "comma-separated list of elements to exclude")
	cmd.PersistentFlags().StringVarP(&flags.matchPath, "match-path", "p", "", "include only paths that match this regular expression")
	cmd.PersistentFlags().StringVarP(&flags.filterExtension, "filter-extension", "", "", "exclude paths and operations with an OpenAPI Extension matching this regular expression")
	cmd.PersistentFlags().IntVarP(&flags.circularReferenceCounter, "max-circular-dep", "", 5, "maximum allowed number of circular dependencies between objects in OpenAPI specs")
	cmd.PersistentFlags().StringVarP(&flags.prefixBase, "prefix-base", "", "", "add this prefix to paths in base-spec before comparison")
	cmd.PersistentFlags().StringVarP(&flags.prefixRevision, "prefix-revision", "", "", "add this prefix to paths in revised-spec before comparison")
	cmd.PersistentFlags().StringVarP(&flags.stripPrefixBase, "strip-prefix-base", "", "", "strip this prefix from paths in base-spec before comparison")
	cmd.PersistentFlags().StringVarP(&flags.stripPrefixRevision, "strip-prefix-revision", "", "", "strip this prefix from paths in revised-spec before comparison")
	cmd.PersistentFlags().BoolVarP(&flags.includePathParams, "include-path-params", "", false, "include path parameter names in endpoint matching")
	cmd.PersistentFlags().BoolVarP(&flags.failOnDiff, "fail-on-diff", "", false, "exit with return code 1 when any change is found")

	return &cmd
}

func runSummary(flags *DiffFlags, stdout io.Writer) (bool, *ReturnError) {

	openapi3.CircularReferenceCounter = flags.circularReferenceCounter

	diffReport, _, err := calcDiff(flags)
	if err != nil {
		return false, err
	}

	if err := outputSummary(stdout, diffReport, flags.format); err != nil {
		return false, err
	}

	return flags.failOnDiff && !diffReport.Empty(), nil
}

func outputSummary(stdout io.Writer, diffReport *diff.Diff, format string) *ReturnError {
	// formatter lookup
	formatter, err := formatters.Lookup(format, formatters.DefaultFormatterOpts())
	if err != nil {
		return getErrUnsupportedSummaryFormat(format)
	}

	// render
	bytes, err := formatter.RenderSummary(nil, diffReport, nil, formatters.RenderOpts{})
	if err != nil {
		return getErrFailedPrint("summary "+format, err)
	}

	// print output
	_, _ = fmt.Fprintf(stdout, "%s\n", bytes)

	return nil
}
