package checker

import (
	"fmt"

	"github.com/tufin/oasdiff/diff"
)

func ResponsePropertyOneOffUpdated(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config Config) Changes {
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

					if mediaTypeDiff.SchemaDiff.OneOfDiff != nil && mediaTypeDiff.SchemaDiff.OneOfDiff.Added > 0 {
						result = append(result, ApiChange{
							Id:    "response-body-one-of-added",
							Level: INFO,
							Text: fmt.Sprintf(
								config.i18n("response-body-one-of-added"),
								ColorizedValue(mediaTypeDiff.SchemaDiff.OneOfDiff.AddedSchemas.String()),
								responseStatus),
							Operation:   operation,
							OperationId: operationItem.Revision.OperationID,
							Path:        path,
							Source:      source,
						})
					}

					if mediaTypeDiff.SchemaDiff.OneOfDiff != nil && mediaTypeDiff.SchemaDiff.OneOfDiff.Deleted > 0 {
						result = append(result, ApiChange{
							Id:    "response-body-one-of-removed",
							Level: INFO,
							Text: fmt.Sprintf(
								config.i18n("response-body-one-of-removed"),
								ColorizedValue(mediaTypeDiff.SchemaDiff.OneOfDiff.DeletedSchemas.String()),
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
							if propertyDiff.OneOfDiff == nil {
								return
							}

							if propertyDiff.OneOfDiff.Added > 0 {
								result = append(result, ApiChange{
									Id:    "response-property-one-of-added",
									Level: INFO,
									Text: fmt.Sprintf(
										config.i18n("response-property-one-of-added"),
										ColorizedValue(propertyDiff.OneOfDiff.AddedSchemas.String()),
										ColorizedValue(propertyFullName(propertyPath, propertyName)),
										responseStatus),
									Operation:   operation,
									OperationId: operationItem.Revision.OperationID,
									Path:        path,
									Source:      source,
								})
							}

							if propertyDiff.OneOfDiff.Deleted > 0 {
								result = append(result, ApiChange{
									Id:    "response-property-one-of-removed",
									Level: INFO,
									Text: fmt.Sprintf(
										config.i18n("response-property-one-of-removed"),
										ColorizedValue(propertyDiff.OneOfDiff.DeletedSchemas.String()),
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
