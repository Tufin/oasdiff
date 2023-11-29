package formatters

import (
	"bytes"
	"fmt"
	"html/template"

	_ "embed"

	"github.com/tufin/oasdiff/checker"
	"github.com/tufin/oasdiff/diff"
	"github.com/tufin/oasdiff/load"
	"github.com/tufin/oasdiff/report"
)

type HTMLFormatter struct {
	notImplementedFormatter
	Localizer checker.Localizer
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

type TemplateData struct {
	APIChanges      ChangesByEndpoint
	BaseVersion     string
	RevisionVersion string
}

func (f HTMLFormatter) RenderChangelog(changes checker.Changes, opts RenderOpts, specInfoPair *load.SpecInfoPair) ([]byte, error) {
	tmpl := template.Must(template.New("changelog").Parse(changelog))

	var out bytes.Buffer
	if err := tmpl.Execute(&out, TemplateData{GroupChanges(changes, f.Localizer), specInfoPair.GetBaseVersion(), specInfoPair.GetRevisionVersion()}); err != nil {
		return nil, err
	}

	return out.Bytes(), nil
}

func (f HTMLFormatter) SupportedOutputs() []Output {
	return []Output{OutputDiff, OutputChangelog}
}
