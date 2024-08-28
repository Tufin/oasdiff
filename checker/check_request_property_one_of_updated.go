package checker

import (
	"github.com/tufin/oasdiff/diff"
)

const (
	RequestBodyOneOfAddedId       = "request-body-one-of-added"
	RequestBodyOneOfRemovedId     = "request-body-one-of-removed"
	RequestPropertyOneOfAddedId   = "request-property-one-of-added"
	RequestPropertyOneOfRemovedId = "request-property-one-of-removed"
)

func RequestPropertyOneOfUpdatedCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config *Config) Changes {
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
			for mediaType, mediaTypeDiff := range modifiedMediaTypes {
				if mediaTypeDiff.SchemaDiff == nil {
					continue
				}

				if mediaTypeDiff.SchemaDiff.OneOfDiff != nil && len(mediaTypeDiff.SchemaDiff.OneOfDiff.Added) > 0 {
					result = append(result, NewApiChange(
						RequestBodyOneOfAddedId,
						config,
						[]any{mediaTypeDiff.SchemaDiff.OneOfDiff.Added.String(), mediaType},
						"",
						operationsSources,
						operationItem.Revision,
						operation,
						path,
					))
				}

				if mediaTypeDiff.SchemaDiff.OneOfDiff != nil && len(mediaTypeDiff.SchemaDiff.OneOfDiff.Deleted) > 0 {
					result = append(result, NewApiChange(
						RequestBodyOneOfRemovedId,
						config,
						[]any{mediaTypeDiff.SchemaDiff.OneOfDiff.Deleted.String(), mediaType},
						"",
						operationsSources,
						operationItem.Revision,
						operation,
						path,
					))
				}

				CheckModifiedPropertiesDiff(
					mediaTypeDiff.SchemaDiff,
					func(propertyPath string, propertyName string, propertyDiff *diff.SchemaDiff, parent *diff.SchemaDiff) {
						if propertyDiff.OneOfDiff == nil {
							return
						}

						propName := propertyFullName(propertyPath, propertyName)

						if len(propertyDiff.OneOfDiff.Added) > 0 {
							result = append(result, NewApiChange(
								RequestPropertyOneOfAddedId,
								config,
								[]any{propertyDiff.OneOfDiff.Added.String(), propName, mediaType},
								"",
								operationsSources,
								operationItem.Revision,
								operation,
								path,
							))
						}

						if len(propertyDiff.OneOfDiff.Deleted) > 0 {
							result = append(result, NewApiChange(
								RequestPropertyOneOfRemovedId,
								config,
								[]any{propertyDiff.OneOfDiff.Deleted.String(), propName, mediaType},
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
