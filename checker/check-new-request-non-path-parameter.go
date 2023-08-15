package checker

import (
	"github.com/tufin/oasdiff/diff"
)

func NewRequestNonPathParameterCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config Config) Changes {
	result := make(Changes, 0)
	if diffReport.PathsDiff == nil {
		return result
	}
	for path, pathItem := range diffReport.PathsDiff.Modified {
		if pathItem.OperationsDiff == nil {
			continue
		}
		for operation, operationItem := range pathItem.OperationsDiff.Modified {
			if operationItem.ParametersDiff == nil {
				continue
			}
			for paramLocation, paramItems := range operationItem.ParametersDiff.Added {
				if paramLocation == "path" {
					// it is processed in the separate check NewRequestPathParameterCheck
					continue
				}

				for _, paramName := range paramItems {
					for _, param := range operationItem.Revision.Parameters {
						if param.Value.Name == paramName {
							id := "new-required-request-parameter"
							level := ERR
							if !param.Value.Required {
								id = "new-optional-request-parameter"
								level = INFO
							}
							source := (*operationsSources)[operationItem.Revision]
							result = append(result, ApiChange{
								Id:          id,
								Level:       level,
								Text:        config.Localize(id, ColorizedValue(paramLocation), ColorizedValue(paramName)),
								Operation:   operation,
								OperationId: operationItem.Revision.OperationID,
								Path:        path,
								Source:      source,
							})

							break
						}
					}
				}
			}
		}
	}
	return result
}
