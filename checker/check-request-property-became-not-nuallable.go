package checker

import (
	"github.com/tufin/oasdiff/diff"
)

const (
	RequestBodyBecomeNotNullableId     = "request-body-became-not-nullable"
	RequestBodyBecomeNullableId        = "request-body-became-nullable"
	RequestPropertyBecomeNotNullableId = "request-property-became-not-nullable"
	RequestPropertyBecomeNullableId    = "request-property-became-nullable"
)

func RequestPropertyBecameNotNullableCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config *Config) Changes {
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

				if mediaTypeDiff.SchemaDiff.NullableDiff != nil {
					if mediaTypeDiff.SchemaDiff.NullableDiff.From == true {
						result = append(result, ApiChange{
							Id:        RequestBodyBecomeNotNullableId,
							Level:     ERR,
							Operation: operation,
							Path:      path,
							Source:    source,

							OperationId: operationItem.Revision.OperationID,
						})
					} else if mediaTypeDiff.SchemaDiff.NullableDiff.To == true {
						result = append(result, ApiChange{
							Id:        RequestBodyBecomeNullableId,
							Level:     INFO,
							Operation: operation,
							Path:      path,
							Source:    source,

							OperationId: operationItem.Revision.OperationID,
						})
					}
				}

				CheckModifiedPropertiesDiff(
					mediaTypeDiff.SchemaDiff,
					func(propertyPath string, propertyName string, propertyDiff *diff.SchemaDiff, parent *diff.SchemaDiff) {
						nullableDiff := propertyDiff.NullableDiff
						if nullableDiff == nil {
							return
						}

						propName := propertyFullName(propertyPath, propertyName)

						if nullableDiff.From == true {
							result = append(result, ApiChange{
								Id:        RequestPropertyBecomeNotNullableId,
								Level:     ERR,
								Args:      []any{propName},
								Operation: operation,
								Path:      path,
								Source:    source,

								OperationId: operationItem.Revision.OperationID,
							})
						} else if nullableDiff.To == true {
							result = append(result, ApiChange{
								Id:        RequestPropertyBecomeNullableId,
								Level:     INFO,
								Args:      []any{propName},
								Operation: operation,
								Path:      path,
								Source:    source,

								OperationId: operationItem.Revision.OperationID,
							})
						}

					})
			}
		}
	}
	return result
}
