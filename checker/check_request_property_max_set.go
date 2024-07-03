package checker

import (
	"github.com/tufin/oasdiff/diff"
)

const (
	RequestBodyMaxSetId     = "request-body-max-set"
	RequestPropertyMaxSetId = "request-property-max-set"
)

func RequestPropertyMaxSetCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config *Config) Changes {
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

			modifiedMediaTypes := operationItem.RequestBodyDiff.ContentDiff.MediaTypeModified
			for _, mediaTypeDiff := range modifiedMediaTypes {
				if mediaTypeDiff.SchemaDiff != nil && mediaTypeDiff.SchemaDiff.MaxDiff != nil {
					maxDiff := mediaTypeDiff.SchemaDiff.MaxDiff
					if maxDiff.From == nil &&
						maxDiff.To != nil {
						result = append(result, NewApiChange(
							RequestBodyMaxSetId,
							WARN,
							[]any{maxDiff.To},
							commentId(RequestBodyMaxSetId),
							operationsSources,
							operationItem.Revision,
							operation,
							path,
						))
					}
				}

				CheckModifiedPropertiesDiff(
					mediaTypeDiff.SchemaDiff,
					func(propertyPath string, propertyName string, propertyDiff *diff.SchemaDiff, parent *diff.SchemaDiff) {
						maxDiff := propertyDiff.MaxDiff
						if maxDiff == nil {
							return
						}
						if maxDiff.From != nil ||
							maxDiff.To == nil {
							return
						}
						if propertyDiff.Revision.ReadOnly {
							return
						}

						result = append(result, NewApiChange(
							RequestPropertyMaxSetId,
							WARN,
							[]any{propertyFullName(propertyPath, propertyName), maxDiff.To},
							commentId(RequestPropertyMaxSetId),
							operationsSources,
							operationItem.Revision,
							operation,
							path,
						))
					})
			}
		}
	}
	return result
}
