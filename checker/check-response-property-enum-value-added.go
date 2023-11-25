package checker

import (
	"github.com/tufin/oasdiff/diff"
)

const (
	ResponsePropertyEnumValueAddedId          = "response-property-enum-value-added"
	ResponseWriteOnlyPropertyEnumValueAddedId = "response-write-only-property-enum-value-added"
)

func ResponsePropertyEnumValueAddedCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config Config) Changes {
	result := make(Changes, 0)
	if diffReport.PathsDiff == nil {
		return result
	}
	for path, pathItem := range diffReport.PathsDiff.Modified {
		if pathItem.OperationsDiff == nil {
			continue
		}
		for operation, operationItem := range pathItem.OperationsDiff.Modified {
			if operationItem.ResponsesDiff == nil || operationItem.ResponsesDiff.Modified == nil {
				continue
			}
			source := (*operationsSources)[operationItem.Revision]
			for responseStatus, responseDiff := range operationItem.ResponsesDiff.Modified {
				if responseDiff == nil ||
					responseDiff.ContentDiff == nil ||
					responseDiff.ContentDiff.MediaTypeModified == nil {
					continue
				}
				modifiedMediaTypes := responseDiff.ContentDiff.MediaTypeModified
				for _, mediaTypeDiff := range modifiedMediaTypes {
					CheckModifiedPropertiesDiff(
						mediaTypeDiff.SchemaDiff,
						func(propertyPath string, propertyName string, propertyDiff *diff.SchemaDiff, parent *diff.SchemaDiff) {
							enumDiff := propertyDiff.EnumDiff
							if enumDiff == nil || enumDiff.Added == nil {
								return
							}

							id := ResponsePropertyEnumValueAddedId
							level := WARN
							comment := config.Localize(comment(ResponsePropertyEnumValueAddedId))

							if propertyDiff.Revision.WriteOnly {
								// Document write-only enum update
								id = ResponseWriteOnlyPropertyEnumValueAddedId
								level = INFO
								comment = ""
							}

							propName := propertyFullName(propertyPath, propertyName)

							for _, enumVal := range enumDiff.Added {
								result = append(result, ApiChange{
									Id:          id,
									Level:       level,
									Text:        config.Localize(id, enumVal, ColorizedValue(propName), ColorizedValue(responseStatus)),
									Args:        []any{enumVal, propName, responseStatus},
									Comment:     comment,
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
