package checker

type Change interface {
	GetSection() string
	IsBreaking() bool
	GetId() string
	GetText(l Localizer) string
	GetArgs() []any
	GetUncolorizedText(l Localizer) string
	GetComment(l Localizer) string
	GetLevel() Level
	GetOperation() string
	GetOperationId() string
	GetPath() string
	GetSource() string
	GetSourceFile() string
	GetSourceLine() int
	GetSourceLineEnd() int
	GetSourceColumn() int
	GetSourceColumnEnd() int
	MatchIgnore(ignorePath, ignoreLine string, l Localizer) bool
	SingleLineError(l Localizer, colorMode ColorMode) string
	MultiLineError(l Localizer, colorMode ColorMode) string
}
