package formatters

import (
	"fmt"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/tufin/oasdiff/checker"
	"github.com/tufin/oasdiff/diff"
	"github.com/tufin/oasdiff/report"
)

type HTMLFormatter struct {
}

func (f HTMLFormatter) RenderDiff(diff *diff.Diff, opts RenderOpts) ([]byte, error) {
	reportAsString, err := report.GetHTMLReportAsString(diff)
	if err != nil {
		return nil, fmt.Errorf("failed to generate HTML report: %w", err)
	}

	return []byte(reportAsString), nil
}

func (f HTMLFormatter) RenderSummary(diff *diff.Diff, opts RenderOpts) ([]byte, error) {
	return nil, fmt.Errorf("not implemented")
}

func (f HTMLFormatter) RenderBreakingChanges(changes checker.Changes, opts RenderOpts) ([]byte, error) {
	return nil, fmt.Errorf("not implemented")
}

func (f HTMLFormatter) RenderChangelog(changes checker.Changes, opts RenderOpts) ([]byte, error) {
	return nil, fmt.Errorf("not implemented")
}

func (f HTMLFormatter) RenderChecks(rules []checker.BackwardCompatibilityRule, opts RenderOpts) ([]byte, error) {
	return nil, fmt.Errorf("not implemented")
}

func (f HTMLFormatter) RenderFlatten(spec *openapi3.T, opts RenderOpts) ([]byte, error) {
	return nil, fmt.Errorf("not implemented")
}

func (f HTMLFormatter) SupportedOutputs() []Output {
	return []Output{OutputDiff}
}
