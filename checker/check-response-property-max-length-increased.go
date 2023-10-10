package checker

import (
	"github.com/tufin/oasdiff/diff"
)

const (
	ResponseBodyMaxLengthIncreasedId     = "response-body-max-length-increased"
	ResponsePropertyMaxLengthIncreasedId = "response-property-max-length-increased"
)

func ResponsePropertyMaxLengthIncreasedCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config Config) Changes {
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
					if mediaTypeDiff.SchemaDiff != nil && mediaTypeDiff.SchemaDiff.MaxLengthDiff != nil {
						maxLengthDiff := mediaTypeDiff.SchemaDiff.MaxLengthDiff
						if maxLengthDiff.From != nil &&
							maxLengthDiff.To != nil {
							if IsIncreasedValue(maxLengthDiff) {
								result = append(result, ApiChange{
									Id:          ResponseBodyMaxLengthIncreasedId,
									Level:       ERR,
									Text:        config.Localize(ResponseBodyMaxLengthIncreasedId, ColorizedValue(maxLengthDiff.From), ColorizedValue(maxLengthDiff.To)),
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
							if maxLengthDiff.To == nil ||
								maxLengthDiff.From == nil {
								return
							}
							if !IsIncreasedValue(maxLengthDiff) {
								return
							}

							if propertyDiff.Revision.WriteOnly {
								return
							}

							result = append(result, ApiChange{
								Id:          ResponsePropertyMaxLengthIncreasedId,
								Level:       ERR,
								Text:        config.Localize(ResponsePropertyMaxLengthIncreasedId, ColorizedValue(propertyFullName(propertyPath, propertyName)), ColorizedValue(maxLengthDiff.From), ColorizedValue(maxLengthDiff.To), ColorizedValue(responseStatus)),
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
