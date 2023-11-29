package checker_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/checker"
)

var apiChange = checker.ApiChange{
	Id:              "change_id",
	Args:            []any{},
	Comment:         "comment_id",
	Level:           checker.ERR,
	Operation:       "GET",
	OperationId:     "123",
	Path:            "/test",
	Source:          "source",
	SourceFile:      "sourceFile",
	SourceLine:      1,
	SourceLineEnd:   2,
	SourceColumn:    3,
	SourceColumnEnd: 4,
}

func TestApiChange(t *testing.T) {
	require.Equal(t, "change_id", apiChange.GetId())
	require.Equal(t, "comment", apiChange.GetComment(MockLocalizer))
	require.Equal(t, checker.ERR, apiChange.GetLevel())
	require.Equal(t, "GET", apiChange.GetOperation())
	require.Equal(t, "123", apiChange.GetOperationId())
	require.Equal(t, "/test", apiChange.GetPath())
	require.Equal(t, "sourceFile", apiChange.GetSourceFile())
	require.Equal(t, 1, apiChange.GetSourceLine())
	require.Equal(t, 2, apiChange.GetSourceLineEnd())
	require.Equal(t, 3, apiChange.GetSourceColumn())
	require.Equal(t, 4, apiChange.GetSourceColumnEnd())
	require.Equal(t, "error at source, in API GET /test This is a breaking change. [change_id]. comment", apiChange.SingleLineError(MockLocalizer, checker.ColorNever))
}

func MockLocalizer(originalKey string, args ...interface{}) string {
	switch originalKey {
	case "change_id":
		return "This is a breaking change."
	case "comment_id":
		return "comment"
	default:
		return originalKey
	}

}

func TestApiChange_MatchIgnore(t *testing.T) {
	require.True(t, apiChange.MatchIgnore("/test", "error at source, in api get /test this is a breaking change. [change_id]. comment", MockLocalizer))
}

func TestApiChange_PrettyPiped(t *testing.T) {
	piped := true
	save := checker.SetPipedOutput(&piped)
	defer checker.SetPipedOutput(save)
	require.Equal(t, "error at source, in API GET /test This is a breaking change. [change_id]. comment", apiChange.SingleLineError(MockLocalizer, checker.ColorAuto))
}
