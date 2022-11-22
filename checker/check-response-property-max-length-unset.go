package checker

import (
	"fmt"

	"github.com/tufin/oasdiff/diff"
)

func ResponsePropertyMaxLengthUnsetCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config BackwardCompatibilityCheckConfig) []BackwardCompatibilityError {
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
					if mediaTypeDiff.SchemaDiff != nil && mediaTypeDiff.SchemaDiff.MaxLengthDiff != nil {
						maxLengthDiff := mediaTypeDiff.SchemaDiff.MaxLengthDiff
						if maxLengthDiff.From != nil &&
							maxLengthDiff.To == nil {
							result = append(result, BackwardCompatibilityError{
								Id:        "response-body-max-length-unset",
								Level:     ERR,
								Text:      fmt.Sprintf("the response's body maxLength was unset from %s", ColorizedValue(maxLengthDiff.From)),
								Operation: operation,
								Path:      path,
								Source:    source,
								ToDo:      "Add to exceptions-list.md",
							})
						}
					}

					CheckModifiedPropertiesDiff(
						mediaTypeDiff.SchemaDiff,
						func(propertyPath string, propertyName string, propertyDiff *diff.SchemaDiff, parent *diff.SchemaDiff) {
							maxLengthDiff := propertyDiff.MaxLengthDiff
							if maxLengthDiff == nil {
								return
							}
							if maxLengthDiff.To != nil ||
								maxLengthDiff.From == nil {
								return
							}
							if propertyDiff.Revision.Value.WriteOnly {
								return
							}

							result = append(result, BackwardCompatibilityError{
								Id:        "response-property-max-length-unset",
								Level:     ERR,
								Text:      fmt.Sprintf("the %s response property's maxLength was unset from %s for the response status %s", ColorizedValue(propertyFullName(propertyPath, propertyName)), ColorizedValue(maxLengthDiff.From), ColorizedValue(responseStatus)),
								Operation: operation,
								Path:      path,
								Source:    source,
								ToDo:      "Add to exceptions-list.md",
							})
						})
				}

			}

		}
	}
	return result
}
