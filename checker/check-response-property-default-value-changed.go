package checker

import (
	"fmt"

	"github.com/tufin/oasdiff/diff"
)

func ResponsePropertyDefaultValueChangedCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config Config) Changes {
	result := make(Changes, 0)
	if diffReport.PathsDiff == nil {
		return result
	}
	for path, pathItem := range diffReport.PathsDiff.Modified {
		if pathItem.OperationsDiff == nil {
			continue
		}
		for operation, operationItem := range pathItem.OperationsDiff.Modified {
			source := (*operationsSources)[operationItem.Revision]
			if operationItem.ResponsesDiff == nil || operationItem.ResponsesDiff.Modified == nil {
				continue
			}

			for responseStatus, responseDiff := range operationItem.ResponsesDiff.Modified {
				if responseDiff.ContentDiff == nil ||
					responseDiff.ContentDiff.MediaTypeModified == nil {
					continue
				}

				modifiedMediaTypes := responseDiff.ContentDiff.MediaTypeModified
				for mediaType, mediaTypeDiff := range modifiedMediaTypes {
					if mediaTypeDiff.SchemaDiff != nil && mediaTypeDiff.SchemaDiff.DefaultDiff != nil {
						defaultValueDiff := mediaTypeDiff.SchemaDiff.DefaultDiff
						result = append(result, ApiChange{
							Id:          "response-body-default-value-changed",
							Level:       INFO,
							Text:        fmt.Sprintf(config.i18n("response-body-default-value-changed"), ColorizedValue(mediaType), empty2none(defaultValueDiff.From), empty2none(defaultValueDiff.To), ColorizedValue(responseStatus)),
							Operation:   operation,
							OperationId: operationItem.Revision.OperationID,
							Path:        path,
							Source:      source,
						})
					}

					CheckModifiedPropertiesDiff(
						mediaTypeDiff.SchemaDiff,
						func(propertyPath string, propertyName string, propertyDiff *diff.SchemaDiff, parent *diff.SchemaDiff) {
							if propertyDiff == nil || propertyDiff.Revision == nil || propertyDiff.Revision.Value == nil || propertyDiff.DefaultDiff == nil {
								return
							}

							defaultValueDiff := propertyDiff.DefaultDiff

							result = append(result, ApiChange{
								Id:          "response-property-default-value-changed",
								Level:       INFO,
								Text:        fmt.Sprintf(config.i18n("response-property-default-value-changed"), ColorizedValue(propertyName), empty2none(defaultValueDiff.From), empty2none(defaultValueDiff.To), ColorizedValue(responseStatus)),
								Operation:   operation,
								OperationId: operationItem.Revision.OperationID,
								Path:        path,
								Source:      source,
							})
						})
				}
			}
		}
	}
	return result
}
