package formatters_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/checker"
	"github.com/tufin/oasdiff/formatters"
)

func TestHtmlFormatter_RenderDiff(t *testing.T) {
	formatter := formatters.HTMLFormatter{}

	out, err := formatter.RenderDiff(nil, formatters.RenderOpts{})
	require.NoError(t, err)
	require.Equal(t, string(out), "<p>No changes</p>\n")
}

func TestHtmlFormatter_RenderChangelog(t *testing.T) {
	formatter := formatters.HTMLFormatter{}

	testChanges := checker.Changes{
		checker.ApiChange{
			Path:      "/test",
			Operation: "GET",
			Id:        "change_id",
			Text:      "This is a breaking change.",
			Level:     checker.ERR,
		},
	}

	out, err := formatter.RenderChangelog(testChanges, formatters.RenderOpts{}, nil)
	require.NoError(t, err)
	require.NotEmpty(t, string(out))
}

func TestHtmlFormatter_NotImplemented(t *testing.T) {
	formatter := formatters.HTMLFormatter{}

	var err error
	_, err = formatter.RenderBreakingChanges(checker.Changes{}, formatters.RenderOpts{})
	assert.Error(t, err)

	_, err = formatter.RenderChecks(formatters.Checks{}, formatters.RenderOpts{})
	assert.Error(t, err)

	_, err = formatter.RenderFlatten(nil, formatters.RenderOpts{})
	assert.Error(t, err)

	_, err = formatter.RenderSummary(nil, formatters.RenderOpts{})
	assert.Error(t, err)
}
