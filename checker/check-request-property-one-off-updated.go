package checker

import (
	"fmt"

	"github.com/tufin/oasdiff/diff"
)

func RequestPropertyOneOffUpdated(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config Config) Changes {
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

				if mediaTypeDiff.SchemaDiff.OneOfDiff != nil && mediaTypeDiff.SchemaDiff.OneOfDiff.Added > 0 {
					result = append(result, ApiChange{
						Id:    "request-body-one-of-added",
						Level: INFO,
						Text: fmt.Sprintf(
							config.i18n("request-body-one-of-added"),
							ColorizedValue(mediaTypeDiff.SchemaDiff.OneOfDiff.AddedSchemas.String())),
						Operation:   operation,
						OperationId: operationItem.Revision.OperationID,
						Path:        path,
						Source:      source,
					})
				}

				if mediaTypeDiff.SchemaDiff.OneOfDiff != nil && mediaTypeDiff.SchemaDiff.OneOfDiff.Deleted > 0 {
					result = append(result, ApiChange{
						Id:    "request-body-one-of-removed",
						Level: ERR,
						Text: fmt.Sprintf(
							config.i18n("request-body-one-of-removed"),
							ColorizedValue(mediaTypeDiff.SchemaDiff.OneOfDiff.DeletedSchemas.String())),
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

						if propertyDiff.OneOfDiff.Added > 0 {
							result = append(result, ApiChange{
								Id:    "request-property-one-of-added",
								Level: INFO,
								Text: fmt.Sprintf(
									config.i18n("request-property-one-of-added"),
									ColorizedValue(propertyDiff.OneOfDiff.AddedSchemas.String()),
									ColorizedValue(propertyFullName(propertyPath, propertyName))),
								Operation:   operation,
								OperationId: operationItem.Revision.OperationID,
								Path:        path,
								Source:      source,
							})
						}

						if propertyDiff.OneOfDiff.Deleted > 0 {
							result = append(result, ApiChange{
								Id:    "request-property-one-of-removed",
								Level: ERR,
								Text: fmt.Sprintf(
									config.i18n("request-property-one-of-removed"),
									ColorizedValue(propertyDiff.OneOfDiff.DeletedSchemas.String()),
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
