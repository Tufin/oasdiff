package diff_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/diff"
)

func TestBreaking_Same(t *testing.T) {
	require.False(t, d(t, diff.NewConfig(), 1, 1).Breaking())
}

func TestBreaking_DeletedPaths(t *testing.T) {
	require.True(t, d(t, diff.NewConfig(), 1, 2).Breaking())
}
