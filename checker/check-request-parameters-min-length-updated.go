package checker

import (
	"github.com/tufin/oasdiff/diff"
)

const (
	RequestParameterMinLengthIncreasedId = "request-parameter-min-length-increased"
	RequestParameterMinLengthDecreasedId = "request-parameter-min-length-decreased"
)

func RequestParameterMinLengthUpdatedCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config *Config) Changes {
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
					minLengthDiff := paramDiff.SchemaDiff.MinLengthDiff
					if minLengthDiff == nil {
						continue
					}
					if minLengthDiff.From == nil ||
						minLengthDiff.To == nil {
						continue
					}

					id := RequestParameterMinLengthIncreasedId
					level := ERR
					if IsDecreasedValue(minLengthDiff) {
						id = RequestParameterMinLengthDecreasedId
						level = INFO
					}

					source := (*operationsSources)[operationItem.Revision]

					result = append(result, ApiChange{
						Id:          id,
						Level:       level,
						Args:        []any{paramLocation, paramName, minLengthDiff.From, minLengthDiff.To},
						Operation:   operation,
						OperationId: operationItem.Revision.OperationID,
						Path:        path,
						Source:      source,
					})
				}
			}
		}
	}
	return result
}
