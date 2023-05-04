package lint_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/lint"
)

func TestPathParam(t *testing.T) {

	const source = "../data/openapi-test2.yaml"
	errs := lint.Run(*lint.NewConfig([]lint.Check{lint.PathParamsCheck}), source, loadFrom(t, source))
	require.Len(t, errs, 2)
}
