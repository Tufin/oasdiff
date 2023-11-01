package formatters_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/checker"
	"github.com/tufin/oasdiff/formatters"
)

func TestHtmlFormatter_RenderBreakingChanges_Normal(t *testing.T) {
	formatter := formatters.HTMLFormatter{}

	out, err := formatter.RenderDiff(nil, formatters.RenderOpts{})
	require.NoError(t, err)
	require.Equal(t, string(out), "<p>No changes</p>\n")

}

func TestHtmlFormatter_RenderBreakingChanges_NotImplemented(t *testing.T) {
	formatter := formatters.HTMLFormatter{}

	testChanges := checker.Changes{}

	var err error
	_, err = formatter.RenderBreakingChanges(testChanges, formatters.RenderOpts{})
	assert.Error(t, err)

	_, err = formatter.RenderChangelog(testChanges, formatters.RenderOpts{})
	assert.Error(t, err)

	_, err = formatter.RenderChecks([]formatters.Check{}, formatters.RenderOpts{})
	assert.Error(t, err)

	_, err = formatter.RenderFlatten(nil, formatters.RenderOpts{})
	assert.Error(t, err)

	_, err = formatter.RenderSummary(nil, formatters.RenderOpts{})
	assert.Error(t, err)
}
