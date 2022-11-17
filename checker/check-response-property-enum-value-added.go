package checker

import (
	"fmt"

	"github.com/tufin/oasdiff/diff"
)

func ResponsePropertyEnumValueAddedCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config BackwardCompatibilityCheckConfig) []BackwardCompatibilityError {
	result := make([]BackwardCompatibilityError, 0)
	if diffReport.PathsDiff == nil {
		return result
	}
	for path, pathItem := range diffReport.PathsDiff.Modified {
		if pathItem.OperationsDiff == nil {
			continue
		}
		for operation, operationItem := range pathItem.OperationsDiff.Modified {
			if operationItem.ResponsesDiff == nil || operationItem.ResponsesDiff.Modified == nil {
				continue
			}
			source := (*operationsSources)[operationItem.Revision]
			for responseStatus, responseDiff := range operationItem.ResponsesDiff.Modified {
				if responseDiff == nil ||
					responseDiff.ContentDiff == nil ||
					responseDiff.ContentDiff.MediaTypeModified == nil {
					continue
				}
				modifiedMediaTypes := responseDiff.ContentDiff.MediaTypeModified
				for _, mediaTypeDiff := range modifiedMediaTypes {
					CheckModifiedPropertiesDiff(
						mediaTypeDiff.SchemaDiff,
						func(propertyPath string, propertyName string, propertyDiff *diff.SchemaDiff, parent *diff.SchemaDiff) {
							enumDiff := propertyDiff.EnumDiff
							if enumDiff == nil || enumDiff.Added == nil {
								return
							}
							if propertyDiff.Revision.Value.WriteOnly {
								return
							}

							for _, enumVal := range enumDiff.Added {
								result = append(result, BackwardCompatibilityError{
									Id:        "response-property-enum-value-added",
									Level:     WARN,
									Text:      fmt.Sprintf("added the new '%s' enum value the %s response property for the response status %s", enumVal, ColorizedValue(propertyFullName(propertyPath, propertyName)), ColorizedValue(responseStatus)),
									Comment:   "Adding new enum values to response could be unexpected for clients, use x-extensible-enum instead.",
									Operation: operation,
									Path:      path,
									Source:    source,
									ToDo:      "Add to exceptions-list.md",
								})
							}
						})
				}

			}

		}
	}
	return result
}
