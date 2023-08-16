package checker

type Change interface {
	GetId() string
	GetText() string
	GetComment() string
	GetLevel() Level
	GetOperation() string
	GetOperationId() string
	GetPath() string

	MatchIgnore(ignorePath, ignoreLine string) bool
	LocalizedError(l Localizer) string
	PrettyErrorText(l Localizer) string
	Error() string
}
