package checker

import (
	"fmt"

	"github.com/tufin/oasdiff/diff"
	"golang.org/x/exp/slices"
)

func NewRequiredRequestHeaderPropertyCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap) []BackwardCompatibilityError {
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

					if paramDiff.SchemaDiff.PropertiesDiff != nil {
						for _, newPropertyName := range paramDiff.SchemaDiff.PropertiesDiff.Added {
							if paramDiff.SchemaDiff.Revision.Value.Properties[newPropertyName].Value.ReadOnly {
								continue
							}
							if !slices.Contains[string](paramDiff.SchemaDiff.Revision.Value.Required, newPropertyName) {
								continue
							}
							result = append(result, BackwardCompatibilityError{
								Id:        "new-required-request-header-property",
								Level:     ERR,
								Text:      fmt.Sprintf("added the new required %s request header's property %s", ColorizedValue(paramName), ColorizedValue(newPropertyName)),
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
								propertiesDiff := propertyDiff.PropertiesDiff
								if propertiesDiff == nil {
									return
								}
								for _, newPropertyName := range propertiesDiff.Added {
									if propertyDiff.Revision.Value.Properties[newPropertyName].Value.ReadOnly {
										continue
									}
									if !slices.Contains[string](propertyDiff.Revision.Value.Required, newPropertyName) {
										continue
									}
		
									result = append(result, BackwardCompatibilityError{
										Id:        "new-required-request-header-property",
										Level:     ERR,
										Text:      fmt.Sprintf("added the new required %s request header's property %s", ColorizedValue(paramName), ColorizedValue(propertyFullName(propertyPath, propertyFullName(propertyName, newPropertyName)))),
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
