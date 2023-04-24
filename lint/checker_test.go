package lint_test

import (
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/lint"
	"github.com/tufin/oasdiff/load"
)

func loadFrom(t *testing.T, path string) *load.OpenAPISpecInfo {
	t.Helper()

	loader := openapi3.NewLoader()
	oas, err := loader.LoadFromFile(path)
	require.NoError(t, err)
	return &load.OpenAPISpecInfo{Spec: oas, Url: path}
}

func TestRun(t *testing.T) {

	const source = "../data/lint/openapi.yaml"
	require.Empty(t, lint.Run(*lint.DefaultConfig(), source, loadFrom(t, source)))
}
