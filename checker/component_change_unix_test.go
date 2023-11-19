package checker_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/checker"
)

func TestComponentChange_PrettyNotPipedUnix(t *testing.T) {
	piped := false
	save := checker.SetPipedOutput(&piped)
	defer checker.SetPipedOutput(save)
	require.Equal(t, "\x1b[31merror\x1b[0m\t[\x1b[33mid\x1b[0m] \t\n\tin components/component\n\t\ttext\n\t\tcomment", componentChange.PrettyErrorText(checker.NewDefaultLocalizer()))
}
