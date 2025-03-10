package formatters_test

import (
	"testing"

	"github.com/oasdiff/oasdiff/formatters"
	"github.com/stretchr/testify/require"
)

func TestTypes(t *testing.T) {
	require.Equal(t, formatters.GetSupportedFormats(), []string{"yaml", "json", "text", "markup", "markdown", "singleline", "html", "githubactions", "junit", "sarif"})
}
