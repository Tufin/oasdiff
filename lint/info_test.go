package lint_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/lint"
)

func TestInfoCheck(t *testing.T) {

	const source = "../data/lint/openapi-info.yaml"
	errs := lint.Run(*lint.NewConfig([]lint.Check{lint.InfoCheck}), source, loadFrom(t, source))
	require.Len(t, errs, 1)
	require.Equal(t, "info-version-missing", errs[0].Id)
	require.Equal(t, lint.LEVEL_WARN, errs[0].Level)
	require.True(t, strings.Contains(errs[0].Text, "license URL is missing"))
	require.Equal(t, source, errs[0].Source)
}
