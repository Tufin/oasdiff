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

			modifiedMediaTypes := operationItem.RequestBodyDiff.ContentDiff.MediaTypeModified
			for mediaType, mediaTypeDiff := range modifiedMediaTypes {
				if mediaTypeDiff.SchemaDiff != nil && mediaTypeDiff.SchemaDiff.DefaultDiff != nil {
					defaultValueDiff := mediaTypeDiff.SchemaDiff.DefaultDiff
					result = append(result, ApiChange{
						Id:          "request-body-default-value-changed",
						Level:       INFO,
						Text:        fmt.Sprintf(config.i18n("request-body-default-value-changed"), ColorizedValue(mediaType), empty2none(defaultValueDiff.From), empty2none(defaultValueDiff.To)),
						Operation:   operation,
						OperationId: operationItem.Revision.OperationID,
						Path:        path,
						Source:      source,
					})
				}

				CheckModifiedPropertiesDiff(
					mediaTypeDiff.SchemaDiff,
					func(propertyPath string, propertyName string, propertyDiff *diff.SchemaDiff, parent *diff.SchemaDiff) {
						if propertyDiff == nil || propertyDiff.DefaultDiff == nil {
							return
						}

						defaultValueDiff := propertyDiff.DefaultDiff

						result = append(result, ApiChange{
							Id:          "request-property-default-value-changed",
							Level:       INFO,
							Text:        fmt.Sprintf(config.i18n("request-property-default-value-changed"), ColorizedValue(propertyName), empty2none(defaultValueDiff.From), empty2none(defaultValueDiff.To)),
							Operation:   operation,
							OperationId: operationItem.Revision.OperationID,
							Path:        path,
							Source:      source,
						})
					})
			}
		}
	}
	return result
}
