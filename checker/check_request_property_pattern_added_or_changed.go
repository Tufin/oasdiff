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

						propName := propertyFullName(propertyPath, propertyName)

						if patternDiff.To == "" {
							result = append(result, NewApiChange(
								RequestPropertyPatternRemovedId,
								INFO,
								[]any{patternDiff.From, propName},
								"",
								operationsSources,
								operationItem.Revision,
								operation,
								path,
							))
						} else if patternDiff.From == "" {
							result = append(result, NewApiChange(
								RequestPropertyPatternAddedId,
								WARN,
								[]any{patternDiff.To, propName},
								PatternChangedCommentId,
								operationsSources,
								operationItem.Revision,
								operation,
								path,
							))
						} else {
							level := WARN
							comment := PatternChangedCommentId
							if patternDiff.To == ".*" {
								level = INFO
								comment = ""
							}
							result = append(result, NewApiChange(
								RequestPropertyPatternChangedId,
								level,
								[]any{propName, patternDiff.From, patternDiff.To},
								comment,
								operationsSources,
								operationItem.Revision,
								operation,
								path,
							))
						}
					})
			}
		}
	}
	return result
}
