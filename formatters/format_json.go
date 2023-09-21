package formatters

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/tufin/oasdiff/checker"
	"github.com/tufin/oasdiff/diff"
)

type JSONFormatter struct {
	Localizer checker.Localizer
}

func (f JSONFormatter) RenderDiff(diff *diff.Diff, changes checker.Changes, opts RenderOpts) ([]byte, error) {
	return printJSON(diff)
}

func (f JSONFormatter) RenderSummary(checks []checker.BackwardCompatibilityCheck, diff *diff.Diff, changes checker.Changes, opts RenderOpts) ([]byte, error) {
	return printJSON(diff.GetSummary())
}

func (f JSONFormatter) RenderBreakingChanges(checks []checker.BackwardCompatibilityCheck, diff *diff.Diff, changes checker.Changes, opts RenderOpts) ([]byte, error) {
	return printJSON(changes)
}

func (f JSONFormatter) RenderChangelog(changes checker.Changes, opts RenderOpts) ([]byte, error) {
	return printJSON(changes)
}

func (f JSONFormatter) RenderChecks(rules []checker.BackwardCompatibilityRule, opts RenderOpts) ([]byte, error) {
	return printJSON(rules)
}

func (f JSONFormatter) SupportedOutputs() []string {
	return []string{"diff", "summary", "breaking-changes", "changelog", "checks"}
}

func printJSON(output interface{}) ([]byte, error) {
	if reflect.ValueOf(output).IsNil() {
		return nil, nil
	}

	bytes, err := json.Marshal(output)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSON: %w", err)
	}

	return StripANSIEscapeCodes(bytes), nil
}
