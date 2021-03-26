package text_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/diff"
	"github.com/tufin/oasdiff/text"
)

func l(t *testing.T, v int) *openapi3.Swagger {
	loader := openapi3.NewSwaggerLoader()
	oas, err := loader.LoadSwaggerFromFile(fmt.Sprintf("../data/openapi-test%d.yaml", v))
	require.NoError(t, err)
	return oas
}

func d(t *testing.T, config *diff.Config, v1, v2 int) *diff.Diff {
	d, err := diff.Get(config, l(t, v1), l(t, v2))
	require.NoError(t, err)
	return d
}

func TestDiff1(t *testing.T) {
	var buf bytes.Buffer
	text.Print(d(t, &diff.Config{}, 3, 5), &buf)

	require.Contains(t, buf.String(), "GET /api/{domain}/{project}/install-command")
}

func TestDiff2(t *testing.T) {
	var buf bytes.Buffer
	text.Print(d(t, &diff.Config{}, 5, 3), &buf)

	require.Contains(t, buf.String(), "GET /api/{domain}/{project}/install-command")
}

func TestDiff3(t *testing.T) {
	var buf bytes.Buffer
	text.Print(d(t, &diff.Config{}, 1, 3), &buf)

	require.Contains(t, buf.String(), "GET /api/{domain}/{project}/install-command")
}

func TestNoDiff(t *testing.T) {
	var buf bytes.Buffer
	text.Print(d(t, &diff.Config{}, 3, 3), &buf)

	require.Equal(t, buf.String(), "No changes\n")
}
