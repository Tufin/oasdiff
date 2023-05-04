package lint_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/lint"
)

func TestPathParam_Missing(t *testing.T) {

	const source = "../data/openapi-test2.yaml"
	errs := lint.Run(*lint.NewConfig([]lint.Check{lint.PathParamsCheck}), source, loadFrom(t, source))
	require.Len(t, errs, 2)
	for _, err := range errs {
		require.Equal(t, "path-param-missing", err.Id)
	}
}

func TestPathParam_OK(t *testing.T) {

	const source = "../data/param-rename/method-base.yaml"
	errs := lint.Run(*lint.NewConfig([]lint.Check{lint.PathParamsCheck}), source, loadFrom(t, source))
	require.Empty(t, errs)
}
