package checker_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/checker"
)

var componentChange = checker.ComponentChange{
	Id:              "id",
	Text:            "text",
	Comment:         "comment",
	Level:           checker.ERR,
	Source:          "source",
	Component:       "component",
	SourceFile:      "sourceFile",
	SourceLine:      1,
	SourceLineEnd:   2,
	SourceColumn:    3,
	SourceColumnEnd: 4,
}

func TestComponentChange(t *testing.T) {
	require.Equal(t, componentChange.GetComment(), "comment")
	require.Equal(t, componentChange.GetOperationId(), "")
	require.Equal(t, componentChange.GetSourceFile(), "sourceFile")
	require.Equal(t, componentChange.GetSourceLine(), 1)
	require.Equal(t, componentChange.GetSourceLineEnd(), 2)
	require.Equal(t, componentChange.GetSourceColumn(), 3)
	require.Equal(t, componentChange.GetSourceColumnEnd(), 4)
	require.Equal(t, componentChange.LocalizedError(checker.NewDefaultLocalizer()), "error, in components/component text [id]. comment")
	require.Equal(t, componentChange.PrettyErrorText(checker.NewDefaultLocalizer()), "error, in components/component text [id]. comment")
	require.Equal(t, componentChange.Error(), "error, in components/component text [id]. comment")
}

func TestComponentChange_MatchIgnore(t *testing.T) {
	require.True(t, componentChange.MatchIgnore("", "error, in components/component text [id]. comment"))
}

func TestComponentChange_PrettyPiped(t *testing.T) {
	piped := true
	save := checker.SetPipedOutput(&piped)
	defer checker.SetPipedOutput(save)
	require.Equal(t, componentChange.PrettyErrorText(checker.NewDefaultLocalizer()), "error, in components/component text [id]. comment")
}

func TestComponentChange_PrettyNotPiped(t *testing.T) {
	piped := false
	save := checker.SetPipedOutput(&piped)
	defer checker.SetPipedOutput(save)
	require.Equal(t, componentChange.PrettyErrorText(checker.NewDefaultLocalizer()), "\x1b[31merror\x1b[0m\t[\x1b[33mid\x1b[0m] \t\n\tin components/component\n\t\ttext\n\t\tcomment")
}
