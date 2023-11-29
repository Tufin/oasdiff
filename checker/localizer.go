package checker

import (
	"fmt"

	"github.com/tufin/oasdiff/checker/localizations"
)

type Localizer func(key string, args ...interface{}) string

func NewDefaultLocalizer() Localizer {
	return NewLocalizer(localizations.LangDefault)
}

func NewLocalizer(locale string) Localizer {
	locales := localizations.New(locale, localizations.LangDefault)

	return func(originalKey string, args ...interface{}) string {
		key := "messages." + originalKey
		pattern := locales.Get(key)

		// if key not found, return original key
		// TODO: improve localizations to return error when key not found
		if pattern == key {
			return originalKey
		}

		return fmt.Sprintf(pattern, args...)
	}
}
