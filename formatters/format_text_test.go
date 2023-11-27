package formatters_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/checker"
	"github.com/tufin/oasdiff/formatters"
)

var textFormatter = formatters.TEXTFormatter{
	Localizer: MockLocalizer,
}

func TestTextFormatter_RenderBreakingChanges(t *testing.T) {
	testChanges := checker.Changes{
		checker.ComponentChange{
			Id:        "change_id",
			Level:     checker.ERR,
			Component: "test",
		},
	}

	out, err := textFormatter.RenderBreakingChanges(testChanges, formatters.RenderOpts{})
	require.NoError(t, err)
	require.Equal(t, "1 breaking changes: 1 error, 0 warning\nerror, in components/test This is a breaking change. [change_id]. \n\n", string(out))
}

func TestTextFormatter_RenderChangelog(t *testing.T) {
	testChanges := checker.Changes{
		checker.ComponentChange{
			Id:        "change_id",
			Level:     checker.ERR,
			Component: "test",
		},
	}

	out, err := textFormatter.RenderChangelog(testChanges, formatters.RenderOpts{}, nil)
	require.NoError(t, err)
	require.Equal(t, "1 changes: 1 error, 0 warning, 0 info\nerror, in components/test This is a breaking change. [change_id]. \n\n", string(out))
}

func TestTextFormatter_RenderChecks(t *testing.T) {
	checks := formatters.Checks{
		{
			Id:          "change_id",
			Level:       "info",
			Description: "This is a breaking change.",
			Required:    true,
		},
	}

	out, err := textFormatter.RenderChecks(checks, formatters.RenderOpts{})
	require.NoError(t, err)
	require.Equal(t, string(out), "ID        DESCRIPTION                LEVEL\nchange_id This is a breaking change. info\n")
}

func TestTextFormatter_RenderDiff(t *testing.T) {
	out, err := textFormatter.RenderDiff(nil, formatters.RenderOpts{})
	require.NoError(t, err)
	require.Equal(t, string(out), "No changes\n")
}

func TestTextFormatter_NotImplemented(t *testing.T) {
	var err error

	_, err = textFormatter.RenderFlatten(nil, formatters.RenderOpts{})
	assert.Error(t, err)

	_, err = textFormatter.RenderSummary(nil, formatters.RenderOpts{})
	assert.Error(t, err)
}
