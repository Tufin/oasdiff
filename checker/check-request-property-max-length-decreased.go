package checker

import (
	"fmt"

	"github.com/tufin/oasdiff/diff"
)

func RequestPropertyMaxLengthDecreasedCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap) []BackwardCompatibilityError {
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
				if mediaTypeDiff.SchemaDiff != nil && mediaTypeDiff.SchemaDiff.MaxLengthDiff != nil {
					maxLengthDiff := mediaTypeDiff.SchemaDiff.MaxLengthDiff
					if maxLengthDiff.From != nil &&
						maxLengthDiff.To != nil {
						if IsDecreasedValue(maxLengthDiff) {
							result = append(result, BackwardCompatibilityError{
								Id:        "request-body-max-length-decreased",
								Level:     ERR,
								Text:      fmt.Sprintf("the request's body maxLength was decreased to '%s'", maxLengthDiff.To),
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
						maxLengthDiff := propertyDiff.MaxLengthDiff
						if maxLengthDiff == nil {
							return
						}
						if maxLengthDiff.From == nil ||
							maxLengthDiff.To == nil {
							return
						}
						if propertyDiff.Revision.Value.ReadOnly {
							return
						}
						if !IsDecreasedValue(maxLengthDiff) {
							return
						}

						result = append(result, BackwardCompatibilityError{
							Id:        "request-property-max-length-decreased",
							Level:     ERR,
							Text:      fmt.Sprintf("the %s request property's maxLength was decreased to %s", ColorizedValue(propertyFullName(propertyPath, propertyName)), maxLengthDiff.To),
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
