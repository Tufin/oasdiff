package formatters_test

import (
	"fmt"
	"html/template"
	"testing"

	_ "embed"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/checker"
	"github.com/tufin/oasdiff/formatters"
)

func MockLocalizer(id string, args ...interface{}) string {
	switch id {
	case "change_id":
		return "This is a breaking change."
	case "warning_id":
		return "This is a warning."
	case "notice_id":
		return "This is a notice."
	case "change_two_lines_id":
		return "This is a breaking change.\nThis is a second line."
	case "total-errors":
		return fmt.Sprintf("%d breaking changes: %d %s, %d %s\n", args...)
	case "total-changes":
		return fmt.Sprintf("%d changes: %d %s, %d %s, %d %s\n", args...)
	default:
		return id
	}
}

var htmlFormatter = formatters.HTMLFormatter{
	Localizer: MockLocalizer,
}

func TestHtmlLookup(t *testing.T) {
	f, err := formatters.Lookup(string(formatters.FormatHTML), formatters.DefaultFormatterOpts())
	require.NoError(t, err)
	require.IsType(t, formatters.HTMLFormatter{}, f)
}

func TestHtmlFormatter_RenderDiff(t *testing.T) {
	out, err := htmlFormatter.RenderDiff(nil, formatters.NewRenderOpts())
	require.NoError(t, err)
	require.Equal(t, string(out), "<p>No changes</p>\n")
}

func TestHtmlFormatter_RenderChangelog(t *testing.T) {
	testChanges := checker.Changes{
		checker.ApiChange{
			Path:      "/test",
			Operation: "GET",
			Id:        "change_id",
			Level:     checker.ERR,
		},
	}

	out, err := htmlFormatter.RenderChangelog(testChanges, formatters.NewRenderOpts(), nil)
	require.NoError(t, err)
	require.NotEmpty(t, string(out))
}

func TestHtmlFormatter_NotImplemented(t *testing.T) {
	var err error

	_, err = htmlFormatter.RenderChecks(formatters.Checks{}, formatters.NewRenderOpts())
	assert.Error(t, err)

	_, err = htmlFormatter.RenderFlatten(nil, formatters.NewRenderOpts())
	assert.Error(t, err)

	_, err = htmlFormatter.RenderSummary(nil, formatters.NewRenderOpts())
	assert.Error(t, err)
}

//go:embed templates/changelog.html
var changelogHtml string

func TestExecuteHtmlTemplate_Err(t *testing.T) {
	tmpl := template.Must(template.New("changelog").Parse(changelogHtml))
	tmpl.Tree = nil
	_, err := formatters.ExecuteHtmlTemplate(tmpl, nil, nil)
	assert.Error(t, err)
}
