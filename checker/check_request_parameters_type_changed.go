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

								level, comment := checkRequestParameterPropertyTypeChanged(typeDiff, formatDiff, schemaDiff)

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

/*
checkRequestParameterPropertyTypeChanged checks the level of the change in the request parameter property type
Explanation:
Objects can be passed in the request parameters, for example, the following calls are equivalent:
PHP style: GET http://localhost:8080/api/tickets?params[id]=123&params[color]=green
JSON: GET http://localhost:8080/api/tickets?params={"id":"123","color":"green"}

The "params" object has two properties: "id" and "color", both with type "string", but note that the "id" values are actually numbers.
Imagine that the OpenAPI type of property "id" was changed from "number" to "string".
In the first example, the change is non-breaking, because the PHP format for numbers and strings is the same.
But in the second example, the change is breaking, because the JSON format requires quotes for strings.
*/
func checkRequestParameterPropertyTypeChanged(typeDiff *diff.StringsDiff, formatDiff *diff.ValueDiff, schemaDiff *diff.SchemaDiff) (Level, string) {

	// try with JSON format
	isBreakingAsJson := breakingTypeFormatChangedInRequestProperty(typeDiff, formatDiff, "application/json", schemaDiff)

	// try with non-JSON format
	isBreakingAsNonJson := breakingTypeFormatChangedInRequestProperty(typeDiff, formatDiff, "", schemaDiff)

	// if the JSON and not breaking as non-JSON formats don't agree, it's a warning
	if isBreakingAsJson != isBreakingAsNonJson {
		return WARN, RequestParameterPropertyTypeChangedCommentId
	}

	// if both are breaking it's an error
	if isBreakingAsJson {
		return ERR, ""
	}

	// if niether are breaking it's an informational change
	return INFO, ""
}
