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
						result = append(result, NewApiChange(
							RequestBodyBecomeNotNullableId,
							ERR,
							nil,
							"",
							operationsSources,
							operationItem.Revision,
							operation,
							path,
						))
					} else if mediaTypeDiff.SchemaDiff.NullableDiff.To == true {
						result = append(result, NewApiChange(
							RequestBodyBecomeNullableId,
							INFO,
							nil,
							"",
							operationsSources,
							operationItem.Revision,
							operation,
							path,
						))
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
							result = append(result, NewApiChange(
								RequestPropertyBecomeNotNullableId,
								ERR,
								[]any{propName},
								"",
								operationsSources,
								operationItem.Revision,
								operation,
								path,
							))
						} else if nullableDiff.To == true {
							result = append(result, NewApiChange(
								RequestPropertyBecomeNullableId,
								INFO,
								[]any{propName},
								"",
								operationsSources,
								operationItem.Revision,
								operation,
								path,
							))
						}

					})
			}
		}
	}
	return result
}
