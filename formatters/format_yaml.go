package formatters

import (
	"fmt"
	"reflect"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/tufin/oasdiff/checker"
	"github.com/tufin/oasdiff/diff"
	"gopkg.in/yaml.v3"
)

type YAMLFormatter struct {
}

func (f YAMLFormatter) RenderDiff(diff *diff.Diff, opts RenderOpts) ([]byte, error) {
	return printYAML(diff)
}

func (f YAMLFormatter) RenderSummary(diff *diff.Diff, opts RenderOpts) ([]byte, error) {
	return printYAML(diff.GetSummary())
}

func (f YAMLFormatter) RenderBreakingChanges(changes checker.Changes, opts RenderOpts) ([]byte, error) {
	return printYAML(changes)
}

func (f YAMLFormatter) RenderChangelog(changes checker.Changes, opts RenderOpts) ([]byte, error) {
	return printYAML(changes)
}

func (f YAMLFormatter) RenderChecks(checks []Check, opts RenderOpts) ([]byte, error) {
	return printYAML(checks)
}

func (f YAMLFormatter) RenderFlatten(spec *openapi3.T, opts RenderOpts) ([]byte, error) {
	return printYAML(spec)
}

func (f YAMLFormatter) SupportedOutputs() []Output {
	return []Output{OutputDiff, OutputSummary, OutputBreaking, OutputChangelog, OutputChecks, OutputFlatten}
}

func printYAML(output interface{}) ([]byte, error) {
	if reflect.ValueOf(output).IsNil() {
		return nil, nil
	}

	bytes, err := yaml.Marshal(output)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal YAML: %w", err)
	}

	return StripANSIEscapeCodes(bytes), nil
}
