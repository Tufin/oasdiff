package checker

import (
	"fmt"

	"github.com/tufin/oasdiff/diff"
)

func RequestPropertyBecameRequiredCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config BackwardCompatibilityCheckConfig) []BackwardCompatibilityError {
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

			if operationItem.RequestBodyDiff == nil ||
				operationItem.RequestBodyDiff.ContentDiff == nil ||
				operationItem.RequestBodyDiff.ContentDiff.MediaTypeModified == nil {
				continue
			}
			modifiedMediaTypes := operationItem.RequestBodyDiff.ContentDiff.MediaTypeModified
			for _, mediaTypeDiff := range modifiedMediaTypes {
				if mediaTypeDiff.SchemaDiff == nil {
					continue
				}

				if mediaTypeDiff.SchemaDiff.RequiredDiff != nil {
					for _, changedRequiredPropertyName := range mediaTypeDiff.SchemaDiff.RequiredDiff.Added {
						if mediaTypeDiff.SchemaDiff.Base.Value.Properties[changedRequiredPropertyName] == nil {
							// it is a new property, checked by the new-required-request-property check
							continue
						}
						if mediaTypeDiff.SchemaDiff.Revision.Value.Properties[changedRequiredPropertyName] == nil {
							continue
						}
						if mediaTypeDiff.SchemaDiff.Revision.Value.Properties[changedRequiredPropertyName].Value.ReadOnly {
							continue
						}
						result = append(result, BackwardCompatibilityError{
							Id:          "request-property-became-required",
							Level:       ERR,
							Text:        fmt.Sprintf(config.i18n("request-property-became-required"), ColorizedValue(changedRequiredPropertyName)),
							Operation:   operation,
							OperationId: operationItem.Revision.OperationID,
							Path:        path,
							Source:      source,
						})
					}
				}

				CheckModifiedPropertiesDiff(
					mediaTypeDiff.SchemaDiff,
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
							if propertyDiff.Base.Value.Properties[changedRequiredPropertyName] == nil {
								// it is a new property, checked by the new-required-request-property check
								return
							}
							result = append(result, BackwardCompatibilityError{
								Id:          "request-property-became-required",
								Level:       ERR,
								Text:        fmt.Sprintf(config.i18n("request-property-became-required"), ColorizedValue(propertyFullName(propertyPath, propertyFullName(propertyName, changedRequiredPropertyName)))),
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
	return result
}
