package checker

import (
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/tufin/oasdiff/diff"
	"golang.org/x/exp/slices"
)

func ResponseOptionalPropertyUpdatedCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config Config) Changes {
	result := make(Changes, 0)
	if diffReport.PathsDiff == nil {
		return result
	}
	for path, pathItem := range diffReport.PathsDiff.Modified {
		if pathItem.OperationsDiff == nil {
			continue
		}
		for operation, operationItem := range pathItem.OperationsDiff.Modified {
			source := (*operationsSources)[operationItem.Revision]

			if operationItem.ResponsesDiff == nil {
				continue
			}

			for responseStatus, responseDiff := range operationItem.ResponsesDiff.Modified {
				if responseDiff.ContentDiff == nil ||
					responseDiff.ContentDiff.MediaTypeModified == nil {
					continue
				}

				modifiedMediaTypes := responseDiff.ContentDiff.MediaTypeModified
				for _, mediaTypeDiff := range modifiedMediaTypes {
					CheckDeletedPropertiesDiff(
						mediaTypeDiff.SchemaDiff,
						func(propertyPath string, propertyName string, propertyItem *openapi3.Schema, parent *diff.SchemaDiff) {
							level := WARN
							id := "response-optional-property-removed"
							if propertyItem.WriteOnly {
								level = INFO
								id = "response-optional-write-only-property-removed"
							}
							if slices.Contains(parent.Base.Value.Required, propertyName) {
								// covered by response-required-property-removed
								return
							}
							result = append(result, ApiChange{
								Id:          id,
								Level:       level,
								Text:        config.Localize(id, ColorizedValue(propertyFullName(propertyPath, propertyName)), ColorizedValue(responseStatus)),
								Operation:   operation,
								OperationId: operationItem.Revision.OperationID,
								Path:        path,
								Source:      source,
							})
						})
					CheckAddedPropertiesDiff(
						mediaTypeDiff.SchemaDiff,
						func(propertyPath string, propertyName string, propertyItem *openapi3.Schema, parent *diff.SchemaDiff) {
							id := "response-optional-property-added"
							if propertyItem.WriteOnly {
								id = "response-optional-write-only-property-added"
							}
							if slices.Contains(parent.Revision.Value.Required, propertyName) {
								// covered by response-required-property-added
								return
							}
							result = append(result, ApiChange{
								Id:          id,
								Level:       INFO,
								Text:        config.Localize(id, ColorizedValue(propertyFullName(propertyPath, propertyName)), ColorizedValue(responseStatus)),
								Operation:   operation,
								OperationId: operationItem.Revision.OperationID,
								Path:        path,
								Source:      source,
							})
						})

				}
			}
		}
	}
	return result
}
