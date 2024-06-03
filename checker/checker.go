package checker

//go:generate go-localize -input localizations_src -output localizations
import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/tufin/oasdiff/diff"
	"github.com/tufin/oasdiff/load"
)

const (
	APIStabilityDecreasedId = "api-stability-decreased"
)

func CheckBackwardCompatibility(config *Config, diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap) Changes {
	return CheckBackwardCompatibilityUntilLevel(config, diffReport, operationsSources, WARN)
}

func CheckBackwardCompatibilityUntilLevel(config *Config, diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, level Level) Changes {
	result := make(Changes, 0)

	if diffReport == nil {
		return result
	}

	result = removeDraftAndAlphaOperationsDiffs(config, diffReport, result, operationsSources)

	for _, check := range config.Checks {
		if check == nil {
			continue
		}
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

func removeDraftAndAlphaOperationsDiffs(config *Config, diffReport *diff.Diff, result Changes, operationsSources *diff.OperationsSourcesMap) Changes {
	if diffReport.PathsDiff == nil {
		return result
	}
	// remove draft and alpha paths diffs delete
	iPath := 0
	for _, path := range diffReport.PathsDiff.Deleted {
		ignore := true
		pathDiff := diffReport.PathsDiff
		for operation, operationItem := range pathDiff.Base.Value(path).Operations() {
			baseStability, err := getStabilityLevel(pathDiff.Base.Value(path).Operations()[operation].Extensions)
			source := (*operationsSources)[pathDiff.Base.Value(path).Operations()[operation]]
			if err != nil {
				result = newParsingError(config, result, err, operation, operationItem, path, source)
				continue
			}
			if !(baseStability == STABILITY_DRAFT || baseStability == STABILITY_ALPHA) {
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
				result = newParsingError(config, result, err, operation, operationItem, path, source)
				continue
			}
			if !(baseStability == STABILITY_DRAFT || baseStability == STABILITY_ALPHA) {
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
					Id:          ParseErrorId,
					Args:        []any{err.Error()},
					Level:       ERR,
					Operation:   operation,
					OperationId: operationItem.Revision.OperationID,
					Path:        path,
					Source:      load.NewSource(source),
				})
				continue
			}
			revisionStability, err := getStabilityLevel(pathDiff.Revision.Operations()[operation].Extensions)
			if err != nil {
				source := (*operationsSources)[pathDiff.Revision.Operations()[operation]]
				result = append(result, ApiChange{
					Id:          ParseErrorId,
					Args:        []any{err.Error()},
					Level:       ERR,
					Operation:   operation,
					OperationId: operationItem.Revision.OperationID,
					Path:        path,
					Source:      load.NewSource(source),
				})
				continue
			}
			source := (*operationsSources)[pathDiff.Revision.Operations()[operation]]
			if baseStability == STABILITY_STABLE && revisionStability != STABILITY_STABLE ||
				baseStability == STABILITY_BETA && revisionStability != STABILITY_BETA && revisionStability != STABILITY_STABLE ||
				baseStability == STABILITY_ALPHA && revisionStability != STABILITY_ALPHA && revisionStability != STABILITY_BETA && revisionStability != STABILITY_STABLE ||
				revisionStability == "" && baseStability != "" {
				result = append(result, ApiChange{
					Id:          APIStabilityDecreasedId,
					Args:        []any{baseStability, revisionStability},
					Level:       ERR,
					Operation:   operation,
					OperationId: operationItem.Revision.OperationID,
					Path:        path,
					Source:      load.NewSource(source),
				})
				continue
			}
			if revisionStability == STABILITY_DRAFT || revisionStability == STABILITY_ALPHA {
				delete(pathDiff.OperationsDiff.Modified, operation)
			}
		}
	}
	return result
}

func newParsingError(config *Config,
	result Changes,
	err error,
	operation string,
	operationItem *openapi3.Operation,
	path string,
	source string) Changes {
	result = append(result, ApiChange{
		Id:          ParseErrorId,
		Args:        []any{err.Error()},
		Level:       ERR,
		Operation:   operation,
		OperationId: operationItem.OperationID,
		Path:        path,
		Source:      load.NewSource(source),
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

	if stabilityLevel != STABILITY_DRAFT &&
		stabilityLevel != STABILITY_ALPHA &&
		stabilityLevel != STABILITY_BETA &&
		stabilityLevel != STABILITY_STABLE {
		return "", fmt.Errorf("invalid x-stability-level: %q", stabilityLevel)
	}

	return stabilityLevel, nil
}
