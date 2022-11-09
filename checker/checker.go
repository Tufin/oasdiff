package checker

import (
	"fmt"
	"os"

	"github.com/TwiN/go-color"
	"github.com/tufin/oasdiff/diff"
)

const (
	ERR  = 0
	WARN = 1
)

type BackwardCompatibilityError struct {
	Id        string `json:"id,omitempty" yaml:"id,omitempty"`
	Text      string `json:"text,omitempty" yaml:"text,omitempty"`
	Level     int    `json:"level,omitempty" yaml:"level,omitempty"`
	Operation string `json:"operation,omitempty" yaml:"operation,omitempty"`
	Path      string `json:"path,omitempty" yaml:"path,omitempty"`
	Source    string `json:"source,omitempty" yaml:"source,omitempty"`
	ToDo      string `json:"source,omitempty" yaml:"source,omitempty"`
}

type BackwardCompatibilityCheck func(diff *diff.Diff, operationsSources *diff.OperationsSourcesMap) []BackwardCompatibilityError

func (r *BackwardCompatibilityError) Error() string {
	var levelName string
	switch r.Level {
	case ERR:
		levelName = "error"
	case WARN:
		levelName = "warning"
	default:
		levelName = "issue"
	}
	return fmt.Sprintf("%s at %s, in API %s %s %s [%s]", levelName, r.Source, r.Operation, r.Path, r.Text, r.Id)
}

func (r *BackwardCompatibilityError) ColorizedError() string {
	fi, _ := os.Stdout.Stat()

	if (fi.Mode() & os.ModeCharDevice) == 0 {
		return r.Error()
	}

	var levelName string
	switch r.Level {
	case ERR:
		levelName = color.InRed("error")
	case WARN:
		levelName = color.InPurple("warning")
	default:
		levelName = color.InGray("issue")
	}
	return fmt.Sprintf("%s at %s, in API %s %s %s [%s]", levelName, r.Source, color.InGreen(r.Operation), color.InGreen(r.Path), r.Text, color.InYellow(r.Id))
}

func CheckBackwardCompatibility(checks []BackwardCompatibilityCheck, diff *diff.Diff, operationsSources *diff.OperationsSourcesMap) []BackwardCompatibilityError {
	result := make([]BackwardCompatibilityError, 0)

	for _, check := range checks {
		errs := check(diff, operationsSources)
		result = append(result, errs...)
	}

	return result
}
