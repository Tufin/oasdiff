package checker

import (
	"fmt"

	"github.com/tufin/oasdiff/diff"
)

func ResponsePropertyAllOfUpdated(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config Config) Changes {
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

					if mediaTypeDiff.SchemaDiff.AllOfDiff != nil && mediaTypeDiff.SchemaDiff.AllOfDiff.Added > 0 {
						result = append(result, ApiChange{
							Id:    "response-body-all-of-added",
							Level: INFO,
							Text: fmt.Sprintf(
								config.i18n("response-body-all-of-added"),
								ColorizedValue(mediaTypeDiff.SchemaDiff.AllOfDiff.AddedSchemas.String()),
								responseStatus),
							Operation:   operation,
							OperationId: operationItem.Revision.OperationID,
							Path:        path,
							Source:      source,
						})
					}

					if mediaTypeDiff.SchemaDiff.AllOfDiff != nil && mediaTypeDiff.SchemaDiff.AllOfDiff.Deleted > 0 {
						result = append(result, ApiChange{
							Id:    "response-body-all-of-removed",
							Level: INFO,
							Text: fmt.Sprintf(
								config.i18n("response-body-all-of-removed"),
								ColorizedValue(mediaTypeDiff.SchemaDiff.AllOfDiff.DeletedSchemas.String()),
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
							if propertyDiff.AllOfDiff == nil {
								return
							}

							if propertyDiff.AllOfDiff.Added > 0 {
								result = append(result, ApiChange{
									Id:    "response-property-all-of-added",
									Level: INFO,
									Text: fmt.Sprintf(
										config.i18n("response-property-all-of-added"),
										ColorizedValue(propertyDiff.AllOfDiff.AddedSchemas.String()),
										ColorizedValue(propertyFullName(propertyPath, propertyName)),
										responseStatus),
									Operation:   operation,
									OperationId: operationItem.Revision.OperationID,
									Path:        path,
									Source:      source,
								})
							}

							if propertyDiff.AllOfDiff.Deleted > 0 {
								result = append(result, ApiChange{
									Id:    "response-property-all-of-removed",
									Level: INFO,
									Text: fmt.Sprintf(
										config.i18n("response-property-all-of-removed"),
										ColorizedValue(propertyDiff.AllOfDiff.DeletedSchemas.String()),
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
