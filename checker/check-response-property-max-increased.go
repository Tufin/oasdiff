package checker

import (
	"fmt"

	"github.com/tufin/oasdiff/diff"
)

func ResponsePropertyMaxIncreasedCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config Config) Changes {
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
			for responseStatus, responseDiff := range operationItem.ResponsesDiff.Modified {
				if responseDiff == nil ||
					responseDiff.ContentDiff == nil ||
					responseDiff.ContentDiff.MediaTypeModified == nil {
					continue
				}
				modifiedMediaTypes := responseDiff.ContentDiff.MediaTypeModified
				for _, mediaTypeDiff := range modifiedMediaTypes {
					if mediaTypeDiff.SchemaDiff != nil && mediaTypeDiff.SchemaDiff.MaxDiff != nil {
						maxDiff := mediaTypeDiff.SchemaDiff.MaxDiff
						if maxDiff.From != nil &&
							maxDiff.To != nil {
							if isIncreasedValue(maxDiff) {
								result = append(result, ApiChange{
									Id:          "response-body-max-increased",
									Level:       ERR,
									Text:        fmt.Sprintf(config.i18n("response-body-max-increased"), colorizedValue(maxDiff.From), colorizedValue(maxDiff.To)),
									Operation:   operation,
									OperationId: operationItem.Revision.OperationID,
									Path:        path,
									Source:      source,
								})
							}
						}
					}

					checkModifiedPropertiesDiff(
						mediaTypeDiff.SchemaDiff,
						func(propertyPath string, propertyName string, propertyDiff *diff.SchemaDiff, parent *diff.SchemaDiff) {
							maxDiff := propertyDiff.MaxDiff
							if maxDiff == nil {
								return
							}
							if maxDiff.To == nil ||
								maxDiff.From == nil {
								return
							}
							if !isIncreasedValue(maxDiff) {
								return
							}

							if propertyDiff.Revision.Value.WriteOnly {
								return
							}

							result = append(result, ApiChange{
								Id:          "response-property-max-increased",
								Level:       ERR,
								Text:        fmt.Sprintf(config.i18n("response-property-max-increased"), colorizedValue(propertyFullName(propertyPath, propertyName)), colorizedValue(maxDiff.From), colorizedValue(maxDiff.To), colorizedValue(responseStatus)),
								Operation:   operation,
								OperationId: operationItem.Revision.OperationID,
								Path:        path,
								Source:      source,
							})
						})
				}

			}

		}
	}
	return result
}
