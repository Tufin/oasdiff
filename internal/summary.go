package internal

import (
	"io"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/spf13/cobra"
	"github.com/tufin/oasdiff/diff"
)

func getSummaryCmd() *cobra.Command {
	flags := DiffFlags{}

	cmd := cobra.Command{
		Use:   "summary base-spec revised-spec [flags]",
		Short: "Generate a diff summary",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {

			flags.base = args[0]
			flags.revision = args[1]

			// by now flags have been parsed successfully so we don't need to show usage on any errors
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
	cmd.PersistentFlags().VarP(newEnumValue([]string{FormatYAML, FormatJSON}, FormatYAML, &flags.format), "format", "f", "output format: yaml or json")
	cmd.PersistentFlags().VarP(newEnumSliceValue(diff.ExcludeDiffOptions, nil, &flags.excludeElements), "exclude-elements", "", "comma-separated list of elements to exclude")
	cmd.PersistentFlags().StringVarP(&flags.matchPath, "match-path", "", "", "include only paths that match this regular expression")
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
	switch format {
	case FormatYAML:
		if err := printYAML(stdout, diffReport.GetSummary()); err != nil {
			return getErrFailedPrint("summary", err)
		}
	case FormatJSON:

		if err := printJSON(stdout, diffReport.GetSummary()); err != nil {
			return getErrFailedPrint("summary", err)
		}
	default:
		return getErrUnsupportedDiffFormat(format)
	}

	return nil

}
