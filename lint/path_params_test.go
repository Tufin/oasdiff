package lint_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/lint"
)

func TestPathParam_PathOK(t *testing.T) {
	const source = "../data/lint/path-params/path.yaml"
	errs := lint.Run(lint.NewConfig([]lint.Check{lint.PathParamsCheck}), source, loadFrom(t, source))
	require.Empty(t, errs)
}

func TestPathParam_MethodOK(t *testing.T) {
	const source = "../data/lint/path-params/method.yaml"
	errs := lint.Run(lint.NewConfig([]lint.Check{lint.PathParamsCheck}), source, loadFrom(t, source))
	require.Empty(t, errs)
}

func TestPathParam_MethodExtra(t *testing.T) {
	const source = "../data/lint/path-params/method-extra.yaml"
	errs := lint.Run(lint.NewConfig([]lint.Check{lint.PathParamsCheck}), source, loadFrom(t, source))
	require.Len(t, errs, 1)
	require.Equal(t, "path-param-extra", errs[0].Id)
	require.Equal(t, "path parameter \"bookId\" appears in the parameters section of the operation but is missing in the URL: GET /books", errs[0].Text)
	require.Equal(t, lint.LEVEL_ERROR, errs[0].Level)
	require.Equal(t, source, errs[0].Source)
}

func TestPathParam_PathExtra(t *testing.T) {
	const source = "../data/lint/path-params/path-extra.yaml"
	errs := lint.Run(lint.NewConfig([]lint.Check{lint.PathParamsCheck}), source, loadFrom(t, source))
	require.Len(t, errs, 1)
	require.Equal(t, "path-param-extra", errs[0].Id)
	require.Equal(t, "path parameter \"bookId\" appears in the parameters section of the path but is missing in the URL: GET /books", errs[0].Text)
	require.Equal(t, lint.LEVEL_ERROR, errs[0].Level)
	require.Equal(t, source, errs[0].Source)
}

func TestPathParam_PathMissing(t *testing.T) {
	const source = "../data/lint/path-params/path-missing.yaml"
	errs := lint.Run(lint.NewConfig([]lint.Check{lint.PathParamsCheck}), source, loadFrom(t, source))
	require.Len(t, errs, 1)
	require.Equal(t, "path-param-missing", errs[0].Id)
	require.Equal(t, "path parameter \"bookId\" appears in the URL path but is missing from the parameters section of the path and operation: GET /books/{bookId}", errs[0].Text)
	require.Equal(t, lint.LEVEL_WARN, errs[0].Level)
	require.Equal(t, source, errs[0].Source)
}

func TestPathParam_Duplicate(t *testing.T) {
	const source = "../data/lint/path-params/duplicate.yaml"
	errs := lint.Run(lint.NewConfig([]lint.Check{lint.PathParamsCheck}), source, loadFrom(t, source))
	require.Len(t, errs, 1)
	require.Equal(t, "path-param-duplicate", errs[0].Id)
	require.Equal(t, "path parameter \"bookId\" is defined both in path and in operation: GET /books/{bookId}", errs[0].Text)
	require.Equal(t, lint.LEVEL_WARN, errs[0].Level)
	require.Equal(t, source, errs[0].Source)
}

func TestPathParam_NotRequired(t *testing.T) {
	const source = "../data/lint/path-params/not-required.yaml"
	errs := lint.Run(lint.NewConfig([]lint.Check{lint.PathParamsCheck}), source, loadFrom(t, source))
	require.Len(t, errs, 1)
	require.Equal(t, "path-param-not-required", errs[0].Id)
	require.Equal(t, "path parameter \"bookId\" should have required=true: GET /books/{bookId}", errs[0].Text)
	require.Equal(t, lint.LEVEL_ERROR, errs[0].Level)
	require.Equal(t, source, errs[0].Source)
}
