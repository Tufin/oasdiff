package lint_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/lint"
)

func TestRequiredParam_PathOK(t *testing.T) {
	const source = "../data/lint/required-params/path.yaml"
	errs := lint.Run(*lint.NewConfig([]lint.Check{lint.RequiredParamsCheck}), source, loadFrom(t, source))
	require.Empty(t, errs)
}

func TestRequiredParam_WithDefault(t *testing.T) {
	const source = "../data/lint/required-params/with_default.yaml"
	errs := lint.Run(*lint.NewConfig([]lint.Check{lint.RequiredParamsCheck}), source, loadFrom(t, source))
	require.Empty(t, errs)
}
