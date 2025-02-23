package localizations_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/checker/localizations"
)

func TestLang_Exists(t *testing.T) {
	require.Equal(t, []string{localizations.LangEn, localizations.LangRu}, localizations.GetSupportedLanguages())
}
