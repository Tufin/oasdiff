package checker

import (
	"github.com/tufin/oasdiff/diff"
	"github.com/tufin/oasdiff/load"
)

const (
	RequestPropertyEnumValueRemovedId = "request-property-enum-value-removed"
	RequestPropertyEnumValueAddedId   = "request-property-enum-value-added"
)

func RequestPropertyEnumValueUpdatedCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config *Config) Changes {
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
						enumDiff := propertyDiff.EnumDiff
						if enumDiff == nil {
							return
						}

						propName := propertyFullName(propertyPath, propertyName)

						for _, enumVal := range enumDiff.Deleted {
							result = append(result, ApiChange{
								Id:          RequestPropertyEnumValueRemovedId,
								Level:       conditionalError(!propertyDiff.Revision.ReadOnly, INFO),
								Args:        []any{enumVal, propName},
								Operation:   operation,
								OperationId: operationItem.Revision.OperationID,
								Path:        path,
								Source:      load.NewSource(source),
							})
						}

						for _, enumVal := range enumDiff.Added {
							result = append(result, ApiChange{
								Id:          RequestPropertyEnumValueAddedId,
								Level:       INFO,
								Args:        []any{enumVal, propName},
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
