package diff_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/diff"
)

func TestModifiedEndpoints(t *testing.T) {
	m := diff.ModifiedEndpoints{diff.Endpoint{}: nil}
	require.Len(t, m.ToEndpoints(), 1)
}

func TestModifiedEndpoints_Empty(t *testing.T) {
	m := diff.ModifiedEndpoints{}
	require.Empty(t, m.ToEndpoints())
}
