package formatters_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/checker"
	"github.com/tufin/oasdiff/formatters"
)

func TestJsonFormatter_RenderBreakingChanges(t *testing.T) {
	formatter := formatters.JSONFormatter{}

	testChanges := checker.Changes{
		checker.ComponentChange{
			Id:    "change_id",
			Text:  "This is a breaking change.",
			Level: checker.ERR,
		},
	}

	out, err := formatter.RenderBreakingChanges(testChanges, formatters.RenderOpts{})
	require.NoError(t, err)
	require.Equal(t, string(out), "[{\"id\":\"change_id\",\"text\":\"This is a breaking change.\",\"level\":3}]")
}

func TestJsonFormatter_RenderChangelog(t *testing.T) {
	formatter := formatters.JSONFormatter{}

	testChanges := checker.Changes{
		checker.ComponentChange{
			Id:    "change_id",
			Text:  "This is a breaking change.",
			Level: checker.ERR,
		},
	}

	out, err := formatter.RenderChangelog(testChanges, formatters.RenderOpts{})
	require.NoError(t, err)
	require.Equal(t, string(out), "[{\"id\":\"change_id\",\"text\":\"This is a breaking change.\",\"level\":3}]")
}

func TestJsonFormatter_RenderChecks(t *testing.T) {
	formatter := formatters.JSONFormatter{}

	checks := []formatters.Check{
		{
			Id:          "change_id",
			Level:       "info",
			Description: "This is a breaking change.",
			Required:    true,
		},
	}

	out, err := formatter.RenderChecks(checks, formatters.RenderOpts{})
	require.NoError(t, err)
	require.Equal(t, string(out), "[{\"id\":\"change_id\",\"level\":\"info\",\"description\":\"This is a breaking change.\",\"reuired\":true}]")
}

func TestJsonFormatter_RenderDiff(t *testing.T) {
	formatter := formatters.JSONFormatter{}

	out, err := formatter.RenderDiff(nil, formatters.RenderOpts{})
	require.NoError(t, err)
	require.Empty(t, string(out))
}

func TestJsonFormatter_RenderFlatten(t *testing.T) {
	formatter := formatters.JSONFormatter{}

	out, err := formatter.RenderFlatten(nil, formatters.RenderOpts{})
	require.NoError(t, err)
	require.Empty(t, string(out))
}

func TestJsonFormatter_RenderSummary(t *testing.T) {
	formatter := formatters.JSONFormatter{}

	out, err := formatter.RenderSummary(nil, formatters.RenderOpts{})
	require.NoError(t, err)
	require.Equal(t, string(out), `{"diff":false}`)
}
