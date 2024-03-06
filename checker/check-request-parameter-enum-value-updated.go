package checker

import (
	"github.com/tufin/oasdiff/diff"
)

const (
	RequestParameterEnumValueAddedId   = "request-parameter-enum-value-added"
	RequestParameterEnumValueRemovedId = "request-parameter-enum-value-removed"
)

func RequestParameterEnumValueUpdatedCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config *Config) Changes {
	changeGetter := newApiChangeGetter(config, operationsSources)
	result := make(Changes, 0)
	if diffReport.PathsDiff == nil {
		return result
	}
	for path, pathItem := range diffReport.PathsDiff.Modified {
		if pathItem.OperationsDiff == nil {
			continue
		}
		for operation, operationItem := range pathItem.OperationsDiff.Modified {
			if operationItem.ParametersDiff == nil {
				continue
			}
			if operationItem.ParametersDiff.Modified == nil {
				continue
			}
			for paramLocation, paramItems := range operationItem.ParametersDiff.Modified {
				for paramName, paramItem := range paramItems {
					if paramItem.SchemaDiff == nil {
						continue
					}
					enumDiff := paramItem.SchemaDiff.EnumDiff
					if enumDiff == nil {
						continue
					}
					for _, enumVal := range enumDiff.Deleted {
						result = append(result, changeGetter(
							RequestParameterEnumValueRemovedId,
							ERR,
							[]any{enumVal, paramLocation, paramName},
							"",
							operation,
							operationItem.Revision,
							path,
							operationItem.Revision,
						))
					}
					for _, enumVal := range enumDiff.Added {
						result = append(result, changeGetter(
							RequestParameterEnumValueAddedId,
							INFO,
							[]any{enumVal, paramLocation, paramName},
							"",
							operation,
							operationItem.Revision,
							path,
							operationItem.Revision,
						))
					}
				}
			}
		}
	}
	return result
}
