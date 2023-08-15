package checker

import (
	"github.com/tufin/oasdiff/diff"
)

func RequestPropertyAllOfUpdated(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config Config) Changes {
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

				if mediaTypeDiff.SchemaDiff.AllOfDiff != nil && len(mediaTypeDiff.SchemaDiff.AllOfDiff.Added) > 0 {
					result = append(result, ApiChange{
						Id:    "request-body-all-of-added",
						Level: ERR,
						Text: config.Localize(
							"request-body-all-of-added",
							ColorizedValue(mediaTypeDiff.SchemaDiff.AllOfDiff.Added.String())),
						Operation:   operation,
						OperationId: operationItem.Revision.OperationID,
						Path:        path,
						Source:      source,
					})
				}

				if mediaTypeDiff.SchemaDiff.AllOfDiff != nil && len(mediaTypeDiff.SchemaDiff.AllOfDiff.Deleted) > 0 {
					result = append(result, ApiChange{
						Id:    "request-body-all-of-removed",
						Level: WARN,
						Text: config.Localize(
							"request-body-all-of-removed",
							ColorizedValue(mediaTypeDiff.SchemaDiff.AllOfDiff.Deleted.String())),
						Operation:   operation,
						OperationId: operationItem.Revision.OperationID,
						Path:        path,
						Source:      source,
					})
				}

				CheckModifiedPropertiesDiff(
					mediaTypeDiff.SchemaDiff,
					func(propertyPath string, propertyName string, propertyDiff *diff.SchemaDiff, parent *diff.SchemaDiff) {
						if propertyDiff.AllOfDiff == nil {
							return
						}

						if len(propertyDiff.AllOfDiff.Added) > 0 {
							result = append(result, ApiChange{
								Id:    "request-property-all-of-added",
								Level: ERR,
								Text: config.Localize(
									"request-property-all-of-added",
									ColorizedValue(propertyDiff.AllOfDiff.Added.String()),
									ColorizedValue(propertyFullName(propertyPath, propertyName))),
								Operation:   operation,
								OperationId: operationItem.Revision.OperationID,
								Path:        path,
								Source:      source,
							})
						}

						if len(propertyDiff.AllOfDiff.Deleted) > 0 {
							result = append(result, ApiChange{
								Id:    "request-property-all-of-removed",
								Level: WARN,
								Text: config.Localize(
									"request-property-all-of-removed",
									ColorizedValue(propertyDiff.AllOfDiff.Deleted.String()),
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
