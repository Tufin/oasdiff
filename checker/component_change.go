package checker

import (
	"fmt"

	"github.com/TwiN/go-color"
	"github.com/tufin/oasdiff/checker/localizations"
)

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

func PrettyLevelText(level Level) string {
	var levelName string
	switch level {
	case ERR:
		levelName = color.InRed("error")
	case WARN:
		levelName = color.InPurple("warning")
	case INFO:
		levelName = color.InCyan("info")
	default:
		levelName = color.InGray("issue")
	}

	return levelName
}

func (r ComponentChange) PrettyErrorText(l localizations.Localizer) string {
	if IsPipedOutput() {
		return r.LocalizedError(l)
	}

	levelName := PrettyLevelText(r.Level)
	comment := ""
	if r.Comment != "" {
		comment = fmt.Sprintf("\n\t\t%s", r.Comment)
	}
	return fmt.Sprintf("%s\t[%s] %s %s\t\n\t%s components\n\t\t%s%s", levelName, color.InYellow(r.Id), l.Get("messages.at"), r.Source, l.Get("messages.in"), r.Text, comment)
}

func (r ComponentChange) Error() string {
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
