package checker

import (
	"github.com/tufin/oasdiff/diff"
)

const responsePropertyBecameNullableId = "response-property-became-nullable"
const responseBodyBecameNullableId = "response-body-became-nullable"

func ResponsePropertyBecameNullableCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config Config) Changes {
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

					if mediaTypeDiff.SchemaDiff.NullableDiff != nil && mediaTypeDiff.SchemaDiff.NullableDiff.To == true {
						result = append(result, ApiChange{
							Id:          responseBodyBecameNullableId,
							Level:       ERR,
							Text:        config.Localize(responseBodyBecameNullableId),
							Operation:   operation,
							OperationId: operationItem.Revision.OperationID,
							Path:        path,
							Source:      source,
						})
					}

					CheckModifiedPropertiesDiff(
						mediaTypeDiff.SchemaDiff,
						func(propertyPath string, propertyName string, propertyDiff *diff.SchemaDiff, parent *diff.SchemaDiff) {
							nullableDiff := propertyDiff.NullableDiff
							if nullableDiff == nil {
								return
							}
							if nullableDiff.To != true {
								return
							}

							result = append(result, ApiChange{
								Id:          responsePropertyBecameNullableId,
								Level:       ERR,
								Text:        config.Localize(responsePropertyBecameNullableId, ColorizedValue(propertyFullName(propertyPath, propertyName)), ColorizedValue(responseStatus)),
								Operation:   operation,
								OperationId: operationItem.Revision.OperationID,
								Path:        path,
								Source:      source,
							})
						})
				}
			}
		}
	}
	return result
}
