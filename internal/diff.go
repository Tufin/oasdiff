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
		Use:   "diff",
		Short: "Generate a diff report",
		PreRun: func(cmd *cobra.Command, args []string) {
			if returnErr := flags.validate(); returnErr != nil {
				exit(false, returnErr, cmd.ErrOrStderr())
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			failEmpty, returnErr := runDiff(&flags, cmd.OutOrStdout())
			exit(failEmpty, returnErr, cmd.ErrOrStderr())
		},
	}

	cmd.PersistentFlags().BoolVarP(&flags.composed, "composed", "c", false, "work in 'composed' mode, compare paths in all specs matching base and revision globs")
	cmd.PersistentFlags().StringVarP(&flags.base, "base", "b", "", "path or URL (or a glob in Composed mode) of original OpenAPI spec in YAML or JSON format")
	cmd.PersistentFlags().StringVarP(&flags.revision, "revision", "r", "", "path or URL (or a glob in Composed mode) of revised OpenAPI spec in YAML or JSON format")
	cmd.PersistentFlags().StringVarP(&flags.format, "format", "f", "yaml", "output format: yaml, json, text or html")
	cmd.PersistentFlags().StringSliceVarP(&flags.excludeElements, "exclude-elements", "", nil, "comma-separated list of elements to exclude from diff")
	cmd.PersistentFlags().StringVarP(&flags.matchPath, "match-path", "", "", "include only paths that match this regular expression")
	cmd.PersistentFlags().StringVarP(&flags.filterExtension, "filter-extension", "", "", "exclude paths and operations with an OpenAPI Extension matching this regular expression")
	cmd.PersistentFlags().IntVarP(&flags.circularReferenceCounter, "max-circular-dep", "", 5, "maximum allowed number of circular dependencies between objects in OpenAPI specs")
	cmd.PersistentFlags().StringVarP(&flags.prefixBase, "prefix-base", "", "", "add this prefix to paths in 'base' spec before comparison")
	cmd.PersistentFlags().StringVarP(&flags.prefixRevision, "prefix-revision", "", "", "add this prefix to paths in 'revision' spec before comparison")
	cmd.PersistentFlags().StringVarP(&flags.stripPrefixBase, "strip-prefix-base", "", "", "strip this prefix from paths in 'base' spec before comparison")
	cmd.PersistentFlags().StringVarP(&flags.stripPrefixRevision, "strip-prefix-revision", "", "", "strip this prefix from paths in 'revision' spec before comparison")
	cmd.PersistentFlags().BoolVarP(&flags.matchPathParams, "match-path-params", "", false, "include path parameter names in endpoint matching")

	return &cmd
}

func runDiff(flags *DiffFlags, stdout io.Writer) (bool, *ReturnError) {

	openapi3.CircularReferenceCounter = flags.circularReferenceCounter

	config := flags.toConfig()

	var diffReport *diff.Diff

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true

	if flags.composed {
		var err *ReturnError
		if diffReport, _, err = composedDiff(loader, flags.base, flags.revision, config); err != nil {
			return false, err
		}
	} else {
		var err *ReturnError
		if diffReport, _, err = normalDiff(loader, flags.base, flags.revision, config); err != nil {
			return false, err
		}
	}

	if flags.summary {
		if err := printYAML(stdout, diffReport.GetSummary()); err != nil {
			return false, getErrFailedPrint("summary", err)
		}
		return failEmpty(flags.failOnDiff, diffReport.Empty()), nil
	}

	return failEmpty(flags.failOnDiff, diffReport.Empty()), handleDiff(stdout, diffReport, flags.format)
}

func handleDiff(stdout io.Writer, diffReport *diff.Diff, format string) *ReturnError {
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

func failEmpty(failOnDiff, diffEmpty bool) bool {
	return failOnDiff && !diffEmpty
}
