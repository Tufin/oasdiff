package checker

import (
	"fmt"

	"github.com/tufin/oasdiff/diff"
)

func RequestPropertyMinSetCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap) []BackwardCompatibilityError {
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
				if mediaTypeDiff.SchemaDiff != nil && mediaTypeDiff.SchemaDiff.MinDiff != nil {
					minDiff := mediaTypeDiff.SchemaDiff.MinDiff
					if minDiff.From == nil &&
						minDiff.To != nil {
						result = append(result, BackwardCompatibilityError{
							Id:        "request-body-min-set",
							Level:     WARN,
							Text:      fmt.Sprintf("the request's body min was set to '%s'", minDiff.To),
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
						minDiff := propertyDiff.MinDiff
						if minDiff == nil {
							return
						}
						if minDiff.From != nil ||
							minDiff.To == nil {
							return
						}
						if propertyDiff.Revision.Value.ReadOnly {
							return
						}

						result = append(result, BackwardCompatibilityError{
							Id:        "request-property-min-set",
							Level:     WARN,
							Text:      fmt.Sprintf("the %s request property's min was set to '%s'", ColorizedValue(propertyFullName(propertyPath, propertyName)), minDiff.To),
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
