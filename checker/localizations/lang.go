package localizations

const (
	LangDefault = LangEn
	LangEn      = "en"
	LangRu      = "ru"
)

func GetSupportedLanguages() []string {
	return []string{LangEn, LangRu}
}
