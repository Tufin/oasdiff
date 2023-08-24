package checker

import (
	"github.com/tufin/oasdiff/diff"
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
								Id:          "request-body-max-decreased",
								Level:       ERR,
								Text:        config.Localize("request-body-max-decreased", ColorizedValue(maxDiff.To)),
								Operation:   operation,
								OperationId: operationItem.Revision.OperationID,
								Path:        path,
								Source:      source,
							})
						} else {
							result = append(result, ApiChange{
								Id:          "request-body-max-increased",
								Level:       INFO,
								Text:        config.Localize("request-body-max-increased", ColorizedValue(maxDiff.From), ColorizedValue(maxDiff.To)),
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
						if IsDecreasedValue(maxDiff) {
							level := ERR
							if propertyDiff.Revision.ReadOnly {
								level = INFO
							}

							result = append(result, ApiChange{
								Id:          "request-property-max-decreased",
								Level:       level,
								Text:        config.Localize("request-property-max-decreased", ColorizedValue(propertyFullName(propertyPath, propertyName)), ColorizedValue(maxDiff.To)),
								Operation:   operation,
								OperationId: operationItem.Revision.OperationID,
								Path:        path,
								Source:      source,
							})
						} else {
							result = append(result, ApiChange{
								Id:          "request-property-max-increased",
								Level:       INFO,
								Text:        config.Localize("request-property-max-increased", ColorizedValue(propertyFullName(propertyPath, propertyName)), ColorizedValue(maxDiff.From), ColorizedValue(maxDiff.To)),
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
