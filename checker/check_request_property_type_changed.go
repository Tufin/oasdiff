package checker

import (
	"github.com/tufin/oasdiff/diff"
	"github.com/tufin/oasdiff/load"
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
			source := (*operationsSources)[operationItem.Revision]

			modifiedMediaTypes := operationItem.RequestBodyDiff.ContentDiff.MediaTypeModified
			for mediaType, mediaTypeDiff := range modifiedMediaTypes {
				if mediaTypeDiff.SchemaDiff == nil {
					continue
				}

				schemaDiff := mediaTypeDiff.SchemaDiff
				typeDiff := schemaDiff.TypeDiff
				formatDiff := schemaDiff.FormatDiff

				if !typeDiff.Empty() || !formatDiff.Empty() {

					result = append(result, ApiChange{
						Id:          RequestBodyTypeChangedId,
						Level:       conditionalError(breakingTypeFormatChangedInRequestProperty(typeDiff, formatDiff, mediaType, schemaDiff), INFO),
						Args:        []any{getBaseType(schemaDiff), getBaseFormat(schemaDiff), getRevisionType(schemaDiff), getRevisionFormat(schemaDiff)},
						Operation:   operation,
						OperationId: operationItem.Revision.OperationID,
						Path:        path,
						Source:      load.NewSource(source),
					})
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

							result = append(result, ApiChange{
								Id:          RequestPropertyTypeChangedId,
								Level:       conditionalError(breakingTypeFormatChangedInRequestProperty(typeDiff, formatDiff, mediaType, schemaDiff), INFO),
								Args:        []any{propertyFullName(propertyPath, propertyName), getBaseType(schemaDiff), getBaseFormat(schemaDiff), getRevisionType(schemaDiff), getRevisionFormat(schemaDiff)},
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
	return result
}
