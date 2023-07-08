package checker

import "github.com/tufin/oasdiff/checker/localizations"

type Change interface {
	GetId() string
	GetText() string
	GetComment() string
	GetLevel() Level
	GetOperation() string
	GetOperationId() string
	GetPath() string

	MatchIgnore(ignorePath, ignoreLine string) bool
	LocalizedError(l localizations.Localizer) string
	PrettyErrorText(l localizations.Localizer) string
	Error() string
}
