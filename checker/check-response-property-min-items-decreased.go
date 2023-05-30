package checker

import (
	"fmt"

	"github.com/tufin/oasdiff/diff"
)

func ResponsePropertyMinItemsDecreasedCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config BackwardCompatibilityCheckConfig) []BackwardCompatibilityError {
	result := make([]BackwardCompatibilityError, 0)
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
					if mediaTypeDiff.SchemaDiff != nil && mediaTypeDiff.SchemaDiff.MinItemsDiff != nil {
						minItemsDiff := mediaTypeDiff.SchemaDiff.MinItemsDiff
						if minItemsDiff.From != nil &&
							minItemsDiff.To != nil {
							if IsDecreasedValue(minItemsDiff) {
								result = append(result, BackwardCompatibilityError{
									Id:          "response-body-min-items-decreased",
									Level:       ERR,
									Text:        fmt.Sprintf(config.i18n("response-body-min-items-decreased"), ColorizedValue(minItemsDiff.From), ColorizedValue(minItemsDiff.To)),
									Operation:   operation,
									OperationId: operationItem.Revision.OperationID,
									Path:        path,
									Source:      source,
								})
							}
						}
					}

					CheckModifiedPropertiesDiff(
						mediaTypeDiff.SchemaDiff,
						func(propertyPath string, propertyName string, propertyDiff *diff.SchemaDiff, parent *diff.SchemaDiff) {
							minItemsDiff := propertyDiff.MinItemsDiff
							if minItemsDiff == nil {
								return
							}
							if minItemsDiff.To == nil ||
								minItemsDiff.From == nil {
								return
							}
							if !IsDecreasedValue(minItemsDiff) {
								return
							}

							if propertyDiff.Revision.Value.WriteOnly {
								return
							}

							result = append(result, BackwardCompatibilityError{
								Id:          "response-property-min-items-decreased",
								Level:       ERR,
								Text:        fmt.Sprintf(config.i18n("response-property-min-items-decreased"), ColorizedValue(propertyFullName(propertyPath, propertyName)), ColorizedValue(minItemsDiff.From), ColorizedValue(minItemsDiff.To), ColorizedValue(responseStatus)),
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
