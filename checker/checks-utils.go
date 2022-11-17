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
		return fmt.Sprintf("'%s'", arg)
	}
	if arg == "" {
		return "''"
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

func CheckModifiedPropertiesDiff(schemaDiff *diff.SchemaDiff, processor func(propertyPath string, propertyName string, propertyItem *diff.SchemaDiff, propertyParentItem *diff.SchemaDiff)) {
	if schemaDiff.PropertiesDiff != nil && schemaDiff.PropertiesDiff.Modified != nil {
		for topName, topDiff := range schemaDiff.PropertiesDiff.Modified {
			processModifiedPropertiesDiff("", topName, topDiff, schemaDiff, processor)
		}
	}
}

func processModifiedPropertiesDiff(propertyPath string, propertyName string, schemaDiff *diff.SchemaDiff, parentDiff *diff.SchemaDiff, processor func(propertyPath string, propertyName string, propertyItem *diff.SchemaDiff, propertyParentItem *diff.SchemaDiff)) {
	if schemaDiff == nil {
		return
	}

	if propertyName != "" {
		processor(propertyPath, propertyName, schemaDiff, parentDiff)
	}

	if propertyName != "" {
		if propertyPath == "" {
			propertyPath = propertyName
		} else {
			propertyPath = propertyPath + "/" + propertyName
		}
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

func CheckAddedPropertiesDiff(schemaDiff *diff.SchemaDiff, processor func(propertyPath string, propertyName string, propertyItem *openapi3.Schema, propertyParentDiff *diff.SchemaDiff)) {
	if schemaDiff == nil {
		return
	}
	processAddedPropertiesDiff("", "", schemaDiff, nil, processor)
}

func processAddedPropertiesDiff(propertyPath string, propertyName string, schemaDiff *diff.SchemaDiff, parentDiff *diff.SchemaDiff, processor func(propertyPath string, propertyName string, propertyItem *openapi3.Schema, propertyParentDiff *diff.SchemaDiff)) {
	if propertyName != "" {
		if propertyPath == "" {
			propertyPath = propertyName
		} else {
			propertyPath = propertyPath + "/" + propertyName
		}
	}

	if schemaDiff.AllOfDiff != nil {
		for k, v := range schemaDiff.AllOfDiff.Modified {
			processAddedPropertiesDiff(fmt.Sprintf("%s/allOf[%s]", propertyPath, k), "", v, schemaDiff, processor)
		}
	}
	if schemaDiff.AnyOfDiff != nil {
		for k, v := range schemaDiff.AnyOfDiff.Modified {
			processAddedPropertiesDiff(fmt.Sprintf("%s/anyOf[%s]", propertyPath, k), "", v, schemaDiff, processor)
		}
	}
	if schemaDiff.PropertiesDiff != nil {
		for _, v := range schemaDiff.PropertiesDiff.Added {
			processor(propertyPath, v, schemaDiff.Revision.Value.Properties[v].Value, schemaDiff)
		}
		for i, v := range schemaDiff.PropertiesDiff.Modified {
			processAddedPropertiesDiff(propertyPath, i, v, schemaDiff, processor)
		}
	}
}

func CheckDeletedPropertiesDiff(schemaDiff *diff.SchemaDiff, processor func(propertyPath string, propertyName string, propertyItem *openapi3.Schema, propertyParentDiff *diff.SchemaDiff)) {
	if schemaDiff == nil {
		return
	}

	processDeletedPropertiesDiff("", "", schemaDiff, nil, processor)
}

func processDeletedPropertiesDiff(propertyPath string, propertyName string, schemaDiff *diff.SchemaDiff, parentDiff *diff.SchemaDiff, processor func(propertyPath string, propertyName string, propertyItem *openapi3.Schema, propertyParentDiff *diff.SchemaDiff)) {
	if propertyName != "" {
		if propertyPath == "" {
			propertyPath = propertyName
		} else {
			propertyPath = propertyPath + "/" + propertyName
		}
	}

	if schemaDiff.AllOfDiff != nil {
		for k, v := range schemaDiff.AllOfDiff.Modified {
			processDeletedPropertiesDiff(fmt.Sprintf("%s/allOf[%s]", propertyPath, k), "", v, schemaDiff, processor)
		}
	}
	if schemaDiff.AnyOfDiff != nil {
		for k, v := range schemaDiff.AnyOfDiff.Modified {
			processDeletedPropertiesDiff(fmt.Sprintf("%s/anyOf[%s]", propertyPath, k), "", v, schemaDiff, processor)
		}
	}
	if schemaDiff.PropertiesDiff != nil {
		for _, v := range schemaDiff.PropertiesDiff.Deleted {
			processor(propertyPath, v, schemaDiff.Base.Value.Properties[v].Value, schemaDiff)
		}
		for i, v := range schemaDiff.PropertiesDiff.Modified {
			processDeletedPropertiesDiff(propertyPath, i, v, schemaDiff, processor)
		}
	}
}

func IsIncreased(from interface{}, to interface{}) bool {
	fromUint64, ok := from.(uint64)
	toUint64, okTo := to.(uint64)
	if ok && okTo {
		return fromUint64 < toUint64
	}
	fromFloat64, ok := from.(float64)
	toFloat64, okTo := to.(float64)
	if ok && okTo {
		return fromFloat64 < toFloat64
	}
	return false
}

func IsIncreasedValue(diff *diff.ValueDiff) bool {
	return IsIncreased(diff.From, diff.To)
}

func IsDecreasedValue(diff *diff.ValueDiff) bool {
	return IsDecreased(diff.From, diff.To)
}

func IsDecreased(from interface{}, to interface{}) bool {
	fromUint64, ok := from.(uint64)
	toUint64, okTo := to.(uint64)
	if ok && okTo {
		return fromUint64 > toUint64
	}
	fromFloat64, ok := from.(float64)
	toFloat64, okTo := to.(float64)
	if ok && okTo {
		return fromFloat64 > toFloat64
	}
	return false
}

func empty2none(a interface{}) interface{} {
	if a == nil || a == "" {
		return "none"
	}
	return a
}
