package lint_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/lint"
)

func TestRequirePropertiesCheck_OK(t *testing.T) {

	const source = "../data/lint/required-properties/ok.yaml"
	errs := lint.Run(lint.NewConfig([]lint.Check{lint.SchemaCheck}), source, loadFrom(t, source))
	require.Empty(t, errs)
}

func TestRequirePropertiesCheck_Extra(t *testing.T) {

	const source = "../data/lint/required-properties/extra.yaml"
	errs := lint.Run(lint.NewConfig([]lint.Check{lint.SchemaCheck}), source, loadFrom(t, source))
	require.Len(t, errs, 1)
	require.Equal(t, "extra_required_props", errs[0].Id)
	require.Equal(t, lint.LEVEL_ERROR, errs[0].Level)
	require.Equal(t, source, errs[0].Source)
}
