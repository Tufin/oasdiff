package checker

import (
	"github.com/tufin/oasdiff/diff"
)

func ResponsePropertyMinDecreasedCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config Config) Changes {
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
					if mediaTypeDiff.SchemaDiff != nil && mediaTypeDiff.SchemaDiff.MinDiff != nil {
						minDiff := mediaTypeDiff.SchemaDiff.MinDiff
						if minDiff.From != nil &&
							minDiff.To != nil {
							if IsDecreasedValue(minDiff) {
								result = append(result, ApiChange{
									Id:          "response-body-min-decreased",
									Level:       ERR,
									Text:        config.Localize("response-body-min-decreased", ColorizedValue(minDiff.From), ColorizedValue(minDiff.To)),
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
							if minDiff.To == nil ||
								minDiff.From == nil {
								return
							}
							if !IsDecreasedValue(minDiff) {
								return
							}

							if propertyDiff.Revision.WriteOnly {
								return
							}

							result = append(result, ApiChange{
								Id:          "response-property-min-decreased",
								Level:       ERR,
								Text:        config.Localize("response-property-min-decreased", ColorizedValue(propertyFullName(propertyPath, propertyName)), ColorizedValue(minDiff.From), ColorizedValue(minDiff.To), ColorizedValue(responseStatus)),
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
