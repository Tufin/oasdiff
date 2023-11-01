package utils_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/utils"
)

func Test_StringList(t *testing.T) {
	l := utils.StringList{}
	require.True(t, l.Empty())
	require.Equal(t, l.String(), "")
	l.Set("b,a")
	l.Sort()
	require.True(t, l.Contains("a"))
	require.True(t, l.Contains("b"))
	require.True(t, l.Minus(l).ToStringSet().Empty())
}
