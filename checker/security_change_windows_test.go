package checker_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/checker"
)

func TestSecurityChange_PrettyNotPipedWindows(t *testing.T) {
	piped := false
	save := checker.SetPipedOutput(&piped)
	defer checker.SetPipedOutput(save)
	require.Equal(t, "error\t[change_id] \t\n\tin security\n\t\tThis is a breaking change.\n\t\tcomment", securityChange.MultiLineError(MockLocalizer, checker.ColorAuto))
}
