package checker

import (
	"fmt"
	"strings"

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
							if propertyDiff == nil || propertyDiff.Revision == nil || propertyDiff.Revision.Value == nil {
								return
							}

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

	if typeDiff != nil {
		return !isTypeOK(typeDiff, mediaType)
	}

	if formatDiff != nil {
		return !isFormatOK(schemaDiff, formatDiff)
	}

	return false
}

func isTypeOK(typeDiff *diff.ValueDiff, mediaType string) bool {
	return (typeDiff.From == "number" && typeDiff.To == "integer") ||
		(typeDiff.From == "string" && !isJsonMediaType(mediaType) && mediaType != "application/xml") // string can change to anything, unless it's json or xml
}

func isFormatOK(schemaDiff *diff.SchemaDiff, formatDiff *diff.ValueDiff) bool {

	switch schemaDiff.Revision.Value.Type {
	case "number":
		return formatDiff.From == "double" && formatDiff.To == "float"
	case "integer":
		return (formatDiff.From == "int64" && formatDiff.To == "int32") ||
			(formatDiff.From == "bigint" && formatDiff.To == "int32") ||
			(formatDiff.From == "bigint" && formatDiff.To == "int64")
	case "string":
		return (formatDiff.From == "date-time" && formatDiff.To == "date" ||
			formatDiff.From == "date-time" && formatDiff.To == "time")
	}

	return false
}

func isJsonMediaType(mediaType string) bool {
	return mediaType == "application/json" ||
		(strings.HasPrefix(mediaType, "application/vnd.") && strings.HasSuffix(mediaType, "+json"))
}
