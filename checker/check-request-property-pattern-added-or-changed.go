package checker

import (
	"github.com/tufin/oasdiff/diff"
)

const (
	RequestPropertyPatternRemovedId = "request-property-pattern-removed"
	RequestPropertyPatternAddedId   = "request-property-pattern-added"
	RequestPropertyPatternChangedId = "request-property-pattern-changed"
)

func RequestPropertyPatternUpdatedCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config *Config) Changes {
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
						propName := propertyFullName(propertyPath, propertyName)

						if patternDiff.To == "" {
							result = append(result, ApiChange{
								Id:          RequestPropertyPatternRemovedId,
								Level:       INFO,
								Args:        []any{patternDiff.From, propName},
								Operation:   operation,
								OperationId: operationItem.Revision.OperationID,
								Path:        path,
								Source:      source,
							})
						} else if patternDiff.From == "" {
							result = append(result, ApiChange{
								Id:          RequestPropertyPatternAddedId,
								Level:       WARN,
								Args:        []any{patternDiff.To, propName},
								Comment:     config.Localize("pattern-changed-warn-comment"),
								Operation:   operation,
								OperationId: operationItem.Revision.OperationID,
								Path:        path,
								Source:      source,
							})
						} else {
							level := WARN
							comment := config.Localize("pattern-changed-warn-comment")
							if patternDiff.To == ".*" {
								level = INFO
								comment = ""
							}
							result = append(result, ApiChange{
								Id:          RequestPropertyPatternChangedId,
								Level:       level,
								Args:        []any{propName, patternDiff.From, patternDiff.To},
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
