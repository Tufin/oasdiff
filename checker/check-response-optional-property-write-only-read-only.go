package checker

import (
	"fmt"

	"github.com/tufin/oasdiff/diff"
	"golang.org/x/exp/slices"
)

const (
	ResponseOptionalPropertyBecameNonWriteOnlyCheckId = "response-optional-property-became-not-write-only"
	ResponseOptionalPropertyBecameWriteOnlyCheckId    = "response-optional-property-became-write-only"
	ResponseOptionalPropertyBecameReadOnlyCheckId     = "response-optional-property-became-read-only"
	ResponseOptionalPropertyBecameNonReadOnlyCheckId  = "response-optional-property-became-not-read-only"
)

func ResponseOptionalPropertyWriteOnlyReadOnlyCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config Config) Changes {
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
					if mediaTypeDiff.SchemaDiff == nil {
						continue
					}

					checkModifiedPropertiesDiff(
						mediaTypeDiff.SchemaDiff,
						func(propertyPath string, propertyName string, propertyDiff *diff.SchemaDiff, parent *diff.SchemaDiff) {
							writeOnlyDiff := propertyDiff.WriteOnlyDiff
							if writeOnlyDiff == nil {
								return
							}
							if parent.Revision.Value.Properties[propertyName] == nil {
								// removed properties processed by the ResponseOptionalPropertyUpdatedCheck check
								return
							}
							if slices.Contains(parent.Base.Value.Required, propertyName) {
								// skip required properties - checked at ResponseRequiredPropertyWriteOnlyReadOnlyCheck
								return
							}

							id := ResponseOptionalPropertyBecameNonWriteOnlyCheckId

							if writeOnlyDiff.To == true {
								id = ResponseOptionalPropertyBecameWriteOnlyCheckId
							}

							result = append(result, ApiChange{
								Id:          id,
								Level:       INFO,
								Text:        fmt.Sprintf(config.i18n(id), colorizedValue(propertyFullName(propertyPath, propertyName)), colorizedValue(responseStatus)),
								Operation:   operation,
								OperationId: operationItem.Revision.OperationID,
								Path:        path,
								Source:      source,
							})
						})

					checkModifiedPropertiesDiff(
						mediaTypeDiff.SchemaDiff,
						func(propertyPath string, propertyName string, propertyDiff *diff.SchemaDiff, parent *diff.SchemaDiff) {
							readOnlyDiff := propertyDiff.ReadOnlyDiff
							if readOnlyDiff == nil {
								return
							}
							if parent.Revision.Value.Properties[propertyName] == nil {
								// removed properties processed by the ResponseOptionalPropertyUpdatedCheck check
								return
							}
							if slices.Contains(parent.Base.Value.Required, propertyName) {
								// skip non-optional properties
								return
							}

							id := ResponseOptionalPropertyBecameNonReadOnlyCheckId

							if readOnlyDiff.To == true {
								id = ResponseOptionalPropertyBecameReadOnlyCheckId
							}

							result = append(result, ApiChange{
								Id:          id,
								Level:       INFO,
								Text:        fmt.Sprintf(config.i18n(id), colorizedValue(propertyFullName(propertyPath, propertyName)), colorizedValue(responseStatus)),
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
