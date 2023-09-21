package internal

import (
	"fmt"
	"io"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/spf13/cobra"
	"github.com/tufin/oasdiff/checker"
	"github.com/tufin/oasdiff/diff"
	"github.com/tufin/oasdiff/flatten"
	"github.com/tufin/oasdiff/formatters"
	"github.com/tufin/oasdiff/load"
)

func getDiffCmd() *cobra.Command {

	flags := DiffFlags{}

	cmd := cobra.Command{
		Use:   "diff base revision [flags]",
		Short: "Generate a diff report",
		Long: `Generate a diff report between base and revision specs.
Base and revision can be a path to a file or a URL.
In 'composed' mode, base and revision can be a glob and oasdiff will compare matching endpoints between the two sets of files.
`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {

			flags.base = args[0]
			flags.revision = args[1]

			// by now flags have been parsed successfully, so we don't need to show usage on any errors
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
	enumWithOptions(&cmd, newEnumValue(formatters.SupportedFormatsByContentType("diff"), string(formatters.FormatYAML), &flags.format), "format", "f", "output format")
	cmd.PersistentFlags().VarP(newEnumSliceValue(diff.ExcludeDiffOptions, nil, &flags.excludeElements), "exclude-elements", "e", "comma-separated list of elements to exclude")
	cmd.PersistentFlags().StringVarP(&flags.matchPath, "match-path", "p", "", "include only paths that match this regular expression")
	cmd.PersistentFlags().StringVarP(&flags.filterExtension, "filter-extension", "", "", "exclude paths and operations with an OpenAPI Extension matching this regular expression")
	cmd.PersistentFlags().IntVarP(&flags.circularReferenceCounter, "max-circular-dep", "", 5, "maximum allowed number of circular dependencies between objects in OpenAPI specs")
	cmd.PersistentFlags().StringVarP(&flags.prefixBase, "prefix-base", "", "", "add this prefix to paths in base-spec before comparison")
	cmd.PersistentFlags().StringVarP(&flags.prefixRevision, "prefix-revision", "", "", "add this prefix to paths in revised-spec before comparison")
	cmd.PersistentFlags().StringVarP(&flags.stripPrefixBase, "strip-prefix-base", "", "", "strip this prefix from paths in base-spec before comparison")
	cmd.PersistentFlags().StringVarP(&flags.stripPrefixRevision, "strip-prefix-revision", "", "", "strip this prefix from paths in revised-spec before comparison")
	cmd.PersistentFlags().BoolVarP(&flags.includePathParams, "include-path-params", "", false, "include path parameter names in endpoint matching")
	cmd.PersistentFlags().BoolVarP(&flags.flatten, "flatten", "", false, "merge subschemas under allOf before diff")
	cmd.PersistentFlags().BoolVarP(&flags.failOnDiff, "fail-on-diff", "o", false, "exit with return code 1 when any change is found")

	return &cmd
}

func runDiff(flags *DiffFlags, stdout io.Writer) (bool, *ReturnError) {

	openapi3.CircularReferenceCounter = flags.circularReferenceCounter

	if flags.format == string(formatters.FormatJSON) {
		flags.excludeElements = append(flags.excludeElements, diff.ExcludeEndpointsOption)
	}

	diffReport, _, err := calcDiff(flags)
	if err != nil {
		return false, err
	}

	if err := outputDiff(stdout, nil, diffReport, flags.format); err != nil {
		return false, err
	}

	return flags.failOnDiff && !diffReport.Empty(), nil
}

func outputDiff(stdout io.Writer, checks []checker.BackwardCompatibilityCheck, diffReport *diff.Diff, format string) *ReturnError {
	// formatter lookup
	formatter, err := formatters.Lookup(format, formatters.DefaultFormatterOpts())
	if err != nil {
		return getErrUnsupportedDiffFormat(format)
	}

	// render
	bytes, err := formatter.RenderDiff(diffReport, nil, formatters.RenderOpts{})
	if err != nil {
		return getErrFailedPrint("diff "+format, err)
	}

	// print output
	_, _ = fmt.Fprintf(stdout, "%s\n", bytes)

	return nil
}

func calcDiff(flags Flags) (*diff.Diff, *diff.OperationsSourcesMap, *ReturnError) {

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true

	if flags.getComposed() {
		return composedDiff(loader, flags)
	}

	return normalDiff(loader, flags)
}

func normalDiff(loader load.Loader, flags Flags) (*diff.Diff, *diff.OperationsSourcesMap, *ReturnError) {
	s1, err := load.LoadSpecInfo(loader, flags.getBase())
	if err != nil {
		return nil, nil, getErrFailedToLoadSpec("base", flags.getBase(), err)
	}

	s2, err := load.LoadSpecInfo(loader, flags.getRevision())
	if err != nil {
		return nil, nil, getErrFailedToLoadSpec("revision", flags.getRevision(), err)
	}

	if flags.getFlatten() {
		if err := mergeAllOf("base", []*load.SpecInfo{s1}); err != nil {
			return nil, nil, err
		}

		if err := mergeAllOf("revision", []*load.SpecInfo{s2}); err != nil {
			return nil, nil, err
		}
	}

	diffReport, operationsSources, err := diff.GetWithOperationsSourcesMap(flags.toConfig(), s1, s2)
	if err != nil {
		return nil, nil, getErrDiffFailed(err)
	}

	return diffReport, operationsSources, nil
}

func composedDiff(loader load.Loader, flags Flags) (*diff.Diff, *diff.OperationsSourcesMap, *ReturnError) {
	s1, err := load.FromGlob(loader, flags.getBase())
	if err != nil {
		return nil, nil, getErrFailedToLoadSpecs("base", flags.getBase(), err)
	}

	s2, err := load.FromGlob(loader, flags.getRevision())
	if err != nil {
		return nil, nil, getErrFailedToLoadSpecs("revision", flags.getRevision(), err)
	}

	if flags.getFlatten() {
		if err := mergeAllOf("base", s1); err != nil {
			return nil, nil, err
		}

		if err := mergeAllOf("revision", s2); err != nil {
			return nil, nil, err
		}
	}

	diffReport, operationsSources, err := diff.GetPathsDiff(flags.toConfig(), s1, s2)
	if err != nil {
		return nil, nil, getErrDiffFailed(err)
	}

	return diffReport, operationsSources, nil
}

func mergeAllOf(title string, specInfos []*load.SpecInfo) *ReturnError {

	var err error

	for _, specInfo := range specInfos {
		if specInfo.Spec, err = flatten.MergeSpec(specInfo.Spec); err != nil {
			return getErrFailedToFlattenSpec(title, specInfo.Url, err)
		}
	}

	return nil
}
