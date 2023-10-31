package checker

import (
	"fmt"
	"strings"

	"github.com/TwiN/go-color"
)

// ComponentChange represnts a change in the Components Section: https://swagger.io/docs/specification/components/
type ComponentChange struct {
	Id      string `json:"id,omitempty" yaml:"id,omitempty"`
	Text    string `json:"text,omitempty" yaml:"text,omitempty"`
	Comment string `json:"comment,omitempty" yaml:"comment,omitempty"`
	Level   Level  `json:"level" yaml:"level"`
	Source  string `json:"source,omitempty" yaml:"source,omitempty"`
}

func (c ComponentChange) MatchIgnore(ignorePath, ignoreLine string) bool {
	return strings.Contains(ignoreLine, strings.ToLower(GetUncolorizedText(c))) &&
		strings.Contains(ignoreLine, "components")
}

func (c ComponentChange) GetId() string {
	return c.Id
}

func (c ComponentChange) GetText() string {
	return c.Text
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

func (c ComponentChange) LocalizedError(l Localizer) string {
	return fmt.Sprintf("%s, %s components %s [%s]. %s", c.Level, l("in"), c.Text, c.Id, c.Comment)
}

func (c ComponentChange) PrettyErrorText(l Localizer) string {
	if IsPipedOutput() {
		return c.LocalizedError(l)
	}

	comment := ""
	if c.Comment != "" {
		comment = fmt.Sprintf("\n\t\t%s", c.Comment)
	}
	return fmt.Sprintf("%s\t[%s] \t\n\t%s components\n\t\t%s%s", c.Level.PrettyString(), color.InYellow(c.Id), l("in"), c.Text, comment)
}

func (c ComponentChange) Error() string {
	return fmt.Sprintf("%s, in components %s [%s]. %s", c.Level, c.Text, c.Id, c.Comment)
}
