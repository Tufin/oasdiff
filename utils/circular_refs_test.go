package utils_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/utils"
)

func Test_VisitedRefs(t *testing.T) {
	v := utils.VisitedRefs{}
	require.False(t, v.IsVisited("test"))
	v.Add("test")
	require.True(t, v.IsVisited("test"))
	v.Remove("test")
	require.False(t, v.IsVisited("test"))
}
