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
			source := (*operationsSources)[operationItem.Revision]
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
					from, ok := paramItem.SchemaDiff.ExtensionsDiff.Modified[diff.XExtensibleEnumExtension].From.(json.RawMessage)
					if !ok {
						continue
					}
					to, ok := paramItem.SchemaDiff.ExtensionsDiff.Modified[diff.XExtensibleEnumExtension].To.(json.RawMessage)
					if !ok {
						continue
					}
					var fromSlice []string
					if err := json.Unmarshal(from, &fromSlice); err != nil {
						result = append(result, ApiChange{
							Id:          UnparsableParameterFromXExtensibleEnumId,
							Level:       ERR,
							Args:        []any{paramLocation, paramName},
							Operation:   operation,
							OperationId: operationItem.Revision.OperationID,
							Path:        path,
							Source:      source,
						})
						continue
					}
					var toSlice []string
					if err := json.Unmarshal(to, &toSlice); err != nil {
						result = append(result, ApiChange{
							Id:          UnparsableParameterToXExtensibleEnumId,
							Level:       ERR,
							Args:        []any{paramLocation, paramName},
							Operation:   operation,
							OperationId: operationItem.Revision.OperationID,
							Path:        path,
							Source:      source,
						})
						continue
					}

					deletedVals := make([]string, 0)
					for _, fromVal := range fromSlice {
						if !slices.Contains(toSlice, fromVal) {
							deletedVals = append(deletedVals, fromVal)
						}
					}

					for _, enumVal := range deletedVals {
						result = append(result, ApiChange{
							Id:          RequestParameterXExtensibleEnumValueRemovedId,
							Level:       ERR,
							Args:        []any{enumVal, paramLocation, paramName},
							Operation:   operation,
							OperationId: operationItem.Revision.OperationID,
							Path:        path,
							Source:      source,
						})
					}
				}
			}
		}
	}
	return result
}
