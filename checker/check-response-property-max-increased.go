package checker

import (
	"github.com/tufin/oasdiff/diff"
)

const (
	ResponseBodyMaxIncreasedId     = "response-body-max-increased"
	ResponsePropertyMaxIncreasedId = "response-property-max-increased"
)

func ResponsePropertyMaxIncreasedCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config *Config) Changes {
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
					if mediaTypeDiff.SchemaDiff != nil && mediaTypeDiff.SchemaDiff.MaxDiff != nil {
						maxDiff := mediaTypeDiff.SchemaDiff.MaxDiff
						if maxDiff.From != nil &&
							maxDiff.To != nil {
							if IsIncreasedValue(maxDiff) {
								result = append(result, ApiChange{
									Id:          ResponseBodyMaxIncreasedId,
									Level:       ERR,
									Text:        config.Localize(ResponseBodyMaxIncreasedId, ColorizedValue(maxDiff.From), ColorizedValue(maxDiff.To)),
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
							maxDiff := propertyDiff.MaxDiff
							if maxDiff == nil {
								return
							}
							if maxDiff.To == nil ||
								maxDiff.From == nil {
								return
							}
							if !IsIncreasedValue(maxDiff) {
								return
							}

							if propertyDiff.Revision.WriteOnly {
								return
							}

							result = append(result, ApiChange{
								Id:          ResponsePropertyMaxIncreasedId,
								Level:       ERR,
								Text:        config.Localize(ResponsePropertyMaxIncreasedId, ColorizedValue(propertyFullName(propertyPath, propertyName)), ColorizedValue(maxDiff.From), ColorizedValue(maxDiff.To), ColorizedValue(responseStatus)),
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
