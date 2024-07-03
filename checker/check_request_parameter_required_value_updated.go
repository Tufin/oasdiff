package checker

import (
	"github.com/tufin/oasdiff/diff"
)

const (
	RequestParameterBecomeRequiredId = "request-parameter-became-required"
	RequestParameterBecomeOptionalId = "request-parameter-became-optional"
)

func RequestParameterRequiredValueUpdatedCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config *Config) Changes {
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
			if operationItem.ParametersDiff.Modified == nil {
				continue
			}
			for paramLocation, paramItems := range operationItem.ParametersDiff.Modified {
				for paramName, paramItem := range paramItems {
					requiredDiff := paramItem.RequiredDiff
					if requiredDiff == nil {
						continue
					}

					id := RequestParameterBecomeRequiredId
					level := ERR

					if requiredDiff.To != true {
						id = RequestParameterBecomeOptionalId
						level = INFO
					}

					result = append(result, NewApiChange(
						id,
						level,
						[]any{paramLocation, paramName},
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
