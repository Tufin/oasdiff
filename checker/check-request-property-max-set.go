package checker

import (
	"fmt"

	"github.com/tufin/oasdiff/diff"
)

func RequestPropertyMaxSetCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap) []BackwardCompatibilityError {
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
					if maxDiff.From == nil &&
						maxDiff.To != nil {
						result = append(result, BackwardCompatibilityError{
							Id:        "request-body-max-set",
							Level:     WARN,
							Text:      fmt.Sprintf("the request's body max was set to '%s'", maxDiff.To),
							Comment:   "It is warn because sometimes it is required to be set. But good clients should be checked to support this restriction before such change in specification.",
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
						maxDiff := propertyDiff.MaxDiff
						if maxDiff == nil {
							return
						}
						if maxDiff.From != nil ||
							maxDiff.To == nil {
							return
						}
						if propertyDiff.Revision.Value.ReadOnly {
							return
						}

						result = append(result, BackwardCompatibilityError{
							Id:        "request-property-max-set",
							Level:     WARN,
							Text:      fmt.Sprintf("the %s request property's max was set to '%s'", ColorizedValue(propertyFullName(propertyPath, propertyName)), maxDiff.To),
							Comment:   "It is warn because sometimes it is required to be set. But good clients should be checked to support this restriction before such change in specification.",
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
