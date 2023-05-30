package checker

import (
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/tufin/oasdiff/diff"
)

func APIAddedCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config BackwardCompatibilityCheckConfig) []BackwardCompatibilityError {
	result := make([]BackwardCompatibilityError, 0)
	if diffReport.PathsDiff == nil {
		return result
	}

	appendErr := func(path, opName string, opConfig *openapi3.Operation) {
		result = append(result, BackwardCompatibilityError{
			Id:          "api-path-added",
			Level:       INFO,
			Text:        config.i18n("api-path-added"),
			Operation:   opName,
			OperationId: opConfig.OperationID,
			Path:        path,
			Source:      (*operationsSources)[opConfig],
		})
	}

	for _, path := range diffReport.PathsDiff.Added {
		for opName, op := range diffReport.PathsDiff.Revision[path].Operations() {
			appendErr(path, opName, op)
		}
	}

	for path := range diffReport.PathsDiff.Modified {
		for opName, op := range diffReport.PathsDiff.Revision[path].Operations() {
			if _, foundInBase := diffReport.PathsDiff.Base[path].Operations()[opName]; foundInBase {
				continue
			}
			appendErr(path, opName, op)
		}
	}

	return result
}
