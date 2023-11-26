package formatters_test

import (
	"testing"

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
	default:
		return ""
	}
}

var htmlFormatter = formatters.HTMLFormatter{
	Localizer: MockLocalizer,
}

func TestHtmlFormatter_RenderDiff(t *testing.T) {
	out, err := htmlFormatter.RenderDiff(nil, formatters.RenderOpts{})
	require.NoError(t, err)
	require.Equal(t, string(out), "<p>No changes</p>\n")
}

func TestHtmlFormatter_RenderChangelog(t *testing.T) {
	testChanges := checker.Changes{
		checker.ApiChange{
			Path:      "/test",
			Operation: "GET",
			Id:        "change_id",
			Text:      "This is a breaking change.",
			Level:     checker.ERR,
		},
	}

	out, err := htmlFormatter.RenderChangelog(testChanges, formatters.RenderOpts{}, nil)
	require.NoError(t, err)
	require.NotEmpty(t, string(out))
}

func TestHtmlFormatter_NotImplemented(t *testing.T) {
	var err error
	_, err = htmlFormatter.RenderBreakingChanges(checker.Changes{}, formatters.RenderOpts{})
	assert.Error(t, err)

	_, err = htmlFormatter.RenderChecks(formatters.Checks{}, formatters.RenderOpts{})
	assert.Error(t, err)

	_, err = htmlFormatter.RenderFlatten(nil, formatters.RenderOpts{})
	assert.Error(t, err)

	_, err = htmlFormatter.RenderSummary(nil, formatters.RenderOpts{})
	assert.Error(t, err)
}
