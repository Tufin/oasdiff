package checker

type Change interface {
	IsBreaking() bool
	GetId() string
	GetText(l Localizer) string
	GetUncolorizedText(l Localizer) string
	GetComment() string
	GetLevel() Level
	GetOperation() string
	GetOperationId() string
	GetPath() string
	GetSourceFile() string
	GetSourceLine() int
	GetSourceLineEnd() int
	GetSourceColumn() int
	GetSourceColumnEnd() int
	MatchIgnore(ignorePath, ignoreLine string, l Localizer) bool
	LocalizedError(l Localizer) string
	PrettyErrorText(l Localizer) string
}
