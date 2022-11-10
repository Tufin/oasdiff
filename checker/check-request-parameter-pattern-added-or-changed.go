package checker

import (
	"fmt"

	"github.com/tufin/oasdiff/diff"
)

func RequestParameterPatternAddedOrChangedCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, diffBC *BCDiff) []BackwardCompatibilityError {
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
					if paramItem.SchemaDiff.PatternDiff == nil {
						continue
					}
					if paramItem.SchemaDiff.PatternDiff.To == "" ||
						paramItem.SchemaDiff.PatternDiff.To == ".*" {
						continue
					}
				
					source := (*operationsSources)[operationItem.Revision]

					if paramItem.SchemaDiff.PatternDiff.From == "" {
						result = append(result, BackwardCompatibilityError{
							Id:        "request-parameter-pattern-added",
							Level:     ERR,
							Text:      fmt.Sprintf("added the pattern '%s' for the %s request parameter %s", paramItem.SchemaDiff.PatternDiff.To, ColorizedValue(paramLocation), ColorizedValue(paramName)),
							Operation: operation,
							Path:      path,
							Source:    source,
							ToDo:      "Add to exceptions-list.md",
						})
					} else {
						result = append(result, BackwardCompatibilityError{
							Id:        "request-parameter-pattern-changed",
							Level:     WARN,
							Text:      fmt.Sprintf("changed the pattern for the %s request parameter %s from '%s' to '%s'", ColorizedValue(paramLocation), ColorizedValue(paramName), paramItem.SchemaDiff.PatternDiff.From, paramItem.SchemaDiff.PatternDiff.To),
							Operation: operation,
							Path:      path,
							Source:    source,
							ToDo:      "Add to exceptions-list.md",
						})
					}

					opDiff := diffBC.AddModifiedOperation(path, operation)
					if opDiff.ParametersDiff == nil {
						opDiff.ParametersDiff = &diff.ParametersDiff{}
					}
					if opDiff.ParametersDiff.Modified == nil {
						opDiff.ParametersDiff.Modified = make(diff.ParamDiffByLocation)
					}
					if opDiff.ParametersDiff.Modified[paramLocation] == nil {
						opDiff.ParametersDiff.Modified[paramLocation] = make(diff.ParamDiffs)
					}
					if opDiff.ParametersDiff.Modified[paramLocation][paramName] == nil {
						opDiff.ParametersDiff.Modified[paramLocation][paramName] = &diff.ParameterDiff{}
					}
					if opDiff.ParametersDiff.Modified[paramLocation][paramName].SchemaDiff == nil {
						opDiff.ParametersDiff.Modified[paramLocation][paramName].SchemaDiff = &diff.SchemaDiff{}
					}
					if opDiff.ParametersDiff.Modified[paramLocation][paramName].SchemaDiff.PatternDiff == nil {
						opDiff.ParametersDiff.Modified[paramLocation][paramName].SchemaDiff.PatternDiff = paramItem.SchemaDiff.PatternDiff
					}
				}
			}
		}
	}
	return result
}
