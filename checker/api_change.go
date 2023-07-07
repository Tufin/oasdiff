package checker

import (
	"fmt"
	"strings"

	"github.com/TwiN/go-color"
	"github.com/tufin/oasdiff/checker/localizations"
)

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

func (r ApiChange) getUncolorizedText() string {
	uncolorizedText := strings.ReplaceAll(r.Text, color.Bold, "")
	return strings.ReplaceAll(uncolorizedText, color.Reset, "")
}

func (r ApiChange) MatchIgnore(ignorePath, ignoreLine string) bool {
	return ignorePath == strings.ToLower(r.Path) &&
		strings.Contains(ignoreLine, strings.ToLower(r.Operation+" "+r.Path)) &&
		strings.Contains(ignoreLine, strings.ToLower(r.getUncolorizedText()))
}

func (r ApiChange) GetId() string {
	return r.Id
}

func (r ApiChange) GetText() string {
	return r.Text
}

func (r ApiChange) GetComment() string {
	return r.Comment
}

func (r ApiChange) GetLevel() Level {
	return r.Level
}

func (r ApiChange) GetOperation() string {
	return r.Operation
}

func (r ApiChange) GetOperationId() string {
	return r.OperationId
}

func (r ApiChange) GetPath() string {
	return r.Path
}

func (r ApiChange) LocalizedError(l localizations.Localizer) string {
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

func (r ApiChange) PrettyErrorText(l localizations.Localizer) string {
	if IsPipedOutput() {
		return r.LocalizedError(l)
	}

	levelName := PrettyLevelText(r.Level)
	comment := ""
	if r.Comment != "" {
		comment = fmt.Sprintf("\n\t\t%s", r.Comment)
	}
	return fmt.Sprintf("%s\t[%s] %s %s\t\n\t%s API %s %s\n\t\t%s%s", levelName, color.InYellow(r.Id), l.Get("messages.at"), r.Source, l.Get("messages.in"), color.InGreen(r.Operation), color.InGreen(r.Path), r.Text, comment)
}

func (r ApiChange) Error() string {
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
