package lint_test

import (
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/oasdiff/oasdiff/lint"
	"github.com/oasdiff/oasdiff/load"
	"github.com/stretchr/testify/require"
)

func loadFrom(t *testing.T, path string) *load.SpecInfo {
	t.Helper()

	loader := openapi3.NewLoader()
	oas, err := loader.LoadFromFile(path)
	require.NoError(t, err)
	return &load.SpecInfo{Spec: oas, Url: path}
}

func TestRun(t *testing.T) {

	const source = "../data/lint/openapi.yaml"
	require.Empty(t, lint.Run(lint.DefaultConfig(), source, loadFrom(t, source)))
}
