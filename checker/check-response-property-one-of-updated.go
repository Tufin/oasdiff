package checker

import (
	"fmt"

	"github.com/tufin/oasdiff/diff"
)

func ResponsePropertyOneOfUpdated(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config Config) Changes {
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
							Id:    "response-body-one-of-added",
							Level: INFO,
							Text: fmt.Sprintf(
								config.i18n("response-body-one-of-added"),
								colorizedValue(mediaTypeDiff.SchemaDiff.OneOfDiff.Added.String()),
								responseStatus),
							Operation:   operation,
							OperationId: operationItem.Revision.OperationID,
							Path:        path,
							Source:      source,
						})
					}

					if mediaTypeDiff.SchemaDiff.OneOfDiff != nil && len(mediaTypeDiff.SchemaDiff.OneOfDiff.Deleted) > 0 {
						result = append(result, ApiChange{
							Id:    "response-body-one-of-removed",
							Level: INFO,
							Text: fmt.Sprintf(
								config.i18n("response-body-one-of-removed"),
								colorizedValue(mediaTypeDiff.SchemaDiff.OneOfDiff.Deleted.String()),
								responseStatus),
							Operation:   operation,
							OperationId: operationItem.Revision.OperationID,
							Path:        path,
							Source:      source,
						})
					}

					checkModifiedPropertiesDiff(
						mediaTypeDiff.SchemaDiff,
						func(propertyPath string, propertyName string, propertyDiff *diff.SchemaDiff, parent *diff.SchemaDiff) {
							if propertyDiff.OneOfDiff == nil {
								return
							}

							if len(propertyDiff.OneOfDiff.Added) > 0 {
								result = append(result, ApiChange{
									Id:    "response-property-one-of-added",
									Level: INFO,
									Text: fmt.Sprintf(
										config.i18n("response-property-one-of-added"),
										colorizedValue(propertyDiff.OneOfDiff.Added.String()),
										colorizedValue(propertyFullName(propertyPath, propertyName)),
										responseStatus),
									Operation:   operation,
									OperationId: operationItem.Revision.OperationID,
									Path:        path,
									Source:      source,
								})
							}

							if len(propertyDiff.OneOfDiff.Deleted) > 0 {
								result = append(result, ApiChange{
									Id:    "response-property-one-of-removed",
									Level: INFO,
									Text: fmt.Sprintf(
										config.i18n("response-property-one-of-removed"),
										colorizedValue(propertyDiff.OneOfDiff.Deleted.String()),
										colorizedValue(propertyFullName(propertyPath, propertyName)),
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
