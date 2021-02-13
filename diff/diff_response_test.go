package diff_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/diff"
)

func TestGetDiffResponse_Diff(t *testing.T) {
	require.Equal(t,
		&diff.DiffSummary{
			Diff:              true,
			AddedEndpoints:    0,
			DeletedEndpoints:  1,
			ModifiedEndpoints: 1,
		},
		diff.GetDiffResponse(l(t, 1), l(t, 2), "", "").DiffSummary)
}

func TestGetDiffResponse_NoDiff(t *testing.T) {
	s := l(t, 1)

	require.Equal(t,
		&diff.DiffSummary{
			Diff:              false,
			DeletedEndpoints:  0,
			ModifiedEndpoints: 0,
		},
		diff.GetDiffResponse(s, s, "", "").DiffSummary)
}

func TestGetDiffResponse_Prefix(t *testing.T) {
	require.Equal(t,
		&diff.DiffSummary{
			Diff:              true,
			DeletedEndpoints:  0,
			ModifiedEndpoints: 1,
		},
		diff.GetDiffResponse(l(t, 4), l(t, 2), "/prefix", "").DiffSummary)
}
