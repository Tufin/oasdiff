package checker

import (
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/tufin/oasdiff/diff"
)

func APIAddedCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config Config) Changes {
	result := make(Changes, 0)
	if diffReport.PathsDiff == nil {
		return result
	}

	appendErr := func(path, opName string, opConfig *openapi3.Operation) {
		result = append(result, ApiChange{
			Id:          "endpoint-added",
			Level:       INFO,
			Text:        config.Localize("endpoint-added"),
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

	for path, pathDiff := range diffReport.PathsDiff.Modified {
		for opName, op := range pathDiff.Revision.Operations() {
			if _, foundInBase := pathDiff.Base.Operations()[opName]; foundInBase {
				continue
			}
			appendErr(path, opName, op)
		}
	}

	return result
}
