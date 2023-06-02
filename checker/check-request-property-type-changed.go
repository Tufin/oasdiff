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
					schemaDiff := mediaTypeDiff.SchemaDiff
					typeDiff := schemaDiff.TypeDiff
					formatDiff := schemaDiff.FormatDiff
					if breakingTypeFormatChangedInRequestProperty(typeDiff, formatDiff, mediaType, schemaDiff) {
						typeDiff, formatDiff = fillEmptyTypeAndFormatDiffs(typeDiff, schemaDiff, formatDiff)
						result = append(result, BackwardCompatibilityError{
							Id:          "request-body-type-changed",
							Level:       ERR,
							Text:        fmt.Sprintf(config.i18n("request-body-type-changed"), empty2none(typeDiff.From), empty2none(formatDiff.From), empty2none(typeDiff.To), empty2none(formatDiff.To)),
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
						if breakingTypeFormatChangedInRequestProperty(typeDiff, formatDiff, mediaType, schemaDiff) {
							typeDiff, formatDiff = fillEmptyTypeAndFormatDiffs(typeDiff, schemaDiff, formatDiff)
							result = append(result, BackwardCompatibilityError{
								Id:          "request-property-type-changed",
								Level:       ERR,
								Text:        fmt.Sprintf(config.i18n("request-property-type-changed"), ColorizedValue(propertyFullName(propertyPath, propertyName)), empty2none(typeDiff.From), empty2none(formatDiff.From), empty2none(typeDiff.To), empty2none(formatDiff.To)),
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
	return result
}

func fillEmptyTypeAndFormatDiffs(typeDiff *diff.ValueDiff, schemaDiff *diff.SchemaDiff, formatDiff *diff.ValueDiff) (*diff.ValueDiff, *diff.ValueDiff) {
	if typeDiff == nil {
		typeDiff = &diff.ValueDiff{From: schemaDiff.Revision.Value.Type, To: schemaDiff.Revision.Value.Type}
	}
	if formatDiff == nil {
		formatDiff = &diff.ValueDiff{From: schemaDiff.Revision.Value.Format, To: schemaDiff.Revision.Value.Format}
	}
	return typeDiff, formatDiff
}

func breakingTypeFormatChangedInRequestProperty(typeDiff *diff.ValueDiff, formatDiff *diff.ValueDiff, mediaType string, schemaDiff *diff.SchemaDiff) bool {
	return (typeDiff != nil || formatDiff != nil) && (typeDiff == nil || typeDiff != nil &&
		!(typeDiff.From == "integer" && typeDiff.To == "number") &&
		!(typeDiff.To == "string" && mediaType != "application/json" && mediaType != "application/xml")) &&
		(formatDiff == nil || formatDiff != nil && formatDiff.To != nil && formatDiff.To != "" &&
			!(schemaDiff.Revision.Value.Type == "string" &&
				(formatDiff.From == "date" && formatDiff.To == "date-time" ||
					formatDiff.From == "time" && formatDiff.To == "date-time")) &&
			!(schemaDiff.Revision.Value.Type == "number" &&
				(formatDiff.From == "float" && formatDiff.To == "double")) &&
			!(schemaDiff.Revision.Value.Type == "integer" &&
				(formatDiff.From == "int32" && formatDiff.To == "int64" ||
					formatDiff.From == "int32" && formatDiff.To == "bigint" ||
					formatDiff.From == "int64" && formatDiff.To == "bigint")))
}
