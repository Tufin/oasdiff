package checker

import (
	"fmt"

	"github.com/tufin/oasdiff/diff"
)

func RequestPropertyMaxDecreasedCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap) []BackwardCompatibilityError {
	result := make([]BackwardCompatibilityError, 0)
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
				if mediaTypeDiff.SchemaDiff != nil && mediaTypeDiff.SchemaDiff.MaxDiff != nil {
					maxDiff := mediaTypeDiff.SchemaDiff.MaxDiff
					if maxDiff.From != nil &&
						maxDiff.To != nil {
						if IsDecreasedValue(maxDiff) {
							result = append(result, BackwardCompatibilityError{
								Id:        "request-body-max-decreased",
								Level:     ERR,
								Text:      fmt.Sprintf("the request's body max was decreased to %s", maxDiff.To),
								Operation: operation,
								Path:      path,
								Source:    source,
								ToDo:      "Add to exceptions-list.md",
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
						if maxDiff.From == nil ||
							maxDiff.To == nil {
							return
						}
						if propertyDiff.Revision.Value.ReadOnly {
							return
						}
						if !IsDecreasedValue(maxDiff) {
							return
						}

						result = append(result, BackwardCompatibilityError{
							Id:        "request-property-max-decreased",
							Level:     ERR,
							Text:      fmt.Sprintf("the %s request property's max was decreased to '%s'", ColorizedValue(propertyFullName(propertyPath, propertyName)), maxDiff.To),
							Operation: operation,
							Path:      path,
							Source:    source,
							ToDo:      "Add to exceptions-list.md",
						})
					})
			}
		}
	}
	return result
}
