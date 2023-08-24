package checker

import (
	"github.com/tufin/oasdiff/diff"
)

func RequestPropertyMinItemsSetCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config Config) Changes {
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
			for _, mediaTypeDiff := range modifiedMediaTypes {
				if mediaTypeDiff.SchemaDiff != nil && mediaTypeDiff.SchemaDiff.MinItemsDiff != nil {
					minItemsDiff := mediaTypeDiff.SchemaDiff.MinItemsDiff
					if minItemsDiff.From == nil &&
						minItemsDiff.To != nil {
						result = append(result, ApiChange{
							Id:          "request-body-min-items-set",
							Level:       WARN,
							Text:        config.Localize("request-body-min-items-set", ColorizedValue(minItemsDiff.To)),
							Comment:     config.Localize("request-body-min-items-set-comment"),
							Operation:   operation,
							OperationId: operationItem.Revision.OperationID,
							Path:        path,
							Source:      source,
						})
					}
				}

				CheckModifiedPropertiesDiff(
					mediaTypeDiff.SchemaDiff,
					func(propertyPath string, propertyName string, propertyDiff *diff.SchemaDiff, parent *diff.SchemaDiff) {
						minItemsDiff := propertyDiff.MinItemsDiff
						if minItemsDiff == nil {
							return
						}
						if minItemsDiff.From != nil ||
							minItemsDiff.To == nil {
							return
						}
						if propertyDiff.Revision.ReadOnly {
							return
						}

						result = append(result, ApiChange{
							Id:          "request-property-min-items-set",
							Level:       WARN,
							Text:        config.Localize("request-property-min-items-set", ColorizedValue(propertyFullName(propertyPath, propertyName)), ColorizedValue(minItemsDiff.To)),
							Comment:     config.Localize("request-property-min-items-set-comment"),
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
