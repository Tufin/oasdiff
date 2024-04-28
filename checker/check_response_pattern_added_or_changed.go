package checker

import (
	"github.com/tufin/oasdiff/diff"
	"github.com/tufin/oasdiff/load"
)

const (
	ResponsePropertyPatternAddedId   = "response-property-pattern-added"
	ResponsePropertyPatternChangedId = "response-property-pattern-changed"
	ResponsePropertyPatternRemovedId = "response-property-pattern-removed"
)

func ResponsePatternAddedOrChangedCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config *Config) Changes {
	result := make(Changes, 0)
	if diffReport.PathsDiff == nil {
		return result
	}
	for path, pathItem := range diffReport.PathsDiff.Modified {
		if pathItem.OperationsDiff == nil {
			continue
		}
		for operation, operationItem := range pathItem.OperationsDiff.Modified {
			source := (*operationsSources)[operationItem.Revision]

			if operationItem.ResponsesDiff == nil {
				continue
			}

			for responseStatus, responseDiff := range operationItem.ResponsesDiff.Modified {
				if responseDiff.ContentDiff == nil ||
					responseDiff.ContentDiff.MediaTypeModified == nil {
					continue
				}

				modifiedMediaTypes := responseDiff.ContentDiff.MediaTypeModified
				for _, mediaTypeDiff := range modifiedMediaTypes {
					if mediaTypeDiff.SchemaDiff == nil {
						continue
					}

					CheckModifiedPropertiesDiff(
						mediaTypeDiff.SchemaDiff,
						func(propertyPath string, propertyName string, propertyDiff *diff.SchemaDiff, parent *diff.SchemaDiff) {
							patternDiff := propertyDiff.PatternDiff
							if patternDiff == nil {
								return
							}

							propName := propertyFullName(propertyPath, propertyName)

							id := ResponsePropertyPatternChangedId
							args := []any{propName, patternDiff.From, patternDiff.To, responseStatus}
							if patternDiff.To == "" || patternDiff.To == nil {
								id = ResponsePropertyPatternRemovedId
								args = []any{propName, patternDiff.From, responseStatus}
							} else if patternDiff.From == "" || patternDiff.From == nil {
								id = ResponsePropertyPatternAddedId
								args = []any{propName, patternDiff.To, responseStatus}
							}

							result = append(result, ApiChange{
								Id:          id,
								Level:       INFO,
								Args:        args,
								Operation:   operation,
								OperationId: operationItem.Revision.OperationID,
								Path:        path,
								Source:      load.NewSource(source),
							})
						})
				}
			}
		}
	}
	return result
}
