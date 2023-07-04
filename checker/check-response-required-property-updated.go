package checker

import (
	"fmt"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/tufin/oasdiff/diff"
	"golang.org/x/exp/slices"
)

const (
	ResponseRequiredPropertyRemovedCheckId          = "response-required-property-removed"
	ResponseRequiredWriteOnlyPropertyRemovedCheckId = "response-required-write-only-property-removed"
	ResponseRequiredPropertyAddedCheckId            = "response-required-property-added"
	ResponseRequiredWriteOnlyPropertyAddedCheckId   = "response-required-write-only-property-added"
)

func ResponseRequiredPropertyUpdatedCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config BackwardCompatibilityCheckConfig) IBackwardCompatibilityErrors {
	result := make(IBackwardCompatibilityErrors, 0)
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
							level := ERR
							id := ResponseRequiredPropertyRemovedCheckId
							if propertyItem.WriteOnly {
								level = INFO
								id = ResponseRequiredWriteOnlyPropertyRemovedCheckId
							}
							if !slices.Contains(parent.Base.Value.Required, propertyName) {
								// Covered by response-optional-property-removed
								return
							}
							result = append(result, BackwardCompatibilityError{
								Id:          id,
								Level:       level,
								Text:        fmt.Sprintf(config.i18n(id), ColorizedValue(propertyFullName(propertyPath, propertyName)), ColorizedValue(responseStatus)),
								Operation:   operation,
								OperationId: operationItem.Revision.OperationID,
								Path:        path,
								Source:      source,
							})
						})
					CheckAddedPropertiesDiff(
						mediaTypeDiff.SchemaDiff,
						func(propertyPath string, propertyName string, propertyItem *openapi3.Schema, parent *diff.SchemaDiff) {
							id := ResponseRequiredPropertyAddedCheckId
							if propertyItem.WriteOnly {
								id = ResponseRequiredWriteOnlyPropertyAddedCheckId
							}
							if !slices.Contains(parent.Revision.Value.Required, propertyName) {
								// Covered by response-optional-property-added
								return
							}
							result = append(result, BackwardCompatibilityError{
								Id:          id,
								Level:       INFO,
								Text:        fmt.Sprintf(config.i18n(id), ColorizedValue(propertyFullName(propertyPath, propertyName)), ColorizedValue(responseStatus)),
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
