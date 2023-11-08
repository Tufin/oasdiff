package formatters

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/tufin/oasdiff/checker"
	"github.com/tufin/oasdiff/diff"
)

type JSONFormatter struct {
}

func (f JSONFormatter) RenderDiff(diff *diff.Diff, opts RenderOpts) ([]byte, error) {
	return printJSON(diff)
}

func (f JSONFormatter) RenderSummary(diff *diff.Diff, opts RenderOpts) ([]byte, error) {
	return printJSON(diff.GetSummary())
}

func (f JSONFormatter) RenderBreakingChanges(changes checker.Changes, opts RenderOpts) ([]byte, error) {
	return printJSON(changes)
}

func (f JSONFormatter) RenderChangelog(changes checker.Changes, opts RenderOpts) ([]byte, error) {
	return printJSON(changes)
}

func (f JSONFormatter) RenderChecks(checks Checks, opts RenderOpts) ([]byte, error) {
	return printJSON(checks)
}

func (f JSONFormatter) RenderFlatten(spec *openapi3.T, opts RenderOpts) ([]byte, error) {
	return printJSON(spec)
}

func (f JSONFormatter) SupportedOutputs() []Output {
	return []Output{OutputDiff, OutputSummary, OutputBreaking, OutputChangelog, OutputChecks, OutputFlatten}
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
