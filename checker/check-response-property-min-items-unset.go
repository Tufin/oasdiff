package checker

import (
	"fmt"

	"github.com/tufin/oasdiff/diff"
)

func ResponsePropertyMinItemsUnsetCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap) []BackwardCompatibilityError {
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
							minItemsDiff.To == nil {
							result = append(result, BackwardCompatibilityError{
								Id:        "response-body-min-items-unset",
								Level:     ERR,
								Text:      fmt.Sprintf("the response's body minItems was unset from %s", minItemsDiff.From),
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
							minItemsDiff := propertyDiff.MinItemsDiff
							if minItemsDiff == nil {
								return
							}
							if minItemsDiff.To != nil ||
								minItemsDiff.From == nil {
								return
							}
							if propertyDiff.Revision.Value.WriteOnly {
								return
							}

							result = append(result, BackwardCompatibilityError{
								Id:        "response-property-min-items-unset",
								Level:     ERR,
								Text:      fmt.Sprintf("the %s response property's minItems was unset from %s for the response status %s", ColorizedValue(propertyFullName(propertyPath, propertyName)), minItemsDiff.From, ColorizedValue(responseStatus)),
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
