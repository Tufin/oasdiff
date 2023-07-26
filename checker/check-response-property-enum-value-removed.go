package checker

import (
	"fmt"

	"github.com/tufin/oasdiff/diff"
)

const responsePropertyEnumValueRemovedId = "response-property-enum-value-removed"

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
					checkModifiedPropertiesDiff(
						mediaTypeDiff.SchemaDiff,
						func(propertyPath string, propertyName string, propertyDiff *diff.SchemaDiff, parent *diff.SchemaDiff) {
							enumDiff := propertyDiff.EnumDiff
							if enumDiff == nil || enumDiff.Deleted == nil {
								return
							}

							for _, enumVal := range enumDiff.Deleted {
								result = append(result, ApiChange{
									Id:          responsePropertyEnumValueRemovedId,
									Level:       config.getLogLevel(responsePropertyEnumValueRemovedId, INFO),
									Text:        fmt.Sprintf(config.i18n(responsePropertyEnumValueRemovedId), enumVal, colorizedValue(propertyFullName(propertyPath, propertyName)), colorizedValue(responseStatus)),
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
