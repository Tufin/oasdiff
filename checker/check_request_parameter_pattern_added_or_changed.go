package checker

import (
	"github.com/tufin/oasdiff/diff"
	"github.com/tufin/oasdiff/load"
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
					source := (*operationsSources)[operationItem.Revision]

					if patternDiff.From == "" {
						result = append(result, ApiChange{
							Id:          RequestParameterPatternAddedId,
							Level:       WARN,
							Args:        []any{patternDiff.To, paramLocation, paramName},
							Comment:     PatternChangedCommentId,
							Operation:   operation,
							OperationId: operationItem.Revision.OperationID,
							Path:        path,
							Source:      load.NewSource(source),
						})
					} else if patternDiff.To == "" {
						result = append(result, ApiChange{
							Id:          RequestParameterPatternRemovedId,
							Level:       INFO,
							Args:        []any{patternDiff.From, paramLocation, paramName},
							Operation:   operation,
							OperationId: operationItem.Revision.OperationID,
							Path:        path,
							Source:      load.NewSource(source),
						})
					} else {
						level := WARN
						comment := PatternChangedCommentId
						if patternDiff.To == ".*" {
							level = INFO
							comment = ""
						}
						result = append(result, ApiChange{
							Id:          RequestParameterPatternChangedId,
							Level:       level,
							Args:        []any{paramLocation, paramName, patternDiff.From, patternDiff.To},
							Comment:     comment,
							Operation:   operation,
							OperationId: operationItem.Revision.OperationID,
							Path:        path,
							Source:      load.NewSource(source),
						})
					}
				}
			}
		}
	}
	return result
}
