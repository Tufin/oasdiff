package checker

import (
	"fmt"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/tufin/oasdiff/diff"
	"golang.org/x/exp/slices"
)

func NewRequiredRequestPropertyCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, diffBC *BCDiff) []BackwardCompatibilityError {
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
			for mediaType, mediaTypeDiff := range modifiedMediaTypes {
				if mediaTypeDiff.SchemaDiff == nil {
					continue
				}
				if mediaTypeDiff.SchemaDiff.PropertiesDiff == nil {
					continue
				}
				for _, topPropertyName := range mediaTypeDiff.SchemaDiff.PropertiesDiff.Added {
					processSchemaProperties(
						"",
						topPropertyName,
						mediaTypeDiff.SchemaDiff.Revision.Value.Properties[topPropertyName].Value,
						mediaTypeDiff.SchemaDiff.Revision.Value,
						func(propertyPath string, propertyName string, propertyItem *openapi3.Schema, parent *openapi3.Schema) {
							if !propertyItem.ReadOnly &&
								slices.Contains(parent.Required, propertyName) {
								source := (*operationsSources)[operationItem.Revision]
								propertyFullName := propertyName
								if propertyPath != "" {
									propertyFullName = propertyPath + "/" + propertyFullName
								}
								result = append(result, BackwardCompatibilityError{
									Id:        "new-required-request-property",
									Level:     ERR,
									Text:      fmt.Sprintf("added new required request property %s", ColorizedValue(propertyFullName)),
									Operation: operation,
									Path:      path,
									Source:    source,
									ToDo:      "Add to exceptions-list.md",
								})

								opDiff := diffBC.AddModifiedOperation(path, operation)
								if opDiff.RequestBodyDiff == nil {
									opDiff.RequestBodyDiff = &diff.RequestBodyDiff{}
								}
								if opDiff.RequestBodyDiff.ContentDiff == nil {
									opDiff.RequestBodyDiff.ContentDiff = &diff.ContentDiff{}
								}
								if opDiff.RequestBodyDiff.ContentDiff.MediaTypeModified == nil {
									opDiff.RequestBodyDiff.ContentDiff.MediaTypeModified = make(diff.ModifiedMediaTypes)
								}
								if opDiff.RequestBodyDiff.ContentDiff.MediaTypeModified[mediaType] == nil {
									opDiff.RequestBodyDiff.ContentDiff.MediaTypeModified[mediaType] = &diff.MediaTypeDiff{}
								}
								mediaTypeBCDiff := opDiff.RequestBodyDiff.ContentDiff.MediaTypeModified[mediaType]
								if mediaTypeBCDiff.SchemaDiff == nil {
									mediaTypeBCDiff.SchemaDiff = &diff.SchemaDiff{}
								}
								if mediaTypeBCDiff.SchemaDiff.PropertiesDiff == nil {
									mediaTypeBCDiff.SchemaDiff.PropertiesDiff = &diff.SchemasDiff{}
								}
								if mediaTypeBCDiff.SchemaDiff.PropertiesDiff.Added == nil {
									mediaTypeBCDiff.SchemaDiff.PropertiesDiff.Added = make(diff.StringList, 0)
								}
								items := mediaTypeBCDiff.SchemaDiff.PropertiesDiff.Added.ToStringSet()
								items.Add(propertyFullName)
								mediaTypeBCDiff.SchemaDiff.PropertiesDiff.Added = items.ToStringList()
							}
						})
				}
			}
		}
	}
	return result
}

func processSchemaProperties(propertyPath string, propertyName string, schema *openapi3.Schema, parent *openapi3.Schema, processor func(propertyPath string, propertyName string, propertyItem *openapi3.Schema, propertyParentItem *openapi3.Schema)) {
	if propertyName != "" {
		processor(propertyPath, propertyName, schema, parent)
	}

	if propertyPath == "" {
		propertyPath = propertyName
	} else {
		propertyPath = propertyPath + "/" + propertyName
	}

	if schema.AllOf != nil {
		for i, v := range schema.AllOf {
			processSchemaProperties(fmt.Sprintf("%s/allOf[%d]", propertyPath, i), "", v.Value, schema, processor)
		}
	}
	if schema.AnyOf != nil {
		for i, v := range schema.AnyOf {
			processSchemaProperties(fmt.Sprintf("%s/anyOf[%d]", propertyPath, i), "", v.Value, schema, processor)
		}
	}
	if schema.Properties != nil {
		for i, v := range schema.Properties {
			processSchemaProperties(propertyPath, i, v.Value, schema, processor)
		}
	}
}
