package checker

import (
	"fmt"
	"strings"

	"github.com/TwiN/go-color"
)

// ComponentChange represnts a change in the Components Section: https://swagger.io/docs/specification/components/
type ComponentChange struct {
	Id        string `json:"id,omitempty" yaml:"id,omitempty"`
	Text      string `json:"text,omitempty" yaml:"text,omitempty"`
	Args      []any  `json:"-" yaml:"-"`
	Comment   string `json:"comment,omitempty" yaml:"comment,omitempty"`
	Level     Level  `json:"level" yaml:"level"`
	Source    string `json:"source,omitempty" yaml:"source,omitempty"`
	Component string `json:"component,omitempty" yaml:"component,omitempty"`

	SourceFile      string `json:"-" yaml:"-"`
	SourceLine      int    `json:"-" yaml:"-"`
	SourceLineEnd   int    `json:"-" yaml:"-"`
	SourceColumn    int    `json:"-" yaml:"-"`
	SourceColumnEnd int    `json:"-" yaml:"-"`
}

func (c ComponentChange) IsBreaking() bool {
	return c.GetLevel().IsBreaking()
}

func (c ComponentChange) MatchIgnore(ignorePath, ignoreLine string, l Localizer) bool {
	return strings.Contains(ignoreLine, strings.ToLower(c.GetUncolorizedText(l))) &&
		strings.Contains(ignoreLine, "components")
}

func (c ComponentChange) GetId() string {
	return c.Id
}

func (c ComponentChange) GetText(l Localizer) string {
	return l(c.Id, ColorizedValues(c.Args)...)
}

func (c ComponentChange) GetArgs() []any {
	return c.Args
}

func (c ComponentChange) GetUncolorizedText(l Localizer) string {
	return l(c.Id, QuotedValues(c.Args)...)
}

func (c ComponentChange) GetComment() string {
	return c.Comment
}

func (c ComponentChange) GetLevel() Level {
	return c.Level
}

func (ComponentChange) GetOperation() string {
	return ""
}

func (ComponentChange) GetOperationId() string {
	return ""
}

func (ComponentChange) GetPath() string {
	return ""
}

func (c ComponentChange) GetSource() string {
	return c.Source
}

func (c ComponentChange) GetSourceFile() string {
	return c.SourceFile
}

func (c ComponentChange) GetSourceLine() int {
	return c.SourceLine
}

func (c ComponentChange) GetSourceLineEnd() int {
	return c.SourceLineEnd
}

func (c ComponentChange) GetSourceColumn() int {
	return c.SourceColumn
}

func (c ComponentChange) GetSourceColumnEnd() int {
	return c.SourceColumnEnd
}

func (c ComponentChange) LocalizedError(l Localizer) string {
	return fmt.Sprintf("%s, %s components/%s %s [%s]. %s", c.Level, l("in"), c.Component, c.Text, c.Id, c.Comment)
}

func (c ComponentChange) PrettyErrorText(l Localizer) string {
	if IsPipedOutput() {
		return c.LocalizedError(l)
	}

	comment := ""
	if c.Comment != "" {
		comment = fmt.Sprintf("\n\t\t%s", c.Comment)
	}
	return fmt.Sprintf("%s\t[%s] \t\n\t%s components/%s\n\t\t%s%s", c.Level.PrettyString(), color.InYellow(c.Id), l("in"), c.Component, c.Text, comment)
}
