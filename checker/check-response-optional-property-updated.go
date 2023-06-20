package checker

import (
	"fmt"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/tufin/oasdiff/diff"
	"golang.org/x/exp/slices"
)

func ResponseOptionalPropertyUpdatedCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config BackwardCompatibilityCheckConfig) []BackwardCompatibilityError {
	result := make([]BackwardCompatibilityError, 0)
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
					comment := ""
					CheckDeletedPropertiesDiff(
						mediaTypeDiff.SchemaDiff,
						func(propertyPath string, propertyName string, propertyItem *openapi3.Schema, parent *diff.SchemaDiff) {
							level := WARN
							if propertyItem.WriteOnly {
								level = INFO
								comment = "This is a non breaking change because the property is write only"
							}
							if slices.Contains(parent.Base.Value.Required, propertyName) {
								// covered by response-required-property-removed
								return
							}
							result = append(result, BackwardCompatibilityError{
								Id:          "response-optional-property-removed",
								Level:       level,
								Text:        fmt.Sprintf(config.i18n("response-optional-property-removed"), ColorizedValue(propertyFullName(propertyPath, propertyName)), ColorizedValue(responseStatus)),
								Operation:   operation,
								OperationId: operationItem.Revision.OperationID,
								Path:        path,
								Source:      source,
								Comment:     comment,
							})
						})
					CheckAddedPropertiesDiff(
						mediaTypeDiff.SchemaDiff,
						func(propertyPath string, propertyName string, propertyItem *openapi3.Schema, parent *diff.SchemaDiff) {
							if propertyItem.WriteOnly {
								comment = "The property is write only"
							}
							if slices.Contains(parent.Base.Value.Required, propertyName) {
								// covered by response-required-property-added
								return
							}
							result = append(result, BackwardCompatibilityError{
								Id:          "response-optional-property-added",
								Level:       INFO,
								Text:        fmt.Sprintf(config.i18n("response-optional-property-added"), ColorizedValue(propertyFullName(propertyPath, propertyName)), ColorizedValue(responseStatus)),
								Operation:   operation,
								OperationId: operationItem.Revision.OperationID,
								Path:        path,
								Source:      source,
								Comment:     comment,
							})
						})

				}
			}
		}
	}
	return result
}
