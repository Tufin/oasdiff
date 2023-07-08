package checker

import (
	"fmt"

	"github.com/tufin/oasdiff/diff"
)

const requestPropertyBecameNotNullableId = "request-property-became-not-nullable"
const requestBodyBecameNotNullableId = "request-body-became-not-nullable"

func RequestPropertyBecameNotNullableCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config Config) Changes {
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

				if mediaTypeDiff.SchemaDiff.NullableDiff != nil && mediaTypeDiff.SchemaDiff.NullableDiff.From == true {
					result = append(result, ApiChange{
						Id:        requestBodyBecameNotNullableId,
						Level:     ERR,
						Text:      config.i18n(requestBodyBecameNotNullableId),
						Operation: operation,
						Path:      path,
						Source:    source,
					})
				}

				CheckModifiedPropertiesDiff(
					mediaTypeDiff.SchemaDiff,
					func(propertyPath string, propertyName string, propertyDiff *diff.SchemaDiff, parent *diff.SchemaDiff) {
						nullableDiff := propertyDiff.NullableDiff
						if nullableDiff == nil {
							return
						}
						if nullableDiff.From != true {
							return
						}

						result = append(result, ApiChange{
							Id:        requestPropertyBecameNotNullableId,
							Level:     ERR,
							Text:      fmt.Sprintf(config.i18n(requestPropertyBecameNotNullableId), ColorizedValue(propertyFullName(propertyPath, propertyName))),
							Operation: operation,
							Path:      path,
							Source:    source,
						})
					})
			}
		}
	}
	return result
}
