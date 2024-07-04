package checker

import (
	"github.com/tufin/oasdiff/diff"
)

const (
	RequestBodyMaxDecreasedId     = "request-body-max-decreased"
	RequestBodyMaxIncreasedId     = "request-body-max-increased"
	RequestPropertyMaxDecreasedId = "request-property-max-decreased"
	RequestPropertyMaxIncreasedId = "request-property-max-increased"
)

func RequestPropertyMaxDecreasedCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config *Config) Changes {
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
					if maxDiff.From != nil &&
						maxDiff.To != nil {
						if IsDecreasedValue(maxDiff) {
							result = append(result, NewApiChange(
								RequestBodyMaxDecreasedId,
								ERR,
								[]any{maxDiff.To},
								"",
								operationsSources,
								operationItem.Revision,
								operation,
								path,
							))
						} else {
							result = append(result, NewApiChange(
								RequestBodyMaxIncreasedId,
								INFO,
								[]any{maxDiff.From, maxDiff.To},
								"",
								operationsSources,
								operationItem.Revision,
								operation,
								path,
							))
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

						propName := propertyFullName(propertyPath, propertyName)

						if IsDecreasedValue(maxDiff) {
							result = append(result, NewApiChange(
								RequestPropertyMaxDecreasedId,
								conditionalError(!propertyDiff.Revision.ReadOnly, INFO),
								[]any{propName, maxDiff.To},
								"",
								operationsSources,
								operationItem.Revision,
								operation,
								path,
							))
						} else {
							result = append(result, NewApiChange(
								RequestPropertyMaxIncreasedId,
								INFO,
								[]any{propName, maxDiff.From, maxDiff.To},
								"",
								operationsSources,
								operationItem.Revision,
								operation,
								path,
							))
						}

					})
			}
		}
	}
	return result
}
