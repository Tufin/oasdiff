package checker

import (
	"encoding/json"

	"github.com/tufin/oasdiff/diff"
	"github.com/tufin/oasdiff/load"
	"golang.org/x/exp/slices"
)

const (
	UnparseablePropertyFromXExtensibleEnumId     = "unparseable-property-from-x-extensible-enum"
	UnparseablePropertyToXExtensibleEnumId       = "unparseable-property-to-x-extensible-enum"
	RequestPropertyXExtensibleEnumValueRemovedId = "request-property-x-extensible-enum-value-removed"
)

func RequestPropertyXExtensibleEnumValueRemovedCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config *Config) Changes {
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
				CheckModifiedPropertiesDiff(
					mediaTypeDiff.SchemaDiff,
					func(propertyPath string, propertyName string, propertyDiff *diff.SchemaDiff, parent *diff.SchemaDiff) {
						if propertyDiff.ExtensionsDiff == nil {
							return
						}
						if propertyDiff.ExtensionsDiff.Modified == nil {
							return
						}
						if propertyDiff.ExtensionsDiff.Modified[diff.XExtensibleEnumExtension] == nil {
							return
						}
						from, ok := diff.GetJsonOrigValue(propertyDiff.ExtensionsDiff.Modified[diff.XExtensibleEnumExtension])
						if !ok {
							return
						}
						to, ok := diff.GetJsonOrigValue(propertyDiff.ExtensionsDiff.Modified[diff.XExtensibleEnumExtension])
						if !ok {
							return
						}
						var fromSlice []string
						if err := json.Unmarshal(from, &fromSlice); err != nil {
							result = append(result, ApiChange{
								Id:          UnparseablePropertyFromXExtensibleEnumId,
								Level:       ERR,
								Args:        []any{propertyFullName(propertyPath, propertyName)},
								Operation:   operation,
								OperationId: operationItem.Revision.OperationID,
								Path:        path,
								Source:      load.NewSource(source),
							})
							return
						}
						var toSlice []string
						if err := json.Unmarshal(to, &toSlice); err != nil {
							result = append(result, ApiChange{
								Id:          UnparseablePropertyToXExtensibleEnumId,
								Level:       ERR,
								Args:        []any{propertyFullName(propertyPath, propertyName)},
								Operation:   operation,
								OperationId: operationItem.Revision.OperationID,
								Path:        path,
								Source:      load.NewSource(source),
							})
							return
						}

						deletedVals := make([]string, 0)
						for _, fromVal := range fromSlice {
							if !slices.Contains(toSlice, fromVal) {
								deletedVals = append(deletedVals, fromVal)
							}
						}

						if propertyDiff.Revision.ReadOnly {
							return
						}
						for _, enumVal := range deletedVals {
							result = append(result, ApiChange{
								Id:          RequestPropertyXExtensibleEnumValueRemovedId,
								Level:       ERR,
								Args:        []any{enumVal, propertyFullName(propertyPath, propertyName)},
								Operation:   operation,
								OperationId: operationItem.Revision.OperationID,
								Path:        path,
								Source:      load.NewSource(source),
							})
						}
					})
			}
		}
	}
	return result
}
