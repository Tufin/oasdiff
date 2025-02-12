package formatters_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/formatters"
)

func TestUnsupportedLookup(t *testing.T) {
	_, err := formatters.Lookup(string("invalid"), formatters.DefaultFormatterOpts())
	require.Error(t, err)
}

func TestDiffOutputFormats(t *testing.T) {
	supportedFormats := formatters.SupportedFormatsByContentType(formatters.OutputDiff)
	assert.Len(t, supportedFormats, 6)
	assert.Contains(t, supportedFormats, string(formatters.FormatYAML))
	assert.Contains(t, supportedFormats, string(formatters.FormatJSON))
	assert.Contains(t, supportedFormats, string(formatters.FormatText))
	assert.Contains(t, supportedFormats, string(formatters.FormatMarkup))
	assert.Contains(t, supportedFormats, string(formatters.FormatMarkdown))
	assert.Contains(t, supportedFormats, string(formatters.FormatHTML))
}

func TestSummaryOutputFormats(t *testing.T) {
	supportedFormats := formatters.SupportedFormatsByContentType(formatters.OutputSummary)
	assert.Len(t, supportedFormats, 2)
	assert.Contains(t, supportedFormats, string(formatters.FormatYAML))
	assert.Contains(t, supportedFormats, string(formatters.FormatJSON))
}

func TestChangelogOutputFormats(t *testing.T) {
	supportedFormats := formatters.SupportedFormatsByContentType(formatters.OutputChangelog)
	assert.Len(t, supportedFormats, 9)
	assert.Contains(t, supportedFormats, string(formatters.FormatYAML))
	assert.Contains(t, supportedFormats, string(formatters.FormatJSON))
	assert.Contains(t, supportedFormats, string(formatters.FormatText))
	assert.Contains(t, supportedFormats, string(formatters.FormatMarkup))
	assert.Contains(t, supportedFormats, string(formatters.FormatMarkdown))
	assert.Contains(t, supportedFormats, string(formatters.FormatSingleLine))
	assert.Contains(t, supportedFormats, string(formatters.FormatHTML))
	assert.Contains(t, supportedFormats, string(formatters.FormatGithubActions))
	assert.Contains(t, supportedFormats, string(formatters.FormatJUnit))
}
