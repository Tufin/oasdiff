package utils_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/utils"
)

func TestStringList(t *testing.T) {
	l := utils.StringList{}
	require.True(t, l.Empty())
	require.Equal(t, l.String(), "")
	require.NoError(t, l.Set("b,a"))
	require.Equal(t, l.String(), "b, a")
	l = l.Sort()
	require.Equal(t, l.String(), "a, b")
	require.True(t, l.Contains("a"))
	require.True(t, l.Contains("b"))
	require.True(t, l.Minus(l).ToStringSet().Empty())
}

func TestStringList_CartesianProduct(t *testing.T) {
	l1 := utils.StringList{"a", "b", "c"}
	l2 := utils.StringList{"x", "y"}
	require.Equal(t, 6, len(l1.CartesianProduct(l2)))
	require.Equal(t, utils.StringPair{"b", "y"}, l1.CartesianProduct(l2)[3])
}

func TestStringList_Is(t *testing.T) {
	l1 := utils.StringList{"a", "b", "c"}
	require.False(t, l1.Is("x"))
}

func TestStringList_IsNil(t *testing.T) {
	var l1 *utils.StringList
	require.False(t, l1.Is("x"))
}

func TestStringList_Comtains(t *testing.T) {
	l1 := utils.StringList{"a", "b", "c"}
	require.False(t, l1.Contains("x"))
}
