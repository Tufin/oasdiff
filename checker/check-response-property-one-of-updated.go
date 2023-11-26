package checker

import (
	"github.com/tufin/oasdiff/diff"
)

const (
	ResponseBodyOneOfAddedId       = "response-body-one-of-added"
	ResponseBodyOneOfRemovedId     = "response-body-one-of-removed"
	ResponsePropertyOneOfAddedId   = "response-property-one-of-added"
	ResponsePropertyOneOfRemovedId = "response-property-one-of-removed"
)

func ResponsePropertyOneOfUpdated(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config *Config) Changes {
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

					if mediaTypeDiff.SchemaDiff.OneOfDiff != nil && len(mediaTypeDiff.SchemaDiff.OneOfDiff.Added) > 0 {
						result = append(result, ApiChange{
							Id:    ResponseBodyOneOfAddedId,
							Level: INFO,
							Text: config.Localize(
								ResponseBodyOneOfAddedId,
								ColorizedValue(mediaTypeDiff.SchemaDiff.OneOfDiff.Added.String()),
								responseStatus),
							Args:        []any{mediaTypeDiff.SchemaDiff.OneOfDiff.Added.String(), responseStatus},
							Operation:   operation,
							OperationId: operationItem.Revision.OperationID,
							Path:        path,
							Source:      source,
						})
					}

					if mediaTypeDiff.SchemaDiff.OneOfDiff != nil && len(mediaTypeDiff.SchemaDiff.OneOfDiff.Deleted) > 0 {
						result = append(result, ApiChange{
							Id:    ResponseBodyOneOfRemovedId,
							Level: INFO,
							Text: config.Localize(
								ResponseBodyOneOfRemovedId,
								ColorizedValue(mediaTypeDiff.SchemaDiff.OneOfDiff.Deleted.String()),
								responseStatus),
							Args:        []any{mediaTypeDiff.SchemaDiff.OneOfDiff.Deleted.String(), responseStatus},
							Operation:   operation,
							OperationId: operationItem.Revision.OperationID,
							Path:        path,
							Source:      source,
						})
					}

					CheckModifiedPropertiesDiff(
						mediaTypeDiff.SchemaDiff,
						func(propertyPath string, propertyName string, propertyDiff *diff.SchemaDiff, parent *diff.SchemaDiff) {
							if propertyDiff.OneOfDiff == nil {
								return
							}

							propName := propertyFullName(propertyPath, propertyName)

							if len(propertyDiff.OneOfDiff.Added) > 0 {
								result = append(result, ApiChange{
									Id:    ResponsePropertyOneOfAddedId,
									Level: INFO,
									Text: config.Localize(
										ResponsePropertyOneOfAddedId,
										ColorizedValue(propertyDiff.OneOfDiff.Added.String()),
										ColorizedValue(propName),
										responseStatus),
									Args:        []any{propertyDiff.OneOfDiff.Added.String(), propName, responseStatus},
									Operation:   operation,
									OperationId: operationItem.Revision.OperationID,
									Path:        path,
									Source:      source,
								})
							}

							if len(propertyDiff.OneOfDiff.Deleted) > 0 {
								result = append(result, ApiChange{
									Id:    ResponsePropertyOneOfRemovedId,
									Level: INFO,
									Text: config.Localize(
										ResponsePropertyOneOfRemovedId,
										ColorizedValue(propertyDiff.OneOfDiff.Deleted.String()),
										ColorizedValue(propName),
										responseStatus),
									Args:        []any{propertyDiff.OneOfDiff.Deleted.String(), propName, responseStatus},
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
