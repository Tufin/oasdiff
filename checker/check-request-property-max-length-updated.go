package checker

import (
	"github.com/tufin/oasdiff/diff"
)

const (
	RequestBodyMaxLengthDecreasedId     = "request-body-max-length-decreased"
	RequestBodyMaxLengthIncreasedId     = "request-body-max-length-increased"
	RequestPropertyMaxLengthDecreasedId = "request-property-max-length-decreased"
	RequestPropertyMaxLengthIncreasedId = "request-property-max-length-increased"
)

func RequestPropertyMaxLengthUpdatedCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config *Config) Changes {
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
				if mediaTypeDiff.SchemaDiff != nil && mediaTypeDiff.SchemaDiff.MaxLengthDiff != nil {
					maxLengthDiff := mediaTypeDiff.SchemaDiff.MaxLengthDiff
					if maxLengthDiff.From != nil &&
						maxLengthDiff.To != nil {
						if IsDecreasedValue(maxLengthDiff) {
							result = append(result, ApiChange{
								Id:          RequestBodyMaxLengthDecreasedId,
								Level:       ERR,
								Args:        []any{maxLengthDiff.To},
								Operation:   operation,
								OperationId: operationItem.Revision.OperationID,
								Path:        path,
								Source:      source,
							})
						} else {
							result = append(result, ApiChange{
								Id:          RequestBodyMaxLengthIncreasedId,
								Level:       INFO,
								Args:        []any{maxLengthDiff.From, maxLengthDiff.To},
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
						maxLengthDiff := propertyDiff.MaxLengthDiff
						if maxLengthDiff == nil {
							return
						}
						if maxLengthDiff.From == nil ||
							maxLengthDiff.To == nil {
							return
						}

						propName := propertyFullName(propertyPath, propertyName)

						if IsDecreasedValue(maxLengthDiff) {
							result = append(result, ApiChange{
								Id:          RequestPropertyMaxLengthDecreasedId,
								Level:       ConditionalError(!propertyDiff.Revision.ReadOnly, INFO),
								Args:        []any{propName, maxLengthDiff.To},
								Operation:   operation,
								OperationId: operationItem.Revision.OperationID,
								Path:        path,
								Source:      source,
							})
						} else {
							result = append(result, ApiChange{
								Id:          RequestPropertyMaxLengthIncreasedId,
								Level:       INFO,
								Args:        []any{propName, maxLengthDiff.From, maxLengthDiff.To},
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
