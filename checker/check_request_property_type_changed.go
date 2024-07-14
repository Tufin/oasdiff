package checker

import (
	"github.com/tufin/oasdiff/diff"
)

const (
	RequestBodyTypeGeneralizedId     = "request-body-type-generalized"
	RequestBodyTypeChangedId         = "request-body-type-changed"
	RequestPropertyTypeGeneralizedId = "request-property-type-generalized"
	RequestPropertyTypeChangedId     = "request-property-type-changed"
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

					id := RequestBodyTypeGeneralizedId
					level := INFO

					if breakingTypeFormatChangedInRequestProperty(typeDiff, formatDiff, mediaType, schemaDiff) {
						id = RequestBodyTypeChangedId
						level = ERR
					}

					result = append(result, NewApiChange(
						id,
						level,
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

							id := RequestPropertyTypeGeneralizedId
							level := INFO

							if breakingTypeFormatChangedInRequestProperty(typeDiff, formatDiff, mediaType, schemaDiff) {
								id = RequestPropertyTypeChangedId
								level = ERR
							}

							result = append(result, NewApiChange(
								id,
								level,
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
