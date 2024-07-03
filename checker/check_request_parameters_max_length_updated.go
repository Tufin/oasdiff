package checker

import (
	"github.com/tufin/oasdiff/diff"
)

const (
	RequestParameterMaxLengthDecreasedId = "request-parameter-max-length-decreased"
	RequestParameterMaxLengthIncreasedId = "request-parameter-max-length-increased"
)

func RequestParameterMaxLengthUpdatedCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config *Config) Changes {
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
			for paramLocation, paramDiffs := range operationItem.ParametersDiff.Modified {
				for paramName, paramDiff := range paramDiffs {
					if paramDiff.SchemaDiff == nil {
						continue
					}
					maxLengthDiff := paramDiff.SchemaDiff.MaxLengthDiff
					if maxLengthDiff == nil {
						continue
					}
					if maxLengthDiff.From == nil ||
						maxLengthDiff.To == nil {
						continue
					}

					id := RequestParameterMaxLengthDecreasedId
					level := ERR
					if !IsDecreasedValue(maxLengthDiff) {
						id = RequestParameterMaxLengthIncreasedId
						level = INFO
					}

					result = append(result, NewApiChange(
						id,
						level,
						[]any{paramLocation, paramName, maxLengthDiff.From, maxLengthDiff.To},
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
