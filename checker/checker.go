package checker

//go:generate go-localize -input localizations_src -output localizations
import (
	"encoding/json"
	"fmt"
	"os"
	"sort"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/tufin/oasdiff/diff"
)

type BackwardCompatibilityCheck func(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config Config) Changes

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

func CheckBackwardCompatibility(config Config, diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap) Changes {
	return CheckBackwardCompatibilityUntilLevel(config, diffReport, operationsSources, WARN)
}

func CheckBackwardCompatibilityUntilLevel(config Config, diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, level Level) Changes {
	result := make(Changes, 0)

	if diffReport == nil {
		return result
	}

	result = removeDraftAndAlphaOperationsDiffs(diffReport, result, operationsSources)

	for _, check := range config.Checks {
		errs := check(diffReport, operationsSources, config)
		result = append(result, errs...)
	}

	filteredResult := make(Changes, 0)
	for _, change := range result {
		if change.GetLevel() >= level {
			filteredResult = append(filteredResult, change)
		}
	}

	sort.Sort(filteredResult)
	return filteredResult
}

func removeDraftAndAlphaOperationsDiffs(diffReport *diff.Diff, result Changes, operationsSources *diff.OperationsSourcesMap) Changes {
	if diffReport.PathsDiff == nil {
		return result
	}
	// remove draft and alpha paths diffs delete
	iPath := 0
	for _, path := range diffReport.PathsDiff.Deleted {
		ignore := true
		pathDiff := diffReport.PathsDiff
		for operation, operationItem := range pathDiff.Base[path].Operations() {
			baseStability, err := getStabilityLevel(pathDiff.Base[path].Operations()[operation].Extensions)
			source := (*operationsSources)[pathDiff.Base[path].Operations()[operation]]
			if err != nil {
				result = newParsingError(result, err, operation, operationItem, path, source)
				continue
			}
			if !(baseStability == "draft" || baseStability == "alpha") {
				ignore = false
				break
			}
		}
		if !ignore {
			diffReport.PathsDiff.Deleted[iPath] = path
			iPath++
		}
	}
	diffReport.PathsDiff.Deleted = diffReport.PathsDiff.Deleted[:iPath]

	// remove draft and alpha paths diffs modified
	for path, pathDiff := range diffReport.PathsDiff.Modified {
		if pathDiff.OperationsDiff == nil {
			continue
		}
		// remove draft and alpha operations diffs deleted
		iOperation := 0
		for _, operation := range pathDiff.OperationsDiff.Deleted {
			operationItem := pathDiff.Base.Operations()[operation]
			baseStability, err := getStabilityLevel(operationItem.Extensions)
			source := (*operationsSources)[pathDiff.Base.Operations()[operation]]
			if err != nil {
				result = newParsingError(result, err, operation, operationItem, path, source)
				continue
			}
			if !(baseStability == "draft" || baseStability == "alpha") {
				pathDiff.OperationsDiff.Deleted[iOperation] = operation
				iOperation++
			}
		}
		pathDiff.OperationsDiff.Deleted = pathDiff.OperationsDiff.Deleted[:iOperation]

		// remove draft and alpha operations diffs modified
		for operation, operationItem := range pathDiff.OperationsDiff.Modified {
			baseStability, err := getStabilityLevel(pathDiff.Base.Operations()[operation].Extensions)
			if err != nil {
				source := (*operationsSources)[pathDiff.Base.Operations()[operation]]
				result = append(result, ApiChange{
					Id:          "parsing-error",
					Level:       ERR,
					Text:        fmt.Sprintf("parsing error %s", err.Error()),
					Operation:   operation,
					OperationId: operationItem.Revision.OperationID,
					Path:        path,
					Source:      source,
				})
				continue
			}
			revisionStability, err := getStabilityLevel(pathDiff.Revision.Operations()[operation].Extensions)
			if err != nil {
				source := (*operationsSources)[pathDiff.Revision.Operations()[operation]]
				result = append(result, ApiChange{
					Id:          "parsing-error",
					Level:       ERR,
					Text:        fmt.Sprintf("parsing error %s", err.Error()),
					Operation:   operation,
					OperationId: operationItem.Revision.OperationID,
					Path:        path,
					Source:      source,
				})
				continue
			}
			source := (*operationsSources)[pathDiff.Revision.Operations()[operation]]
			if baseStability == "stable" && revisionStability != "stable" ||
				baseStability == "beta" && revisionStability != "beta" && revisionStability != "stable" ||
				baseStability == "alpha" && revisionStability != "alpha" && revisionStability != "beta" && revisionStability != "stable" ||
				revisionStability == "" && baseStability != "" {
				result = append(result, ApiChange{
					Id:          "api-stability-decreased",
					Level:       ERR,
					Text:        fmt.Sprintf("stability level decreased from '%s' to '%s'", baseStability, revisionStability),
					Operation:   operation,
					OperationId: operationItem.Revision.OperationID,
					Path:        path,
					Source:      source,
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

func newParsingError(result Changes,
	err error,
	operation string,
	operationItem *openapi3.Operation,
	path string,
	source string) Changes {
	result = append(result, ApiChange{
		Id:          "parsing-error",
		Level:       ERR,
		Text:        fmt.Sprintf("parsing error %s", err.Error()),
		Operation:   operation,
		OperationId: operationItem.OperationID,
		Path:        path,
		Source:      source,
	})
	return result
}

func getStabilityLevel(i map[string]interface{}) (string, error) {
	if i == nil || i[diff.XStabilityLevelExtension] == nil {
		return "", nil
	}
	var stabilityLevel string

	stabilityLevel, ok := i[diff.XStabilityLevelExtension].(string)
	if !ok {
		jsonStability, ok := i[diff.XStabilityLevelExtension].(json.RawMessage)
		if !ok {
			return "", fmt.Errorf("unparseable x-stability-level")
		}
		err := json.Unmarshal(jsonStability, &stabilityLevel)
		if err != nil {
			return "", fmt.Errorf("unparseable x-stability-level")
		}
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
