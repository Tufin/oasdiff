package checker

import (
	"fmt"

	"github.com/tufin/oasdiff/diff"
)

const PatternChangedWarnComment = "It is the warning because it is difficult to automatically analyze if the new pattern is a superset of the previous pattern(e.g. changed from '[0-9]+' to '[0-9]*')"

func RequestPropertyPatternAddedOrChangedCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap) []BackwardCompatibilityError {
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
			modifiedMediaTypes := operationItem.RequestBodyDiff.ContentDiff.MediaTypeModified
			for _, mediaTypeDiff := range modifiedMediaTypes {
				CheckModifiedPropertiesDiff(
					mediaTypeDiff.SchemaDiff,
					func(propertyPath string, propertyName string, propertyDiff *diff.SchemaDiff, parent *diff.SchemaDiff) {
						patternDiff := propertyDiff.PatternDiff
						if patternDiff == nil {
							return
						}
						if patternDiff.To == "" ||
							patternDiff.To == ".*" {
							return
						}

						source := (*operationsSources)[operationItem.Revision]

						if patternDiff.From == "" {
							result = append(result, BackwardCompatibilityError{
								Id:        "request-property-pattern-added",
								Level:     WARN,
								Text:      fmt.Sprintf("added the pattern '%s' for the request property %s", patternDiff.To, ColorizedValue(propertyFullName(propertyPath, propertyName))),
								Comment:   PatternChangedWarnComment,
								Operation: operation,
								Path:      path,
								Source:    source,
								ToDo:      "Add to exceptions-list.md",
							})
						} else {
							result = append(result, BackwardCompatibilityError{
								Id:        "request-property-pattern-changed",
								Level:     WARN,
								Text:      fmt.Sprintf("changed the pattern for the request property %s from '%s' to '%s'", ColorizedValue(propertyFullName(propertyPath, propertyName)), patternDiff.From, patternDiff.To),
								Comment:   PatternChangedWarnComment,
								Operation: operation,
								Path:      path,
								Source:    source,
								ToDo:      "Add to exceptions-list.md",
							})
						}
					})
			}
		}
	}
	return result
}
