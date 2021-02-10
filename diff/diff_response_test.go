package diff_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/diff"
	"github.com/tufin/oasdiff/load"
)

func TestGetDiffResponse(t *testing.T) {
	s1, err := load.LoadPath(test1)
	require.NoError(t, err)

	s2, err := load.LoadPath(test2)
	require.NoError(t, err)

	require.Equal(t,
		&diff.DiffSummary{
			Diff:              true,
			MissingEndpoints:  1,
			ModifiedEndpoints: 1,
		},
		diff.GetDiffResponse(s1, s2, "", "").DiffSummary)
}
