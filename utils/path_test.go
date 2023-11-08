package utils_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/utils"
)

func Test_NormalizeTemplatedPath(t *testing.T) {
	normalizedPath, numOfParams, paramNames := utils.NormalizeTemplatedPath("/person/{personName}")
	require.Equal(t, "/person/{}", normalizedPath)
	require.Equal(t, uint(0x1), numOfParams)
	require.Equal(t, []string{"personName"}, paramNames)
}
