package checker

import (
	"fmt"

	"github.com/tufin/oasdiff/diff"
)

func RequestHeaderPropertyBecameRequiredCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config BackwardCompatibilityCheckConfig) []BackwardCompatibilityError {
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
							if paramDiff.SchemaDiff.Revision.Value.Properties[changedRequiredPropertyName] == nil {
								continue
							}
							if paramDiff.SchemaDiff.Revision.Value.Properties[changedRequiredPropertyName].Value.ReadOnly {
								continue
							}

							if paramDiff.SchemaDiff.Base.Value.Properties[changedRequiredPropertyName] == nil {
								// new added required properties processed via the new-required-request-header-property check
								continue
							}

							result = append(result, BackwardCompatibilityError{
								Id:          "request-header-property-became-required",
								Level:       ERR,
								Text:        fmt.Sprintf(config.i18n("request-header-property-became-required"), ColorizedValue(paramName), ColorizedValue(changedRequiredPropertyName)),
								Operation:   operation,
								OperationId: operationItem.Revision.OperationID,
								Path:        path,
								Source:      source,
							})
						}
					}

					CheckModifiedPropertiesDiff(
						paramDiff.SchemaDiff,
						func(propertyPath string, propertyName string, propertyDiff *diff.SchemaDiff, parent *diff.SchemaDiff) {
							requiredDiff := propertyDiff.RequiredDiff
							if requiredDiff == nil {
								return
							}
							for _, changedRequiredPropertyName := range requiredDiff.Added {
								if propertyDiff.Revision.Value.Properties[changedRequiredPropertyName] == nil {
									continue
								}
								if propertyDiff.Revision.Value.Properties[changedRequiredPropertyName].Value.ReadOnly {
									continue
								}
								result = append(result, BackwardCompatibilityError{
									Id:          "request-header-property-became-required",
									Level:       ERR,
									Text:        fmt.Sprintf(config.i18n("request-header-property-became-required"), ColorizedValue(paramName), ColorizedValue(propertyFullName(propertyPath, propertyFullName(propertyName, changedRequiredPropertyName)))),
									Operation:   operation,
									OperationId: operationItem.Revision.OperationID,
									Path:        path,
									Source:      source,
								})
							}
						})
				}
			}
		}
	}
	return result
}
