package checker

import (
	"fmt"

	"github.com/tufin/oasdiff/diff"
)

func RequestPropertyTypeChangedCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config Config) Changes {
	result := make(Changes, 0)
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

					if !typeDiff.Empty() || !formatDiff.Empty() {
						typeDiff, formatDiff = fillEmptyTypeAndFormatDiffs(typeDiff, schemaDiff, formatDiff)

						result = append(result, ApiChange{
							Id:          "request-body-type-changed",
							Level:       ConditionalError(breakingTypeFormatChangedInRequestProperty(typeDiff, formatDiff, mediaType, schemaDiff)),
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

						if !typeDiff.Empty() || !formatDiff.Empty() {
							typeDiff, formatDiff = fillEmptyTypeAndFormatDiffs(typeDiff, schemaDiff, formatDiff)
							result = append(result, ApiChange{
								Id:          "request-property-type-changed",
								Level:       ConditionalError(breakingTypeFormatChangedInRequestProperty(typeDiff, formatDiff, mediaType, schemaDiff)),
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

	if typeDiff != nil {
		return !isTypeContained(typeDiff.To, typeDiff.From, mediaType)
	}

	if formatDiff != nil {
		return !isFormatContained(schemaDiff.Revision.Value.Type, formatDiff.To, formatDiff.From)
	}

	return false
}
