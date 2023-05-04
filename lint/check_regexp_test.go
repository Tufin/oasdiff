package lint_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/lint"
)

func TestRegexCheck(t *testing.T) {

	const source = "../data/lint/openapi-invalid-regex.yaml"
	errs := lint.Run(*lint.NewConfig([]lint.Check{lint.RegexCheck}), source, loadFrom(t, source))
	require.Len(t, errs, 1)
	require.Equal(t, "invalid-regex-pattern", errs[0].Id)
}

func TestRegexCheck_Embedded(t *testing.T) {

	const source = "../data/lint/openapi-invalid-regex-embedded.yaml"
	errs := lint.Run(*lint.NewConfig([]lint.Check{lint.RegexCheck}), source, loadFrom(t, source))
	require.Len(t, errs, 7)
	for i := range errs {
		require.Equal(t, "invalid-regex-pattern", errs[i].Id)
	}
}
