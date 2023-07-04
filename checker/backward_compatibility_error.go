package checker

import (
	"fmt"

	"github.com/TwiN/go-color"
	"github.com/tufin/oasdiff/checker/localizations"
)

type BackwardCompatibilityError struct {
	Id          string `json:"id,omitempty" yaml:"id,omitempty"`
	Text        string `json:"text,omitempty" yaml:"text,omitempty"`
	Comment     string `json:"comment,omitempty" yaml:"comment,omitempty"`
	Level       Level  `json:"level" yaml:"level"`
	Operation   string `json:"operation,omitempty" yaml:"operation,omitempty"`
	OperationId string `json:"operationId,omitempty" yaml:"operationId,omitempty"`
	Path        string `json:"path,omitempty" yaml:"path,omitempty"`
	Source      string `json:"source,omitempty" yaml:"source,omitempty"`
}

func (r *BackwardCompatibilityError) LocalizedError(l localizations.Localizer) string {
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
	return fmt.Sprintf("%s %s %s, %s API %s %s %s [%s]. %s", levelName, l.Get("messages.at"), r.Source, l.Get("messages.in"), r.Operation, r.Path, r.Text, r.Id, r.Comment)
}

func (r *BackwardCompatibilityError) PrettyErrorText(l localizations.Localizer) string {
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
	return fmt.Sprintf("%s\t[%s] %s %s\t\n\t%s API %s %s\n\t\t%s%s", levelName, color.InYellow(r.Id), l.Get("messages.at"), r.Source, l.Get("messages.in"), color.InGreen(r.Operation), color.InGreen(r.Path), r.Text, comment)
}

func (r *BackwardCompatibilityError) Error() string {
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
	return fmt.Sprintf("%s at %s, in API %s %s %s [%s]. %s", levelName, r.Source, r.Operation, r.Path, r.Text, r.Id, r.Comment)
}
