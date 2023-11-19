package checker_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/checker"
)

var apiChange = checker.ApiChange{
	Id:              "id",
	Text:            "text",
	Comment:         "comment",
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
	require.Equal(t, "id", apiChange.GetId())
	require.Equal(t, "text", apiChange.GetText())
	require.Equal(t, "comment", apiChange.GetComment())
	require.Equal(t, checker.ERR, apiChange.GetLevel())
	require.Equal(t, "GET", apiChange.GetOperation())
	require.Equal(t, "123", apiChange.GetOperationId())
	require.Equal(t, "/test", apiChange.GetPath())
	require.Equal(t, "sourceFile", apiChange.GetSourceFile())
	require.Equal(t, 1, apiChange.GetSourceLine())
	require.Equal(t, 2, apiChange.GetSourceLineEnd())
	require.Equal(t, 3, apiChange.GetSourceColumn())
	require.Equal(t, 4, apiChange.GetSourceColumnEnd())
	require.Equal(t, "error at source, in API GET /test text [id]. comment", apiChange.LocalizedError(checker.NewDefaultLocalizer()))
	require.Equal(t, "error at source, in API GET /test text [id]. comment", apiChange.Error())
}

func TestApiChange_MatchIgnore(t *testing.T) {
	require.True(t, apiChange.MatchIgnore("/test", "error at source, in api get /test text [id]. comment"))
}

func TestApiChange_PrettyPiped(t *testing.T) {
	piped := true
	save := checker.SetPipedOutput(&piped)
	defer checker.SetPipedOutput(save)
	require.Equal(t, "error at source, in API GET /test text [id]. comment", apiChange.PrettyErrorText(checker.NewDefaultLocalizer()))
}
