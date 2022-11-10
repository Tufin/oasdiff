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

type BackwardCompatibilityCheck func(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, diffBC *BCDiff) []BackwardCompatibilityError

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
	if IsPipedOutput() {
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
	return fmt.Sprintf("%s at %s in API %s %s %s [%s]", levelName, r.Source, color.InGreen(r.Operation), color.InGreen(r.Path), r.Text, color.InYellow(r.Id))
}

var pipedOutput *bool
func IsPipedOutput() bool {
	if pipedOutput != nil {
		return *pipedOutput
	}
	fi, _ := os.Stdout.Stat()
	a := (fi.Mode() & os.ModeCharDevice) == 0
	pipedOutput = &a
	return *pipedOutput
}

func (r *BackwardCompatibilityError) PrettyError() string {
	if IsPipedOutput() {
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
	return fmt.Sprintf("%s\t[%s] at %s\t\n\tin API %s %s\n\t\t%s", levelName, color.InYellow(r.Id), r.Source, color.InGreen(r.Operation), color.InGreen(r.Path), r.Text)
}

func CheckBackwardCompatibility(checks []BackwardCompatibilityCheck, diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap) ([]BackwardCompatibilityError, diff.Diff) {
	result := make([]BackwardCompatibilityError, 0)
	diffBC := BCDiff{}

	for _, check := range checks {
		errs := check(diffReport, operationsSources, &diffBC)
		result = append(result, errs...)
	}

	return result, diffBC.Diff
}

type BCDiff struct {
	diff.Diff
}

func (d *BCDiff) AddModifiedOperation(path string, operation string) *diff.MethodDiff {
	pathDiff := d.AddModifiedPath(path)
	if pathDiff.OperationsDiff == nil {
		pathDiff.OperationsDiff = &diff.OperationsDiff{}
	}
	if pathDiff.OperationsDiff.Modified == nil {
		pathDiff.OperationsDiff.Modified = make(diff.ModifiedOperations)
	}
	if pathDiff.OperationsDiff.Modified[operation] == nil {
		pathDiff.OperationsDiff.Modified[operation] = &diff.MethodDiff{}
	}
	return pathDiff.OperationsDiff.Modified[operation]
}

func (d *BCDiff) AddModifiedPath(path string) *diff.PathDiff {
	if d.PathsDiff == nil {
		d.PathsDiff = &diff.PathsDiff{}
	}
	if d.PathsDiff.Modified == nil {
		d.PathsDiff.Modified = make(diff.ModifiedPaths)
	}
	if d.PathsDiff.Modified[path] == nil {
		d.PathsDiff.Modified[path] = &diff.PathDiff{}
	}
	return d.PathsDiff.Modified[path]
}

func ColorizedValue(arg string) string {
	if IsPipedOutput() {
		return arg
	}
	return color.InBold(arg)
}
