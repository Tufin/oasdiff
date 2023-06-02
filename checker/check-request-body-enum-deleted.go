package checker

import (
	"fmt"

	"github.com/tufin/oasdiff/diff"
)

const requestBodyEnumRemovedId = "request-body-enum-value-removed"

func RequestBodyEnumValueRemovedCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config BackwardCompatibilityCheckConfig) []BackwardCompatibilityError {
	result := make([]BackwardCompatibilityError, 0)
	if diffReport.PathsDiff == nil {
		return result
	}
	for path, pathItem := range diffReport.PathsDiff.Modified {
		if pathItem.OperationsDiff == nil {
			continue
		}
		for operation, operationItem := range pathItem.OperationsDiff.Modified {
			if operationItem.RequestBodyDiff == nil {
				continue
			}
			if operationItem.RequestBodyDiff.ContentDiff == nil {
				continue
			}
			if operationItem.RequestBodyDiff.ContentDiff.MediaTypeModified == nil {
				continue
			}

			mediaTypeChanges := operationItem.RequestBodyDiff.ContentDiff.MediaTypeModified

			source := (*operationsSources)[operationItem.Revision]
			for _, mediaTypeItem := range mediaTypeChanges {
				if mediaTypeItem.SchemaDiff == nil {
					continue
				}

				enumDiff := mediaTypeItem.SchemaDiff.EnumDiff
				if enumDiff == nil || enumDiff.Deleted == nil {
					continue
				}
				for _, enumVal := range enumDiff.Deleted {
					result = append(result, BackwardCompatibilityError{
						Id:          requestBodyEnumRemovedId,
						Level:       ERR,
						Text:        fmt.Sprintf(config.i18n(requestBodyEnumRemovedId), ColorizedValue(enumVal)),
						Operation:   operation,
						OperationId: operationItem.Revision.OperationID,
						Path:        path,
						Source:      source,
					})
				}
			}
		}
	}
	return result
}
