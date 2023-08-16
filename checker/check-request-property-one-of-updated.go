package checker

import (
	"github.com/tufin/oasdiff/diff"
)

func RequestPropertyOneOfUpdated(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config Config) Changes {
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
						Id:    "request-body-one-of-added",
						Level: INFO,
						Text: config.Localize(
							"request-body-one-of-added",
							ColorizedValue(mediaTypeDiff.SchemaDiff.OneOfDiff.Added.String())),
						Operation:   operation,
						OperationId: operationItem.Revision.OperationID,
						Path:        path,
						Source:      source,
					})
				}

				if mediaTypeDiff.SchemaDiff.OneOfDiff != nil && len(mediaTypeDiff.SchemaDiff.OneOfDiff.Deleted) > 0 {
					result = append(result, ApiChange{
						Id:    "request-body-one-of-removed",
						Level: ERR,
						Text: config.Localize(
							"request-body-one-of-removed",
							ColorizedValue(mediaTypeDiff.SchemaDiff.OneOfDiff.Deleted.String())),
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

						if len(propertyDiff.OneOfDiff.Added) > 0 {
							result = append(result, ApiChange{
								Id:    "request-property-one-of-added",
								Level: INFO,
								Text: config.Localize(
									"request-property-one-of-added",
									ColorizedValue(propertyDiff.OneOfDiff.Added.String()),
									ColorizedValue(propertyFullName(propertyPath, propertyName))),
								Operation:   operation,
								OperationId: operationItem.Revision.OperationID,
								Path:        path,
								Source:      source,
							})
						}

						if len(propertyDiff.OneOfDiff.Deleted) > 0 {
							result = append(result, ApiChange{
								Id:    "request-property-one-of-removed",
								Level: ERR,
								Text: config.Localize(
									"request-property-one-of-removed",
									ColorizedValue(propertyDiff.OneOfDiff.Deleted.String()),
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
