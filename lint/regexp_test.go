package lint_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/lint"
)

func TestRegexCheck(t *testing.T) {

	const source = "../data/lint/regex/openapi-invalid-regex.yaml"
	errs := lint.Run(lint.NewConfig([]lint.Check{lint.SchemaCheck}), source, loadFrom(t, source))
	require.Len(t, errs, 1)
	require.Equal(t, "invalid-regex-pattern", errs[0].Id)
	require.Equal(t, lint.LEVEL_ERROR, errs[0].Level)
	require.Equal(t, source, errs[0].Source)
}

func TestRegexCheck_Embedded(t *testing.T) {

	const source = "../data/lint/regex/openapi-invalid-regex-embedded.yaml"
	errs := lint.Run(lint.NewConfig([]lint.Check{lint.SchemaCheck}), source, loadFrom(t, source))
	require.Len(t, errs, 7)
	for i := range errs {
		require.Equal(t, "invalid-regex-pattern", errs[i].Id)
		require.Equal(t, lint.LEVEL_ERROR, errs[i].Level)
		require.Equal(t, source, errs[i].Source)
	}
}

func TestRegexCheck_Circular(t *testing.T) {

	const source = "../data/circular2.yaml"
	errs := lint.Run(lint.NewConfig([]lint.Check{lint.SchemaCheck}), source, loadFrom(t, source))
	require.Empty(t, errs)
}
