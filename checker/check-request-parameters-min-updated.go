package checker

import (
	"github.com/tufin/oasdiff/diff"
	"github.com/tufin/oasdiff/load"
)

const (
	RequestParameterMinIncreasedId = "request-parameter-min-increased"
	RequestParameterMinDecreasedId = "request-parameter-min-decreased"
)

func RequestParameterMinUpdatedCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config *Config) Changes {
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
					minDiff := paramDiff.SchemaDiff.MinDiff
					if minDiff == nil {
						continue
					}
					if minDiff.From == nil ||
						minDiff.To == nil {
						continue
					}

					source := (*operationsSources)[operationItem.Revision]

					id := RequestParameterMinIncreasedId
					level := ERR
					if !IsIncreasedValue(minDiff) {
						id = RequestParameterMinDecreasedId
						level = INFO
					}

					result = append(result, ApiChange{
						Id:          id,
						Level:       level,
						Args:        []any{paramLocation, paramName, minDiff.From, minDiff.To},
						Operation:   operation,
						OperationId: operationItem.Revision.OperationID,
						Path:        path,
						Source:      load.NewSource(source),
					})
				}
			}
		}
	}
	return result
}
