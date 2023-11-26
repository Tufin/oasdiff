package checker

import (
	"github.com/tufin/oasdiff/diff"
)

const (
	ResponseBodyMinItemsDecreasedId     = "response-body-min-items-decreased"
	ResponsePropertyMinItemsDecreasedId = "response-property-min-items-decreased"
)

func ResponsePropertyMinItemsDecreasedCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config *Config) Changes {
	result := make(Changes, 0)
	if diffReport.PathsDiff == nil {
		return result
	}
	for path, pathItem := range diffReport.PathsDiff.Modified {
		if pathItem.OperationsDiff == nil {
			continue
		}
		for operation, operationItem := range pathItem.OperationsDiff.Modified {
			if operationItem.ResponsesDiff == nil || operationItem.ResponsesDiff.Modified == nil {
				continue
			}
			source := (*operationsSources)[operationItem.Revision]
			for responseStatus, responseDiff := range operationItem.ResponsesDiff.Modified {
				if responseDiff == nil ||
					responseDiff.ContentDiff == nil ||
					responseDiff.ContentDiff.MediaTypeModified == nil {
					continue
				}
				modifiedMediaTypes := responseDiff.ContentDiff.MediaTypeModified
				for _, mediaTypeDiff := range modifiedMediaTypes {
					if mediaTypeDiff.SchemaDiff != nil && mediaTypeDiff.SchemaDiff.MinItemsDiff != nil {
						minItemsDiff := mediaTypeDiff.SchemaDiff.MinItemsDiff
						if minItemsDiff.From != nil &&
							minItemsDiff.To != nil {
							if IsDecreasedValue(minItemsDiff) {
								result = append(result, ApiChange{
									Id:          ResponseBodyMinItemsDecreasedId,
									Level:       ERR,
									Text:        config.Localize(ResponseBodyMinItemsDecreasedId, ColorizedValue(minItemsDiff.From), ColorizedValue(minItemsDiff.To)),
									Args:        []any{},
									Operation:   operation,
									OperationId: operationItem.Revision.OperationID,
									Path:        path,
									Source:      source,
								})
							}
						}
					}

					CheckModifiedPropertiesDiff(
						mediaTypeDiff.SchemaDiff,
						func(propertyPath string, propertyName string, propertyDiff *diff.SchemaDiff, parent *diff.SchemaDiff) {
							minItemsDiff := propertyDiff.MinItemsDiff
							if minItemsDiff == nil {
								return
							}
							if minItemsDiff.To == nil ||
								minItemsDiff.From == nil {
								return
							}
							if !IsDecreasedValue(minItemsDiff) {
								return
							}

							if propertyDiff.Revision.WriteOnly {
								return
							}

							result = append(result, ApiChange{
								Id:          ResponsePropertyMinItemsDecreasedId,
								Level:       ERR,
								Text:        config.Localize(ResponsePropertyMinItemsDecreasedId, ColorizedValue(propertyFullName(propertyPath, propertyName)), ColorizedValue(minItemsDiff.From), ColorizedValue(minItemsDiff.To), ColorizedValue(responseStatus)),
								Args:        []any{},
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
