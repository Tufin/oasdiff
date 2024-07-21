package formatters

import (
	"bytes"
	"text/template"

	_ "embed"

	"github.com/tufin/oasdiff/checker"
	"github.com/tufin/oasdiff/diff"
	"github.com/tufin/oasdiff/load"
	"github.com/tufin/oasdiff/report"
)

type MarkupFormatter struct {
	notImplementedFormatter
	Localizer checker.Localizer
}

func newMarkupFormatter(l checker.Localizer) MarkupFormatter {
	return MarkupFormatter{
		Localizer: l,
	}
}

func (f MarkupFormatter) RenderDiff(diff *diff.Diff, opts RenderOpts) ([]byte, error) {
	return []byte(report.GetTextReportAsString(diff)), nil
}

//go:embed templates/changelog.md
var changelogMarkdown string

func (f MarkupFormatter) RenderChangelog(changes checker.Changes, opts RenderOpts, specInfoPair *load.SpecInfoPair) ([]byte, error) {
	tmpl := template.Must(template.New("changelog").Parse(changelogMarkdown))
	return ExecuteTextTemplate(tmpl, GroupChanges(changes, f.Localizer), specInfoPair)
}

func ExecuteTextTemplate(tmpl *template.Template, changes ChangesByEndpoint, specInfoPair *load.SpecInfoPair) ([]byte, error) {
	var out bytes.Buffer
	if err := tmpl.Execute(&out, TemplateData{changes, specInfoPair.GetBaseVersion(), specInfoPair.GetRevisionVersion()}); err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}

func (f MarkupFormatter) SupportedOutputs() []Output {
	return []Output{OutputDiff, OutputChangelog}
}
