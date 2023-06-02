package checker

import (
	"fmt"

	"github.com/tufin/oasdiff/diff"
)

func ResponsePropertyMaxIncreasedCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config BackwardCompatibilityCheckConfig) []BackwardCompatibilityError {
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
					if mediaTypeDiff.SchemaDiff != nil && mediaTypeDiff.SchemaDiff.MaxDiff != nil {
						maxDiff := mediaTypeDiff.SchemaDiff.MaxDiff
						if maxDiff.From != nil &&
							maxDiff.To != nil {
							if IsIncreasedValue(maxDiff) {
								result = append(result, BackwardCompatibilityError{
									Id:          "response-body-max-increased",
									Level:       ERR,
									Text:        fmt.Sprintf(config.i18n("response-body-max-increased"), ColorizedValue(maxDiff.From), ColorizedValue(maxDiff.To)),
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
							maxDiff := propertyDiff.MaxDiff
							if maxDiff == nil {
								return
							}
							if maxDiff.To == nil ||
								maxDiff.From == nil {
								return
							}
							if !IsIncreasedValue(maxDiff) {
								return
							}

							if propertyDiff.Revision.Value.WriteOnly {
								return
							}

							result = append(result, BackwardCompatibilityError{
								Id:          "response-property-max-increased",
								Level:       ERR,
								Text:        fmt.Sprintf(config.i18n("response-property-max-increased"), ColorizedValue(propertyFullName(propertyPath, propertyName)), ColorizedValue(maxDiff.From), ColorizedValue(maxDiff.To), ColorizedValue(responseStatus)),
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
