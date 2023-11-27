package checker

import (
	"github.com/tufin/oasdiff/diff"
)

const (
	RequestBodyMinItemsSetId     = "request-body-min-items-set"
	RequestPropertyMinItemsSetId = "request-property-min-items-set"
)

func RequestPropertyMinItemsSetCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config *Config) Changes {
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
							Id:          RequestBodyMinItemsSetId,
							Level:       WARN,
							Args:        []any{minItemsDiff.To},
							Comment:     config.Localize(comment(RequestBodyMinItemsSetId)),
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
							Id:          RequestPropertyMinItemsSetId,
							Level:       WARN,
							Args:        []any{propertyFullName(propertyPath, propertyName), minItemsDiff.To},
							Comment:     config.Localize(comment(RequestPropertyMinItemsSetId)),
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
