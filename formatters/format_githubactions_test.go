package formatters_test

import (
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tufin/oasdiff/checker"
	"github.com/tufin/oasdiff/formatters"
)

func TestGitHubActionsFormatter_RenderBreakingChanges_OneFailure(t *testing.T) {
	// prepare formatter and test changes
	formatter := formatters.GitHubActionsFormatter{}
	testChanges := checker.Changes{
		checker.ApiChange{
			Id:        "change_id",
			Text:      "This is a breaking change.",
			Level:     checker.ERR,
			Operation: http.MethodGet,
			Path:      "/api/test",
			Source:    "openapi.yaml",
		},
	}

	// check output
	output, err := formatter.RenderBreakingChanges(testChanges, formatters.RenderOpts{})
	assert.NoError(t, err)
	expectedOutput := "::error title=change_id::at openapi.yaml, in API GET /api/test This is a breaking change.\n"
	assert.Equal(t, expectedOutput, string(output))
}

func TestGitHubActionsFormatter_RenderBreakingChanges_MultipleLevels(t *testing.T) {
	// prepare formatter and test changes
	formatter := formatters.GitHubActionsFormatter{}
	testChanges := checker.Changes{
		checker.ApiChange{
			Id:        "change_id",
			Text:      "This is a breaking change.",
			Level:     checker.ERR,
			Operation: http.MethodGet,
			Path:      "/api/test",
			Source:    "openapi.yaml",
		},
		checker.ApiChange{
			Id:        "change_id",
			Text:      "This is a warning.",
			Level:     checker.WARN,
			Operation: http.MethodGet,
			Path:      "/api/test",
			Source:    "openapi.yaml",
		},
		checker.ApiChange{
			Id:        "change_id",
			Text:      "This is a notice.",
			Level:     checker.INFO,
			Operation: http.MethodGet,
			Path:      "/api/test",
			Source:    "openapi.yaml",
		},
	}

	// check output
	output, err := formatter.RenderBreakingChanges(testChanges, formatters.RenderOpts{})
	assert.NoError(t, err)
	expectedOutput := "::error title=change_id::at openapi.yaml, in API GET /api/test This is a breaking change.\n::warning title=change_id::at openapi.yaml, in API GET /api/test This is a warning.\n::notice title=change_id::at openapi.yaml, in API GET /api/test This is a notice.\n"
	assert.Equal(t, expectedOutput, string(output))
}

func TestGitHubActionsFormatter_RenderBreakingChanges_MultilineText(t *testing.T) {
	// prepare formatter and test changes
	formatter := formatters.GitHubActionsFormatter{}
	testChanges := checker.Changes{
		checker.ApiChange{
			Id:        "change_id",
			Text:      "This is a breaking change.\nThis is a second line.",
			Level:     checker.ERR,
			Operation: http.MethodGet,
			Path:      "/api/test",
			Source:    "openapi.yaml",
		},
	}

	// check output
	output, err := formatter.RenderBreakingChanges(testChanges, formatters.RenderOpts{})
	assert.NoError(t, err)
	expectedOutput := "::error title=change_id::at openapi.yaml, in API GET /api/test This is a breaking change.%0AThis is a second line.\n"
	assert.Equal(t, expectedOutput, string(output))
}

func TestGitHubActionsFormatter_RenderBreakingChanges_FileLocation(t *testing.T) {
	// prepare formatter and test changes
	formatter := formatters.GitHubActionsFormatter{}
	testChanges := checker.Changes{
		checker.ApiChange{
			Id:              "change_id",
			Text:            "This is a breaking change.",
			Level:           checker.ERR,
			Operation:       http.MethodGet,
			Path:            "/api/test",
			Source:          "openapi.yaml",
			SourceFile:      "openapi.json",
			SourceLine:      20,
			SourceLineEnd:   25,
			SourceColumn:    5,
			SourceColumnEnd: 10,
		},
	}

	// check output
	output, err := formatter.RenderBreakingChanges(testChanges, formatters.RenderOpts{})
	assert.NoError(t, err)
	expectedOutput := "::error title=change_id,file=openapi.json,col=6,endColumn=11,line=21,endLine=26::at openapi.yaml, in API GET /api/test This is a breaking change.\n"
	assert.Equal(t, expectedOutput, string(output))
}

func TestGitHubActionsFormatter_RenderBreakingChanges_JobOutputParameters(t *testing.T) {
	// temp file to mock GITHUB_OUTPUT
	tempFile, err := os.CreateTemp("", "github-output")
	assert.NoError(t, err)
	defer os.Remove(tempFile.Name())
	_ = os.Setenv("GITHUB_OUTPUT", tempFile.Name())

	// prepare formatter and test changes
	formatter := formatters.GitHubActionsFormatter{}
	testChanges := checker.Changes{
		checker.ApiChange{
			Id:        "change_id",
			Text:      "This is a breaking change.",
			Level:     checker.ERR,
			Operation: http.MethodGet,
			Path:      "/api/test",
			Source:    "openapi.yaml",
		},
		checker.ApiChange{
			Id:        "change_id",
			Text:      "This is a second breaking change.",
			Level:     checker.ERR,
			Operation: http.MethodGet,
			Path:      "/api/test",
			Source:    "openapi.yaml",
		},
		checker.ApiChange{
			Id:        "change_id",
			Text:      "This is a warning.",
			Level:     checker.WARN,
			Operation: http.MethodGet,
			Path:      "/api/test",
			Source:    "openapi.yaml",
		},
		checker.ApiChange{
			Id:        "change_id",
			Text:      "This is a notice.",
			Level:     checker.INFO,
			Operation: http.MethodGet,
			Path:      "/api/test",
			Source:    "openapi.yaml",
		},
	}

	// check output
	output, err := formatter.RenderBreakingChanges(testChanges, formatters.RenderOpts{})
	assert.NoError(t, err)
	_ = os.Unsetenv("GITHUB_OUTPUT")
	expectedOutput := "::error title=change_id::at openapi.yaml, in API GET /api/test This is a breaking change.\n::error title=change_id::at openapi.yaml, in API GET /api/test This is a second breaking change.\n::warning title=change_id::at openapi.yaml, in API GET /api/test This is a warning.\n::notice title=change_id::at openapi.yaml, in API GET /api/test This is a notice.\n"
	assert.Equal(t, expectedOutput, string(output))

	// check job output parameters (NOTE: order of parameters is not guaranteed)
	outputFile, err := os.ReadFile(tempFile.Name())
	assert.NoError(t, err)
	assert.Contains(t, string(outputFile), "error_count=2\n")
	assert.Contains(t, string(outputFile), "warning_count=1\n")
	assert.Contains(t, string(outputFile), "info_count=1\n")
}

func TestGitHubActionsFormatterr_NotImplemented(t *testing.T) {
	formatter := formatters.GitHubActionsFormatter{}

	var err error
	_, err = formatter.RenderDiff(nil, formatters.RenderOpts{})
	assert.Error(t, err)

	_, err = formatter.RenderSummary(nil, formatters.RenderOpts{})
	assert.Error(t, err)

	_, err = formatter.RenderChangelog(nil, formatters.RenderOpts{})
	assert.Error(t, err)

	_, err = formatter.RenderChecks(nil, formatters.RenderOpts{})
	assert.Error(t, err)

	_, err = formatter.RenderFlatten(nil, formatters.RenderOpts{})
	assert.Error(t, err)
}
