package checker

import (
	"github.com/tufin/oasdiff/diff"
)

func RequestPropertyAnyOfUpdated(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config Config) Changes {
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
						Id:    "request-body-any-of-added",
						Level: INFO,
						Text: config.Localize(
							"request-body-any-of-added",
							ColorizedValue(mediaTypeDiff.SchemaDiff.AnyOfDiff.Added.String())),
						Operation:   operation,
						OperationId: operationItem.Revision.OperationID,
						Path:        path,
						Source:      source,
					})
				}

				if mediaTypeDiff.SchemaDiff.AnyOfDiff != nil && len(mediaTypeDiff.SchemaDiff.AnyOfDiff.Deleted) > 0 {
					result = append(result, ApiChange{
						Id:    "request-body-any-of-removed",
						Level: ERR,
						Text: config.Localize(
							"request-body-any-of-removed",
							ColorizedValue(mediaTypeDiff.SchemaDiff.AnyOfDiff.Deleted.String())),
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
								Id:    "request-property-any-of-added",
								Level: INFO,
								Text: config.Localize(
									"request-property-any-of-added",
									ColorizedValue(propertyDiff.AnyOfDiff.Added.String()),
									ColorizedValue(propertyFullName(propertyPath, propertyName))),
								Operation:   operation,
								OperationId: operationItem.Revision.OperationID,
								Path:        path,
								Source:      source,
							})
						}

						if len(propertyDiff.AnyOfDiff.Deleted) > 0 {
							result = append(result, ApiChange{
								Id:    "request-property-any-of-removed",
								Level: ERR,
								Text: config.Localize(
									"request-property-any-of-removed",
									ColorizedValue(propertyDiff.AnyOfDiff.Deleted.String()),
									ColorizedValue(propertyFullName(propertyPath, propertyName))),
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
