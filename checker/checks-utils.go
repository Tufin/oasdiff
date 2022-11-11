package checker

import (
	"fmt"
	"strings"

	"github.com/TwiN/go-color"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/tufin/oasdiff/diff"
)

func propertyFullName(propertyPath string, propertyNames ...string) string {
	propertyFullName := strings.Join(propertyNames, "/")
	if propertyPath != "" {
		propertyFullName = propertyPath + "/" + propertyFullName
	}
	return propertyFullName
}

func ColorizedValue(arg string) string {
	if IsPipedOutput() {
		return arg
	}
	return color.InBold(arg)
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

func processModifiedPropertiesDiff(propertyPath string, propertyName string, schemaDiff *diff.SchemaDiff, parentDiff *diff.SchemaDiff, processor func(propertyPath string, propertyName string, propertyItem *diff.SchemaDiff, propertyParentItem *diff.SchemaDiff)) {
	if propertyName != "" {
		processor(propertyPath, propertyName, schemaDiff, parentDiff)
	}

	if propertyPath == "" {
		propertyPath = propertyName
	} else {
		propertyPath = propertyPath + "/" + propertyName
	}

	if schemaDiff.AllOfDiff != nil {
		for k, v := range schemaDiff.AllOfDiff.Modified {
			processModifiedPropertiesDiff(fmt.Sprintf("%s/allOf[%s]", propertyPath, k), "", v, schemaDiff, processor)
		}
	}
	if schemaDiff.AnyOfDiff != nil {
		for k, v := range schemaDiff.AnyOfDiff.Modified {
			processModifiedPropertiesDiff(fmt.Sprintf("%s/anyOf[%s]", propertyPath, k), "", v, schemaDiff, processor)
		}
	}
	if schemaDiff.PropertiesDiff != nil {
		for i, v := range schemaDiff.PropertiesDiff.Modified {
			processModifiedPropertiesDiff(propertyPath, i, v, schemaDiff, processor)
		}
	}
}
