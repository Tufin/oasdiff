package diff_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/diff"
)

func TestFilterByRegex(t *testing.T) {
	diffResult := diff.DiffResult{
		&diff.PathsDiff{
			AddedEndpoints:    []string{"a"},
			DeletedEndpoints:  []string{"ab"},
			ModifiedEndpoints: diff.ModifiedEndpoints{"abc": &diff.EndpointDiff{}},
		},
	}

	diffResult.FilterByRegex("ab")

	require.Empty(t, diffResult.PathsDiff.AddedEndpoints)
	require.Equal(t, []string{"ab"}, diffResult.PathsDiff.DeletedEndpoints)
	require.Contains(t, diffResult.PathsDiff.ModifiedEndpoints, "abc")
}
