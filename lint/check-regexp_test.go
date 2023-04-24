package lint_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/lint"
)

func TestRegexCheck(t *testing.T) {

	const source = "../data/lint/openapi-invalid-regex.yaml"
	require.Len(t, lint.Run(*lint.NewConfig([]lint.Check{lint.RegexCheck}),
		source, loadFrom(t, source)), 1)
}
