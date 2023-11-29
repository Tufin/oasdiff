package checker_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/checker"
)

func TestComponentChange_PrettyNotPipedWindows(t *testing.T) {
	piped := false
	save := checker.SetPipedOutput(&piped)
	defer checker.SetPipedOutput(save)
	require.Equal(t, "error\t[change_id] \t\n\tin components/component\n\t\tThis is a breaking change.\n\t\tcomment", componentChange.MultiLineError(MockLocalizer, checker.ColorAuto))
}
