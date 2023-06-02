package checker

import (
	"fmt"

	"github.com/tufin/oasdiff/diff"
)

func RequestParameterPatternAddedOrChangedCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config BackwardCompatibilityCheckConfig) []BackwardCompatibilityError {
	result := make([]BackwardCompatibilityError, 0)
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
					if paramItem.SchemaDiff == nil {
						continue
					}
					patternDiff := paramItem.SchemaDiff.PatternDiff
					if patternDiff == nil {
						continue
					}
					if patternDiff.To == "" ||
						patternDiff.To == ".*" {
						continue
					}

					source := (*operationsSources)[operationItem.Revision]

					if patternDiff.From == "" {
						result = append(result, BackwardCompatibilityError{
							Id:          "request-parameter-pattern-added",
							Level:       WARN,
							Text:        fmt.Sprintf(config.i18n("request-parameter-pattern-added"), patternDiff.To, ColorizedValue(paramLocation), ColorizedValue(paramName)),
							Comment:     config.i18n("pattern-changed-warn-comment"),
							Operation:   operation,
							OperationId: operationItem.Revision.OperationID,
							Path:        path,
							Source:      source,
						})
					} else {
						result = append(result, BackwardCompatibilityError{
							Id:          "request-parameter-pattern-changed",
							Level:       WARN,
							Text:        fmt.Sprintf(config.i18n("request-parameter-pattern-changed"), ColorizedValue(paramLocation), ColorizedValue(paramName), patternDiff.From, patternDiff.To),
							Comment:     config.i18n("pattern-changed-warn-comment"),
							Operation:   operation,
							OperationId: operationItem.Revision.OperationID,
							Path:        path,
							Source:      source,
						})
					}
				}
			}
		}
	}
	return result
}
