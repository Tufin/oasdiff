package checker

import (
	"github.com/tufin/oasdiff/diff"
)

const (
	RequestBodyMaxDecreasedId     = "request-body-max-decreased"
	RequestBodyMaxIncreasedId     = "request-body-max-increased"
	RequestPropertyMaxDecreasedId = "request-property-max-decreased"
	RequestPropertyMaxIncreasedId = "request-property-max-increased"
)

func RequestPropertyMaxDecreasedCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config Config) Changes {
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
				if mediaTypeDiff.SchemaDiff != nil && mediaTypeDiff.SchemaDiff.MaxDiff != nil {
					maxDiff := mediaTypeDiff.SchemaDiff.MaxDiff
					if maxDiff.From != nil &&
						maxDiff.To != nil {
						if IsDecreasedValue(maxDiff) {
							result = append(result, ApiChange{
								Id:          RequestBodyMaxDecreasedId,
								Level:       ERR,
								Text:        config.Localize(RequestBodyMaxDecreasedId, ColorizedValue(maxDiff.To)),
								Args:        []any{maxDiff.To},
								Operation:   operation,
								OperationId: operationItem.Revision.OperationID,
								Path:        path,
								Source:      source,
							})
						} else {
							result = append(result, ApiChange{
								Id:          RequestBodyMaxIncreasedId,
								Level:       INFO,
								Text:        config.Localize(RequestBodyMaxIncreasedId, ColorizedValue(maxDiff.From), ColorizedValue(maxDiff.To)),
								Args:        []any{maxDiff.From, maxDiff.To},
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
						if maxDiff.From == nil ||
							maxDiff.To == nil {
							return
						}

						propName := propertyFullName(propertyPath, propertyName)

						if IsDecreasedValue(maxDiff) {
							result = append(result, ApiChange{
								Id:          RequestPropertyMaxDecreasedId,
								Level:       ConditionalError(!propertyDiff.Revision.ReadOnly, INFO),
								Text:        config.Localize(RequestPropertyMaxDecreasedId, ColorizedValue(propName), ColorizedValue(maxDiff.To)),
								Args:        []any{propName, maxDiff.To},
								Operation:   operation,
								OperationId: operationItem.Revision.OperationID,
								Path:        path,
								Source:      source,
							})
						} else {
							result = append(result, ApiChange{
								Id:          RequestPropertyMaxIncreasedId,
								Level:       INFO,
								Text:        config.Localize(RequestPropertyMaxIncreasedId, ColorizedValue(propName), ColorizedValue(maxDiff.From), ColorizedValue(maxDiff.To)),
								Args:        []any{propName, maxDiff.From, maxDiff.To},
								Operation:   operation,
								OperationId: operationItem.Revision.OperationID,
								Path:        path,
								Source:      source,
							})
						}

					})
			}
		}
	}
	return result
}
