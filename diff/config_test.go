package diff_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/diff"
	"github.com/tufin/oasdiff/utils"
)

func TestConfig_Valid(t *testing.T) {
	excludeElements := utils.StringList{"examples", "description", "endpoints", "title", "summary"}
	require.Empty(t, diff.ValidateExcludeElements(excludeElements))
}

func TestConfig_Invalid(t *testing.T) {
	excludeElements := utils.StringList{"y", "x"}
	require.Equal(t, utils.StringList{"x", "y"}, diff.ValidateExcludeElements(excludeElements))
}
func TestConfig_Default(t *testing.T) {
	c := diff.NewConfig()
	c.SetExcludeElements(utils.StringSet{}, false, false, false)
	require.False(t, c.IsExcludeExamples())
	require.False(t, c.IsExcludeDescription())
	require.False(t, c.IsExcludeEndpoints())
	require.False(t, c.IsExcludeTitle())
	require.False(t, c.IsExcludeSummary())
}

func TestConfig_BC(t *testing.T) {
	c := diff.NewConfig()
	c.SetExcludeElements(utils.StringSet{}, true, true, true)
	require.True(t, c.IsExcludeExamples())
	require.True(t, c.IsExcludeDescription())
	require.True(t, c.IsExcludeEndpoints())
	require.False(t, c.IsExcludeTitle())
	require.False(t, c.IsExcludeSummary())
}

func TestConfig_ExcludeElements(t *testing.T) {
	c := diff.NewConfig()
	excludeElements := utils.StringList{"examples", "description", "endpoints", "title", "summary"}
	c.SetExcludeElements(excludeElements.ToStringSet(), false, false, false)
	require.True(t, c.IsExcludeExamples())
	require.True(t, c.IsExcludeDescription())
	require.True(t, c.IsExcludeEndpoints())
	require.True(t, c.IsExcludeTitle())
	require.True(t, c.IsExcludeSummary())
}
