package checker

import (
	"github.com/tufin/oasdiff/diff"
)

const (
	ResponseBodyDefaultValueAddedId       = "response-body-default-value-added"
	ResponseBodyDefaultValueRemovedId     = "response-body-default-value-removed"
	ResponseBodyDefaultValueChangedId     = "response-body-default-value-changed"
	ResponsePropertyDefaultValueAddedId   = "response-property-default-value-added"
	ResponsePropertyDefaultValueRemovedId = "response-property-default-value-removed"
	ResponsePropertyDefaultValueChangedId = "response-property-default-value-changed"
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
					Args:        []any{},
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
							appendResultItem(ResponseBodyDefaultValueAddedId, ColorizedValue(mediaType), ColorizedValue(defaultValueDiff.To), ColorizedValue(responseStatus))
						} else if defaultValueDiff.To == nil {
							appendResultItem(ResponseBodyDefaultValueRemovedId, ColorizedValue(mediaType), ColorizedValue(defaultValueDiff.From), ColorizedValue(responseStatus))
						} else {
							appendResultItem(ResponseBodyDefaultValueChangedId, ColorizedValue(mediaType), empty2none(defaultValueDiff.From), empty2none(defaultValueDiff.To), ColorizedValue(responseStatus))
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
								appendResultItem(ResponsePropertyDefaultValueAddedId, ColorizedValue(propertyName), empty2none(defaultValueDiff.To), ColorizedValue(responseStatus))
							} else if defaultValueDiff.To == nil {
								appendResultItem(ResponsePropertyDefaultValueRemovedId, ColorizedValue(propertyName), empty2none(defaultValueDiff.From), ColorizedValue(responseStatus))
							} else {
								appendResultItem(ResponsePropertyDefaultValueChangedId, ColorizedValue(propertyName), empty2none(defaultValueDiff.From), empty2none(defaultValueDiff.To), ColorizedValue(responseStatus))
							}
						})
				}
			}
		}
	}
	return result
}
