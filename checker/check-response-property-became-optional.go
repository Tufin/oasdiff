package checker

import (
	"fmt"

	"github.com/tufin/oasdiff/diff"
)

func ResponsePropertyBecameOptionalCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap) []BackwardCompatibilityError {
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

			for _, responseDiff := range operationItem.ResponsesDiff.Modified {
				if responseDiff.ContentDiff == nil ||
					responseDiff.ContentDiff.MediaTypeModified == nil {
					continue
				}

				modifiedMediaTypes := responseDiff.ContentDiff.MediaTypeModified
				for _, mediaTypeDiff := range modifiedMediaTypes {
					if mediaTypeDiff.SchemaDiff == nil {
						continue
					}
	
					if mediaTypeDiff.SchemaDiff.RequiredDiff != nil {
						for _, changedRequiredPropertyName := range mediaTypeDiff.SchemaDiff.RequiredDiff.Deleted {
							if mediaTypeDiff.SchemaDiff.Revision.Value.Properties[changedRequiredPropertyName].Value.WriteOnly {
								continue
							}
							result = append(result, BackwardCompatibilityError{
								Id:        "response-property-became-optional",
								Level:     ERR,
								Text:      fmt.Sprintf("the response property %s became optional", ColorizedValue(changedRequiredPropertyName)),
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
								for _, changedRequiredPropertyName := range requiredDiff.Deleted {
									if propertyDiff.Base.Value.Properties[changedRequiredPropertyName].Value.WriteOnly {
										continue
									}
									result = append(result, BackwardCompatibilityError{
										Id:        "response-property-became-optional",
										Level:     ERR,
										Text:      fmt.Sprintf("the response property %s became optional", ColorizedValue(propertyFullName(propertyPath, propertyFullName(propertyName, changedRequiredPropertyName)))),
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
	}
	return result
}
