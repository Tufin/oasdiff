package checker

import (
	"fmt"

	"github.com/tufin/oasdiff/diff"
)

func UncheckedRequestAllOfWarnCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config BackwardCompatibilityCheckConfig) []BackwardCompatibilityError {
	result := make([]BackwardCompatibilityError, 0)
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
				processModifiedPropertiesAllOfDiff(
					"",
					"",
					mediaTypeDiff.SchemaDiff,
					nil,
					func(propertyPath string, propertyName string, allOfDiff *diff.SchemaListDiff, parent *diff.SchemaDiff) {
						if allOfDiff.Added > 0 && allOfDiff.Deleted > 0 {
							source := (*operationsSources)[operationItem.Revision]
							result = append(result, BackwardCompatibilityError{
								Id:          "request-allOf-modified",
								Level:       WARN,
								Text:        fmt.Sprintf(config.i18n("request-allOf-modified"), ColorizedValue(propertyFullName(propertyPath, propertyName))),
								Comment:     config.i18n("request-allOf-modified-comment"),
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

func UncheckedResponseAllOfWarnCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config BackwardCompatibilityCheckConfig) []BackwardCompatibilityError {
	result := make([]BackwardCompatibilityError, 0)
	if diffReport.PathsDiff == nil {
		return result
	}
	for path, pathItem := range diffReport.PathsDiff.Modified {
		if pathItem.OperationsDiff == nil {
			continue
		}
		for operation, operationItem := range pathItem.OperationsDiff.Modified {
			if operationItem.ResponsesDiff == nil {
				continue
			}
			for responseStatus, responseDiff := range operationItem.ResponsesDiff.Modified {
				if responseDiff == nil ||
					responseDiff.ContentDiff == nil ||
					responseDiff.ContentDiff.MediaTypeModified == nil {
					continue
				}
				modifiedMediaTypes := responseDiff.ContentDiff.MediaTypeModified
				for _, mediaTypeDiff := range modifiedMediaTypes {
					processModifiedPropertiesAllOfDiff(
						"",
						"",
						mediaTypeDiff.SchemaDiff,
						nil,
						func(propertyPath string, propertyName string, allOfDiff *diff.SchemaListDiff, parent *diff.SchemaDiff) {
							if allOfDiff.Added > 0 && allOfDiff.Deleted > 0 {
								source := (*operationsSources)[operationItem.Revision]
								result = append(result, BackwardCompatibilityError{
									Id:          "response-allOf-modified",
									Level:       WARN,
									Text:        fmt.Sprintf("modified allOf for the response property %s for status %s", ColorizedValue(propertyFullName(propertyPath, propertyName)), ColorizedValue(responseStatus)),
									Comment:     "It is a warn because it is very difficult to check that allOf changed correctly without breaking changes",
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
	}
	return result
}

func processModifiedPropertiesAllOfDiff(propertyPath string, propertyName string, schemaDiff *diff.SchemaDiff, parentDiff *diff.SchemaDiff, processor func(propertyPath string, propertyName string, allOfDoff *diff.SchemaListDiff, propertyParentItem *diff.SchemaDiff)) {
	if schemaDiff == nil {
		return
	}

	if propertyName != "" {
		if propertyPath == "" {
			propertyPath = propertyName
		} else {
			propertyPath = propertyPath + "/" + propertyName
		}
	}

	if schemaDiff.AllOfDiff != nil {
		processor(propertyName, propertyName, schemaDiff.AllOfDiff, schemaDiff)
		for k, v := range schemaDiff.AllOfDiff.Modified {
			processModifiedPropertiesAllOfDiff(fmt.Sprintf("%s/allOf[%s]", propertyPath, k), "", v, schemaDiff, processor)
		}
	}
	if schemaDiff.AnyOfDiff != nil {
		for k, v := range schemaDiff.AnyOfDiff.Modified {
			processModifiedPropertiesAllOfDiff(fmt.Sprintf("%s/anyOf[%s]", propertyPath, k), "", v, schemaDiff, processor)
		}
	}
	if schemaDiff.PropertiesDiff != nil {
		for i, v := range schemaDiff.PropertiesDiff.Modified {
			processModifiedPropertiesAllOfDiff(propertyPath, i, v, schemaDiff, processor)
		}
	}
}
