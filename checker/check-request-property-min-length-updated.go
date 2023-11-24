package checker

import (
	"github.com/tufin/oasdiff/diff"
)

const (
	RequestBodyMinLengthIncreasedId     = "request-body-min-length-increased"
	RequestBodyMinLengthDecreasedId     = "request-body-min-length-decreased"
	RequestPropertyMinLengthIncreasedId = "request-property-min-length-increased"
	RequestPropertyMinLengthDecreasedId = "request-property-min-length-decreased"
)

func RequestPropertyMinLengthUpdatedCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config Config) Changes {
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
				if mediaTypeDiff.SchemaDiff != nil && mediaTypeDiff.SchemaDiff.MinLengthDiff != nil {
					minLengthDiff := mediaTypeDiff.SchemaDiff.MinLengthDiff
					if minLengthDiff.From != nil &&
						minLengthDiff.To != nil {
						if IsIncreasedValue(minLengthDiff) {
							result = append(result, ApiChange{
								Id:          RequestBodyMinLengthIncreasedId,
								Level:       ERR,
								Text:        config.Localize(RequestBodyMinLengthIncreasedId, ColorizedValue(minLengthDiff.From), ColorizedValue(minLengthDiff.To)),
								Args:        []any{minLengthDiff.From, minLengthDiff.To},
								Operation:   operation,
								OperationId: operationItem.Revision.OperationID,
								Path:        path,
								Source:      source,
							})
						} else {
							result = append(result, ApiChange{
								Id:          RequestBodyMinLengthDecreasedId,
								Level:       INFO,
								Text:        config.Localize(RequestBodyMinLengthDecreasedId, ColorizedValue(minLengthDiff.From), ColorizedValue(minLengthDiff.To)),
								Args:        []any{minLengthDiff.From, minLengthDiff.To},
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
						minLengthDiff := propertyDiff.MinLengthDiff
						if minLengthDiff == nil {
							return
						}
						if minLengthDiff.From == nil ||
							minLengthDiff.To == nil {
							return
						}

						fullName := propertyFullName(propertyPath, propertyName)

						if IsDecreasedValue(minLengthDiff) {
							result = append(result, ApiChange{
								Id:          RequestPropertyMinLengthDecreasedId,
								Level:       INFO,
								Text:        config.Localize(RequestPropertyMinLengthDecreasedId, ColorizedValue(fullName), ColorizedValue(minLengthDiff.From), ColorizedValue(minLengthDiff.To)),
								Args:        []any{fullName, minLengthDiff.From, minLengthDiff.To},
								Operation:   operation,
								OperationId: operationItem.Revision.OperationID,
								Path:        path,
								Source:      source,
							})
						} else {
							result = append(result, ApiChange{
								Id:          RequestPropertyMinLengthIncreasedId,
								Level:       ERR,
								Text:        config.Localize(RequestPropertyMinLengthIncreasedId, ColorizedValue(fullName), ColorizedValue(minLengthDiff.From), ColorizedValue(minLengthDiff.To)),
								Args:        []any{fullName, minLengthDiff.From, minLengthDiff.To},
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
