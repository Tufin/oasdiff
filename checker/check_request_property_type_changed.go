package checker

import (
	"github.com/tufin/oasdiff/diff"
)

const (
	RequestBodyTypeChangedId     = "request-body-type-changed"
	RequestPropertyTypeChangedId = "request-property-type-changed"
)

func RequestPropertyTypeChangedCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config *Config) Changes {
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

			modifiedMediaTypes := operationItem.RequestBodyDiff.ContentDiff.MediaTypeModified
			for mediaType, mediaTypeDiff := range modifiedMediaTypes {
				if mediaTypeDiff.SchemaDiff == nil {
					continue
				}

				schemaDiff := mediaTypeDiff.SchemaDiff
				typeDiff := schemaDiff.TypeDiff
				formatDiff := schemaDiff.FormatDiff

				if !typeDiff.Empty() || !formatDiff.Empty() {

					result = append(result, NewApiChange(
						RequestBodyTypeChangedId,
						conditionalError(breakingTypeFormatChangedInRequestProperty(typeDiff, formatDiff, mediaType, schemaDiff), INFO),
						[]any{getBaseType(schemaDiff), getBaseFormat(schemaDiff), getRevisionType(schemaDiff), getRevisionFormat(schemaDiff)},
						"",
						operationsSources,
						operationItem.Revision,
						operation,
						path,
					))
				}

				CheckModifiedPropertiesDiff(
					schemaDiff,
					func(propertyPath string, propertyName string, propertyDiff *diff.SchemaDiff, parent *diff.SchemaDiff) {
						if propertyDiff.Revision == nil {
							return
						}
						if propertyDiff.Revision.ReadOnly {
							return
						}
						schemaDiff := propertyDiff
						typeDiff := schemaDiff.TypeDiff
						formatDiff := schemaDiff.FormatDiff

						if !typeDiff.Empty() || !formatDiff.Empty() {

							result = append(result, NewApiChange(
								RequestPropertyTypeChangedId,
								conditionalError(breakingTypeFormatChangedInRequestProperty(typeDiff, formatDiff, mediaType, schemaDiff), INFO),
								[]any{propertyFullName(propertyPath, propertyName), getBaseType(schemaDiff), getBaseFormat(schemaDiff), getRevisionType(schemaDiff), getRevisionFormat(schemaDiff)},
								"",
								operationsSources,
								operationItem.Revision,
								operation,
								path,
							))
						}
					})
			}
		}
	}
	return result
}
