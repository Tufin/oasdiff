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
						schemaDiff := mediaTypeDiff.SchemaDiff
						typeDiff := schemaDiff.TypeDiff
						formatDiff := schemaDiff.FormatDiff
						if breakingTypeFormatChangedInResponseProperty(typeDiff, formatDiff, mediaType, schemaDiff) {
							typeDiff, formatDiff = fillEmptyTypeAndFormatDiffs(typeDiff, schemaDiff, formatDiff)
							result = append(result, BackwardCompatibilityError{
								Id:          "response-body-type-changed",
								Level:       ERR,
								Text:        fmt.Sprintf(config.i18n("response-body-type-changed"), empty2none(typeDiff.From), empty2none(formatDiff.From), empty2none(typeDiff.To), empty2none(formatDiff.To), ColorizedValue(responseStatus)),
								Operation:   operation,
								OperationId: operationItem.Revision.OperationID,
								Path:        path,
								Source:      source,
							})
						}
					}

					CheckModifiedPropertiesDiff(
						mediaTypeDiff.SchemaDiff,
						func(propertyPath string, propertyName string, propertyDiff *diff.SchemaDiff, parent *diff.SchemaDiff) {
							if propertyDiff.Revision.Value.ReadOnly {
								return
							}
							schemaDiff := propertyDiff
							typeDiff := schemaDiff.TypeDiff
							formatDiff := schemaDiff.FormatDiff

							if breakingTypeFormatChangedInResponseProperty(typeDiff, formatDiff, mediaType, schemaDiff) {
								typeDiff, formatDiff = fillEmptyTypeAndFormatDiffs(typeDiff, schemaDiff, formatDiff)
								result = append(result, BackwardCompatibilityError{
									Id:          "response-property-type-changed",
									Level:       ERR,
									Text:        fmt.Sprintf(config.i18n("response-property-type-changed"), empty2none(typeDiff.From), empty2none(formatDiff.From), empty2none(typeDiff.To), empty2none(formatDiff.To), ColorizedValue(responseStatus)),
									Operation:   operation,
									OperationId: operationItem.Revision.OperationID,
									Path:        path,
									Source:      source,
								})
							}
						})
				}
			}
		}
	}
	return result
}

func breakingTypeFormatChangedInResponseProperty(typeDiff *diff.ValueDiff, formatDiff *diff.ValueDiff, mediaType string, schemaDiff *diff.SchemaDiff) bool {
	return (typeDiff != nil || formatDiff != nil) && (typeDiff == nil || typeDiff != nil &&
		!(typeDiff.To == "integer" && typeDiff.From == "number") &&
		!(typeDiff.From == "string" && mediaType != "application/json" && mediaType != "application/xml")) &&
		(formatDiff == nil || formatDiff != nil && formatDiff.From != nil && formatDiff.From != "" &&
			!(schemaDiff.Revision.Value.Type == "number" &&
				(formatDiff.To == "float" && formatDiff.From == "double")) &&
			!(schemaDiff.Revision.Value.Type == "integer" &&
				(formatDiff.To == "int32" && formatDiff.From == "int64" ||
					formatDiff.To == "int32" && formatDiff.From == "bigint" ||
					formatDiff.To == "int64" && formatDiff.From == "bigint")))
}
