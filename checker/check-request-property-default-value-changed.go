package checker

import (
	"fmt"

	"github.com/tufin/oasdiff/diff"
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
					Text:        fmt.Sprintf(config.i18n(messageId), a...),
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
						appendResultItem("request-body-default-value-added", ColorizedValue(mediaType), empty2none(defaultValueDiff.To))
					} else if defaultValueDiff.To == nil {
						appendResultItem("request-body-default-value-removed", ColorizedValue(mediaType), empty2none(defaultValueDiff.From))
					} else {
						appendResultItem("request-body-default-value-changed", ColorizedValue(mediaType), empty2none(defaultValueDiff.From), empty2none(defaultValueDiff.To))
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
							appendResultItem("request-property-default-value-added", ColorizedValue(propertyName), empty2none(defaultValueDiff.To))
						} else if defaultValueDiff.To == nil {
							appendResultItem("request-property-default-value-removed", ColorizedValue(propertyName), empty2none(defaultValueDiff.From))
						} else {
							appendResultItem("request-property-default-value-changed", ColorizedValue(propertyName), empty2none(defaultValueDiff.From), empty2none(defaultValueDiff.To))
						}
					})
			}
		}
	}
	return result
}
