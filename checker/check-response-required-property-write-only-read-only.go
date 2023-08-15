package checker

import (
	"github.com/tufin/oasdiff/diff"
	"golang.org/x/exp/slices"
)

const (
	ResponseRequiredPropertyBecameNonWriteOnlyCheckId = "response-required-property-became-not-write-only"
	ResponseRequiredPropertyBecameWriteOnlyCheckId    = "response-required-property-became-write-only"
	ResponseRequiredPropertyBecameReadOnlyCheckId     = "response-required-property-became-read-only"
	ResponseRequiredPropertyBecameNonReadOnlyCheckId  = "response-required-property-became-not-read-only"
)

func ResponseRequiredPropertyWriteOnlyReadOnlyCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config Config) Changes {
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

					CheckModifiedPropertiesDiff(
						mediaTypeDiff.SchemaDiff,
						func(propertyPath string, propertyName string, propertyDiff *diff.SchemaDiff, parent *diff.SchemaDiff) {
							writeOnlyDiff := propertyDiff.WriteOnlyDiff
							if writeOnlyDiff == nil {
								return
							}
							if parent.Revision.Value.Properties[propertyName] == nil {
								// removed properties processed by the ResponseRequiredPropertyUpdatedCheck check
								return
							}
							if !slices.Contains(parent.Base.Value.Required, propertyName) {
								// skip non-required properties
								return
							}

							id := ResponseRequiredPropertyBecameNonWriteOnlyCheckId
							level := WARN
							comment := config.Localize("response-required-property-became-not-write-only-comment")

							if writeOnlyDiff.To == true {
								id = ResponseRequiredPropertyBecameWriteOnlyCheckId
								level = INFO
								comment = ""
							}

							result = append(result, ApiChange{
								Id:          id,
								Level:       level,
								Text:        config.Localize(id, ColorizedValue(propertyFullName(propertyPath, propertyName)), ColorizedValue(responseStatus)),
								Comment:     comment,
								Operation:   operation,
								OperationId: operationItem.Revision.OperationID,
								Path:        path,
								Source:      source,
							})
						})

					CheckModifiedPropertiesDiff(
						mediaTypeDiff.SchemaDiff,
						func(propertyPath string, propertyName string, propertyDiff *diff.SchemaDiff, parent *diff.SchemaDiff) {
							readOnlyDiff := propertyDiff.ReadOnlyDiff
							if readOnlyDiff == nil {
								return
							}
							if parent.Revision.Value.Properties[propertyName] == nil {
								// removed properties processed by the ResponseRequiredPropertyUpdatedCheck check
								return
							}
							if !slices.Contains(parent.Base.Value.Required, propertyName) {
								// skip non-required properties
								return
							}

							id := ResponseRequiredPropertyBecameNonReadOnlyCheckId
							level := INFO

							if readOnlyDiff.To == true {
								id = ResponseRequiredPropertyBecameReadOnlyCheckId
								level = INFO
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
				}
			}
		}
	}
	return result
}
