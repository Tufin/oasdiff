package checker

import (
	"github.com/tufin/oasdiff/diff"
)

const (
	NewRequiredRequestParameterId = "new-required-request-parameter"
	NewOptionalRequestParameterId = "new-optional-request-parameter"
)

func NewRequestNonPathParameterCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config *Config) Changes {
	changeGetter := newApiChangeGetter(config, operationsSources)
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
							result = append(result, changeGetter(
								id,
								level,
								[]any{paramLocation, paramName},
								"",
								operation,
								operationItem.Revision,
								path,
								operationItem.Revision,
							))

							break
						}
					}
				}
			}
		}
	}
	return result
}
