package load_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/load"
)

func TestSource_NewStdin(t *testing.T) {
	require.True(t, load.NewSource("-").IsStdin())
}

func TestSource_NewFile(t *testing.T) {
	require.True(t, load.NewSource("../spec.yaml").IsFile())
}

func TestSource_String(t *testing.T) {
	require.Equal(t, "stdin", load.NewSource("-").String())
}

func TestSource_OutStdin(t *testing.T) {
	require.Equal(t, `stdin`, load.NewSource("-").Out())
}

func TestSource_Out(t *testing.T) {
	require.Equal(t, `"http://twitter.com"`, load.NewSource("http://twitter.com").Out())
}
