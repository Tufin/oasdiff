package formatters

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tufin/oasdiff/checker"
)

func TestRenderBreakingChanges(t *testing.T) {
	// prepare formatter and test changes
	formatter := GitHubActionsFormatter{}
	testChanges := checker.Changes{
		checker.ComponentChange{
			Id:    "change_id",
			Text:  "This is a breaking change.",
			Level: checker.ERR,
		},
	}

	// check output
	output, err := formatter.RenderBreakingChanges(nil, nil, testChanges, RenderOpts{})
	assert.NoError(t, err)
	expectedOutput := "::error title=change_id::This is a breaking change.\n"
	assert.Equal(t, expectedOutput, string(output))
}

func TestRenderBreakingChangesLevels(t *testing.T) {
	// prepare formatter and test changes
	formatter := GitHubActionsFormatter{}
	testChanges := checker.Changes{
		checker.ComponentChange{
			Id:    "change_id",
			Text:  "This is a breaking change.",
			Level: checker.ERR,
		},
		checker.ComponentChange{
			Id:    "change_id",
			Text:  "This is a warning.",
			Level: checker.WARN,
		},
		checker.ComponentChange{
			Id:    "change_id",
			Text:  "This is a notice.",
			Level: checker.INFO,
		},
	}

	// check output
	output, err := formatter.RenderBreakingChanges(nil, nil, testChanges, RenderOpts{})
	assert.NoError(t, err)
	expectedOutput := "::error title=change_id::This is a breaking change.\n::warning title=change_id::This is a warning.\n::notice title=change_id::This is a notice.\n"
	assert.Equal(t, expectedOutput, string(output))
}

func TestRenderBreakingChangesMultilineText(t *testing.T) {
	// prepare formatter and test changes
	formatter := GitHubActionsFormatter{}
	testChanges := checker.Changes{
		checker.ComponentChange{
			Id:    "change_id",
			Text:  "This is a breaking change.\nThis is a second line.",
			Level: checker.ERR,
		},
	}

	// check output
	output, err := formatter.RenderBreakingChanges(nil, nil, testChanges, RenderOpts{})
	assert.NoError(t, err)
	expectedOutput := "::error title=change_id::This is a breaking change.%0AThis is a second line.\n"
	assert.Equal(t, expectedOutput, string(output))
}

func TestRenderBreakingChangesWithFileAndLine(t *testing.T) {
	// prepare formatter and test changes
	formatter := GitHubActionsFormatter{}
	testChanges := checker.Changes{
		checker.ComponentChange{
			Id:              "change_id",
			Text:            "This is a breaking change.",
			Level:           checker.ERR,
			SourceFile:      "openapi.json",
			SourceLine:      20,
			SourceLineEnd:   25,
			SourceColumn:    5,
			SourceColumnEnd: 10,
		},
	}

	// check output
	output, err := formatter.RenderBreakingChanges(nil, nil, testChanges, RenderOpts{})
	assert.NoError(t, err)
	expectedOutput := "::error title=change_id,file=openapi.json,col=6,endColumn=11,line=21,endLine=26::This is a breaking change.\n"
	assert.Equal(t, expectedOutput, string(output))
}
