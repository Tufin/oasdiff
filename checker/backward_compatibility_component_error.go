package checker

import (
	"fmt"

	"github.com/TwiN/go-color"
	"github.com/tufin/oasdiff/checker/localizations"
)

type BackwardCompatibilityComponentError struct {
	Id      string `json:"id,omitempty" yaml:"id,omitempty"`
	Text    string `json:"text,omitempty" yaml:"text,omitempty"`
	Comment string `json:"comment,omitempty" yaml:"comment,omitempty"`
	Level   Level  `json:"level" yaml:"level"`
	Source  string `json:"source,omitempty" yaml:"source,omitempty"`
}

func (r BackwardCompatibilityComponentError) MatchIgnore(ignorePath, ignoreLine string) bool {
	return false
}

func (r BackwardCompatibilityComponentError) GetId() string {
	return r.Id
}

func (r BackwardCompatibilityComponentError) GetText() string {
	return r.Text
}

func (r BackwardCompatibilityComponentError) GetComment() string {
	return r.Comment
}

func (r BackwardCompatibilityComponentError) GetLevel() Level {
	return r.Level
}

func (r BackwardCompatibilityComponentError) GetOperation() string {
	return ""
}

func (r BackwardCompatibilityComponentError) GetOperationId() string {
	return ""
}

func (r BackwardCompatibilityComponentError) GetPath() string {
	return ""
}

func (r BackwardCompatibilityComponentError) LocalizedError(l localizations.Localizer) string {
	var levelName string
	switch r.Level {
	case ERR:
		levelName = "error"
	case WARN:
		levelName = "warning"
	case INFO:
		levelName = "info"
	default:
		levelName = "issue"
	}
	return fmt.Sprintf("%s %s %s, %s components %s [%s]. %s", levelName, l.Get("messages.at"), r.Source, l.Get("messages.in"), r.Text, r.Id, r.Comment)
}

func (r BackwardCompatibilityComponentError) PrettyErrorText(l localizations.Localizer) string {
	if IsPipedOutput() {
		return r.LocalizedError(l)
	}

	var levelName string
	switch r.Level {
	case ERR:
		levelName = color.InRed("error")
	case WARN:
		levelName = color.InPurple("warning")
	case INFO:
		levelName = color.InCyan("info")
	default:
		levelName = color.InGray("issue")
	}
	comment := ""
	if r.Comment != "" {
		comment = fmt.Sprintf("\n\t\t%s", r.Comment)
	}
	return fmt.Sprintf("%s\t[%s] %s %s\t\n\t%s components\n\t\t%s%s", levelName, color.InYellow(r.Id), l.Get("messages.at"), r.Source, l.Get("messages.in"), r.Text, comment)
}

func (r BackwardCompatibilityComponentError) Error() string {
	var levelName string
	switch r.Level {
	case ERR:
		levelName = "error"
	case WARN:
		levelName = "warning"
	case INFO:
		levelName = "info"
	default:
		levelName = "issue"
	}
	return fmt.Sprintf("%s at %s, in components %s [%s]. %s", levelName, r.Source, r.Text, r.Id, r.Comment)
}
