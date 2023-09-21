package formatters

import (
	"fmt"
	"reflect"

	"github.com/tufin/oasdiff/checker"
	"github.com/tufin/oasdiff/diff"
	"gopkg.in/yaml.v3"
)

type YAMLFormatter struct {
	Localizer checker.Localizer
}

func (f YAMLFormatter) RenderDiff(diff *diff.Diff, changes checker.Changes, opts RenderOpts) ([]byte, error) {
	return printYAML(diff)
}

func (f YAMLFormatter) RenderSummary(checks []checker.BackwardCompatibilityCheck, diff *diff.Diff, changes checker.Changes, opts RenderOpts) ([]byte, error) {
	return printYAML(diff.GetSummary())
}

func (f YAMLFormatter) RenderBreakingChanges(checks []checker.BackwardCompatibilityCheck, diff *diff.Diff, changes checker.Changes, opts RenderOpts) ([]byte, error) {
	return printYAML(changes)
}

func (f YAMLFormatter) RenderChangelog(changes checker.Changes, opts RenderOpts) ([]byte, error) {
	return printYAML(changes)
}

func (f YAMLFormatter) RenderChecks(rules []checker.BackwardCompatibilityRule, opts RenderOpts) ([]byte, error) {
	return printYAML(rules)
}

func (f YAMLFormatter) SupportedOutputs() []string {
	return []string{"diff", "summary", "breaking-changes", "changelog", "checks"}
}

func printYAML(output interface{}) ([]byte, error) {
	if reflect.ValueOf(output).IsNil() {
		return nil, nil
	}

	bytes, err := yaml.Marshal(output)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal YAML: %w", err)
	}

	return bytes, nil
}
