package checker

import (
	"fmt"

	"github.com/tufin/oasdiff/diff"
)

func RequestParameterTypeChangedCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config BackwardCompatibilityCheckConfig) []BackwardCompatibilityError {
	result := make([]BackwardCompatibilityError, 0)
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
					typeDiff := paramDiff.SchemaDiff.TypeDiff
					formatDiff := paramDiff.SchemaDiff.FormatDiff
					if typeDiff == nil && formatDiff == nil {
						continue
					}

					if typeDiff != nil && typeDiff.From == "integer" && typeDiff.To == "number" {
						continue
					}

					if typeDiff != nil && typeDiff.To == "string" {
						// parameters can be changed to string anytime
						continue
					}

					if formatDiff != nil && (formatDiff.To == nil || formatDiff.To == "") {
						continue
					}

					if formatDiff != nil && paramDiff.Revision.Schema.Value.Type == "string" &&
						(formatDiff.From == "date" && formatDiff.To == "date-time" ||
							formatDiff.From == "time" && formatDiff.To == "date-time") {
						continue
					}

					if formatDiff != nil && paramDiff.Revision.Schema.Value.Type == "number" &&
						(formatDiff.From == "float" && formatDiff.To == "double") {
						continue
					}

					if formatDiff != nil && paramDiff.Revision.Schema.Value.Type == "integer" &&
						(formatDiff.From == "int32" && formatDiff.To == "int64" ||
							formatDiff.From == "int32" && formatDiff.To == "bigint" ||
							formatDiff.From == "int64" && formatDiff.To == "bigint") {
						continue
					}
					source := (*operationsSources)[operationItem.Revision]

					if typeDiff == nil {
						typeDiff = &diff.ValueDiff{From: paramDiff.Revision.Schema.Value.Type, To: paramDiff.Revision.Schema.Value.Type}
					}
					if formatDiff == nil {
						formatDiff = &diff.ValueDiff{From: paramDiff.Revision.Schema.Value.Format, To: paramDiff.Revision.Schema.Value.Format}
					}

					result = append(result, BackwardCompatibilityError{
						Id:          "request-parameter-type-changed",
						Level:       ERR,
						Text:        fmt.Sprintf(config.i18n("request-parameter-type-changed"), ColorizedValue(paramLocation), ColorizedValue(paramName), empty2none(typeDiff.From), empty2none(formatDiff.From), empty2none(typeDiff.To), empty2none(formatDiff.To)),
						Operation:   operation,
						OperationId: operationItem.Revision.OperationID,
						Path:        path,
						Source:      source,
					})
				}
			}
		}
	}
	return result
}
