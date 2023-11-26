package checker

import (
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/tufin/oasdiff/diff"
	"golang.org/x/exp/slices"
)

const (
	RequestPropertyRemovedId     = "request-property-removed"
	NewRequiredRequestPropertyId = "new-required-request-property"
	NewOptionalRequestPropertyId = "new-optional-request-property"
)

func RequestPropertyUpdatedCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config *Config) Changes {
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
				CheckDeletedPropertiesDiff(
					mediaTypeDiff.SchemaDiff,
					func(propertyPath string, propertyName string, propertyItem *openapi3.Schema, parent *diff.SchemaDiff) {
						if !propertyItem.ReadOnly {
							source := (*operationsSources)[operationItem.Revision]
							propName := propertyFullName(propertyPath, propertyName)
							result = append(result, ApiChange{
								Id:          RequestPropertyRemovedId,
								Level:       WARN,
								Text:        config.Localize(RequestPropertyRemovedId, ColorizedValue(propName)),
								Args:        []any{propName},
								Operation:   operation,
								OperationId: operationItem.Revision.OperationID,
								Path:        path,
								Source:      source,
							})
						}
					})
				CheckAddedPropertiesDiff(
					mediaTypeDiff.SchemaDiff,
					func(propertyPath string, propertyName string, propertyItem *openapi3.Schema, parent *diff.SchemaDiff) {
						source := (*operationsSources)[operationItem.Revision]
						if propertyItem.ReadOnly {
							return
						}

						propName := propertyFullName(propertyPath, propertyName)

						if slices.Contains(parent.Revision.Required, propertyName) {
							result = append(result, ApiChange{
								Id:          NewRequiredRequestPropertyId,
								Level:       ERR,
								Text:        config.Localize(NewRequiredRequestPropertyId, ColorizedValue(propName)),
								Args:        []any{propName},
								Operation:   operation,
								OperationId: operationItem.Revision.OperationID,
								Path:        path,
								Source:      source,
							})
						} else {
							result = append(result, ApiChange{
								Id:          NewOptionalRequestPropertyId,
								Level:       INFO,
								Text:        config.Localize(NewOptionalRequestPropertyId, ColorizedValue(propName)),
								Args:        []any{propName},
								Operation:   operation,
								OperationId: operationItem.Revision.OperationID,
								Path:        path,
								Source:      source,
							})
						}
					})
			}
		}
	}
	return result
}
