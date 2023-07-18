package checker

import (
	"fmt"

	"github.com/tufin/oasdiff/diff"
)

const (
	ResponsePropertyEnumValueAddedCheckId          = "response-property-enum-value-added"
	ResponseWriteOnlyPropertyEnumValueAddedCheckId = "response-write-only-property-enum-value-added"
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

							id := ResponsePropertyEnumValueAddedCheckId
							level := WARN
							comment := config.i18n("response-property-enum-value-added-comment")

							if propertyDiff.Revision.Value.WriteOnly {
								// Document write-only enum update
								id = ResponseWriteOnlyPropertyEnumValueAddedCheckId
								level = INFO
								comment = ""
							}

							for _, enumVal := range enumDiff.Added {
								result = append(result, ApiChange{
									Id:          id,
									Level:       level,
									Text:        fmt.Sprintf(config.i18n(id), enumVal, ColorizedValue(propertyFullName(propertyPath, propertyName)), ColorizedValue(responseStatus)),
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
