package linter_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/linter"
)

func TestRegexCheck(t *testing.T) {
	require.Len(t, linter.Run(*linter.NewConfig([]linter.Check{linter.RegexCheck}),
		loadFrom(t, "../data/linter/openapi-invalid-regex.yaml")), 1)
}
