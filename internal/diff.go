package internal

import (
	"fmt"
	"io"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/spf13/cobra"
	"github.com/tufin/oasdiff/diff"
	"github.com/tufin/oasdiff/flatten/allof"
	"github.com/tufin/oasdiff/formatters"
	"github.com/tufin/oasdiff/load"
)

func getDiffCmd() *cobra.Command {

	flags := DiffFlags{}

	cmd := cobra.Command{
		Use:   "diff base revision [flags]",
		Short: "Generate a diff report",
		Long:  "Generate a diff report between base and revision specs." + specHelp,
		Args:  getParseArgs(&flags),
		RunE:  getRun(&flags, runDiff),
	}

	cmd.PersistentFlags().BoolVarP(&flags.composed, "composed", "c", false, "work in 'composed' mode, compare paths in all specs matching base and revision globs")
	enumWithOptions(&cmd, newEnumValue(formatters.SupportedFormatsByContentType(formatters.OutputDiff), string(formatters.FormatYAML), &flags.format), "format", "f", "output format")
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
	cmd.PersistentFlags().BoolVarP(&flags.failOnDiff, "fail-on-diff", "o", false, "exit with return code 1 when any change is found")

	return &cmd
}

func runDiff(flags Flags, stdout io.Writer) (bool, *ReturnError) {

	openapi3.CircularReferenceCounter = flags.getCircularReferenceCounter()

	if flags.getFormat() == string(formatters.FormatJSON) {
		flags.addExcludeElements(diff.ExcludeEndpointsOption)
	}

	diffResult, err := calcDiff(flags)
	if err != nil {
		return false, err
	}

	if err := outputDiff(stdout, diffResult.diffReport, flags.getFormat()); err != nil {
		return false, err
	}

	return flags.getFailOnDiff() && !diffResult.diffReport.Empty(), nil
}

func outputDiff(stdout io.Writer, diffReport *diff.Diff, format string) *ReturnError {
	// formatter lookup
	formatter, err := formatters.Lookup(format, formatters.DefaultFormatterOpts())
	if err != nil {
		return getErrUnsupportedDiffFormat(format)
	}

	// render
	bytes, err := formatter.RenderDiff(diffReport, formatters.NewRenderOpts())
	if err != nil {
		return getErrFailedPrint("diff "+format, err)
	}

	// print output
	_, _ = fmt.Fprintf(stdout, "%s\n", bytes)

	return nil
}

func calcDiff(flags Flags) (*diffResult, *ReturnError) {

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true

	if flags.getComposed() {
		return composedDiff(loader, flags)
	}

	return normalDiff(loader, flags)
}

type diffResult struct {
	diffReport        *diff.Diff
	operationsSources *diff.OperationsSourcesMap
	specInfoPair      *load.SpecInfoPair
}

func newDiffResult(d *diff.Diff, o *diff.OperationsSourcesMap, s *load.SpecInfoPair) *diffResult {
	return &diffResult{
		diffReport:        d,
		operationsSources: o,
		specInfoPair:      s,
	}
}

func normalDiff(loader load.Loader, flags Flags) (*diffResult, *ReturnError) {
	s1, err := load.LoadSpecInfo(loader, flags.getBase())
	if err != nil {
		return nil, getErrFailedToLoadSpec("base", flags.getBase(), err)
	}

	s2, err := load.LoadSpecInfo(loader, flags.getRevision())
	if err != nil {
		return nil, getErrFailedToLoadSpec("revision", flags.getRevision(), err)
	}

	if flags.getBase().IsStdin() && flags.getRevision().IsStdin() {
		// io.ReadAll can only read stdin once, so in this edge case, we copy base into revision
		s2.Spec = s1.Spec
	}

	if flags.getFlatten() {
		if err := mergeAllOf("base", []*load.SpecInfo{s1}, flags.getBase()); err != nil {
			return nil, err
		}

		if err := mergeAllOf("revision", []*load.SpecInfo{s2}, flags.getRevision()); err != nil {
			return nil, err
		}
	}

	diffReport, operationsSources, err := diff.GetWithOperationsSourcesMap(flags.toConfig(), s1, s2)
	if err != nil {
		return nil, getErrDiffFailed(err)
	}

	return newDiffResult(diffReport, operationsSources, load.NewSpecInfoPair(s1, s2)), nil
}

func composedDiff(loader load.Loader, flags Flags) (*diffResult, *ReturnError) {
	s1, err := load.FromGlob(loader, flags.getBase().Path)
	if err != nil {
		return nil, getErrFailedToLoadSpecs("base", flags.getBase().Path, err)
	}

	s2, err := load.FromGlob(loader, flags.getRevision().Path)
	if err != nil {
		return nil, getErrFailedToLoadSpecs("revision", flags.getRevision().Path, err)
	}

	if flags.getFlatten() {
		if err := mergeAllOf("base", s1, flags.getBase()); err != nil {
			return nil, err
		}

		if err := mergeAllOf("revision", s2, flags.getRevision()); err != nil {
			return nil, err
		}
	}

	diffReport, operationsSources, err := diff.GetPathsDiff(flags.toConfig(), s1, s2)
	if err != nil {
		return nil, getErrDiffFailed(err)
	}

	return newDiffResult(diffReport, operationsSources, nil), nil
}

func mergeAllOf(title string, specInfos []*load.SpecInfo, source *load.Source) *ReturnError {

	var err error

	for _, specInfo := range specInfos {
		if specInfo.Spec, err = allof.MergeSpec(specInfo.Spec); err != nil {
			return getErrFailedToFlattenSpec(title, source, err)
		}
	}

	return nil
}
