package checker

import (
	"fmt"

	"github.com/tufin/oasdiff/diff"
)

func RequestParameterMinIncreasedCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config Config) Changes {
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

					if !isIncreasedValue(minDiff) {
						continue
					}

					source := (*operationsSources)[operationItem.Revision]

					result = append(result, ApiChange{
						Id:          "request-parameter-min-increased",
						Level:       ERR,
						Text:        fmt.Sprintf(config.i18n("request-parameter-min-increased"), colorizedValue(paramLocation), colorizedValue(paramName), colorizedValue(minDiff.From), colorizedValue(minDiff.To)),
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
