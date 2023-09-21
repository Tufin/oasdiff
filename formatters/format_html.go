package formatters

import (
	"fmt"

	"github.com/tufin/oasdiff/checker"
	"github.com/tufin/oasdiff/diff"
	"github.com/tufin/oasdiff/report"
)

type HTMLFormatter struct {
	Localizer checker.Localizer
}

func (f HTMLFormatter) RenderDiff(diff *diff.Diff, changes checker.Changes, opts RenderOpts) ([]byte, error) {
	reportAsString, err := report.GetHTMLReportAsString(diff)
	if err != nil {
		return nil, fmt.Errorf("failed to generate HTML report: %w", err)
	}

	return []byte(reportAsString), nil
}

func (f HTMLFormatter) RenderSummary(checks []checker.BackwardCompatibilityCheck, diff *diff.Diff, changes checker.Changes, opts RenderOpts) ([]byte, error) {
	return nil, fmt.Errorf("not implemented")
}

func (f HTMLFormatter) RenderBreakingChanges(checks []checker.BackwardCompatibilityCheck, diff *diff.Diff, changes checker.Changes, opts RenderOpts) ([]byte, error) {
	return nil, fmt.Errorf("not implemented")
}

func (f HTMLFormatter) RenderChangelog(changes checker.Changes, opts RenderOpts) ([]byte, error) {
	return nil, fmt.Errorf("not implemented")
}

func (f HTMLFormatter) RenderChecks(rules []checker.BackwardCompatibilityRule, opts RenderOpts) ([]byte, error) {
	return nil, fmt.Errorf("not implemented")
}

func (f HTMLFormatter) SupportedOutputs() []string {
	return []string{"diff"}
}
