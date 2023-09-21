package formatters

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDiffOutputFormats(t *testing.T) {
	supportedFormats := SupportedFormatsByContentType("diff")
	assert.Len(t, supportedFormats, 4)
	assert.Contains(t, supportedFormats, string(FormatYAML))
	assert.Contains(t, supportedFormats, string(FormatJSON))
	assert.Contains(t, supportedFormats, string(FormatText))
	assert.Contains(t, supportedFormats, string(FormatHTML))
}

func TestSummaryOutputFormats(t *testing.T) {
	supportedFormats := SupportedFormatsByContentType("summary")
	assert.Len(t, supportedFormats, 2)
	assert.Contains(t, supportedFormats, string(FormatYAML))
	assert.Contains(t, supportedFormats, string(FormatJSON))
}

func TestChangelogOutputFormats(t *testing.T) {
	supportedFormats := SupportedFormatsByContentType("changelog")
	assert.Len(t, supportedFormats, 3)
	assert.Contains(t, supportedFormats, string(FormatYAML))
	assert.Contains(t, supportedFormats, string(FormatJSON))
	assert.Contains(t, supportedFormats, string(FormatText))
}

func TestBreakingChangesOutputFormats(t *testing.T) {
	supportedFormats := SupportedFormatsByContentType("breaking-changes")
	assert.Len(t, supportedFormats, 4)
	assert.Contains(t, supportedFormats, string(FormatYAML))
	assert.Contains(t, supportedFormats, string(FormatJSON))
	assert.Contains(t, supportedFormats, string(FormatText))
	assert.Contains(t, supportedFormats, string(FormatGithubActions))
}
