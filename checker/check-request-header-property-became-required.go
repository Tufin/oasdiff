package checker

import (
	"fmt"

	"github.com/tufin/oasdiff/diff"
)

func RequestHeaderPropertyBecameRequiredCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap) []BackwardCompatibilityError {
	result := make([]BackwardCompatibilityError, 0)
	if diffReport.PathsDiff == nil {
		return result
	}
	for path, pathItem := range diffReport.PathsDiff.Modified {
		if pathItem.OperationsDiff == nil {
			continue
		}
		for operation, operationItem := range pathItem.OperationsDiff.Modified {
			source := (*operationsSources)[operationItem.Revision]

			if operationItem.ParametersDiff == nil {
				continue
			}

			for paramLocation, paramDiffs := range operationItem.ParametersDiff.Modified {

				if paramLocation != "header" {
					continue
				}

				for paramName, paramDiff := range paramDiffs {
					if paramDiff.SchemaDiff == nil {
						continue
					}

					if paramDiff.SchemaDiff.RequiredDiff != nil {
						for _, changedRequiredPropertyName := range paramDiff.SchemaDiff.RequiredDiff.Added {
							if paramDiff.SchemaDiff.Revision.Value.Properties[changedRequiredPropertyName].Value.ReadOnly {
								continue
							}
							result = append(result, BackwardCompatibilityError{
								Id:        "request-header-property-became-required",
								Level:     ERR,
								Text:      fmt.Sprintf("the %s request header's property %s became required", ColorizedValue(paramName), ColorizedValue(changedRequiredPropertyName)),
								Operation: operation,
								Path:      path,
								Source:    source,
								ToDo:      "Add to exceptions-list.md",
							})
						}
					}
	
					if paramDiff.SchemaDiff.PropertiesDiff == nil {
						continue
					}
	
					for topPropertyName, topPropertyDiff := range paramDiff.SchemaDiff.PropertiesDiff.Modified {
						processModifiedPropertiesDiff(
							"",
							topPropertyName,
							topPropertyDiff,
							nil,
							func(propertyPath string, propertyName string, propertyDiff *diff.SchemaDiff, parent *diff.SchemaDiff) {
								requiredDiff := propertyDiff.RequiredDiff
								if requiredDiff == nil {
									return
								}
								for _, changedRequiredPropertyName := range requiredDiff.Added {
									if propertyDiff.Revision.Value.Properties[changedRequiredPropertyName].Value.ReadOnly {
										continue
									}		
									result = append(result, BackwardCompatibilityError{
										Id:        "request-header-property-became-required",
										Level:     ERR,
										Text:      fmt.Sprintf("the %s request header's property %s became required", ColorizedValue(paramName), ColorizedValue(propertyFullName(propertyPath, propertyFullName(propertyName, changedRequiredPropertyName)))),
										Operation: operation,
										Path:      path,
										Source:    source,
										ToDo:      "Add to exceptions-list.md",
									})							
								}
							})
					}
	

				}

			}
		}
	}
	return result
}
