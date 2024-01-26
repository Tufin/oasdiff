package diff_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/diff"
	"github.com/tufin/oasdiff/utils"
)

func TestParamNamesByLocation_Len(t *testing.T) {
	require.Equal(t, 3, diff.ParamNamesByLocation{
		"query":  utils.StringList{"name"},
		"header": utils.StringList{"id", "organization"},
	}.Len())
}

func TestParamDiffByLocation_Len(t *testing.T) {
	require.Equal(t, 3, diff.ParamDiffByLocation{
		"query":  diff.ParamDiffs{"query": &diff.ParameterDiff{}},
		"header": diff.ParamDiffs{"id": &diff.ParameterDiff{}, "organization": &diff.ParameterDiff{}},
	}.Len())
}
