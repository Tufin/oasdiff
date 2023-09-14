package checker

import (
	"github.com/tufin/oasdiff/diff"
)

func NewRequestNonPathDefaultParameterCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config Config) Changes {
	result := make(Changes, 0)
	if diffReport.PathsDiff == nil || len(diffReport.PathsDiff.Modified) == 0 {
		return result
	}
	for path, pathItem := range diffReport.PathsDiff.Modified {
		if pathItem.ParametersDiff == nil || pathItem.Revision == nil {
			continue
		}

		for paramLoc := range pathItem.ParametersDiff.Added {
			if paramLoc == "path" {
				continue
			}

			for _, param := range pathItem.Revision.Parameters {

				id := "new-required-request-default-parameter-to-existing-path"
				level := ERR
				if !param.Value.Required {
					id = "new-optional-request-default-parameter-to-existing-path"
					level = INFO
				}

				for operation, operationItem := range pathItem.Revision.Operations() {

					// TODO: if base operation had this param individually (not through the path) - continue

					result = append(result, ApiChange{
						Id:          id,
						Level:       level,
						Text:        config.Localize(id, ColorizedValue(paramLoc), ColorizedValue(param.Value.Name)),
						Operation:   operation,
						OperationId: operationItem.OperationID,
						Path:        path,
						Source:      (*operationsSources)[operationItem],
					})
				}
			}
		}
	}
	return result
}
