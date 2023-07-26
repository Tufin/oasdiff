package checker

import (
	"fmt"

	"github.com/tufin/oasdiff/diff"
)

func RequestPropertyPatternUpdatedCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config Config) Changes {
	result := make(Changes, 0)
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

						source := (*operationsSources)[operationItem.Revision]

						if patternDiff.To == "" {
							result = append(result, ApiChange{
								Id:          "request-property-pattern-removed",
								Level:       INFO,
								Text:        fmt.Sprintf(config.i18n("request-property-pattern-removed"), patternDiff.From, ColorizedValue(propertyFullName(propertyPath, propertyName))),
								Operation:   operation,
								OperationId: operationItem.Revision.OperationID,
								Path:        path,
								Source:      source,
							})
						} else if patternDiff.From == "" {
							result = append(result, ApiChange{
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
							level := WARN
							comment := config.i18n("pattern-changed-warn-comment")
							if patternDiff.To == ".*" {
								level = INFO
								comment = ""
							}
							result = append(result, ApiChange{
								Id:          "request-property-pattern-changed",
								Level:       level,
								Text:        fmt.Sprintf(config.i18n("request-property-pattern-changed"), ColorizedValue(propertyFullName(propertyPath, propertyName)), patternDiff.From, patternDiff.To),
								Comment:     comment,
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
