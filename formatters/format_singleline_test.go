package formatters_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/checker"
	"github.com/tufin/oasdiff/formatters"
)

var singleLineFormatter = formatters.SingleLineFormatter{
	Localizer: MockLocalizer,
}

func TestSingleLineFormatter_RenderBreakingChanges(t *testing.T) {
	testChanges := checker.Changes{
		checker.ComponentChange{
			Id:        "change_id",
			Level:     checker.ERR,
			Component: "test",
		},
	}

	out, err := singleLineFormatter.RenderBreakingChanges(testChanges, formatters.NewRenderOpts())
	require.NoError(t, err)
	require.Equal(t, "1 breaking changes: 1 error, 0 warning\nerror, in components/test This is a breaking change. [change_id]. \n\n", string(out))
}

func TestSingleLineFormatter_RenderChangelog(t *testing.T) {
	testChanges := checker.Changes{
		checker.ComponentChange{
			Id:        "change_id",
			Level:     checker.ERR,
			Component: "test",
		},
	}

	out, err := singleLineFormatter.RenderChangelog(testChanges, formatters.NewRenderOpts(), nil)
	require.NoError(t, err)
	require.Equal(t, "1 changes: 1 error, 0 warning, 0 info\nerror, in components/test This is a breaking change. [change_id]. \n\n", string(out))
}

func TestSingleLineFormatter_RenderChecks(t *testing.T) {
	checks := formatters.Checks{
		{
			Id:          "change_id",
			Level:       "info",
			Description: "This is a breaking change.",
			Required:    true,
		},
	}

	out, err := singleLineFormatter.RenderChecks(checks, formatters.NewRenderOpts())
	require.NoError(t, err)
	require.Equal(t, string(out), "ID        DESCRIPTION                LEVEL\nchange_id This is a breaking change. info\n")
}

func TestSingleLineFormatter_RenderDiff(t *testing.T) {
	out, err := singleLineFormatter.RenderDiff(nil, formatters.NewRenderOpts())
	require.NoError(t, err)
	require.Equal(t, string(out), "No changes\n")
}

func TestSingleLineFormatter_NotImplemented(t *testing.T) {
	var err error

	_, err = singleLineFormatter.RenderFlatten(nil, formatters.NewRenderOpts())
	assert.Error(t, err)

	_, err = singleLineFormatter.RenderSummary(nil, formatters.NewRenderOpts())
	assert.Error(t, err)
}
