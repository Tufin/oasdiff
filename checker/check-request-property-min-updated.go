package checker

import (
	"github.com/tufin/oasdiff/diff"
)

const (
	RequestBodyMinIncreasedId     = "request-body-min-increased"
	RequestBodyMinDecreasedId     = "request-body-min-decreased"
	RequestPropertyMinIncreasedId = "request-property-min-increased"
	RequestPropertyMinDecreasedId = "request-property-min-decreased"
)

func RequestPropertyMinIncreasedCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config Config) Changes {
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
				if mediaTypeDiff.SchemaDiff != nil && mediaTypeDiff.SchemaDiff.MinDiff != nil {
					minDiff := mediaTypeDiff.SchemaDiff.MinDiff
					if minDiff.From != nil &&
						minDiff.To != nil {
						if IsIncreasedValue(minDiff) {
							result = append(result, ApiChange{
								Id:          RequestBodyMinIncreasedId,
								Level:       ERR,
								Text:        config.Localize(RequestBodyMinIncreasedId, ColorizedValue(minDiff.To)),
								Args:        []any{},
								Operation:   operation,
								OperationId: operationItem.Revision.OperationID,
								Path:        path,
								Source:      source,
							})
						} else {
							result = append(result, ApiChange{
								Id:          RequestBodyMinDecreasedId,
								Level:       INFO,
								Text:        config.Localize(RequestBodyMinDecreasedId, ColorizedValue(minDiff.From), ColorizedValue(minDiff.To)),
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
						minDiff := propertyDiff.MinDiff
						if minDiff == nil {
							return
						}
						if minDiff.From == nil ||
							minDiff.To == nil {
							return
						}
						if IsIncreasedValue(minDiff) {
							result = append(result, ApiChange{
								Id:          RequestPropertyMinIncreasedId,
								Level:       ConditionalError(!propertyDiff.Revision.ReadOnly, INFO),
								Text:        config.Localize(RequestPropertyMinIncreasedId, ColorizedValue(propertyFullName(propertyPath, propertyName)), ColorizedValue(minDiff.To)),
								Args:        []any{},
								Operation:   operation,
								OperationId: operationItem.Revision.OperationID,
								Path:        path,
								Source:      source,
							})
						} else {
							result = append(result, ApiChange{
								Id:          RequestPropertyMinDecreasedId,
								Level:       INFO,
								Text:        config.Localize(RequestPropertyMinDecreasedId, ColorizedValue(propertyFullName(propertyPath, propertyName)), ColorizedValue(minDiff.From), ColorizedValue(minDiff.To)),
								Args:        []any{},
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
