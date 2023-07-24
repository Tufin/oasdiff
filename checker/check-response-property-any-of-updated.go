package checker

import (
	"fmt"

	"github.com/tufin/oasdiff/diff"
)

func ResponsePropertyAnyOfUpdated(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config Config) Changes {
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

					if mediaTypeDiff.SchemaDiff.AnyOfDiff != nil && mediaTypeDiff.SchemaDiff.AnyOfDiff.Added > 0 {
						result = append(result, ApiChange{
							Id:    "response-body-any-of-added",
							Level: INFO,
							Text: fmt.Sprintf(
								config.i18n("response-body-any-of-added"),
								ColorizedValue(mediaTypeDiff.SchemaDiff.AnyOfDiff.AddedSchemas.String()),
								responseStatus),
							Operation:   operation,
							OperationId: operationItem.Revision.OperationID,
							Path:        path,
							Source:      source,
						})
					}

					if mediaTypeDiff.SchemaDiff.AnyOfDiff != nil && mediaTypeDiff.SchemaDiff.AnyOfDiff.Deleted > 0 {
						result = append(result, ApiChange{
							Id:    "response-body-any-of-removed",
							Level: INFO,
							Text: fmt.Sprintf(
								config.i18n("response-body-any-of-removed"),
								ColorizedValue(mediaTypeDiff.SchemaDiff.AnyOfDiff.DeletedSchemas.String()),
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

							if propertyDiff.AnyOfDiff.Added > 0 {
								result = append(result, ApiChange{
									Id:    "response-property-any-of-added",
									Level: INFO,
									Text: fmt.Sprintf(
										config.i18n("response-property-any-of-added"),
										ColorizedValue(propertyDiff.AnyOfDiff.AddedSchemas.String()),
										ColorizedValue(propertyFullName(propertyPath, propertyName)),
										responseStatus),
									Operation:   operation,
									OperationId: operationItem.Revision.OperationID,
									Path:        path,
									Source:      source,
								})
							}

							if propertyDiff.AnyOfDiff.Deleted > 0 {
								result = append(result, ApiChange{
									Id:    "response-property-any-of-removed",
									Level: INFO,
									Text: fmt.Sprintf(
										config.i18n("response-property-any-of-removed"),
										ColorizedValue(propertyDiff.AnyOfDiff.DeletedSchemas.String()),
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
