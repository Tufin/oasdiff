package formatters_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/checker"
	"github.com/tufin/oasdiff/formatters"
)

var jsonFormatter = formatters.JSONFormatter{
	Localizer: MockLocalizer,
}

func TestJsonFormatter_RenderBreakingChanges(t *testing.T) {
	testChanges := checker.Changes{
		checker.ComponentChange{
			Id:    "change_id",
			Level: checker.ERR,
		},
	}

	out, err := jsonFormatter.RenderBreakingChanges(testChanges, formatters.RenderOpts{})
	require.NoError(t, err)
	require.Equal(t, string(out), "[{\"id\":\"change_id\",\"text\":\"This is a breaking change.\",\"level\":3}]")
}

func TestJsonFormatter_RenderChangelog(t *testing.T) {
	testChanges := checker.Changes{
		checker.ComponentChange{
			Id:    "change_id",
			Level: checker.ERR,
		},
	}

	out, err := jsonFormatter.RenderChangelog(testChanges, formatters.RenderOpts{}, nil)
	require.NoError(t, err)
	require.Equal(t, string(out), "[{\"id\":\"change_id\",\"text\":\"This is a breaking change.\",\"level\":3}]")
}

func TestJsonFormatter_RenderChecks(t *testing.T) {
	checks := formatters.Checks{
		{
			Id:          "change_id",
			Level:       "info",
			Description: "This is a breaking change.",
			Required:    true,
		},
	}

	out, err := jsonFormatter.RenderChecks(checks, formatters.RenderOpts{})
	require.NoError(t, err)
	require.Equal(t, string(out), "[{\"id\":\"change_id\",\"level\":\"info\",\"description\":\"This is a breaking change.\",\"reuired\":true}]")
}

func TestJsonFormatter_RenderDiff(t *testing.T) {
	out, err := jsonFormatter.RenderDiff(nil, formatters.RenderOpts{})
	require.NoError(t, err)
	require.Empty(t, string(out))
}

func TestJsonFormatter_RenderFlatten(t *testing.T) {
	out, err := jsonFormatter.RenderFlatten(nil, formatters.RenderOpts{})
	require.NoError(t, err)
	require.Empty(t, string(out))
}

func TestJsonFormatter_RenderSummary(t *testing.T) {
	out, err := jsonFormatter.RenderSummary(nil, formatters.RenderOpts{})
	require.NoError(t, err)
	require.Equal(t, string(out), `{"diff":false}`)
}
