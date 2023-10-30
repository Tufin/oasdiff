package formatters

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tufin/oasdiff/checker"
)

func TestGitHubActionsFormatter_RenderBreakingChanges_OneFailure(t *testing.T) {
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
	output, err := formatter.RenderBreakingChanges(testChanges, RenderOpts{})
	assert.NoError(t, err)
	expectedOutput := "::error title=change_id::This is a breaking change.\n"
	assert.Equal(t, expectedOutput, string(output))
}

func TestGitHubActionsFormatter_RenderBreakingChanges_MultipleLevels(t *testing.T) {
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
	output, err := formatter.RenderBreakingChanges(testChanges, RenderOpts{})
	assert.NoError(t, err)
	expectedOutput := "::error title=change_id::This is a breaking change.\n::warning title=change_id::This is a warning.\n::notice title=change_id::This is a notice.\n"
	assert.Equal(t, expectedOutput, string(output))
}

func TestGitHubActionsFormatter_RenderBreakingChanges_MultilineText(t *testing.T) {
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
	output, err := formatter.RenderBreakingChanges(testChanges, RenderOpts{})
	assert.NoError(t, err)
	expectedOutput := "::error title=change_id::This is a breaking change.%0AThis is a second line.\n"
	assert.Equal(t, expectedOutput, string(output))
}

func TestGitHubActionsFormatter_RenderBreakingChanges_FileLocation(t *testing.T) {
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
	output, err := formatter.RenderBreakingChanges(testChanges, RenderOpts{})
	assert.NoError(t, err)
	expectedOutput := "::error title=change_id,file=openapi.json,col=6,endColumn=11,line=21,endLine=26::This is a breaking change.\n"
	assert.Equal(t, expectedOutput, string(output))
}

func TestGitHubActionsFormatter_RenderBreakingChanges_JobOutputParameters(t *testing.T) {
	// temp file to mock GITHUB_OUTPUT
	tempFile, err := os.CreateTemp("", "github-output")
	assert.NoError(t, err)
	defer os.Remove(tempFile.Name())
	_ = os.Setenv("GITHUB_OUTPUT", tempFile.Name())

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
			Text:  "This is a second breaking change.",
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
	output, err := formatter.RenderBreakingChanges(testChanges, RenderOpts{})
	assert.NoError(t, err)
	_ = os.Unsetenv("GITHUB_OUTPUT")
	expectedOutput := "::error title=change_id::This is a breaking change.\n::error title=change_id::This is a second breaking change.\n::warning title=change_id::This is a warning.\n::notice title=change_id::This is a notice.\n"
	assert.Equal(t, expectedOutput, string(output))

	// check job output parameters (NOTE: order of parameters is not guaranteed)
	outputFile, err := os.ReadFile(tempFile.Name())
	assert.NoError(t, err)
	assert.Contains(t, string(outputFile), "error_count=2\n")
	assert.Contains(t, string(outputFile), "warning_count=1\n")
	assert.Contains(t, string(outputFile), "info_count=1\n")
}
