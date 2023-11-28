package checker

import (
	"fmt"
	"strings"

	"github.com/TwiN/go-color"
)

// ApiChange represnts a change in the Paths Section of an OpenAPI spec
type ApiChange struct {
	Id          string `json:"id,omitempty" yaml:"id,omitempty"`
	Args        []any
	Comment     string `json:"comment,omitempty" yaml:"comment,omitempty"`
	Level       Level  `json:"level" yaml:"level"`
	Operation   string `json:"operation,omitempty" yaml:"operation,omitempty"`
	OperationId string `json:"operationId,omitempty" yaml:"operationId,omitempty"`
	Path        string `json:"path,omitempty" yaml:"path,omitempty"`
	Source      string `json:"source,omitempty" yaml:"source,omitempty"`

	SourceFile      string `json:"-" yaml:"-"`
	SourceLine      int    `json:"-" yaml:"-"`
	SourceLineEnd   int    `json:"-" yaml:"-"`
	SourceColumn    int    `json:"-" yaml:"-"`
	SourceColumnEnd int    `json:"-" yaml:"-"`
}

func (c ApiChange) IsBreaking() bool {
	return c.GetLevel().IsBreaking()
}

func (c ApiChange) MatchIgnore(ignorePath, ignoreLine string, l Localizer) bool {
	if ignorePath == "" {
		return false
	}
	x := c.GetUncolorizedText(l)
	return ignorePath == strings.ToLower(c.Path) &&
		strings.Contains(ignoreLine, strings.ToLower(c.Operation+" "+c.Path)) &&
		strings.Contains(ignoreLine, strings.ToLower(x))
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
	return c.Source
}

func (c ApiChange) GetSourceFile() string {
	return c.SourceFile
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

func (c ApiChange) LocalizedError(l Localizer) string {
	return fmt.Sprintf("%s %s %s, %s API %s %s %s [%s]. %s", c.Level, l("at"), c.Source, l("in"), c.Operation, c.Path, c.GetText(l), c.Id, c.GetComment(l))
}

func (c ApiChange) PrettyErrorText(l Localizer) string {
	if IsPipedOutput() {
		return c.LocalizedError(l)
	}

	comment := c.GetComment(l)
	if comment != "" {
		comment = fmt.Sprintf("\n\t\t%s", c.GetComment(l))
	}
	return fmt.Sprintf("%s\t[%s] %s %s\t\n\t%s API %s %s\n\t\t%s%s", c.Level.PrettyString(), color.InYellow(c.Id), l("at"), c.Source, l("in"), color.InGreen(c.Operation), color.InGreen(c.Path), c.GetText(l), comment)
}
