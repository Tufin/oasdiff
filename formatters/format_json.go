package formatters

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/tufin/oasdiff/checker"
	"github.com/tufin/oasdiff/diff"
	"github.com/tufin/oasdiff/load"
)

type JSONFormatter struct {
	Localizer checker.Localizer
}

func newJSONFormatter(l checker.Localizer) JSONFormatter {
	return JSONFormatter{
		Localizer: l,
	}
}

func (f JSONFormatter) RenderDiff(diff *diff.Diff, opts RenderOpts) ([]byte, error) {
	return printJSON(diff)
}

func (f JSONFormatter) RenderSummary(diff *diff.Diff, opts RenderOpts) ([]byte, error) {
	return printJSON(diff.GetSummary())
}

func (f JSONFormatter) RenderChangelog(changes checker.Changes, opts RenderOpts, specInfoPair *load.SpecInfoPair) ([]byte, error) {
	return printJSON(NewChanges(changes, f.Localizer))
}

func (f JSONFormatter) RenderChecks(checks Checks, opts RenderOpts) ([]byte, error) {
	return printJSON(checks)
}

func (f JSONFormatter) RenderFlatten(spec *openapi3.T, opts RenderOpts) ([]byte, error) {
	return printJSON(spec)
}

func (f JSONFormatter) SupportedOutputs() []Output {
	return []Output{OutputDiff, OutputSummary, OutputChangelog, OutputChecks, OutputFlatten}
}

func printJSON(output interface{}) ([]byte, error) {
	if reflect.ValueOf(output).IsNil() {
		return nil, nil
	}

	bytes, err := json.Marshal(output)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSON: %w", err)
	}

	return bytes, nil
}
