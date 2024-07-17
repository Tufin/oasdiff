package checker

import (
	"github.com/tufin/oasdiff/diff"
)

const (
	NewRequestPathParameterId = "new-request-path-parameter"
)

func NewRequestPathParameterCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config *Config) Changes {
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
				if paramLocation != "path" {
					continue
				}

				for _, paramName := range paramItems {
					result = append(result, NewApiChange(
						NewRequestPathParameterId,
						config,
						[]any{paramName},
						"",
						operationsSources,
						operationItem.Revision,
						operation,
						path,
					))
				}
			}
		}
	}
	return result
}
