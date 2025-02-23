package localizations_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/checker/localizations"
)

func TestLocalizations_ExistsNoSunstitute(t *testing.T) {
	locales := localizations.New(localizations.LangEn, localizations.LangDefault)
	require.Equal(t, "the %s response's property pattern was changed from %s to %s for the status %s", locales.Get("messages.response-property-pattern-changed"))
}

func TestLocalizations_SetLocal(t *testing.T) {
	locales := localizations.New(localizations.LangEn, localizations.LangDefault).SetLocale(localizations.LangEn)
	require.Equal(t, "the %s response's property pattern was changed from %s to %s for the status %s", locales.Get("messages.response-property-pattern-changed"))
}

func TestLocalizations_SetFallbackLocal(t *testing.T) {
	locales := localizations.New(localizations.LangEn, localizations.LangDefault).SetFallbackLocale(localizations.LangEn)
	require.Equal(t, "the %s response's property pattern was changed from %s to %s for the status %s", locales.Get("messages.response-property-pattern-changed"))
}

func TestLocalizations_SetLocals(t *testing.T) {
	locales := localizations.New(localizations.LangEn, localizations.LangDefault).SetLocales(localizations.LangEn, localizations.LangEn)
	require.Equal(t, "the %s response's property pattern was changed from %s to %s for the status %s", locales.Get("messages.response-property-pattern-changed"))
}

func TestLocalizations_NotExists(t *testing.T) {
	locales := localizations.New(localizations.LangEn, localizations.LangDefault)
	require.Equal(t, "invalid", locales.Get("invalid"))
}

func TestLocalizations_ExistsWithSunstitute(t *testing.T) {
	locales := localizations.New(localizations.LangEn, localizations.LangDefault)
	locales.Localizations = map[string]string{
		"en.messages.response-property-pattern-changed": "{{/* a comment */}}",
	}

	replacements := localizations.Replacements{}
	require.Equal(t, "", locales.Get("messages.response-property-pattern-changed", &replacements))
}
