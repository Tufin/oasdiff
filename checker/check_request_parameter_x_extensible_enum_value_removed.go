package checker

import (
	"encoding/json"

	"github.com/tufin/oasdiff/diff"
	"golang.org/x/exp/slices"
)

const (
	UnparsableParameterFromXExtensibleEnumId      = "unparseable-parameter-from-x-extensible-enum"
	UnparsableParameterToXExtensibleEnumId        = "unparseable-parameter-to-x-extensible-enum"
	RequestParameterXExtensibleEnumValueRemovedId = "request-parameter-x-extensible-enum-value-removed"
)

func RequestParameterXExtensibleEnumValueRemovedCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config *Config) Changes {
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
					if paramItem.SchemaDiff.ExtensionsDiff == nil {
						continue
					}
					if paramItem.SchemaDiff.ExtensionsDiff.Modified == nil {
						continue
					}
					if paramItem.SchemaDiff.ExtensionsDiff.Modified[diff.XExtensibleEnumExtension] == nil {
						continue
					}
					from, ok := paramItem.SchemaDiff.Base.Extensions[diff.XExtensibleEnumExtension].(json.RawMessage)
					if !ok {
						continue
					}
					to, ok := paramItem.SchemaDiff.Revision.Extensions[diff.XExtensibleEnumExtension].(json.RawMessage)
					if !ok {
						continue
					}
					var fromSlice []string
					if err := json.Unmarshal(from, &fromSlice); err != nil {
						result = append(result, NewApiChange(
							UnparsableParameterFromXExtensibleEnumId,
							ERR,
							[]any{paramLocation, paramName},
							"",
							operationsSources,
							operationItem.Revision,
							operation,
							path,
						))
						continue
					}
					var toSlice []string
					if err := json.Unmarshal(to, &toSlice); err != nil {
						result = append(result, NewApiChange(
							UnparsableParameterToXExtensibleEnumId,
							ERR,
							[]any{paramLocation, paramName},
							"",
							operationsSources,
							operationItem.Revision,
							operation,
							path,
						))
						continue
					}

					deletedVals := make([]string, 0)
					for _, fromVal := range fromSlice {
						if !slices.Contains(toSlice, fromVal) {
							deletedVals = append(deletedVals, fromVal)
						}
					}

					for _, enumVal := range deletedVals {
						result = append(result, NewApiChange(
							RequestParameterXExtensibleEnumValueRemovedId,
							ERR,
							[]any{enumVal, paramLocation, paramName},
							"",
							operationsSources,
							operationItem.Revision,
							operation,
							path,
						))
					}
				}
			}
		}
	}
	return result
}
