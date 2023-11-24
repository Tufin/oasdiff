package checker

import (
	"github.com/tufin/oasdiff/diff"
)

const (
	RequestBodyDefaultValueAddedId       = "request-body-default-value-added"
	RequestBodyDefaultValueRemovedId     = "request-body-default-value-removed"
	RequestBodyDefaultValueChangedId     = "request-body-default-value-changed"
	RequestPropertyDefaultValueAddedId   = "request-property-default-value-added"
	RequestPropertyDefaultValueRemovedId = "request-property-default-value-removed"
	RequestPropertyDefaultValueChangedId = "request-property-default-value-changed"
)

func RequestPropertyDefaultValueChangedCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config Config) Changes {
	result := make(Changes, 0)
	if diffReport.PathsDiff == nil {
		return result
	}
	for path, pathItem := range diffReport.PathsDiff.Modified {
		if pathItem.OperationsDiff == nil {
			continue
		}
		for operation, operationItem := range pathItem.OperationsDiff.Modified {
			if operationItem.RequestBodyDiff == nil ||
				operationItem.RequestBodyDiff.ContentDiff == nil ||
				operationItem.RequestBodyDiff.ContentDiff.MediaTypeModified == nil {
				continue
			}
			source := (*operationsSources)[operationItem.Revision]

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

			modifiedMediaTypes := operationItem.RequestBodyDiff.ContentDiff.MediaTypeModified
			for mediaType, mediaTypeDiff := range modifiedMediaTypes {
				if mediaTypeDiff.SchemaDiff != nil && mediaTypeDiff.SchemaDiff.DefaultDiff != nil {
					defaultValueDiff := mediaTypeDiff.SchemaDiff.DefaultDiff

					if defaultValueDiff.From == nil {
						appendResultItem(RequestBodyDefaultValueAddedId, ColorizedValue(mediaType), ColorizedValue(defaultValueDiff.To))
					} else if defaultValueDiff.To == nil {
						appendResultItem(RequestBodyDefaultValueRemovedId, ColorizedValue(mediaType), ColorizedValue(defaultValueDiff.From))
					} else {
						appendResultItem(RequestBodyDefaultValueChangedId, ColorizedValue(mediaType), ColorizedValue(defaultValueDiff.From), ColorizedValue(defaultValueDiff.To))
					}
				}

				CheckModifiedPropertiesDiff(
					mediaTypeDiff.SchemaDiff,
					func(propertyPath string, propertyName string, propertyDiff *diff.SchemaDiff, parent *diff.SchemaDiff) {
						if propertyDiff == nil || propertyDiff.DefaultDiff == nil {
							return
						}

						defaultValueDiff := propertyDiff.DefaultDiff

						if defaultValueDiff.From == nil {
							appendResultItem(RequestPropertyDefaultValueAddedId, ColorizedValue(propertyName), ColorizedValue(defaultValueDiff.To))
						} else if defaultValueDiff.To == nil {
							appendResultItem(RequestPropertyDefaultValueRemovedId, ColorizedValue(propertyName), ColorizedValue(defaultValueDiff.From))
						} else {
							appendResultItem(RequestPropertyDefaultValueChangedId, ColorizedValue(propertyName), ColorizedValue(defaultValueDiff.From), ColorizedValue(defaultValueDiff.To))
						}
					})
			}
		}
	}
	return result
}
