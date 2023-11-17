package checker

import (
	"fmt"
	"strings"

	"github.com/TwiN/go-color"
)

// SecurityChange represents a change in the Security Section (not to be confised with components/securitySchemes)
type SecurityChange struct {
	Id      string `json:"id,omitempty" yaml:"id,omitempty"`
	Text    string `json:"text,omitempty" yaml:"text,omitempty"`
	Comment string `json:"comment,omitempty" yaml:"comment,omitempty"`
	Level   Level  `json:"level" yaml:"level"`
	Source  string `json:"source,omitempty" yaml:"source,omitempty"`

	SourceFile      string `json:"-" yaml:"-"`
	SourceLine      int    `json:"-" yaml:"-"`
	SourceLineEnd   int    `json:"-" yaml:"-"`
	SourceColumn    int    `json:"-" yaml:"-"`
	SourceColumnEnd int    `json:"-" yaml:"-"`
}

func (c SecurityChange) IsBreaking() bool {
	return c.GetLevel().IsBreaking()
}

func (c SecurityChange) MatchIgnore(ignorePath, ignoreLine string) bool {
	return strings.Contains(ignoreLine, strings.ToLower(GetUncolorizedText(c))) &&
		strings.Contains(ignoreLine, "security")
}

func (c SecurityChange) GetId() string {
	return c.Id
}

func (c SecurityChange) GetText() string {
	return c.Text
}

func (c SecurityChange) GetComment() string {
	return c.Comment
}

func (c SecurityChange) GetLevel() Level {
	return c.Level
}

func (r SecurityChange) GetOperation() string {
	return ""
}

func (SecurityChange) GetOperationId() string {
	return ""
}

func (SecurityChange) GetPath() string {
	return ""
}

func (c SecurityChange) GetSourceFile() string {
	return c.SourceFile
}

func (c SecurityChange) GetSourceLine() int {
	return c.SourceLine
}

func (c SecurityChange) GetSourceLineEnd() int {
	return c.SourceLineEnd
}

func (c SecurityChange) GetSourceColumn() int {
	return c.SourceColumn
}

func (c SecurityChange) GetSourceColumnEnd() int {
	return c.SourceColumnEnd
}

func (c SecurityChange) LocalizedError(l Localizer) string {
	return fmt.Sprintf("%s, %s security %s [%s]. %s", c.Level, l("in"), c.Text, c.Id, c.Comment)
}

func (c SecurityChange) PrettyErrorText(l Localizer) string {
	if IsPipedOutput() {
		return c.LocalizedError(l)
	}

	comment := ""
	if c.Comment != "" {
		comment = fmt.Sprintf("\n\t\t%s", c.Comment)
	}
	return fmt.Sprintf("%s\t[%s] \t\n\t%s security\n\t\t%s%s", c.Level.PrettyString(), color.InYellow(c.Id), l("in"), c.Text, comment)
}

func (c SecurityChange) Error() string {
	return fmt.Sprintf("%s, in security %s [%s]. %s", c.Level, c.Text, c.Id, c.Comment)
}
