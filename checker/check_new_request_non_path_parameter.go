package checker

import (
	"github.com/tufin/oasdiff/diff"
	"github.com/tufin/oasdiff/load"
)

const (
	NewRequiredRequestParameterId = "new-required-request-parameter"
	NewOptionalRequestParameterId = "new-optional-request-parameter"
)

func NewRequestNonPathParameterCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config *Config) Changes {
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
							id := NewRequiredRequestParameterId
							level := ERR
							if !param.Value.Required {
								id = NewOptionalRequestParameterId
								level = INFO
							}
							source := (*operationsSources)[operationItem.Revision]
							result = append(result, ApiChange{
								Id:          id,
								Level:       level,
								Args:        []any{paramLocation, paramName},
								Operation:   operation,
								OperationId: operationItem.Revision.OperationID,
								Path:        path,
								Source:      load.NewSource(source),
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
