package internal

import "errors"

type Lang string

const (
	LangEn Lang = "en"
	LangRu Lang = "ru"
)

// String is used both by fmt.Print and by Cobra in help text
func (lang *Lang) String() string {
	return string(*lang)
}

// Set must have pointer receiver so it doesn't change the value of a copy
func (lang *Lang) Set(v string) error {
	switch v {
	case "en", "ru":
		*lang = Lang(v)
		return nil
	default:
		return errors.New(`must be one of "en", or "ru"`)
	}
}

// Type is only used in help text
func (lang *Lang) Type() string {
	return "lang"
}
