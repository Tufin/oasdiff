package checker_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/checker"
)

var securityChange = checker.SecurityChange{
	Id:              "id",
	Text:            "text",
	Comment:         "comment",
	Level:           checker.ERR,
	Source:          "source",
	SourceFile:      "sourceFile",
	SourceLine:      1,
	SourceLineEnd:   2,
	SourceColumn:    3,
	SourceColumnEnd: 4,
}

func TestSecurityChange(t *testing.T) {
	require.Equal(t, securityChange.GetComment(), "comment")
	require.Equal(t, securityChange.GetOperationId(), "")
	require.Equal(t, securityChange.GetSourceFile(), "sourceFile")
	require.Equal(t, securityChange.GetSourceLine(), 1)
	require.Equal(t, securityChange.GetSourceLineEnd(), 2)
	require.Equal(t, securityChange.GetSourceColumn(), 3)
	require.Equal(t, securityChange.GetSourceColumnEnd(), 4)
	require.Equal(t, securityChange.LocalizedError(checker.NewDefaultLocalizer()), "error, in security text [id]. comment")
	require.Equal(t, securityChange.Error(), "error, in security text [id]. comment")
}

func TestSecurityChange_MatchIgnore(t *testing.T) {
	require.True(t, securityChange.MatchIgnore("", "error, in security text [id]. comment"))
}

func TestSecurityChange_PrettyPiped(t *testing.T) {
	piped := true
	save := checker.SetPipedOutput(&piped)
	defer checker.SetPipedOutput(save)
	require.Equal(t, securityChange.PrettyErrorText(checker.NewDefaultLocalizer()), "error, in security text [id]. comment")
}

func TestSecurityChange_PrettyNotPiped(t *testing.T) {
	piped := false
	save := checker.SetPipedOutput(&piped)
	defer checker.SetPipedOutput(save)
	require.Equal(t, securityChange.PrettyErrorText(checker.NewDefaultLocalizer()), "\x1b[31merror\x1b[0m\t[\x1b[33mid\x1b[0m] \t\n\tin security\n\t\ttext\n\t\tcomment")
}
