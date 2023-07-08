package checker

import (
	"fmt"

	"github.com/tufin/oasdiff/diff"
)

func ResponsePropertyBecameOptionalCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config Config) Changes {
	result := make(Changes, 0)
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
						for _, changedRequiredPropertyName := range mediaTypeDiff.SchemaDiff.RequiredDiff.Deleted {
							id := "response-property-became-optional"
							level := ERR
							if mediaTypeDiff.SchemaDiff.Revision.Value.Properties[changedRequiredPropertyName] == nil {
								// removed properties processed by the ResponseRequiredPropertyUpdatedCheck check
								continue
							}
							if mediaTypeDiff.SchemaDiff.Revision.Value.Properties[changedRequiredPropertyName].Value.WriteOnly {
								id = "response-write-only-property-became-optional"
								level = INFO
							}

							result = append(result, ApiChange{
								Id:     id,
								Level:  level,
								Text:   fmt.Sprintf(config.i18n(id), ColorizedValue(changedRequiredPropertyName), ColorizedValue(responseStatus)),
								Source: source,
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
							for _, changedRequiredPropertyName := range requiredDiff.Deleted {
								level := ERR
								id := "response-property-became-optional"

								if propertyDiff.Base.Value.Properties[changedRequiredPropertyName] == nil {
									continue
								}
								if propertyDiff.Base.Value.Properties[changedRequiredPropertyName].Value.WriteOnly {
									level = INFO
									id = "response-write-only-property-became-optional"
								}
								if propertyDiff.Revision.Value.Properties[changedRequiredPropertyName] == nil {
									// removed properties processed by the ResponseRequiredPropertyUpdatedCheck check
									continue
								}
								result = append(result, ApiChange{
									Id:          id,
									Level:       level,
									Text:        fmt.Sprintf(config.i18n(id), ColorizedValue(propertyFullName(propertyPath, propertyFullName(propertyName, changedRequiredPropertyName))), ColorizedValue(responseStatus)),
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
