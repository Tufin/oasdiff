package checker

import (
	"fmt"

	"github.com/TwiN/go-color"
	"github.com/tufin/oasdiff/checker/localizations"
)

// ComponentChange represnts a change in the Components Section: https://swagger.io/docs/specification/components/
type ComponentChange struct {
	Id      string `json:"id,omitempty" yaml:"id,omitempty"`
	Text    string `json:"text,omitempty" yaml:"text,omitempty"`
	Comment string `json:"comment,omitempty" yaml:"comment,omitempty"`
	Level   Level  `json:"level" yaml:"level"`
	Source  string `json:"source,omitempty" yaml:"source,omitempty"`
}

func (r ComponentChange) MatchIgnore(ignorePath, ignoreLine string) bool {
	return false
}

func (r ComponentChange) GetId() string {
	return r.Id
}

func (r ComponentChange) GetText() string {
	return r.Text
}

func (r ComponentChange) GetComment() string {
	return r.Comment
}

func (r ComponentChange) GetLevel() Level {
	return r.Level
}

func (r ComponentChange) GetOperation() string {
	return ""
}

func (r ComponentChange) GetOperationId() string {
	return ""
}

func (r ComponentChange) GetPath() string {
	return ""
}

func (r ComponentChange) LocalizedError(l localizations.Localizer) string {
	return fmt.Sprintf("%s %s %s, %s components %s [%s]. %s", r.Level, l.Get("messages.at"), r.Source, l.Get("messages.in"), r.Text, r.Id, r.Comment)
}

func (r ComponentChange) PrettyErrorText(l localizations.Localizer) string {
	if IsPipedOutput() {
		return r.LocalizedError(l)
	}

	comment := ""
	if r.Comment != "" {
		comment = fmt.Sprintf("\n\t\t%s", r.Comment)
	}
	return fmt.Sprintf("%s\t[%s] %s %s\t\n\t%s components\n\t\t%s%s", r.Level.PrettyString(), color.InYellow(r.Id), l.Get("messages.at"), r.Source, l.Get("messages.in"), r.Text, comment)
}

func (r ComponentChange) Error() string {
	return fmt.Sprintf("%s at %s, in components %s [%s]. %s", r.Level.String(), r.Source, r.Text, r.Id, r.Comment)
}
