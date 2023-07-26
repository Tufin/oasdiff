package checker

import (
	"fmt"

	"github.com/tufin/oasdiff/diff"
)

const responseMediatypeEnumValueRemovedId = "response-mediatype-enum-value-removed"

func ResponseMediaTypeEnumValueRemovedCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config Config) Changes {
	result := make(Changes, 0)
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
				if responseItems.ContentDiff == nil {
					continue
				}

				if responseItems.ContentDiff.MediaTypeModified == nil {
					continue
				}
				for mediaType, mediaTypeItem := range responseItems.ContentDiff.MediaTypeModified {
					if mediaTypeItem.SchemaDiff == nil {
						continue
					}

					enumDiff := mediaTypeItem.SchemaDiff.EnumDiff
					if enumDiff == nil {
						continue
					}

					for _, enumVal := range enumDiff.Deleted {
						result = append(result, ApiChange{
							Id:          responseMediatypeEnumValueRemovedId,
							Level:       config.getLogLevel(responseMediatypeEnumValueRemovedId, ERR),
							Text:        fmt.Sprintf(config.i18n(responseMediatypeEnumValueRemovedId), mediaType, colorizedValue(enumVal)),
							Operation:   operation,
							OperationId: operationItem.Revision.OperationID,
							Path:        path,
							Source:      source,
						})
					}

				}

			}
		}
	}
	return result
}
