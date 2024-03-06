package checker

import "github.com/getkin/kin-openapi/openapi3"

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
	GetTags() map[string]any
	GetSourceFile() string
	GetSourceLine() int
	GetSourceLineEnd() int
	GetSourceColumn() int
	GetSourceColumnEnd() int
	MatchIgnore(ignorePath, ignoreLine string, l Localizer) bool
	SingleLineError(l Localizer, colorMode ColorMode) string
	MultiLineError(l Localizer, colorMode ColorMode) string
}

type ChangeGetter func(id string, defaultLevel Level, args []any, comment string, method string, op *openapi3.Operation, path string, sourceOp *openapi3.Operation) Change
