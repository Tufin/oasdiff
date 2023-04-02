package checker

import (
	"fmt"

	"github.com/tufin/oasdiff/diff"
)

const responseMediatypeEnumValueRemovedId = "response-mediatype-enum-value-removed"

func ResponseMediaTypeEnumValueRemovedCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config BackwardCompatibilityCheckConfig) []BackwardCompatibilityError {
	result := make([]BackwardCompatibilityError, 0)
	if diffReport.PathsDiff == nil {
		return result
	}
	for path, pathItem := range diffReport.PathsDiff.Modified {
		if pathItem.OperationsDiff == nil {
			continue
		}
		for operation, operationItem := range pathItem.OperationsDiff.Modified {
			if operationItem.ResponsesDiff == nil {
				continue
			}
			if operationItem.ResponsesDiff.Modified == nil {
				continue
			}
			source := (*operationsSources)[operationItem.Revision]
			for _, responseItems := range operationItem.ResponsesDiff.Modified {
				for mediaType, mediaTypeItem := range responseItems.ContentDiff.MediaTypeModified {
					if mediaTypeItem.SchemaDiff == nil {
						continue
					}

					enumDiff := mediaTypeItem.SchemaDiff.EnumDiff
					if enumDiff == nil {
						continue
					}

					for _, enumVal := range enumDiff.Deleted {
						result = append(result, BackwardCompatibilityError{
							Id:        responseMediatypeEnumValueRemovedId,
							Level:     ERR,
							Text:      fmt.Sprintf(config.i18n(responseMediatypeEnumValueRemovedId), mediaType, ColorizedValue(enumVal)),
							Operation: operation,
							Path:      path,
							Source:    source,
						})
					}

				}

			}
		}
	}
	return result
}
