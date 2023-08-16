package checker

import (
	"github.com/tufin/oasdiff/diff"
)

func RequestParameterPatternAddedOrChangedCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config Config) Changes {
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
					if paramItem.SchemaDiff == nil {
						continue
					}
					patternDiff := paramItem.SchemaDiff.PatternDiff
					if patternDiff == nil {
						continue
					}
					source := (*operationsSources)[operationItem.Revision]

					if patternDiff.From == "" {
						result = append(result, ApiChange{
							Id:          "request-parameter-pattern-added",
							Level:       WARN,
							Text:        config.Localize("request-parameter-pattern-added", patternDiff.To, ColorizedValue(paramLocation), ColorizedValue(paramName)),
							Comment:     config.Localize("pattern-changed-warn-comment"),
							Operation:   operation,
							OperationId: operationItem.Revision.OperationID,
							Path:        path,
							Source:      source,
						})
					} else if patternDiff.To == "" {
						result = append(result, ApiChange{
							Id:          "request-parameter-pattern-removed",
							Level:       INFO,
							Text:        config.Localize("request-parameter-pattern-removed", patternDiff.From, ColorizedValue(paramLocation), ColorizedValue(paramName)),
							Operation:   operation,
							OperationId: operationItem.Revision.OperationID,
							Path:        path,
							Source:      source,
						})
					} else {
						level := WARN
						comment := config.Localize("pattern-changed-warn-comment")
						if patternDiff.To == ".*" {
							level = INFO
							comment = ""
						}
						result = append(result, ApiChange{
							Id:          "request-parameter-pattern-changed",
							Level:       level,
							Text:        config.Localize("request-parameter-pattern-changed", ColorizedValue(paramLocation), ColorizedValue(paramName), patternDiff.From, patternDiff.To),
							Comment:     comment,
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
