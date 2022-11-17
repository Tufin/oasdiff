package checker

import (
	"fmt"

	"github.com/tufin/oasdiff/diff"
)

func ResponsePropertyTypeChangedCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config BackwardCompatibilityCheckConfig) []BackwardCompatibilityError {
	result := make([]BackwardCompatibilityError, 0)
	if diffReport.PathsDiff == nil {
		return result
	}
	for path, pathItem := range diffReport.PathsDiff.Modified {
		if pathItem.OperationsDiff == nil {
			continue
		}
		for operation, operationItem := range pathItem.OperationsDiff.Modified {
			source := (*operationsSources)[operationItem.Revision]
			if operationItem.ResponsesDiff == nil || operationItem.ResponsesDiff.Modified == nil {
				continue
			}

			for responseStatus, responseDiff := range operationItem.ResponsesDiff.Modified {
				if responseDiff.ContentDiff == nil ||
					responseDiff.ContentDiff.MediaTypeModified == nil {
					continue
				}

				modifiedMediaTypes := responseDiff.ContentDiff.MediaTypeModified
				for mediaType, mediaTypeDiff := range modifiedMediaTypes {
					if mediaTypeDiff.SchemaDiff != nil {
						typeDiff := mediaTypeDiff.SchemaDiff.TypeDiff
						formatDiff := mediaTypeDiff.SchemaDiff.FormatDiff
						if (typeDiff != nil || formatDiff != nil) && (typeDiff == nil || typeDiff != nil &&
							!(typeDiff.To == "integer" && typeDiff.From == "number") &&
							!(typeDiff.From == "string" && mediaType != "application/json" && mediaType != "application/xml")) &&
							(formatDiff == nil || formatDiff != nil && formatDiff.From != nil && formatDiff.From != "" &&
								!(mediaTypeDiff.SchemaDiff.Revision.Value.Type == "number" &&
									(formatDiff.To == "float" && formatDiff.From == "double")) &&
								!(mediaTypeDiff.SchemaDiff.Revision.Value.Type == "integer" &&
									(formatDiff.To == "int32" && formatDiff.From == "int64" ||
										formatDiff.To == "int32" && formatDiff.From == "bigint" ||
										formatDiff.To == "int64" && formatDiff.From == "bigint"))) {
							if typeDiff == nil {
								typeDiff = &diff.ValueDiff{From: mediaTypeDiff.SchemaDiff.Revision.Value.Type, To: mediaTypeDiff.SchemaDiff.Revision.Value.Type}
							}
							if formatDiff == nil {
								formatDiff = &diff.ValueDiff{From: mediaTypeDiff.SchemaDiff.Revision.Value.Format, To: mediaTypeDiff.SchemaDiff.Revision.Value.Format}
							}
							result = append(result, BackwardCompatibilityError{
								Id:        "response-body-type-changed",
								Level:     ERR,
								Text:      fmt.Sprintf("the response's body type/format changed from %s/%s to %s/%s for status %s", empty2none(typeDiff.From), empty2none(formatDiff.From), empty2none(typeDiff.To), empty2none(formatDiff.To), ColorizedValue(responseStatus)),
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
								!(typeDiff.To == "integer" && typeDiff.From == "number") &&
								!(typeDiff.From == "string" && mediaType != "application/json" && mediaType != "application/xml")) &&
								(formatDiff == nil || formatDiff != nil && formatDiff.From != nil && formatDiff.From != "" &&
									!(propertyDiff.Revision.Value.Type == "number" &&
										(formatDiff.To == "float" && formatDiff.From == "double")) &&
									!(propertyDiff.Revision.Value.Type == "integer" &&
										(formatDiff.To == "int32" && formatDiff.From == "int64" ||
											formatDiff.To == "int32" && formatDiff.From == "bigint" ||
											formatDiff.To == "int64" && formatDiff.From == "bigint"))) {
								if typeDiff == nil {
									typeDiff = &diff.ValueDiff{From: propertyDiff.Revision.Value.Type, To: propertyDiff.Revision.Value.Type}
								}
								if formatDiff == nil {
									formatDiff = &diff.ValueDiff{From: propertyDiff.Revision.Value.Format, To: propertyDiff.Revision.Value.Format}
								}
								result = append(result, BackwardCompatibilityError{
									Id:        "response-property-type-changed",
									Level:     ERR,
									Text:      fmt.Sprintf("the response's property type/format changed from %s/%s to %s/%s for status %s", empty2none(typeDiff.From), empty2none(formatDiff.From), empty2none(typeDiff.To), empty2none(formatDiff.To), ColorizedValue(responseStatus)),
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
	}
	return result
}
