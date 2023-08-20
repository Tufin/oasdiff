package checker

import (
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

			appendResultItem := func(messageId string, a ...any) {
				result = append(result, ApiChange{
					Id:          messageId,
					Level:       INFO,
					Text:        config.Localize(messageId, a...),
					Operation:   operation,
					OperationId: operationItem.Revision.OperationID,
					Path:        path,
					Source:      source,
				})
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
						if defaultValueDiff.From == nil {
							appendResultItem("response-body-default-value-added", ColorizedValue(mediaType), ColorizedValue(defaultValueDiff.To), ColorizedValue(responseStatus))
						} else if defaultValueDiff.To == nil {
							appendResultItem("response-body-default-value-removed", ColorizedValue(mediaType), ColorizedValue(defaultValueDiff.From), ColorizedValue(responseStatus))
						} else {
							appendResultItem("response-body-default-value-changed", ColorizedValue(mediaType), empty2none(defaultValueDiff.From), empty2none(defaultValueDiff.To), ColorizedValue(responseStatus))
						}
					}

					CheckModifiedPropertiesDiff(
						mediaTypeDiff.SchemaDiff,
						func(propertyPath string, propertyName string, propertyDiff *diff.SchemaDiff, parent *diff.SchemaDiff) {
							if propertyDiff == nil || propertyDiff.Revision == nil || propertyDiff.DefaultDiff == nil {
								return
							}

							defaultValueDiff := propertyDiff.DefaultDiff
							if defaultValueDiff.From == nil {
								appendResultItem("response-property-default-value-added", ColorizedValue(propertyName), empty2none(defaultValueDiff.To), ColorizedValue(responseStatus))
							} else if defaultValueDiff.To == nil {
								appendResultItem("response-property-default-value-removed", ColorizedValue(propertyName), empty2none(defaultValueDiff.From), ColorizedValue(responseStatus))
							} else {
								appendResultItem("response-property-default-value-changed", ColorizedValue(propertyName), empty2none(defaultValueDiff.From), empty2none(defaultValueDiff.To), ColorizedValue(responseStatus))
							}
						})
				}
			}
		}
	}
	return result
}
