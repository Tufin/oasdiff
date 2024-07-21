package formatters_test

import (
	"testing"
	"text/template"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/checker"
	"github.com/tufin/oasdiff/formatters"
)

var markupFormatter = formatters.MarkupFormatter{
	Localizer: MockLocalizer,
}

func TestMarkupLookup(t *testing.T) {
	f, err := formatters.Lookup(string(formatters.FormatMarkup), formatters.DefaultFormatterOpts())
	require.NoError(t, err)
	require.IsType(t, formatters.MarkupFormatter{}, f)
}

func TestMarkupFormatter_RenderDiff(t *testing.T) {
	out, err := markupFormatter.RenderDiff(nil, formatters.NewRenderOpts())
	require.NoError(t, err)
	require.Equal(t, string(out), "No changes\n")
}

func TestMarkupFormatter_RenderChangelog(t *testing.T) {
	testChanges := checker.Changes{
		checker.ApiChange{
			Path:      "/test",
			Operation: "GET",
			Id:        "change_id",
			Level:     checker.ERR,
		},
	}

	out, err := markupFormatter.RenderChangelog(testChanges, formatters.NewRenderOpts(), nil)
	require.NoError(t, err)
	require.NotEmpty(t, string(out))
}

func TestMarkupFormatter_NotImplemented(t *testing.T) {
	var err error

	_, err = markupFormatter.RenderChecks(formatters.Checks{}, formatters.NewRenderOpts())
	assert.Error(t, err)

	_, err = markupFormatter.RenderFlatten(nil, formatters.NewRenderOpts())
	assert.Error(t, err)

	_, err = markupFormatter.RenderSummary(nil, formatters.NewRenderOpts())
	assert.Error(t, err)
}

func TestExecuteMarkupTemplate_Err(t *testing.T) {
	_, err := formatters.ExecuteTextTemplate(&template.Template{}, nil, nil)
	assert.Error(t, err)
}
