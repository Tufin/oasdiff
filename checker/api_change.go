package checker

import (
	"fmt"
	"strings"

	"github.com/TwiN/go-color"
)

// ApiChange represnts a change in the Paths Section of an OpenAPI spec
type ApiChange struct {
	Id          string `json:"id,omitempty" yaml:"id,omitempty"`
	Text        string `json:"text,omitempty" yaml:"text,omitempty"`
	Comment     string `json:"comment,omitempty" yaml:"comment,omitempty"`
	Level       Level  `json:"level" yaml:"level"`
	Operation   string `json:"operation,omitempty" yaml:"operation,omitempty"`
	OperationId string `json:"operationId,omitempty" yaml:"operationId,omitempty"`
	Path        string `json:"path,omitempty" yaml:"path,omitempty"`
	Source      string `json:"source,omitempty" yaml:"source,omitempty"`
}

func (c ApiChange) getUncolorizedText() string {
	uncolorizedText := strings.ReplaceAll(c.Text, color.Bold, "")
	return strings.ReplaceAll(uncolorizedText, color.Reset, "")
}

func (c ApiChange) MatchIgnore(ignorePath, ignoreLine string) bool {
	return ignorePath == strings.ToLower(c.Path) &&
		strings.Contains(ignoreLine, strings.ToLower(c.Operation+" "+c.Path)) &&
		strings.Contains(ignoreLine, strings.ToLower(c.getUncolorizedText()))
}

func (c ApiChange) GetId() string {
	return c.Id
}

func (c ApiChange) GetText() string {
	return c.Text
}

func (c ApiChange) GetComment() string {
	return c.Comment
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

func (c ApiChange) LocalizedError(l Localizer) string {
	return fmt.Sprintf("%s %s %s, %s API %s %s %s [%s]. %s", c.Level, l("at"), c.Source, l("in"), c.Operation, c.Path, c.Text, c.Id, c.Comment)
}

func (c ApiChange) PrettyErrorText(l Localizer) string {
	if IsPipedOutput() {
		return c.LocalizedError(l)
	}

	comment := ""
	if c.Comment != "" {
		comment = fmt.Sprintf("\n\t\t%s", c.Comment)
	}
	return fmt.Sprintf("%s\t[%s] %s %s\t\n\t%s API %s %s\n\t\t%s%s", c.Level.PrettyString(), color.InYellow(c.Id), l("at"), c.Source, l("in"), color.InGreen(c.Operation), color.InGreen(c.Path), c.Text, comment)
}

func (c ApiChange) Error() string {
	return fmt.Sprintf("%s at %s, in API %s %s %s [%s]. %s", c.Level, c.Source, c.Operation, c.Path, c.Text, c.Id, c.Comment)
}
