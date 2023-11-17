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
	require.Equal(t, apiChange.GetId(), "id")
	require.Equal(t, apiChange.GetText(), "text")
	require.Equal(t, apiChange.GetText(), "text")
	require.Equal(t, apiChange.GetComment(), "comment")
	require.Equal(t, apiChange.GetLevel(), checker.ERR)
	require.Equal(t, apiChange.GetOperation(), "GET")
	require.Equal(t, apiChange.GetOperationId(), "123")
	require.Equal(t, apiChange.GetPath(), "/test")
	require.Equal(t, apiChange.GetSourceFile(), "sourceFile")
	require.Equal(t, apiChange.GetSourceLine(), 1)
	require.Equal(t, apiChange.GetSourceLineEnd(), 2)
	require.Equal(t, apiChange.GetSourceColumn(), 3)
	require.Equal(t, apiChange.GetSourceColumnEnd(), 4)
	require.Equal(t, apiChange.LocalizedError(checker.NewDefaultLocalizer()), "error at source, in API GET /test text [id]. comment")
	require.Equal(t, apiChange.Error(), "error at source, in API GET /test text [id]. comment")
}

func TestApiChange_MatchIgnore(t *testing.T) {
	require.True(t, apiChange.MatchIgnore("/test", "error at source, in api get /test text [id]. comment"))
}

func TestApiChange_PrettyPiped(t *testing.T) {
	piped := true
	save := checker.SetPipedOutput(&piped)
	defer checker.SetPipedOutput(save)
	require.Equal(t, apiChange.PrettyErrorText(checker.NewDefaultLocalizer()), "error at source, in API GET /test text [id]. comment")
}

func TestApiChange_PrettyNotPiped(t *testing.T) {
	piped := false
	save := checker.SetPipedOutput(&piped)
	defer checker.SetPipedOutput(save)
	require.Equal(t, apiChange.PrettyErrorText(checker.NewDefaultLocalizer()), "\x1b[31merror\x1b[0m\t[\x1b[33mid\x1b[0m] at source\t\n\tin API \x1b[32mGET\x1b[0m \x1b[32m/test\x1b[0m\n\t\ttext\n\t\tcomment")
}
