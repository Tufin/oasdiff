package checker

import (
	"fmt"

	"github.com/tufin/oasdiff/diff"
)

func ResponsePropertyBecameRequiredCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config BackwardCompatibilityCheckConfig) []BackwardCompatibilityError {
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

			if operationItem.ResponsesDiff == nil {
				continue
			}

			for responseStatus, responseDiff := range operationItem.ResponsesDiff.Modified {
				if responseDiff.ContentDiff == nil ||
					responseDiff.ContentDiff.MediaTypeModified == nil {
					continue
				}

				modifiedMediaTypes := responseDiff.ContentDiff.MediaTypeModified
				for _, mediaTypeDiff := range modifiedMediaTypes {
					if mediaTypeDiff.SchemaDiff == nil {
						continue
					}

					if mediaTypeDiff.SchemaDiff.RequiredDiff != nil {
						for _, changedRequiredPropertyName := range mediaTypeDiff.SchemaDiff.RequiredDiff.Added {
							comment := ""
							if mediaTypeDiff.SchemaDiff.Revision.Value.Properties[changedRequiredPropertyName] == nil {
								// removed properties processed by the ResponseRequiredPropertyRemovedCheck check
								continue
							}
							if mediaTypeDiff.SchemaDiff.Revision.Value.Properties[changedRequiredPropertyName].Value.WriteOnly {
								comment = "the property is write only"

							}

							result = append(result, BackwardCompatibilityError{
								Id:          "response-property-became-required",
								Level:       INFO,
								Text:        fmt.Sprintf(config.i18n("response-property-became-required"), ColorizedValue(changedRequiredPropertyName), ColorizedValue(responseStatus)),
								Operation:   operation,
								OperationId: operationItem.Revision.OperationID,
								Path:        path,
								Source:      source,
								Comment:     comment,
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
								comment := ""
								if propertyDiff.Base.Value.Properties[changedRequiredPropertyName] == nil {
									continue
								}
								if propertyDiff.Base.Value.Properties[changedRequiredPropertyName].Value.WriteOnly {
									comment = "the property is write only"
								}
								if propertyDiff.Revision.Value.Properties[changedRequiredPropertyName] == nil {
									// removed properties processed by the ResponseRequiredPropertyRemovedCheck check
									continue
								}
								result = append(result, BackwardCompatibilityError{
									Id:          "response-property-became-required",
									Level:       INFO,
									Text:        fmt.Sprintf(config.i18n("response-property-became-required"), ColorizedValue(propertyFullName(propertyPath, propertyFullName(propertyName, changedRequiredPropertyName))), ColorizedValue(responseStatus)),
									Operation:   operation,
									OperationId: operationItem.Revision.OperationID,
									Path:        path,
									Source:      source,
									Comment:     comment,
								})
							}
						})
				}

			}
		}
	}
	return result
}
