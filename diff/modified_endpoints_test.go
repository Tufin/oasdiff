package diff_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/diff"
)

func TestModifiedEndpoints(t *testing.T) {
	m := diff.ModifiedEndpoints{diff.Endpoint{}: nil}
	require.NotEmpty(t, m.ToEndpoints())
}
