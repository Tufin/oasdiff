package checker

import (
	"github.com/tufin/oasdiff/diff"
)

const (
	RequestBodyAnyOfAddedId       = "request-body-any-of-added"
	RequestBodyAnyOfRemovedId     = "request-body-any-of-removed"
	RequestPropertyAnyOfAddedId   = "request-property-any-of-added"
	RequestPropertyAnyOfRemovedId = "request-property-any-of-removed"
)

func RequestPropertyAnyOfUpdatedCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config *Config) Changes {
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

				if mediaTypeDiff.SchemaDiff.AnyOfDiff != nil && len(mediaTypeDiff.SchemaDiff.AnyOfDiff.Added) > 0 {
					result = append(result, ApiChange{
						Id:    RequestBodyAnyOfAddedId,
						Level: INFO,
						Text: config.Localize(
							RequestBodyAnyOfAddedId,
							ColorizedValue(mediaTypeDiff.SchemaDiff.AnyOfDiff.Added.String())),
						Args:        []any{mediaTypeDiff.SchemaDiff.AnyOfDiff.Added.String()},
						Operation:   operation,
						OperationId: operationItem.Revision.OperationID,
						Path:        path,
						Source:      source,
					})
				}

				if mediaTypeDiff.SchemaDiff.AnyOfDiff != nil && len(mediaTypeDiff.SchemaDiff.AnyOfDiff.Deleted) > 0 {
					result = append(result, ApiChange{
						Id:    RequestBodyAnyOfRemovedId,
						Level: ERR,
						Text: config.Localize(
							RequestBodyAnyOfRemovedId,
							ColorizedValue(mediaTypeDiff.SchemaDiff.AnyOfDiff.Deleted.String())),
						Args:        []any{mediaTypeDiff.SchemaDiff.AnyOfDiff.Deleted.String()},
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

						propName := propertyFullName(propertyPath, propertyName)

						if len(propertyDiff.AnyOfDiff.Added) > 0 {
							result = append(result, ApiChange{
								Id:    RequestPropertyAnyOfAddedId,
								Level: INFO,
								Text: config.Localize(
									RequestPropertyAnyOfAddedId,
									ColorizedValue(propertyDiff.AnyOfDiff.Added.String()),
									ColorizedValue(propName)),
								Args:        []any{propertyDiff.AnyOfDiff.Added.String(), propName},
								Operation:   operation,
								OperationId: operationItem.Revision.OperationID,
								Path:        path,
								Source:      source,
							})
						}

						if len(propertyDiff.AnyOfDiff.Deleted) > 0 {
							result = append(result, ApiChange{
								Id:    RequestPropertyAnyOfRemovedId,
								Level: ERR,
								Text: config.Localize(
									RequestPropertyAnyOfRemovedId,
									ColorizedValue(propertyDiff.AnyOfDiff.Deleted.String()),
									ColorizedValue(propName)),
								Args:        []any{propertyDiff.AnyOfDiff.Deleted.String(), propName},
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
