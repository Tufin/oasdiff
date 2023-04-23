package linter_test

import (
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/linter"
	"github.com/tufin/oasdiff/load"
)

func loadFrom(t *testing.T, path string) *load.OpenAPISpecInfo {
	loader := openapi3.NewLoader()
	oas, err := loader.LoadFromFile(path)
	require.NoError(t, err)
	return &load.OpenAPISpecInfo{Spec: oas, Url: path}
}

func TestRun(t *testing.T) {
	require.Empty(t, linter.Run(*linter.DefaultConfig(), loadFrom(t, "../data/linter/openapi.yaml")))
}
