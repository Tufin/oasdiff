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

func ColorizedValue(arg interface{}) string {
	str := interfaceToString(arg)
	if IsPipedOutput() {
		return fmt.Sprintf("'%s'", str)
	}
	return color.InBold(fmt.Sprintf("'%s'", str))
}

func interfaceToString(arg interface{}) string {
	if arg == nil {
		return "undefined"
	}

	if argString, ok := arg.(string); ok {
		return argString
	}

	if argUint64, ok := arg.(uint64); ok {
		return fmt.Sprintf("%d", argUint64)
	}

	if argFloat64, ok := arg.(float64); ok {
		return fmt.Sprintf("%.2f", argFloat64)
	}

	if argBool, ok := arg.(bool); ok {
		return fmt.Sprintf("%t", argBool)
	}

	return fmt.Sprintf("%v", arg)
}

func CheckModifiedPropertiesDiff(schemaDiff *diff.SchemaDiff, processor func(propertyPath string, propertyName string, propertyItem *diff.SchemaDiff, propertyParentItem *diff.SchemaDiff)) {
	if schemaDiff == nil {
		return
	}

	processModifiedPropertiesDiff("", "", schemaDiff, nil, processor)
}

func processModifiedPropertiesDiff(propertyPath string, propertyName string, schemaDiff *diff.SchemaDiff, parentDiff *diff.SchemaDiff, processor func(propertyPath string, propertyName string, propertyItem *diff.SchemaDiff, propertyParentItem *diff.SchemaDiff)) {
	if propertyName != "" || propertyPath != "" {
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

	if schemaDiff.OneOfDiff != nil {
		for k, v := range schemaDiff.OneOfDiff.Modified {
			processModifiedPropertiesDiff(fmt.Sprintf("%s/oneOf[%s]", propertyPath, k), "", v, schemaDiff, processor)
		}
	}

	if schemaDiff.ItemsDiff != nil {
		processModifiedPropertiesDiff(fmt.Sprintf("%s/items", propertyPath), "", schemaDiff.ItemsDiff, schemaDiff, processor)
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

	if schemaDiff.OneOfDiff != nil {
		for k, v := range schemaDiff.OneOfDiff.Modified {
			processAddedPropertiesDiff(fmt.Sprintf("%s/oneOf[%s]", propertyPath, k), "", v, schemaDiff, processor)
		}
	}

	if schemaDiff.ItemsDiff != nil {
		processAddedPropertiesDiff(fmt.Sprintf("%s/items", propertyPath), "", schemaDiff.ItemsDiff, schemaDiff, processor)
	}

	if schemaDiff.PropertiesDiff != nil {
		for _, v := range schemaDiff.PropertiesDiff.Added {
			processor(propertyPath, v, schemaDiff.Revision.Properties[v].Value, schemaDiff)
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

	if schemaDiff.OneOfDiff != nil {
		for k, v := range schemaDiff.OneOfDiff.Modified {
			processDeletedPropertiesDiff(fmt.Sprintf("%s/oneOf[%s]", propertyPath, k), "", v, schemaDiff, processor)
		}
	}

	if schemaDiff.ItemsDiff != nil {
		processDeletedPropertiesDiff(fmt.Sprintf("%s/items", propertyPath), "", schemaDiff.ItemsDiff, schemaDiff, processor)
	}

	if schemaDiff.PropertiesDiff != nil {
		for _, v := range schemaDiff.PropertiesDiff.Deleted {
			processor(propertyPath, v, schemaDiff.Base.Properties[v].Value, schemaDiff)
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
		return ColorizedValue("none")
	}
	return ColorizedValue(a)
}
