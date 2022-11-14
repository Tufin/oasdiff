package checker

import (
	"fmt"

	"github.com/tufin/oasdiff/diff"
)

func RequestPropertyBecameRequiredCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap) []BackwardCompatibilityError {
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

			if operationItem.RequestBodyDiff == nil ||
				operationItem.RequestBodyDiff.ContentDiff == nil ||
				operationItem.RequestBodyDiff.ContentDiff.MediaTypeModified == nil {
				continue
			}
			modifiedMediaTypes := operationItem.RequestBodyDiff.ContentDiff.MediaTypeModified
			for _, mediaTypeDiff := range modifiedMediaTypes {
				if mediaTypeDiff.SchemaDiff == nil {
					continue
				}

				if mediaTypeDiff.SchemaDiff.RequiredDiff != nil {
					for _, changedRequiredPropertyName := range mediaTypeDiff.SchemaDiff.RequiredDiff.Added {
						if mediaTypeDiff.SchemaDiff.Revision.Value.Properties[changedRequiredPropertyName].Value.ReadOnly {
							continue
						}
						result = append(result, BackwardCompatibilityError{
							Id:        "request-property-became-required",
							Level:     ERR,
							Text:      fmt.Sprintf("the request property %s became required", ColorizedValue(changedRequiredPropertyName)),
							Operation: operation,
							Path:      path,
							Source:    source,
							ToDo:      "Add to exceptions-list.md",
						})
					}
				}

				if mediaTypeDiff.SchemaDiff.PropertiesDiff == nil {
					continue
				}

				for topPropertyName, topPropertyDiff := range mediaTypeDiff.SchemaDiff.PropertiesDiff.Modified {
					processModifiedPropertiesDiff(
						"",
						topPropertyName,
						topPropertyDiff,
						nil,
						func(propertyPath string, propertyName string, propertyDiff *diff.SchemaDiff, parent *diff.SchemaDiff) {
							requiredDiff := propertyDiff.RequiredDiff
							if requiredDiff == nil {
								return
							}
							for _, changedRequiredPropertyName := range requiredDiff.Added {
								if propertyDiff.Revision.Value.Properties[changedRequiredPropertyName].Value.ReadOnly {
									continue
								}		
								result = append(result, BackwardCompatibilityError{
									Id:        "request-property-became-required",
									Level:     ERR,
									Text:      fmt.Sprintf("the request property %s became required", ColorizedValue(propertyFullName(propertyPath, propertyFullName(propertyName, changedRequiredPropertyName)))),
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
