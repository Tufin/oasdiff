package checker_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/checker"
)

var componentChange = checker.ComponentChange{
	Id:              "change_id",
	Comment:         "comment",
	Level:           checker.ERR,
	Component:       "component",
	SourceFile:      "sourceFile",
	SourceLine:      1,
	SourceLineEnd:   2,
	SourceColumn:    3,
	SourceColumnEnd: 4,
}

func TestComponentChange(t *testing.T) {
	require.Equal(t, "components", componentChange.GetSection())
	require.Equal(t, "comment", componentChange.GetComment(MockLocalizer))
	require.Equal(t, "", componentChange.GetOperationId())
	require.Equal(t, "", componentChange.GetSource())
	require.Equal(t, "sourceFile", componentChange.GetSourceFile())
	require.Equal(t, 1, componentChange.GetSourceLine())
	require.Equal(t, 2, componentChange.GetSourceLineEnd())
	require.Equal(t, 3, componentChange.GetSourceColumn())
	require.Equal(t, 4, componentChange.GetSourceColumnEnd())
	require.Equal(t, "error, in components/component This is a breaking change. [change_id]. comment", componentChange.SingleLineError(MockLocalizer, checker.ColorNever))
	require.Equal(t, "error, in components/component This is a breaking change. [change_id]. comment", componentChange.SingleLineError(MockLocalizer, checker.ColorNever))
}

func TestComponentChange_MatchIgnore(t *testing.T) {
	require.True(t, componentChange.MatchIgnore("", "error, in components/component this is a breaking change. [change_id]. comment", MockLocalizer))
}

func TestComponentChange_SingleLineError(t *testing.T) {
	require.Equal(t, "error, in components/component This is a breaking change. [change_id]. comment", componentChange.SingleLineError(MockLocalizer, checker.ColorNever))
}

func TestComponentChange_MultiLineError_NoColor(t *testing.T) {
	require.Equal(t, "error, in components/component This is a breaking change. [change_id]. comment", componentChange.SingleLineError(MockLocalizer, checker.ColorNever))
}
