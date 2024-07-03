package checker

import (
	"github.com/tufin/oasdiff/diff"
)

const (
	RequestBodyMinLengthIncreasedId     = "request-body-min-length-increased"
	RequestBodyMinLengthDecreasedId     = "request-body-min-length-decreased"
	RequestPropertyMinLengthIncreasedId = "request-property-min-length-increased"
	RequestPropertyMinLengthDecreasedId = "request-property-min-length-decreased"
)

func RequestPropertyMinLengthUpdatedCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config *Config) Changes {
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
				if mediaTypeDiff.SchemaDiff != nil && mediaTypeDiff.SchemaDiff.MinLengthDiff != nil {
					minLengthDiff := mediaTypeDiff.SchemaDiff.MinLengthDiff
					if minLengthDiff.From != nil &&
						minLengthDiff.To != nil {
						if IsIncreasedValue(minLengthDiff) {
							result = append(result, NewApiChange(
								RequestBodyMinLengthIncreasedId,
								ERR,
								[]any{minLengthDiff.From, minLengthDiff.To},
								"",
								operationsSources,
								operationItem.Revision,
								operation,
								path,
							))
						} else {
							result = append(result, NewApiChange(
								RequestBodyMinLengthDecreasedId,
								INFO,
								[]any{minLengthDiff.From, minLengthDiff.To},
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
						minLengthDiff := propertyDiff.MinLengthDiff
						if minLengthDiff == nil {
							return
						}
						if minLengthDiff.From == nil ||
							minLengthDiff.To == nil {
							return
						}

						propName := propertyFullName(propertyPath, propertyName)

						if IsDecreasedValue(minLengthDiff) {
							result = append(result, NewApiChange(
								RequestPropertyMinLengthDecreasedId,
								INFO,
								[]any{propName, minLengthDiff.From, minLengthDiff.To},
								"",
								operationsSources,
								operationItem.Revision,
								operation,
								path,
							))
						} else {
							result = append(result, NewApiChange(
								RequestPropertyMinLengthIncreasedId,
								ERR,
								[]any{propName, minLengthDiff.From, minLengthDiff.To},
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
