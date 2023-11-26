package checker

import (
	"github.com/tufin/oasdiff/diff"
)

const (
	ResponseBodyMinLengthDecreasedId     = "response-body-min-length-decreased"
	ResponsePropertyMinLengthDecreasedId = "response-property-min-length-decreased"
)

func ResponsePropertyMinLengthDecreasedCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config *Config) Changes {
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
					if mediaTypeDiff.SchemaDiff != nil && mediaTypeDiff.SchemaDiff.MinLengthDiff != nil {
						minLengthDiff := mediaTypeDiff.SchemaDiff.MinLengthDiff
						if minLengthDiff.From != nil &&
							minLengthDiff.To != nil {
							if IsDecreasedValue(minLengthDiff) {
								result = append(result, ApiChange{
									Id:          ResponseBodyMinLengthDecreasedId,
									Level:       ERR,
									Text:        config.Localize(ResponseBodyMinLengthDecreasedId, ColorizedValue(minLengthDiff.From), ColorizedValue(minLengthDiff.To)),
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
							minLengthDiff := propertyDiff.MinLengthDiff
							if minLengthDiff == nil {
								return
							}
							if minLengthDiff.To == nil ||
								minLengthDiff.From == nil {
								return
							}
							if !IsDecreasedValue(minLengthDiff) {
								return
							}

							if propertyDiff.Revision.WriteOnly {
								return
							}

							result = append(result, ApiChange{
								Id:          ResponsePropertyMinLengthDecreasedId,
								Level:       ERR,
								Text:        config.Localize(ResponsePropertyMinLengthDecreasedId, ColorizedValue(propertyFullName(propertyPath, propertyName)), ColorizedValue(minLengthDiff.From), ColorizedValue(minLengthDiff.To), ColorizedValue(responseStatus)),
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
