package formatters_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/formatters"
)

func TestTypes(t *testing.T) {
	require.Equal(t, formatters.GetSupportedFormats(), []string{"yaml", "json", "text", "markup", "markdown", "singleline", "html", "githubactions", "junit", "sarif"})
}
