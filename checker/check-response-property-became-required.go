package checker

import (
	"github.com/tufin/oasdiff/diff"
	"github.com/tufin/oasdiff/load"
)

const (
	ResponsePropertyBecameRequiredId          = "response-property-became-required"
	ResponseWriteOnlyPropertyBecameRequiredId = "response-write-only-property-became-required"
)

func ResponsePropertyBecameRequiredCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config *Config) Changes {
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
						for _, changedRequiredPropertyName := range mediaTypeDiff.SchemaDiff.RequiredDiff.Added {
							id := ResponsePropertyBecameRequiredId
							if mediaTypeDiff.SchemaDiff.Revision.Properties[changedRequiredPropertyName] == nil {
								// removed properties processed by the ResponseRequiredPropertyUpdatedCheck check
								continue
							}
							if mediaTypeDiff.SchemaDiff.Revision.Properties[changedRequiredPropertyName].Value.WriteOnly {
								id = ResponseWriteOnlyPropertyBecameRequiredId
							}

							result = append(result, ApiChange{
								Id:          id,
								Level:       INFO,
								Args:        []any{changedRequiredPropertyName, responseStatus},
								Operation:   operation,
								OperationId: operationItem.Revision.OperationID,
								Path:        path,
								Source:      load.NewSource(source),
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
								id := ResponsePropertyBecameRequiredId
								if propertyDiff.Base.Properties[changedRequiredPropertyName] == nil {
									continue
								}
								if propertyDiff.Base.Properties[changedRequiredPropertyName].Value.WriteOnly {
									id = ResponseWriteOnlyPropertyBecameRequiredId
								}
								if propertyDiff.Revision.Properties[changedRequiredPropertyName] == nil {
									// removed properties processed by the ResponseRequiredPropertyUpdatedCheck check
									continue
								}

								result = append(result, ApiChange{
									Id:          id,
									Level:       INFO,
									Args:        []any{propertyFullName(propertyPath, propertyFullName(propertyName, changedRequiredPropertyName)), responseStatus},
									Operation:   operation,
									OperationId: operationItem.Revision.OperationID,
									Path:        path,
									Source:      load.NewSource(source),
								})
							}
						})
				}

			}
		}
	}
	return result
}
