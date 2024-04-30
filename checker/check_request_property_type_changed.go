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
				if mediaTypeDiff.SchemaDiff != nil {
					schemaDiff := mediaTypeDiff.SchemaDiff
					typeDiff := schemaDiff.TypeDiff
					formatDiff := schemaDiff.FormatDiff

					if !typeDiff.Empty() || !formatDiff.Empty() {

						typeDiff = getDetailedTypeDiff(schemaDiff)
						formatDiff = getDetailedFormatDiff(schemaDiff)

						result = append(result, ApiChange{
							Id:          RequestBodyTypeChangedId,
							Level:       conditionalError(breakingTypeFormatChangedInRequestProperty(typeDiff, formatDiff, mediaType, schemaDiff), INFO),
							Args:        []any{typeDiff.Deleted, formatDiff.From, typeDiff.Added, formatDiff.To},
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

							typeDiff = getDetailedTypeDiff(schemaDiff)
							formatDiff = getDetailedFormatDiff(schemaDiff)

							result = append(result, ApiChange{
								Id:          RequestPropertyTypeChangedId,
								Level:       conditionalError(breakingTypeFormatChangedInRequestProperty(typeDiff, formatDiff, mediaType, schemaDiff), INFO),
								Args:        []any{propertyFullName(propertyPath, propertyName), typeDiff.Deleted, formatDiff.From, typeDiff.Added, formatDiff.To},
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

/*
getDetailedTypeDiff ensures that we have:
1. a non-empty diff, even if type wasn't changed
2. the full list of values in base and revision, rather than the added/deleted values only
*/
func getDetailedTypeDiff(schemaDiff *diff.SchemaDiff) *diff.StringsDiff {
	return &diff.StringsDiff{
		Deleted: schemaDiff.Base.Type.Slice(),
		Added:   schemaDiff.Revision.Type.Slice(),
	}
}

/*
getDetailedFormatDiff ensures that we have:
1. a non-empty diff, even if format wasn't changed
2. the original and revised format values which is implicit because format can only have a single value
*/
func getDetailedFormatDiff(schemaDiff *diff.SchemaDiff) *diff.ValueDiff {
	return &diff.ValueDiff{
		From: schemaDiff.Base.Format,
		To:   schemaDiff.Revision.Format,
	}
}

func breakingTypeFormatChangedInRequestProperty(typeDiff *diff.StringsDiff, formatDiff *diff.ValueDiff, mediaType string, schemaDiff *diff.SchemaDiff) bool {

	if typeDiff != nil {
		return !isTypeContained(typeDiff.Added, typeDiff.Deleted, mediaType)
	}

	if formatDiff != nil {
		return !isFormatContained(schemaDiff.Revision.Type, formatDiff.To, formatDiff.From)
	}

	return false
}
