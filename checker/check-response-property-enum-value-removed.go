package checker

import (
	"github.com/tufin/oasdiff/diff"
)

const (
	ResponsePropertyEnumValueRemovedId = "response-property-enum-value-removed"
)

func ResponseParameterEnumValueRemovedCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config Config) Changes {
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
				for _, mediaTypeDiff := range responseDiff.ContentDiff.MediaTypeModified {
					CheckModifiedPropertiesDiff(
						mediaTypeDiff.SchemaDiff,
						func(propertyPath string, propertyName string, propertyDiff *diff.SchemaDiff, parent *diff.SchemaDiff) {
							enumDiff := propertyDiff.EnumDiff
							if enumDiff == nil || enumDiff.Deleted == nil {
								return
							}

							propName := propertyFullName(propertyPath, propertyName)

							for _, enumVal := range enumDiff.Deleted {
								result = append(result, ApiChange{
									Id:          ResponsePropertyEnumValueRemovedId,
									Level:       config.getLogLevel(ResponsePropertyEnumValueRemovedId, INFO),
									Text:        config.Localize(ResponsePropertyEnumValueRemovedId, enumVal, ColorizedValue(propName), ColorizedValue(responseStatus)),
									Args:        []any{enumVal, propName, responseStatus},
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
