package generator_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/checker/generator"
)

func TestGenerator(t *testing.T) {
	require.NoError(t, generator.Generate())
}
