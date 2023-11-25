package checker

import (
	"github.com/tufin/oasdiff/diff"
)

const (
	ResponsePropertyBecameOptionalId          = "response-property-became-optional"
	ResponseWriteOnlyPropertyBecameOptionalId = "response-write-only-property-became-optional"
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
							id := ResponsePropertyBecameOptionalId
							level := ERR
							if mediaTypeDiff.SchemaDiff.Revision.Properties[changedRequiredPropertyName] == nil {
								// removed properties processed by the ResponseRequiredPropertyUpdatedCheck check
								continue
							}
							if mediaTypeDiff.SchemaDiff.Revision.Properties[changedRequiredPropertyName].Value.WriteOnly {
								id = ResponseWriteOnlyPropertyBecameOptionalId
								level = INFO
							}

							result = append(result, ApiChange{
								Id:          id,
								Level:       level,
								Text:        config.Localize(id, ColorizedValue(changedRequiredPropertyName), ColorizedValue(responseStatus)),
								Args:        []any{changedRequiredPropertyName, responseStatus},
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
							for _, changedRequiredPropertyName := range requiredDiff.Deleted {
								level := ERR
								id := ResponsePropertyBecameOptionalId

								if propertyDiff.Base.Properties[changedRequiredPropertyName] == nil {
									continue
								}
								if propertyDiff.Base.Properties[changedRequiredPropertyName].Value.WriteOnly {
									level = INFO
									id = ResponseWriteOnlyPropertyBecameOptionalId
								}
								if propertyDiff.Revision.Properties[changedRequiredPropertyName] == nil {
									// removed properties processed by the ResponseRequiredPropertyUpdatedCheck check
									continue
								}

								propName := propertyFullName(propertyPath, propertyFullName(propertyName, changedRequiredPropertyName))

								result = append(result, ApiChange{
									Id:          id,
									Level:       level,
									Text:        config.Localize(id, ColorizedValue(propName), ColorizedValue(responseStatus)),
									Args:        []any{propName, responseStatus},
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
