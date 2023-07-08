package checker

import (
	"fmt"

	"github.com/TwiN/go-color"
	"github.com/tufin/oasdiff/checker/localizations"
)

// SecurityChange represnts a change in the Security Section (not to be confised with components/securitySchemes)
type SecurityChange struct {
	Id      string `json:"id,omitempty" yaml:"id,omitempty"`
	Text    string `json:"text,omitempty" yaml:"text,omitempty"`
	Comment string `json:"comment,omitempty" yaml:"comment,omitempty"`
	Level   Level  `json:"level" yaml:"level"`
	Source  string `json:"source,omitempty" yaml:"source,omitempty"`
}

func (SecurityChange) MatchIgnore(ignorePath, ignoreLine string) bool {
	return false
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

func (c SecurityChange) LocalizedError(l localizations.Localizer) string {
	return fmt.Sprintf("%s, %s security %s [%s]. %s", c.Level, l.Get("messages.in"), c.Text, c.Id, c.Comment)
}

func (c SecurityChange) PrettyErrorText(l localizations.Localizer) string {
	if IsPipedOutput() {
		return c.LocalizedError(l)
	}

	comment := ""
	if c.Comment != "" {
		comment = fmt.Sprintf("\n\t\t%s", c.Comment)
	}
	return fmt.Sprintf("%s\t[%s] \t\n\t%s security\n\t\t%s%s", c.Level.PrettyString(), color.InYellow(c.Id), l.Get("messages.in"), c.Text, comment)
}

func (c SecurityChange) Error() string {
	return fmt.Sprintf("%s, in security %s [%s]. %s", c.Level.String(), c.Text, c.Id, c.Comment)
}
