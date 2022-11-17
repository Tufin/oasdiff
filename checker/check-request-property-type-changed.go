package checker

import (
	"fmt"

	"github.com/tufin/oasdiff/diff"
)

func RequestPropertyTypeChangedCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config BackwardCompatibilityCheckConfig) []BackwardCompatibilityError {
	result := make([]BackwardCompatibilityError, 0)
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
			source := (*operationsSources)[operationItem.Revision]

			modifiedMediaTypes := operationItem.RequestBodyDiff.ContentDiff.MediaTypeModified
			for mediaType, mediaTypeDiff := range modifiedMediaTypes {
				if mediaTypeDiff.SchemaDiff != nil {
					typeDiff := mediaTypeDiff.SchemaDiff.TypeDiff
					formatDiff := mediaTypeDiff.SchemaDiff.FormatDiff
					if (typeDiff != nil || formatDiff != nil) && (typeDiff == nil || typeDiff != nil &&
						!(typeDiff.From == "integer" && typeDiff.To == "number") &&
						!(typeDiff.To == "string" && mediaType != "application/json" && mediaType != "application/xml")) &&
						(formatDiff == nil || formatDiff != nil && formatDiff.To != nil && formatDiff.To != "" &&
							!(mediaTypeDiff.SchemaDiff.Revision.Value.Type == "string" &&
								(formatDiff.From == "date" && formatDiff.To == "date-time" ||
									formatDiff.From == "time" && formatDiff.To == "date-time")) &&
							!(mediaTypeDiff.SchemaDiff.Revision.Value.Type == "number" &&
								(formatDiff.From == "float" && formatDiff.To == "double")) &&
							!(mediaTypeDiff.SchemaDiff.Revision.Value.Type == "string" &&
								(formatDiff.From == "decimal" && formatDiff.To == "uuid")) &&
							!(mediaTypeDiff.SchemaDiff.Revision.Value.Type == "integer" &&
								(formatDiff.From == "int32" && formatDiff.To == "int64" ||
									formatDiff.From == "int32" && formatDiff.To == "bigint" ||
									formatDiff.From == "int64" && formatDiff.To == "bigint"))) {
						if typeDiff == nil {
							typeDiff = &diff.ValueDiff{From: mediaTypeDiff.SchemaDiff.Revision.Value.Type, To: mediaTypeDiff.SchemaDiff.Revision.Value.Type}
						}
						if formatDiff == nil {
							formatDiff = &diff.ValueDiff{From: mediaTypeDiff.SchemaDiff.Revision.Value.Format, To: mediaTypeDiff.SchemaDiff.Revision.Value.Format}
						}
						result = append(result, BackwardCompatibilityError{
							Id:        "request-body-type-changed",
							Level:     ERR,
							Text:      fmt.Sprintf("the request's body type/format changed from %s/%s to %s/%s", empty2none(typeDiff.From), empty2none(formatDiff.From), empty2none(typeDiff.To), empty2none(formatDiff.To)),
							Operation: operation,
							Path:      path,
							Source:    source,
							ToDo:      "Add to exceptions-list.md",
						})
					}
				}

				CheckModifiedPropertiesDiff(
					mediaTypeDiff.SchemaDiff,
					func(propertyPath string, propertyName string, propertyDiff *diff.SchemaDiff, parent *diff.SchemaDiff) {
						if propertyDiff.Revision.Value.ReadOnly {
							return
						}
						typeDiff := propertyDiff.TypeDiff
						formatDiff := propertyDiff.FormatDiff

						if (typeDiff != nil || formatDiff != nil) && (typeDiff == nil || typeDiff != nil &&
							!(typeDiff.From == "integer" && typeDiff.To == "number") &&
							!(typeDiff.To == "string" && mediaType != "application/json" && mediaType != "application/xml")) &&
							(formatDiff == nil || formatDiff != nil && formatDiff.To != nil && formatDiff.To != "" &&
								!(propertyDiff.Revision.Value.Type == "string" &&
									(formatDiff.From == "date" && formatDiff.To == "date-time" ||
										formatDiff.From == "time" && formatDiff.To == "date-time")) &&
								!(propertyDiff.Revision.Value.Type == "number" &&
									(formatDiff.From == "float" && formatDiff.To == "double")) &&
								!(propertyDiff.Revision.Value.Type == "string" &&
									(formatDiff.From == "decimal" && formatDiff.To == "uuid")) &&
								!(propertyDiff.Revision.Value.Type == "integer" &&
									(formatDiff.From == "int32" && formatDiff.To == "int64" ||
										formatDiff.From == "int32" && formatDiff.To == "bigint" ||
										formatDiff.From == "int64" && formatDiff.To == "bigint"))) {
							if typeDiff == nil {
								typeDiff = &diff.ValueDiff{From: propertyDiff.Revision.Value.Type, To: propertyDiff.Revision.Value.Type}
							}
							if formatDiff == nil {
								formatDiff = &diff.ValueDiff{From: propertyDiff.Revision.Value.Format, To: propertyDiff.Revision.Value.Format}
							}
							result = append(result, BackwardCompatibilityError{
								Id:        "request-property-type-changed",
								Level:     ERR,
								Text:      fmt.Sprintf("the %s request property type/format changed from %s/%s to %s/%s", ColorizedValue(propertyFullName(propertyPath, propertyName)), empty2none(typeDiff.From), empty2none(formatDiff.From), empty2none(typeDiff.To), empty2none(formatDiff.To)),
								Operation: operation,
								Path:      path,
								Source:    source,
								ToDo:      "Add to exceptions-list.md",
							})
						}
					})
			}
		}
	}
	return result
}
