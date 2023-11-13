package checker

import (
	"github.com/tufin/oasdiff/diff"
)

const (
	RequestBodyMinSetId     = "request-body-min-set"
	RequestPropertyMinSetId = "request-property-min-set"
)

func RequestPropertyMinSetCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config Config) Changes {
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
					if minDiff.From == nil &&
						minDiff.To != nil {
						result = append(result, ApiChange{
							Id:          RequestBodyMinSetId,
							Level:       WARN,
							Text:        config.Localize(RequestBodyMinSetId, ColorizedValue(minDiff.To)),
							Comment:     config.Localize(comment(RequestBodyMinSetId)),
							Operation:   operation,
							OperationId: operationItem.Revision.OperationID,
							Path:        path,
							Source:      source,
						})
					}
				}

				CheckModifiedPropertiesDiff(
					mediaTypeDiff.SchemaDiff,
					func(propertyPath string, propertyName string, propertyDiff *diff.SchemaDiff, parent *diff.SchemaDiff) {
						minDiff := propertyDiff.MinDiff
						if minDiff == nil {
							return
						}
						if minDiff.From != nil ||
							minDiff.To == nil {
							return
						}
						if propertyDiff.Revision.ReadOnly {
							return
						}

						result = append(result, ApiChange{
							Id:          RequestPropertyMinSetId,
							Level:       WARN,
							Text:        config.Localize(RequestPropertyMinSetId, ColorizedValue(propertyFullName(propertyPath, propertyName)), ColorizedValue(minDiff.To)),
							Comment:     config.Localize(comment(RequestPropertyMinSetId)),
							Operation:   operation,
							OperationId: operationItem.Revision.OperationID,
							Path:        path,
							Source:      source,
						})
					})
			}
		}
	}
	return result
}
