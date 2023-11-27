package formatters_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/checker"
	"github.com/tufin/oasdiff/formatters"
)

var yamlFormatter = formatters.YAMLFormatter{
	Localizer: MockLocalizer,
}

func TestYamlFormatter_RenderBreakingChanges(t *testing.T) {
	testChanges := checker.Changes{
		checker.ComponentChange{
			Id:    "change_id",
			Text:  "This is a breaking change.",
			Level: checker.ERR,
		},
	}

	out, err := yamlFormatter.RenderBreakingChanges(testChanges, formatters.RenderOpts{})
	require.NoError(t, err)
	require.Equal(t, string(out), "- id: change_id\n  text: This is a breaking change.\n  level: 3\n")
}

func TestYamlFormatter_RenderChangelog(t *testing.T) {
	testChanges := checker.Changes{
		checker.ComponentChange{
			Id:    "change_id",
			Text:  "This is a breaking change.",
			Level: checker.ERR,
		},
	}

	out, err := yamlFormatter.RenderChangelog(testChanges, formatters.RenderOpts{}, nil)
	require.NoError(t, err)
	require.Equal(t, string(out), "- id: change_id\n  text: This is a breaking change.\n  level: 3\n")
}

func TestYamlFormatter_RenderChecks(t *testing.T) {
	checks := formatters.Checks{
		{
			Id:          "change_id",
			Level:       "info",
			Description: "This is a breaking change.",
			Required:    true,
		},
	}

	out, err := yamlFormatter.RenderChecks(checks, formatters.RenderOpts{})
	require.NoError(t, err)
	require.Equal(t, string(out), "- id: change_id\n  level: info\n  description: This is a breaking change.\n  reuired: true\n")
}

func TestYamlFormatter_RenderDiff(t *testing.T) {
	out, err := yamlFormatter.RenderDiff(nil, formatters.RenderOpts{})
	require.NoError(t, err)
	require.Empty(t, string(out))
}

func TestYamlFormatter_RenderFlatten(t *testing.T) {
	out, err := yamlFormatter.RenderFlatten(nil, formatters.RenderOpts{})
	require.NoError(t, err)
	require.Empty(t, string(out))
}

func TestYamlFormatter_RenderSummary(t *testing.T) {
	out, err := yamlFormatter.RenderSummary(nil, formatters.RenderOpts{})
	require.NoError(t, err)
	require.Equal(t, string(out), "diff: false\n")
}
