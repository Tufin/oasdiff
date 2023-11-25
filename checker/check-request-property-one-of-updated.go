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

func RequestPropertyOneOfUpdatedCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config Config) Changes {
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
			source := (*operationsSources)[operationItem.Revision]

			modifiedMediaTypes := operationItem.RequestBodyDiff.ContentDiff.MediaTypeModified
			for _, mediaTypeDiff := range modifiedMediaTypes {
				if mediaTypeDiff.SchemaDiff == nil {
					continue
				}

				if mediaTypeDiff.SchemaDiff.OneOfDiff != nil && len(mediaTypeDiff.SchemaDiff.OneOfDiff.Added) > 0 {
					result = append(result, ApiChange{
						Id:    RequestBodyOneOfAddedId,
						Level: INFO,
						Text: config.Localize(
							RequestBodyOneOfAddedId,
							ColorizedValue(mediaTypeDiff.SchemaDiff.OneOfDiff.Added.String())),
						Args:        []any{mediaTypeDiff.SchemaDiff.OneOfDiff.Added.String()},
						Operation:   operation,
						OperationId: operationItem.Revision.OperationID,
						Path:        path,
						Source:      source,
					})
				}

				if mediaTypeDiff.SchemaDiff.OneOfDiff != nil && len(mediaTypeDiff.SchemaDiff.OneOfDiff.Deleted) > 0 {
					result = append(result, ApiChange{
						Id:    RequestBodyOneOfRemovedId,
						Level: ERR,
						Text: config.Localize(
							RequestBodyOneOfRemovedId,
							ColorizedValue(mediaTypeDiff.SchemaDiff.OneOfDiff.Deleted.String())),
						Args:        []any{mediaTypeDiff.SchemaDiff.OneOfDiff.Deleted.String()},
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
								Id:    RequestPropertyOneOfAddedId,
								Level: INFO,
								Text: config.Localize(
									RequestPropertyOneOfAddedId,
									ColorizedValue(propertyDiff.OneOfDiff.Added.String()),
									ColorizedValue(propName)),
								Args:        []any{propertyDiff.OneOfDiff.Added.String(), propName},
								Operation:   operation,
								OperationId: operationItem.Revision.OperationID,
								Path:        path,
								Source:      source,
							})
						}

						if len(propertyDiff.OneOfDiff.Deleted) > 0 {
							result = append(result, ApiChange{
								Id:    RequestPropertyOneOfRemovedId,
								Level: ERR,
								Text: config.Localize(
									RequestPropertyOneOfRemovedId,
									ColorizedValue(propertyDiff.OneOfDiff.Deleted.String()),
									ColorizedValue(propName)),
								Args:        []any{propertyDiff.OneOfDiff.Deleted.String(), propName},
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
	return result
}
