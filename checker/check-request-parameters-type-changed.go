package checker

import (
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
					typeDiff := paramDiff.SchemaDiff.TypeDiff
					formatDiff := paramDiff.SchemaDiff.FormatDiff
					if typeDiff == nil && formatDiff == nil {
						continue
					}

					if typeDiff != nil && typeDiff.Deleted.Is("integer") && typeDiff.Added.Is("number") {
						continue
					}

					if typeDiff != nil && typeDiff.Added.Is("string") {
						// parameters can be changed to string anytime
						continue
					}

					if formatDiff != nil && (formatDiff.To == nil || formatDiff.To == "") {
						continue
					}

					if formatDiff != nil && paramDiff.Revision.Schema.Value.Type.Is("string") &&
						(formatDiff.From == "date" && formatDiff.To == "date-time" ||
							formatDiff.From == "time" && formatDiff.To == "date-time") {
						continue
					}

					if formatDiff != nil && paramDiff.Revision.Schema.Value.Type.Is("number") &&
						(formatDiff.From == "float" && formatDiff.To == "double") {
						continue
					}

					if formatDiff != nil && paramDiff.Revision.Schema.Value.Type.Is("integer") &&
						(formatDiff.From == "int32" && formatDiff.To == "int64" ||
							formatDiff.From == "int32" && formatDiff.To == "bigint" ||
							formatDiff.From == "int64" && formatDiff.To == "bigint") {
						continue
					}
					source := (*operationsSources)[operationItem.Revision]

					if typeDiff == nil {
						typeDiff = &diff.StringsDiff{
							Deleted: paramDiff.Base.Schema.Value.Type.Slice(), // TODO: check bug fix
							Added:   paramDiff.Revision.Schema.Value.Type.Slice(),
						}
					}
					if formatDiff == nil {
						formatDiff = &diff.ValueDiff{From: paramDiff.Revision.Schema.Value.Format, To: paramDiff.Revision.Schema.Value.Format}
					}

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
