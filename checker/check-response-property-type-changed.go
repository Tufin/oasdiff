package checker

import (
	"strings"

	"github.com/tufin/oasdiff/diff"
	"github.com/tufin/oasdiff/load"
)

const (
	ResponseBodyTypeChangedId     = "response-body-type-changed"
	ResponsePropertyTypeChangedId = "response-property-type-changed"
)

func ResponsePropertyTypeChangedCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config *Config) Changes {
	result := make(Changes, 0)
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
							result = append(result, ApiChange{
								Id:          ResponseBodyTypeChangedId,
								Level:       ERR,
								Args:        []any{typeDiff.From, formatDiff.From, typeDiff.To, formatDiff.To, responseStatus},
								Operation:   operation,
								OperationId: operationItem.Revision.OperationID,
								Path:        path,
								Source:      load.NewSource(source),
							})
						}
					}

					CheckModifiedPropertiesDiff(
						mediaTypeDiff.SchemaDiff,
						func(propertyPath string, propertyName string, propertyDiff *diff.SchemaDiff, parent *diff.SchemaDiff) {
							if propertyDiff == nil || propertyDiff.Revision == nil {
								return
							}

							schemaDiff := propertyDiff
							typeDiff := schemaDiff.TypeDiff
							formatDiff := schemaDiff.FormatDiff

							if breakingTypeFormatChangedInResponseProperty(typeDiff, formatDiff, mediaType, schemaDiff) {
								typeDiff, formatDiff = fillEmptyTypeAndFormatDiffs(typeDiff, schemaDiff, formatDiff)
								result = append(result, ApiChange{
									Id:          ResponsePropertyTypeChangedId,
									Level:       ERR,
									Args:        []any{typeDiff.From, formatDiff.From, typeDiff.To, formatDiff.To, responseStatus},
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

func breakingTypeFormatChangedInResponseProperty(typeDiff *diff.ValueDiff, formatDiff *diff.ValueDiff, mediaType string, schemaDiff *diff.SchemaDiff) bool {

	if typeDiff != nil {
		return !isTypeContained(typeDiff.From, typeDiff.To, mediaType)
	}

	if formatDiff != nil {
		return !isFormatContained(schemaDiff.Revision.Type, formatDiff.From, formatDiff.To)
	}

	return false
}

// isTypeContained checks if type2 is contained in type1
func isTypeContained(type1, type2 interface{}, mediaType string) bool {
	return (type1 == "number" && type2 == "integer") ||
		(type1 == "string" && !isJsonMediaType(mediaType) && mediaType != "application/xml") // string can change to anything, unless it's json or xml
}

// isFormatContained checks if format2 is contained in format1
func isFormatContained(schemaType string, format1, format2 interface{}) bool {

	switch schemaType {
	case "number":
		return format1 == "double" && format2 == "float"
	case "integer":
		return (format1 == "int64" && format2 == "int32") ||
			(format1 == "bigint" && format2 == "int32") ||
			(format1 == "bigint" && format2 == "int64")
	case "string":
		return (format1 == "date-time" && format2 == "date" ||
			format1 == "date-time" && format2 == "time")
	}

	return false
}

func isJsonMediaType(mediaType string) bool {
	return mediaType == "application/json" ||
		(strings.HasPrefix(mediaType, "application/vnd.") && strings.HasSuffix(mediaType, "+json"))
}
