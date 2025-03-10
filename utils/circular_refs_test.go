package utils_test

import (
	"testing"

	"github.com/oasdiff/oasdiff/utils"
	"github.com/stretchr/testify/require"
)

func Test_VisitedRefs(t *testing.T) {
	v := utils.VisitedRefs{}
	require.False(t, v.IsVisited("test"))
	v.Add("test")
	require.True(t, v.IsVisited("test"))
	v.Remove("test")
	require.False(t, v.IsVisited("test"))
}
