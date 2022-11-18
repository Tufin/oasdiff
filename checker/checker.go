package checker

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/TwiN/go-color"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/tufin/oasdiff/diff"
	"github.com/tufin/oasdiff/load"
)

const (
	ERR  = 0
	WARN = 1
)

const (
	XStabilityLevelExtension = "x-stability-level"
	XExtensibleEnumExtension = "x-extensible-enum"
)

type BackwardCompatibilityError struct {
	Id        string `json:"id,omitempty" yaml:"id,omitempty"`
	Text      string `json:"text,omitempty" yaml:"text,omitempty"`
	Comment   string `json:"comment,omitempty" yaml:"comment,omitempty"`
	Level     int    `json:"level,omitempty" yaml:"level,omitempty"`
	Operation string `json:"operation,omitempty" yaml:"operation,omitempty"`
	Path      string `json:"path,omitempty" yaml:"path,omitempty"`
	Source    string `json:"source,omitempty" yaml:"source,omitempty"`
	ToDo      string `json:"todo,omitempty" yaml:"todo,omitempty"`
}

type BackwardCompatibilityCheck func(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config BackwardCompatibilityCheckConfig) []BackwardCompatibilityError

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
	return fmt.Sprintf("%s at %s, in API %s %s %s [%s]. %s", levelName, r.Source, r.Operation, r.Path, r.Text, r.Id, r.Comment)
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
	return fmt.Sprintf("%s at %s in API %s %s %s [%s]. %s", levelName, r.Source, color.InGreen(r.Operation), color.InGreen(r.Path), r.Text, color.InYellow(r.Id), r.Comment)
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
	comment := ""
	if r.Comment != "" {
		comment = fmt.Sprintf("\n\t\t%s", r.Comment)
	}
	return fmt.Sprintf("%s\t[%s] at %s\t\n\tin API %s %s\n\t\t%s%s", levelName, color.InYellow(r.Id), r.Source, color.InGreen(r.Operation), color.InGreen(r.Path), r.Text, comment)
}

type BackwardCompatibilityCheckConfig struct {
	Checks              []BackwardCompatibilityCheck
	MinSunsetBetaDays   int
	MinSunsetStableDays int
}

func CheckBackwardCompatibility(config BackwardCompatibilityCheckConfig, diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap) []BackwardCompatibilityError {
	result := make([]BackwardCompatibilityError, 0)

	if diffReport == nil {
		return result
	}

	result = removeDraftAndAlphaOperationsDiffs(diffReport, result, operationsSources)

	for _, check := range config.Checks {
		errs := check(diffReport, operationsSources, config)
		result = append(result, errs...)
	}

	return result
}

func removeDraftAndAlphaOperationsDiffs(diffReport *diff.Diff, result []BackwardCompatibilityError, operationsSources *diff.OperationsSourcesMap) []BackwardCompatibilityError {
	if diffReport.PathsDiff == nil {
		return result
	}
	// remove draft and alpha paths diffs modified
	for path, pathDiff := range diffReport.PathsDiff.Modified {
		if pathDiff.OperationsDiff == nil {
			continue
		}
		// remove draft and alpha operations diffs deleted
		iOperation := 0
		for _, operation := range pathDiff.OperationsDiff.Deleted {
			baseStability, err := getStabilityLevel(pathDiff.Base.Operations()[operation].ExtensionProps)
			source := (*operationsSources)[pathDiff.Base.Operations()[operation]]
			if err != nil {
				result = append(result, BackwardCompatibilityError{
					Id:        "parsing-error",
					Level:     ERR,
					Text:      fmt.Sprintf("parsing error %s", err.Error()),
					Operation: operation,
					Path:      path,
					Source:    source,
					ToDo:      "Add to exceptions-list.md",
				})
				continue
			}
			if !(baseStability == "draft" || baseStability == "alpha") {
				pathDiff.OperationsDiff.Deleted[iOperation] = operation
				iOperation++
			}
		}
		pathDiff.OperationsDiff.Deleted = pathDiff.OperationsDiff.Deleted[:iOperation]

		// remove draft and alpha operations diffs modified
		for operation := range pathDiff.OperationsDiff.Modified {
			baseStability, err := getStabilityLevel(pathDiff.Base.Operations()[operation].ExtensionProps)
			if err != nil {
				source := (*operationsSources)[pathDiff.Base.Operations()[operation]]
				result = append(result, BackwardCompatibilityError{
					Id:        "parsing-error",
					Level:     ERR,
					Text:      fmt.Sprintf("parsing error %s", err.Error()),
					Operation: operation,
					Path:      path,
					Source:    source,
					ToDo:      "Add to exceptions-list.md",
				})
				continue
			}
			revisionStability, err := getStabilityLevel(pathDiff.Revision.Operations()[operation].ExtensionProps)
			if err != nil {
				source := (*operationsSources)[pathDiff.Revision.Operations()[operation]]
				result = append(result, BackwardCompatibilityError{
					Id:        "parsing-error",
					Level:     ERR,
					Text:      fmt.Sprintf("parsing error %s", err.Error()),
					Operation: operation,
					Path:      path,
					Source:    source,
					ToDo:      "Add to exceptions-list.md",
				})
				continue
			}
			source := (*operationsSources)[pathDiff.Revision.Operations()[operation]]
			if baseStability == "stable" && revisionStability != "stable" ||
				baseStability == "beta" && revisionStability != "beta" && revisionStability != "stable" ||
				baseStability == "alpha" && revisionStability != "alpha" && revisionStability != "beta" && revisionStability != "stable" ||
				revisionStability == "" && baseStability != "" {
				result = append(result, BackwardCompatibilityError{
					Id:        "api-stability-decreased",
					Level:     ERR,
					Text:      fmt.Sprintf("API stability decreased from '%s' to '%s'", baseStability, revisionStability),
					Operation: operation,
					Path:      path,
					Source:    source,
					ToDo:      "Add to exceptions-list.md",
				})
				continue
			}
			if revisionStability == "draft" || revisionStability == "alpha" {
				delete(pathDiff.OperationsDiff.Modified, operation)
			}
		}
	}
	return result
}

func getStabilityLevel(i openapi3.ExtensionProps) (string, error) {
	if i.Extensions == nil || i.Extensions[XStabilityLevelExtension] == nil {
		return "", nil
	}
	jsonStability, ok := i.Extensions[XStabilityLevelExtension].(json.RawMessage)
	if !ok {
		return "", fmt.Errorf("unparseable x-stability-level")
	}
	var stabilityLevel string
	err := json.Unmarshal(jsonStability, &stabilityLevel)
	if err != nil {
		return "", fmt.Errorf("unparseable x-stability-level")
	}
	return stabilityLevel, nil
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

func (diffBC *BCDiff) AddModifiedParameter(path string, operation string, paramLocation string, paramName string) *diff.ParameterDiff {
	opDiff := diffBC.AddModifiedOperation(path, operation)
	if opDiff.ParametersDiff == nil {
		opDiff.ParametersDiff = &diff.ParametersDiff{}
	}
	if opDiff.ParametersDiff.Modified == nil {
		opDiff.ParametersDiff.Modified = make(diff.ParamDiffByLocation)
	}
	if opDiff.ParametersDiff.Modified[paramLocation] == nil {
		opDiff.ParametersDiff.Modified[paramLocation] = make(diff.ParamDiffs)
	}
	if opDiff.ParametersDiff.Modified[paramLocation][paramName] == nil {
		opDiff.ParametersDiff.Modified[paramLocation][paramName] = &diff.ParameterDiff{}
	}
	return opDiff.ParametersDiff.Modified[paramLocation][paramName]
}

func (diffBC *BCDiff) AddRequestPropertiesDiff(path string, operation string, mediaType string) *diff.SchemasDiff {
	opDiff := diffBC.AddModifiedOperation(path, operation)
	if opDiff.RequestBodyDiff == nil {
		opDiff.RequestBodyDiff = &diff.RequestBodyDiff{}
	}
	if opDiff.RequestBodyDiff.ContentDiff == nil {
		opDiff.RequestBodyDiff.ContentDiff = &diff.ContentDiff{}
	}
	if opDiff.RequestBodyDiff.ContentDiff.MediaTypeModified == nil {
		opDiff.RequestBodyDiff.ContentDiff.MediaTypeModified = make(diff.ModifiedMediaTypes)
	}
	if opDiff.RequestBodyDiff.ContentDiff.MediaTypeModified[mediaType] == nil {
		opDiff.RequestBodyDiff.ContentDiff.MediaTypeModified[mediaType] = &diff.MediaTypeDiff{}
	}
	mediaTypeBCDiff := opDiff.RequestBodyDiff.ContentDiff.MediaTypeModified[mediaType]
	if mediaTypeBCDiff.SchemaDiff == nil {
		mediaTypeBCDiff.SchemaDiff = &diff.SchemaDiff{}
	}
	if mediaTypeBCDiff.SchemaDiff.PropertiesDiff == nil {
		mediaTypeBCDiff.SchemaDiff.PropertiesDiff = &diff.SchemasDiff{}
	}
	return mediaTypeBCDiff.SchemaDiff.PropertiesDiff
}

// LoadOpenAPISpecInfoFromFile loads a LoadOpenAPISpecInfoFromFile from a local file path
func LoadOpenAPISpecInfoFromFile(location string) (*load.OpenAPISpecInfo, error) {
	loader := openapi3.NewLoader()
	s, err := loader.LoadFromFile(location)
	return &load.OpenAPISpecInfo{Spec: s, Url: location}, err
}
