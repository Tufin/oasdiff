package checker

import (
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/tufin/oasdiff/diff"
	"golang.org/x/exp/slices"
)

func RequestPropertyUpdatedCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config Config) Changes {
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
							result = append(result, ApiChange{
								Id:          "request-property-removed",
								Level:       WARN,
								Text:        config.Localize("request-property-removed", ColorizedValue(propertyFullName(propertyPath, propertyName))),
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
						if slices.Contains(parent.Revision.Required, propertyName) {
							result = append(result, ApiChange{
								Id:          "new-required-request-property",
								Level:       ERR,
								Text:        config.Localize("new-required-request-property", ColorizedValue(propertyFullName(propertyPath, propertyName))),
								Operation:   operation,
								OperationId: operationItem.Revision.OperationID,
								Path:        path,
								Source:      source,
							})
						} else {
							result = append(result, ApiChange{
								Id:          "new-optional-request-property",
								Level:       INFO,
								Text:        config.Localize("new-optional-request-property", ColorizedValue(propertyFullName(propertyPath, propertyName))),
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
