package checker

import (
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/tufin/oasdiff/diff"
	"github.com/tufin/oasdiff/load"
)

const (
	RequestParameterTypeChangedId = "request-parameter-type-changed"
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
			for paramLocation, paramDiffs := range operationItem.ParametersDiff.Modified {
				for paramName, paramDiff := range paramDiffs {
					if paramDiff.SchemaDiff == nil {
						continue
					}

					schemaDiff := paramDiff.SchemaDiff
					typeDiff := schemaDiff.TypeDiff
					formatDiff := schemaDiff.FormatDiff

					if typeDiff == nil && formatDiff == nil {
						continue
					}

					if typeAndFormatContained(typeDiff, formatDiff, paramDiff.Revision.Schema.Value.Type) {
						continue
					}

					source := (*operationsSources)[operationItem.Revision]

					typeDiff = getDetailedTypeDiff(schemaDiff)
					formatDiff = getDetailedFormatDiff(schemaDiff)

					result = append(result, ApiChange{
						Id:          RequestParameterTypeChangedId,
						Level:       ERR,
						Args:        []any{paramLocation, paramName, typeDiff.Deleted, formatDiff.From, typeDiff.Added, formatDiff.To},
						Operation:   operation,
						OperationId: operationItem.Revision.OperationID,
						Path:        path,
						Source:      load.NewSource(source),
					})
				}
			}
		}
	}
	return result
}

func typeAndFormatContained(typeDiff *diff.StringsDiff, formatDiff *diff.ValueDiff, revisionType *openapi3.Types) bool {

	if typeDiff != nil && typeDiff.Deleted.Is("integer") && typeDiff.Added.Is("number") {
		return true
	}

	if typeDiff != nil && typeDiff.Added.Is("string") {
		return true
	}

	if formatDiff != nil && (formatDiff.To == nil || formatDiff.To == "") {
		// TODO: is this correct?
		return true
	}

	if formatDiff != nil && revisionType.Is("string") &&
		(formatDiff.From == "date" && formatDiff.To == "date-time" ||
			formatDiff.From == "time" && formatDiff.To == "date-time") {
		return true
	}

	if formatDiff != nil && revisionType.Is("number") &&
		(formatDiff.From == "float" && formatDiff.To == "double") {
		return true
	}

	if formatDiff != nil && revisionType.Is("integer") &&
		(formatDiff.From == "int32" && formatDiff.To == "int64" ||
			formatDiff.From == "int32" && formatDiff.To == "bigint" ||
			formatDiff.From == "int64" && formatDiff.To == "bigint") {
		return true
	}

	return false
}
