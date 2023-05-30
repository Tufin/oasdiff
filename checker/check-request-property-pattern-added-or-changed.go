package checker

import (
	"fmt"

	"github.com/tufin/oasdiff/diff"
)

func RequestPropertyPatternAddedOrChangedCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config BackwardCompatibilityCheckConfig) []BackwardCompatibilityError {
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
								Id:          "request-property-pattern-added",
								Level:       WARN,
								Text:        fmt.Sprintf(config.i18n("request-property-pattern-added"), patternDiff.To, ColorizedValue(propertyFullName(propertyPath, propertyName))),
								Comment:     config.i18n("pattern-changed-warn-comment"),
								Operation:   operation,
								OperationId: operationItem.Revision.OperationID,
								Path:        path,
								Source:      source,
							})
						} else {
							result = append(result, BackwardCompatibilityError{
								Id:          "request-property-pattern-changed",
								Level:       WARN,
								Text:        fmt.Sprintf(config.i18n("request-property-pattern-changed"), ColorizedValue(propertyFullName(propertyPath, propertyName)), patternDiff.From, patternDiff.To),
								Comment:     config.i18n("pattern-changed-warn-comment"),
								Operation:   operation,
								OperationId: operationItem.Revision.OperationID,
								Path:        path,
								Source:      source,
							})
						}
					})
			}
		}
	}
	return result
}
