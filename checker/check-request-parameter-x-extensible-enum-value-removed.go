package checker

import (
	"fmt"

	"github.com/tufin/oasdiff/diff"
	"golang.org/x/exp/slices"
)

func RequestParameterXExtensibleEnumValueRemovedCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap) []BackwardCompatibilityError {
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
			source := (*operationsSources)[operationItem.Revision]
			for paramLocation, paramItems := range operationItem.ParametersDiff.Modified {
				for paramName, paramItem := range paramItems {
					if paramItem.SchemaDiff == nil {
						continue
					}
					if paramItem.SchemaDiff.ExtensionsDiff == nil {
						continue
					}
					if paramItem.SchemaDiff.ExtensionsDiff.Modified == nil {
						continue
					}
					if paramItem.SchemaDiff.ExtensionsDiff.Modified["x-extensible-enum"] == nil {
						continue
					}
					from, ok := paramItem.SchemaDiff.ExtensionsDiff.Modified["x-extensible-enum"].From.([]string)
					if !ok {
						continue
					}
					to, ok := paramItem.SchemaDiff.ExtensionsDiff.Modified["x-extensible-enum"].To.([]string)
					if !ok {
						continue
					}
					deletedVals := make([]string, 0)

					for _, fromVal := range from {
						if !slices.Contains(to, fromVal) {
							deletedVals = append(deletedVals, fromVal)
						}
					}

					for _, enumVal := range deletedVals {
						result = append(result, BackwardCompatibilityError{
							Id:        "request-parameter-x-extensible-enum-value-removed",
							Level:     ERR,
							Text:      fmt.Sprintf("removed the x-extensible-enum value %s for the %s request parameter %s", enumVal, ColorizedValue(paramLocation), ColorizedValue(paramName)),
							Operation: operation,
							Path:      path,
							Source:    source,
							ToDo:      "Add to exceptions-list.md",
						})
					}
				}
			}
		}
	}
	return result
}
