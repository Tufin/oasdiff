package internal_test

import (
	"io"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/internal"
)

func Test_NoArgs(t *testing.T) {
	failOnDiff, diffEmpty, returnErr := internal.Run([]string{}, io.Discard)
	require.False(t, failOnDiff)
	require.False(t, diffEmpty)
	require.Equal(t, 101, returnErr.Code)
}

func Test_OneArg(t *testing.T) {
	failOnDiff, diffEmpty, returnErr := internal.Run([]string{"oasdiff"}, io.Discard)
	require.False(t, failOnDiff)
	require.False(t, diffEmpty)
	require.Equal(t, 101, returnErr.Code)
}

func Test_NoRevision(t *testing.T) {
	failOnDiff, diffEmpty, returnErr := internal.Run([]string{"oasdiff", "-base", "base.yaml"}, io.Discard)
	require.False(t, failOnDiff)
	require.False(t, diffEmpty)
	require.Equal(t, 101, returnErr.Code)
}

func Test_Basic(t *testing.T) {
	failOnDiff, diffEmpty, returnErr := internal.Run([]string{"oasdiff", "-base", "../data/openapi-test1.yaml", "-revision", "../data/openapi-test3.yaml"}, io.Discard)
	require.False(t, failOnDiff)
	require.False(t, diffEmpty)
	require.Nil(t, returnErr)
}
