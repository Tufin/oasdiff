package lint_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/lint"
)

func TestRequiredParam_PathOK(t *testing.T) {
	const source = "../data/lint/required-params/path.yaml"
	errs := lint.Run(lint.NewConfig([]lint.Check{lint.RequiredParamsCheck}), source, loadFrom(t, source))
	require.Empty(t, errs)
}

func TestRequiredParam_PathWithDefault(t *testing.T) {
	const source = "../data/lint/required-params/path_with_default.yaml"
	errs := lint.Run(lint.NewConfig([]lint.Check{lint.RequiredParamsCheck}), source, loadFrom(t, source))
	require.Len(t, errs, 1)
	require.Equal(t, "required-param-with-default", errs[0].Id)
	require.Equal(t, "required path parameter \"bookId\" shouldn't have a default value: /books/{bookId}", errs[0].Text)
	require.Equal(t, lint.LEVEL_ERROR, errs[0].Level)
	require.Equal(t, source, errs[0].Source)
}

func TestRequiredParam_MethodWithDefault(t *testing.T) {
	const source = "../data/lint/required-params/method_with_default.yaml"
	errs := lint.Run(lint.NewConfig([]lint.Check{lint.RequiredParamsCheck}), source, loadFrom(t, source))
	require.Len(t, errs, 1)
	require.Equal(t, "required-param-with-default", errs[0].Id)
	require.Equal(t, "required path parameter \"bookId\" shouldn't have a default value: GET /books/{bookId}", errs[0].Text)
	require.Equal(t, lint.LEVEL_ERROR, errs[0].Level)
	require.Equal(t, source, errs[0].Source)
}
