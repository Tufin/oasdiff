package formatters

import (
	"bytes"
	"fmt"
	"html/template"

	_ "embed"

	"github.com/tufin/oasdiff/checker"
	"github.com/tufin/oasdiff/diff"
	"github.com/tufin/oasdiff/report"
)

type HTMLFormatter struct {
	notImplementedFormatter
}

func (f HTMLFormatter) RenderDiff(diff *diff.Diff, opts RenderOpts) ([]byte, error) {
	reportAsString, err := report.GetHTMLReportAsString(diff)
	if err != nil {
		return nil, fmt.Errorf("failed to generate HTML report: %w", err)
	}

	return []byte(reportAsString), nil
}

//go:embed templates/changelog.html
var changelog string

func (f HTMLFormatter) RenderChangelog(changes checker.Changes, opts RenderOpts) ([]byte, error) {
	tmpl, err := template.New("changelog").Parse(changelog)
	if err != nil {
		return nil, err
	}

	var out bytes.Buffer
	err = tmpl.Execute(&out, changes.Group())
	if err != nil {
		return nil, err
	}

	return out.Bytes(), nil
}

func (f HTMLFormatter) SupportedOutputs() []Output {
	return []Output{OutputDiff, OutputChangelog}
}
