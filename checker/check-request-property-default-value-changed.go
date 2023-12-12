package checker

import (
	"github.com/tufin/oasdiff/diff"
	"github.com/tufin/oasdiff/load"
)

const (
	RequestBodyDefaultValueAddedId       = "request-body-default-value-added"
	RequestBodyDefaultValueRemovedId     = "request-body-default-value-removed"
	RequestBodyDefaultValueChangedId     = "request-body-default-value-changed"
	RequestPropertyDefaultValueAddedId   = "request-property-default-value-added"
	RequestPropertyDefaultValueRemovedId = "request-property-default-value-removed"
	RequestPropertyDefaultValueChangedId = "request-property-default-value-changed"
)

func RequestPropertyDefaultValueChangedCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config *Config) Changes {
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
					Args:        a,
					Operation:   operation,
					OperationId: operationItem.Revision.OperationID,
					Path:        path,
					Source:      load.NewSource(source),
				})
			}

			modifiedMediaTypes := operationItem.RequestBodyDiff.ContentDiff.MediaTypeModified
			for mediaType, mediaTypeDiff := range modifiedMediaTypes {
				if mediaTypeDiff.SchemaDiff != nil && mediaTypeDiff.SchemaDiff.DefaultDiff != nil {
					defaultValueDiff := mediaTypeDiff.SchemaDiff.DefaultDiff

					if defaultValueDiff.From == nil {
						appendResultItem(RequestBodyDefaultValueAddedId, mediaType, defaultValueDiff.To)
					} else if defaultValueDiff.To == nil {
						appendResultItem(RequestBodyDefaultValueRemovedId, mediaType, defaultValueDiff.From)
					} else {
						appendResultItem(RequestBodyDefaultValueChangedId, mediaType, defaultValueDiff.From, defaultValueDiff.To)
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
							appendResultItem(RequestPropertyDefaultValueAddedId, propertyName, defaultValueDiff.To)
						} else if defaultValueDiff.To == nil {
							appendResultItem(RequestPropertyDefaultValueRemovedId, propertyName, defaultValueDiff.From)
						} else {
							appendResultItem(RequestPropertyDefaultValueChangedId, propertyName, defaultValueDiff.From, defaultValueDiff.To)
						}
					})
			}
		}
	}
	return result
}
