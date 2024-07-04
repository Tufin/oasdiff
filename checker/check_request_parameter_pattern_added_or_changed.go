package checker

import (
	"github.com/tufin/oasdiff/diff"
)

const (
	RequestParameterPatternAddedId   = "request-parameter-pattern-added"
	RequestParameterPatternRemovedId = "request-parameter-pattern-removed"
	RequestParameterPatternChangedId = "request-parameter-pattern-changed"
	PatternChangedCommentId          = "pattern-changed-warn-comment"
)

func RequestParameterPatternAddedOrChangedCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config *Config) Changes {
	result := make(Changes, 0)
	if diffReport.PathsDiff == nil {
		return result
	}
	for path, pathItem := range diffReport.PathsDiff.Modified {
		if pathItem.OperationsDiff == nil {
			continue
		}
		for operation, operationItem := range pathItem.OperationsDiff.Modified {
			if operationItem.ParametersDiff == nil {
				continue
			}
			if operationItem.ParametersDiff.Modified == nil {
				continue
			}
			for paramLocation, paramItems := range operationItem.ParametersDiff.Modified {
				for paramName, paramItem := range paramItems {
					if paramItem.SchemaDiff == nil {
						continue
					}
					patternDiff := paramItem.SchemaDiff.PatternDiff
					if patternDiff == nil {
						continue
					}

					if patternDiff.From == "" {
						result = append(result, NewApiChange(
							RequestParameterPatternAddedId,
							WARN,
							[]any{patternDiff.To, paramLocation, paramName},
							PatternChangedCommentId,
							operationsSources,
							operationItem.Revision,
							operation,
							path,
						))
					} else if patternDiff.To == "" {
						result = append(result, NewApiChange(
							RequestParameterPatternRemovedId,
							INFO,
							[]any{patternDiff.From, paramLocation, paramName},
							"",
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
							RequestParameterPatternChangedId,
							level,
							[]any{paramLocation, paramName, patternDiff.From, patternDiff.To},
							comment,
							operationsSources,
							operationItem.Revision,
							operation,
							path,
						))
					}
				}
			}
		}
	}
	return result
}
