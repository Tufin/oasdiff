package formatters_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/checker"
	"github.com/tufin/oasdiff/checker/localizations"
	"github.com/tufin/oasdiff/formatters"
)

func TestTextFormatter_RenderBreakingChanges(t *testing.T) {
	formatter, err := formatters.Lookup("text", formatters.FormatterOpts{
		Language: localizations.LangDefault,
	})
	require.NoError(t, err)

	testChanges := checker.Changes{
		checker.ComponentChange{
			Id:    "change_id",
			Text:  "This is a breaking change.",
			Level: checker.ERR,
		},
	}

	out, err := formatter.RenderBreakingChanges(testChanges, formatters.RenderOpts{})
	require.NoError(t, err)
	require.Equal(t, string(out), "1 breaking changes: 1 error, 0 warning\nerror, in components This is a breaking change. [change_id]. \n\n")
}

func TestTextFormatter_RenderChangelog(t *testing.T) {
	formatter, err := formatters.Lookup("text", formatters.FormatterOpts{
		Language: localizations.LangDefault,
	})
	require.NoError(t, err)

	testChanges := checker.Changes{
		checker.ComponentChange{
			Id:    "change_id",
			Text:  "This is a breaking change.",
			Level: checker.ERR,
		},
	}

	out, err := formatter.RenderChangelog(testChanges, formatters.RenderOpts{})
	require.NoError(t, err)
	require.Equal(t, string(out), "1 changes: 1 error, 0 warning, 0 info\nerror, in components This is a breaking change. [change_id]. \n\n")
}

func TestTextFormatter_RenderChecks(t *testing.T) {
	formatter, err := formatters.Lookup("text", formatters.FormatterOpts{
		Language: localizations.LangDefault,
	})
	require.NoError(t, err)

	checks := formatters.Checks{
		{
			Id:          "change_id",
			Level:       "info",
			Description: "This is a breaking change.",
			Required:    true,
		},
	}

	out, err := formatter.RenderChecks(checks, formatters.RenderOpts{})
	require.NoError(t, err)
	require.Equal(t, string(out), "ID        DESCRIPTION                LEVEL\nchange_id This is a breaking change. info\n")
}

func TestTextFormatter_RenderDiff(t *testing.T) {
	formatter, err := formatters.Lookup("text", formatters.FormatterOpts{
		Language: localizations.LangDefault,
	})
	require.NoError(t, err)

	out, err := formatter.RenderDiff(nil, formatters.RenderOpts{})
	require.NoError(t, err)
	require.Equal(t, string(out), "No changes\n")
}

func TestTextFormatter_NotImplemented(t *testing.T) {
	formatter, err := formatters.Lookup("text", formatters.FormatterOpts{
		Language: localizations.LangDefault,
	})
	require.NoError(t, err)

	_, err = formatter.RenderFlatten(nil, formatters.RenderOpts{})
	assert.Error(t, err)

	_, err = formatter.RenderSummary(nil, formatters.RenderOpts{})
	assert.Error(t, err)
}
