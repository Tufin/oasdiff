package checker

import (
	"fmt"

	"github.com/tufin/oasdiff/diff"
)

func ResponsePropertyMaxLengthIncreasedCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config Config) Changes {
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
					if mediaTypeDiff.SchemaDiff != nil && mediaTypeDiff.SchemaDiff.MaxLengthDiff != nil {
						maxLengthDiff := mediaTypeDiff.SchemaDiff.MaxLengthDiff
						if maxLengthDiff.From != nil &&
							maxLengthDiff.To != nil {
							if isIncreasedValue(maxLengthDiff) {
								result = append(result, ApiChange{
									Id:          "response-body-max-length-increased",
									Level:       ERR,
									Text:        fmt.Sprintf(config.i18n("response-body-max-length-increased"), colorizedValue(maxLengthDiff.From), colorizedValue(maxLengthDiff.To)),
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
							maxLengthDiff := propertyDiff.MaxLengthDiff
							if maxLengthDiff == nil {
								return
							}
							if maxLengthDiff.To == nil ||
								maxLengthDiff.From == nil {
								return
							}
							if !isIncreasedValue(maxLengthDiff) {
								return
							}

							if propertyDiff.Revision.Value.WriteOnly {
								return
							}

							result = append(result, ApiChange{
								Id:          "response-property-max-length-increased",
								Level:       ERR,
								Text:        fmt.Sprintf(config.i18n("response-property-max-length-increased"), colorizedValue(propertyFullName(propertyPath, propertyName)), colorizedValue(maxLengthDiff.From), colorizedValue(maxLengthDiff.To), colorizedValue(responseStatus)),
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
