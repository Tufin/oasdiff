package checker

import (
	"fmt"
	"strings"

	"github.com/TwiN/go-color"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/tufin/oasdiff/diff"
	"github.com/tufin/oasdiff/load"
)

// ApiChange represnts a change in the Paths Section of an OpenAPI spec
type ApiChange struct {
	Id          string
	Args        []any
	Comment     string
	Level       Level
	Operation   string
	OperationId string
	Path        string
	Source      *load.Source
	Tags        map[string]any

	SourceFile      string
	SourceLine      int
	SourceLineEnd   int
	SourceColumn    int
	SourceColumnEnd int
}

func newApiChangeGetter(config *Config, operationsSources *diff.OperationsSourcesMap) ChangeGetter {
	return func(id string, defaultLevel Level, args []any, comment string, method string, op *openapi3.Operation, path string, sourceOp *openapi3.Operation) Change {
		return ApiChange{
			Id:          id,
			Comment:     comment,
			Level:       config.getLogLevel(id, defaultLevel),
			Operation:   method,
			Path:        path,
			OperationId: op.OperationID,
			Source:      load.NewSource((*operationsSources)[sourceOp]),
			Args:        args,
			Tags:        getTags(config, op),
		}
	}
}

func getTags(config *Config, op *openapi3.Operation) map[string]any {
	if config.Tags == nil {
		return nil
	}

	result := map[string]any{}
	for _, tag := range config.Tags {
		if val, ok := op.Extensions[tag]; ok {
			result[tag] = val
		}
	}

	return result
}

func (c ApiChange) GetSection() string {
	return "paths"
}

func (c ApiChange) IsBreaking() bool {
	return c.GetLevel().IsBreaking()
}

func (c ApiChange) MatchIgnore(ignorePath, ignoreLine string, l Localizer) bool {
	if ignorePath == "" {
		return false
	}

	return ignorePath == strings.ToLower(c.Path) &&
		strings.Contains(ignoreLine, strings.ToLower(c.Operation+" "+c.Path)) &&
		strings.Contains(ignoreLine, strings.ToLower(c.GetUncolorizedText(l)))
}

func (c ApiChange) GetId() string {
	return c.Id
}

func (c ApiChange) GetText(l Localizer) string {
	return l(c.Id, colorizedValues(c.Args)...)
}

func (c ApiChange) GetArgs() []any {
	return c.Args
}

func (c ApiChange) GetUncolorizedText(l Localizer) string {
	return l(c.Id, quotedValues(c.Args)...)
}

func (c ApiChange) GetComment(l Localizer) string {
	return l(c.Comment)
}

func (c ApiChange) GetLevel() Level {
	return c.Level
}

func (c ApiChange) GetOperation() string {
	return c.Operation
}

func (c ApiChange) GetOperationId() string {
	return c.OperationId
}

func (c ApiChange) GetPath() string {
	return c.Path
}

func (c ApiChange) GetSource() string {
	return c.Source.String()
}

func (c ApiChange) GetTags() map[string]any {
	return c.Tags
}

func (c ApiChange) GetSourceFile() string {
	if c.SourceFile != "" {
		return c.SourceFile
	}

	if c.Source.IsFile() {
		return c.Source.String()
	}

	return ""
}

func (c ApiChange) GetSourceLine() int {
	return c.SourceLine
}

func (c ApiChange) GetSourceLineEnd() int {
	return c.SourceLineEnd
}

func (c ApiChange) GetSourceColumn() int {
	return c.SourceColumn
}

func (c ApiChange) GetSourceColumnEnd() int {
	return c.SourceColumnEnd
}

func (c ApiChange) SingleLineError(l Localizer, colorMode ColorMode) string {
	const format = "%s %s %s, %s API %s %s %s [%s]. %s"

	if isColorEnabled(colorMode) {
		return fmt.Sprintf(format, c.Level.PrettyString(), l("at"), c.GetSource(), l("in"), color.InGreen(c.Operation), color.InGreen(c.Path), c.GetText(l), color.InYellow(c.Id), c.GetComment(l))
	}

	return fmt.Sprintf(format, c.Level.String(), l("at"), c.GetSource(), l("in"), c.Operation, c.Path, c.GetUncolorizedText(l), c.Id, c.GetComment(l))

}

func (c ApiChange) MultiLineError(l Localizer, colorMode ColorMode) string {
	const format = "%s\t[%s] %s %s\t\n\t%s API %s %s\n\t\t%s%s"

	if isColorEnabled(colorMode) {
		return fmt.Sprintf(format, c.Level.PrettyString(), color.InYellow(c.Id), l("at"), c.GetSource(), l("in"), color.InGreen(c.Operation), color.InGreen(c.Path), c.GetText(l), multiLineComment(c.GetComment(l)))
	}

	return fmt.Sprintf(format, c.Level.String(), c.Id, l("at"), c.GetSource(), l("in"), c.Operation, c.Path, c.GetUncolorizedText(l), multiLineComment(c.GetComment(l)))
}

func multiLineComment(comment string) string {
	if comment == "" {
		return ""
	}
	return fmt.Sprintf("\n\t\t%s", comment)
}
