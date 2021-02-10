package diff_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/diff"
	"github.com/tufin/oasdiff/load"
)

func TestGetDiffResponse_Diff(t *testing.T) {
	s1, err := load.LoadPath(test1)
	require.NoError(t, err)

	s2, err := load.LoadPath(test2)
	require.NoError(t, err)

	require.Equal(t,
		&diff.DiffSummary{
			Diff:              true,
			AddedEndpoints:    0,
			DeletedEndpoints:  1,
			ModifiedEndpoints: 1,
		},
		diff.GetDiffResponse(s1, s2, "", "").DiffSummary)
}

func TestGetDiffResponse_NoDiff(t *testing.T) {
	s, err := load.LoadPath(test1)
	require.NoError(t, err)

	require.Equal(t,
		&diff.DiffSummary{
			Diff:              false,
			DeletedEndpoints:  0,
			ModifiedEndpoints: 0,
		},
		diff.GetDiffResponse(s, s, "", "").DiffSummary)
}

func TestGetDiffResponse_Prefix(t *testing.T) {
	s2, err := load.LoadPath(test2)
	require.NoError(t, err)

	s4, err := load.LoadPath(test4)
	require.NoError(t, err)

	require.Equal(t,
		&diff.DiffSummary{
			Diff:              false,
			DeletedEndpoints:  0,
			ModifiedEndpoints: 0,
		},
		diff.GetDiffResponse(s4, s2, "/prefix", "").DiffSummary)
}
