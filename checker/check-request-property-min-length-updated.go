package checker

import (
	"github.com/tufin/oasdiff/diff"
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
								Id:          "request-body-min-length-increased",
								Level:       ERR,
								Text:        config.Localize("request-body-min-length-increased", ColorizedValue(minLengthDiff.From), ColorizedValue(minLengthDiff.To)),
								Operation:   operation,
								OperationId: operationItem.Revision.OperationID,
								Path:        path,
								Source:      source,
							})
						} else {
							result = append(result, ApiChange{
								Id:          "request-body-min-length-decreased",
								Level:       INFO,
								Text:        config.Localize("request-body-min-length-decreased", ColorizedValue(minLengthDiff.From), ColorizedValue(minLengthDiff.To)),
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

						if IsDecreasedValue(minLengthDiff) {
							result = append(result, ApiChange{
								Id:          "request-property-min-length-decreased",
								Level:       INFO,
								Text:        config.Localize("request-property-min-length-decreased", ColorizedValue(propertyFullName(propertyPath, propertyName)), ColorizedValue(minLengthDiff.From), ColorizedValue(minLengthDiff.To)),
								Operation:   operation,
								OperationId: operationItem.Revision.OperationID,
								Path:        path,
								Source:      source,
							})
						} else {
							result = append(result, ApiChange{
								Id:          "request-property-min-length-increased",
								Level:       ERR,
								Text:        config.Localize("request-property-min-length-increased", ColorizedValue(propertyFullName(propertyPath, propertyName)), ColorizedValue(minLengthDiff.From), ColorizedValue(minLengthDiff.To)),
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
