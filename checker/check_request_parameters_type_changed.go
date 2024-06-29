package checker

import (
	"github.com/tufin/oasdiff/diff"
	"github.com/tufin/oasdiff/load"
)

const (
	RequestParameterTypeChangedId                = "request-parameter-type-changed"
	RequestParameterPropertyTypeChangedId        = "request-parameter-property-type-changed"
	RequestParameterPropertyTypeChangedCommentId = "request-parameter-property-type-changed-warn-comment"
)

func RequestParameterTypeChangedCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config *Config) Changes {
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
			source := (*operationsSources)[operationItem.Revision]

			for paramLocation, paramDiffs := range operationItem.ParametersDiff.Modified {
				for paramName, paramDiff := range paramDiffs {
					if paramDiff.SchemaDiff == nil {
						continue
					}

					schemaDiff := paramDiff.SchemaDiff
					typeDiff := schemaDiff.TypeDiff
					formatDiff := schemaDiff.FormatDiff

					if !typeDiff.Empty() || !formatDiff.Empty() {

						result = append(result, ApiChange{
							Id:          RequestParameterTypeChangedId,
							Level:       conditionalError(breakingTypeFormatChangedInRequestParam(typeDiff, formatDiff, schemaDiff), INFO),
							Args:        []any{paramLocation, paramName, getBaseType(schemaDiff), getBaseFormat(schemaDiff), getRevisionType(schemaDiff), getRevisionFormat(schemaDiff)},
							Operation:   operation,
							OperationId: operationItem.Revision.OperationID,
							Path:        path,
							Source:      load.NewSource(source),
						})
					}

					CheckModifiedPropertiesDiff(
						schemaDiff,
						func(propertyPath string, propertyName string, propertyDiff *diff.SchemaDiff, parent *diff.SchemaDiff) {

							schemaDiff := propertyDiff
							typeDiff := schemaDiff.TypeDiff
							formatDiff := schemaDiff.FormatDiff

							if !typeDiff.Empty() || !formatDiff.Empty() {

								level := WARN
								comment := RequestParameterPropertyTypeChangedCommentId
								if isWarn := breakingTypeFormatChangedInRequestParam(typeDiff, formatDiff, schemaDiff); isWarn {
									level = ERR
									comment = ""
								}

								result = append(result, ApiChange{
									Id:          RequestParameterPropertyTypeChangedId,
									Level:       level,
									Args:        []any{paramLocation, paramName, propertyFullName(propertyPath, propertyName), getBaseType(schemaDiff), getBaseFormat(schemaDiff), getRevisionType(schemaDiff), getRevisionFormat(schemaDiff)},
									Comment:     comment,
									Operation:   operation,
									OperationId: operationItem.Revision.OperationID,
									Path:        path,
									Source:      load.NewSource(source),
								})
							}
						})
				}
			}
		}
	}
	return result
}
