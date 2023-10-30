package checker

import (
	"github.com/tufin/oasdiff/diff"
)

const (
	RequestBodyMaxLengthSetId     = "request-body-max-length-set"
	RequestPropertyMaxLengthSetId = "request-property-max-length-set"
)

func RequestPropertyMaxLengthSetCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config Config) Changes {
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
					if maxLengthDiff.From == nil &&
						maxLengthDiff.To != nil {
						result = append(result, ApiChange{
							Id:          RequestBodyMaxLengthSetId,
							Level:       WARN,
							Text:        config.Localize(RequestBodyMaxLengthSetId, ColorizedValue(maxLengthDiff.To)),
							Comment:     config.Localize(RequestBodyMaxLengthSetId + "-comment"),
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
						maxLengthDiff := propertyDiff.MaxLengthDiff
						if maxLengthDiff == nil {
							return
						}
						if maxLengthDiff.From != nil ||
							maxLengthDiff.To == nil {
							return
						}
						if propertyDiff.Revision.ReadOnly {
							return
						}

						result = append(result, ApiChange{
							Id:          RequestPropertyMaxLengthSetId,
							Level:       WARN,
							Text:        config.Localize(RequestPropertyMaxLengthSetId, ColorizedValue(propertyFullName(propertyPath, propertyName)), ColorizedValue(maxLengthDiff.To)),
							Comment:     config.Localize(RequestPropertyMaxLengthSetId + "-comment"),
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
