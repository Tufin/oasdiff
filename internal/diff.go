package internal

import (
	"fmt"
	"io"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/spf13/cobra"
	"github.com/tufin/oasdiff/diff"
	"github.com/tufin/oasdiff/load"
	"github.com/tufin/oasdiff/report"
)

func getDiffCmd() *cobra.Command {

	flags := DiffFlags{}

	cmd := cobra.Command{
		Use:   "diff base revision [flags]",
		Short: "Generate a diff report",
		Long: `Generate a diff report between base and revision specs.
Base and revision can be a path to a file or a URL.
In 'composed' mode, base and revision can be a glob and oasdiff will compare mathcing endpoints between the two sets of files.
`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {

			flags.base = args[0]
			flags.revision = args[1]

			// by now flags have been parsed successfully so we don't need to show usage on any errors
			cmd.Root().SilenceUsage = true

			failEmpty, err := runDiff(&flags, cmd.OutOrStdout())
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
	cmd.PersistentFlags().VarP(newEnumValue([]string{FormatYAML, FormatJSON, FormatText, FormatHTML}, FormatYAML, &flags.format), "format", "f", "output format: yaml, json, text or html")
	cmd.PersistentFlags().VarP(newEnumSliceValue(diff.ExcludeDiffOptions, nil, &flags.excludeElements), "exclude-elements", "e", "comma-separated list of elements to exclude")
	cmd.PersistentFlags().StringVarP(&flags.matchPath, "match-path", "p", "", "include only paths that match this regular expression")
	cmd.PersistentFlags().StringVarP(&flags.filterExtension, "filter-extension", "", "", "exclude paths and operations with an OpenAPI Extension matching this regular expression")
	cmd.PersistentFlags().IntVarP(&flags.circularReferenceCounter, "max-circular-dep", "", 5, "maximum allowed number of circular dependencies between objects in OpenAPI specs")
	cmd.PersistentFlags().StringVarP(&flags.prefixBase, "prefix-base", "", "", "add this prefix to paths in base-spec before comparison")
	cmd.PersistentFlags().StringVarP(&flags.prefixRevision, "prefix-revision", "", "", "add this prefix to paths in revised-spec before comparison")
	cmd.PersistentFlags().StringVarP(&flags.stripPrefixBase, "strip-prefix-base", "", "", "strip this prefix from paths in base-spec before comparison")
	cmd.PersistentFlags().StringVarP(&flags.stripPrefixRevision, "strip-prefix-revision", "", "", "strip this prefix from paths in revised-spec before comparison")
	cmd.PersistentFlags().BoolVarP(&flags.includePathParams, "include-path-params", "", false, "include path parameter names in endpoint matching")
	cmd.PersistentFlags().BoolVarP(&flags.failOnDiff, "fail-on-diff", "o", false, "exit with return code 1 when any change is found")
	cmd.PersistentFlags().BoolVarP(&flags.mergeAllOf, "merge-all-of", "m", false, "merge subschemas under allOf before diff")

	return &cmd
}

func runDiff(flags *DiffFlags, stdout io.Writer) (bool, *ReturnError) {

	openapi3.CircularReferenceCounter = flags.circularReferenceCounter

	if flags.format == FormatJSON {
		flags.excludeElements = append(flags.excludeElements, diff.ExcludeEndpointsOption)
	}

	diffReport, _, err := calcDiff(flags)
	if err != nil {
		return false, err
	}

	if err := outputDiff(stdout, diffReport, flags.format); err != nil {
		return false, err
	}

	return flags.failOnDiff && !diffReport.Empty(), nil
}

func outputDiff(stdout io.Writer, diffReport *diff.Diff, format string) *ReturnError {
	switch format {
	case FormatYAML:
		if err := printYAML(stdout, diffReport); err != nil {
			return getErrFailedPrint("diff YAML", err)
		}
	case FormatJSON:
		if err := printJSON(stdout, diffReport); err != nil {
			return getErrFailedPrint("diff JSON", err)
		}
	case FormatText:
		fmt.Fprintf(stdout, "%s", report.GetTextReportAsString(diffReport))
	case FormatHTML:
		html, err := report.GetHTMLReportAsString(diffReport)
		if err != nil {
			return getErrFailedGenerateHTML(err)
		}
		fmt.Fprintf(stdout, "%s", html)
	default:
		return getErrUnsupportedDiffFormat(format)
	}

	return nil
}

func calcDiff(flags Flags) (*diff.Diff, *diff.OperationsSourcesMap, *ReturnError) {

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true

	if flags.getComposed() {
		return composedDiff(loader, flags.getBase(), flags.getRevision(), flags.toConfig())
	}

	return normalDiff(loader, flags.getBase(), flags.getRevision(), flags.toConfig())
}

func normalDiff(loader load.Loader, base, revision string, config *diff.Config) (*diff.Diff, *diff.OperationsSourcesMap, *ReturnError) {
	s1, err := load.LoadSpecInfo(loader, base)
	if err != nil {
		return nil, nil, getErrFailedToLoadSpec("base", base, err)
	}
	s2, err := load.LoadSpecInfo(loader, revision)
	if err != nil {
		return nil, nil, getErrFailedToLoadSpec("revision", revision, err)
	}

	diffReport, operationsSources, err := diff.GetWithOperationsSourcesMap(config, s1, s2)
	if err != nil {
		return nil, nil, getErrDiffFailed(err)
	}

	return diffReport, operationsSources, nil
}

func composedDiff(loader load.Loader, base, revision string, config *diff.Config) (*diff.Diff, *diff.OperationsSourcesMap, *ReturnError) {
	s1, err := load.FromGlob(loader, base)
	if err != nil {
		return nil, nil, getErrFailedToLoadSpec("base", base, err)
	}

	s2, err := load.FromGlob(loader, revision)
	if err != nil {
		return nil, nil, getErrFailedToLoadSpec("revision", revision, err)
	}
	diffReport, operationsSources, err := diff.GetPathsDiff(config, s1, s2)
	if err != nil {
		return nil, nil, getErrDiffFailed(err)
	}

	return diffReport, operationsSources, nil
}
