package checker

import (
	"github.com/tufin/oasdiff/diff"
)

const (
	ResponseBodyAnyOfAddedId       = "response-body-any-of-added"
	ResponseBodyAnyOfRemovedId     = "response-body-any-of-removed"
	ResponsePropertyAnyOfAddedId   = "response-property-any-of-added"
	ResponsePropertyAnyOfRemovedId = "response-property-any-of-removed"
)

func ResponsePropertyAnyOfUpdatedCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config Config) Changes {
	result := make(Changes, 0)
	if diffReport.PathsDiff == nil {
		return result
	}

	for path, pathItem := range diffReport.PathsDiff.Modified {
		if pathItem.OperationsDiff == nil {
			continue
		}

		for operation, operationItem := range pathItem.OperationsDiff.Modified {
			if operationItem.ResponsesDiff == nil || operationItem.ResponsesDiff.Modified == nil {
				continue
			}
			source := (*operationsSources)[operationItem.Revision]

			for responseStatus, responsesDiff := range operationItem.ResponsesDiff.Modified {
				if responsesDiff.ContentDiff == nil || responsesDiff.ContentDiff.MediaTypeModified == nil {
					continue
				}

				modifiedMediaTypes := responsesDiff.ContentDiff.MediaTypeModified
				for _, mediaTypeDiff := range modifiedMediaTypes {
					if mediaTypeDiff.SchemaDiff == nil {
						continue
					}

					if mediaTypeDiff.SchemaDiff.AnyOfDiff != nil && len(mediaTypeDiff.SchemaDiff.AnyOfDiff.Added) > 0 {
						result = append(result, ApiChange{
							Id:    ResponseBodyAnyOfAddedId,
							Level: INFO,
							Text: config.Localize(
								ResponseBodyAnyOfAddedId,
								ColorizedValue(mediaTypeDiff.SchemaDiff.AnyOfDiff.Added.String()),
								responseStatus),
							Operation:   operation,
							OperationId: operationItem.Revision.OperationID,
							Path:        path,
							Source:      source,
						})
					}

					if mediaTypeDiff.SchemaDiff.AnyOfDiff != nil && len(mediaTypeDiff.SchemaDiff.AnyOfDiff.Deleted) > 0 {
						result = append(result, ApiChange{
							Id:    ResponseBodyAnyOfRemovedId,
							Level: INFO,
							Text: config.Localize(
								ResponseBodyAnyOfRemovedId,
								ColorizedValue(mediaTypeDiff.SchemaDiff.AnyOfDiff.Deleted.String()),
								responseStatus),
							Operation:   operation,
							OperationId: operationItem.Revision.OperationID,
							Path:        path,
							Source:      source,
						})
					}

					CheckModifiedPropertiesDiff(
						mediaTypeDiff.SchemaDiff,
						func(propertyPath string, propertyName string, propertyDiff *diff.SchemaDiff, parent *diff.SchemaDiff) {
							if propertyDiff.AnyOfDiff == nil {
								return
							}

							if len(propertyDiff.AnyOfDiff.Added) > 0 {
								result = append(result, ApiChange{
									Id:    ResponsePropertyAnyOfAddedId,
									Level: INFO,
									Text: config.Localize(
										ResponsePropertyAnyOfAddedId,
										ColorizedValue(propertyDiff.AnyOfDiff.Added.String()),
										ColorizedValue(propertyFullName(propertyPath, propertyName)),
										responseStatus),
									Operation:   operation,
									OperationId: operationItem.Revision.OperationID,
									Path:        path,
									Source:      source,
								})
							}

							if len(propertyDiff.AnyOfDiff.Deleted) > 0 {
								result = append(result, ApiChange{
									Id:    ResponsePropertyAnyOfRemovedId,
									Level: INFO,
									Text: config.Localize(
										ResponsePropertyAnyOfRemovedId,
										ColorizedValue(propertyDiff.AnyOfDiff.Deleted.String()),
										ColorizedValue(propertyFullName(propertyPath, propertyName)),
										responseStatus),
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
