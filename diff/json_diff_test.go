package diff_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/diff"
)

func TestJsonOperationString_Add(t *testing.T) {
	p := diff.JsonOperation{
		OldValue: "old",
		Value:    "new",
		Type:     "add",
		From:     "from",
		Path:     "path",
	}
	require.Equal(t, "Added path with value: 'new'", p.String())
}

func TestJsonOperationString_Remove(t *testing.T) {
	p := diff.JsonOperation{
		OldValue: "old",
		Value:    "new",
		Type:     "remove",
		From:     "from",
		Path:     "path",
	}
	require.Equal(t, "Removed path with value: 'old'", p.String())
}

func TestJsonOperationString_Replace(t *testing.T) {
	p := diff.JsonOperation{
		OldValue: "old",
		Value:    "new",
		Type:     "replace",
		From:     "from",
		Path:     "path",
	}
	require.Equal(t, "Modified path from 'old' to 'new'", p.String())
}

func TestJsonOperationString_ReplaceNoPath(t *testing.T) {
	p := diff.JsonOperation{
		OldValue: "old",
		Value:    "new",
		Type:     "replace",
		From:     "from",
		Path:     "",
	}
	require.Equal(t, "Modified value from 'old' to 'new'", p.String())
}

func TestJsonOperationString_Unknown(t *testing.T) {
	p := diff.JsonOperation{
		OldValue: "old",
		Value:    "new",
		Type:     "unknown",
		From:     "from",
		Path:     "path",
	}
	require.Equal(t, "unknown path", p.String())
}
