package checker_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/checker"
)

var securityChange = checker.SecurityChange{
	Id:              "change_id",
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
	require.Equal(t, "comment", securityChange.GetComment())
	require.Equal(t, "", securityChange.GetOperationId())
	require.Equal(t, "sourceFile", securityChange.GetSourceFile())
	require.Equal(t, 1, securityChange.GetSourceLine())
	require.Equal(t, 2, securityChange.GetSourceLineEnd())
	require.Equal(t, 3, securityChange.GetSourceColumn())
	require.Equal(t, 4, securityChange.GetSourceColumnEnd())
	require.Equal(t, "error, in security This is a breaking change. [change_id]. comment", securityChange.LocalizedError(MockLocalizer))
}

func TestSecurityChange_MatchIgnore(t *testing.T) {
	require.True(t, securityChange.MatchIgnore("", "error, in security this is a breaking change. [change_id]. comment", MockLocalizer))
}

func TestSecurityChange_PrettyPiped(t *testing.T) {
	piped := true
	save := checker.SetPipedOutput(&piped)
	defer checker.SetPipedOutput(save)
	require.Equal(t, "error, in security This is a breaking change. [change_id]. comment", securityChange.PrettyErrorText(MockLocalizer))
}
