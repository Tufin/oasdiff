package diff_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/diff"
)

func TestFilterByRegex(t *testing.T) {
	diffResult := diff.DiffResult{
		&diff.PathDiff{
			AddedEndpoints:    []string{"a"},
			DeletedEndpoints:  []string{"ab"},
			ModifiedEndpoints: diff.ModifiedEndpoints{"abc": &diff.EndpointDiff{}},
		},
	}

	diffResult.FilterByRegex("ab")

	require.Empty(t, diffResult.PathDiff.AddedEndpoints)
	require.Equal(t, []string{"ab"}, diffResult.PathDiff.DeletedEndpoints)
	require.Contains(t, diffResult.PathDiff.ModifiedEndpoints, "abc")
}
