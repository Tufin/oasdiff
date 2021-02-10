package diff

import (
	"reflect"

	"github.com/getkin/kin-openapi/openapi3"
)

type SchemaDiff struct {
	SchemaAdded     bool       `json:"schemaAdded,omitempty"`
	SchemaDeleted   bool       `json:"schemaDelete,omitempty"`
	ValueAdded      bool       `json:"valueAdded,omitempty"`
	ValueDeleted    bool       `json:"valueDeleted,omitempty"`
	OneOfDiff       bool       `json:"oneOfDiff,omitempty"`
	AnyOfDiff       bool       `json:"anyOfDiff,omitempty"`
	AllOfDiff       bool       `json:"allOfDiff,omitempty"`
	NotDiff         bool       `json:"notDiff,omitempty"`
	TypeDiff        *ValueDiff `json:"typeDiff,omitempty"`
	TitleDiff       *ValueDiff `json:"titleDiff,omitempty"`
	FormatDiff      *ValueDiff `json:"formatDiff,omitempty"`
	DescriptionDiff *ValueDiff `json:"descriptionDiff,omitempty"`
	EnumDiff        bool       `json:"enumDiff,omitempty"`
}

func (schemaDiff SchemaDiff) empty() bool {
	return !schemaDiff.SchemaAdded && !schemaDiff.SchemaDeleted && !schemaDiff.ValueAdded && !schemaDiff.ValueDeleted &&
		!schemaDiff.OneOfDiff && !schemaDiff.AnyOfDiff && !schemaDiff.AllOfDiff && !schemaDiff.NotDiff &&
		schemaDiff.TypeDiff == nil &&
		schemaDiff.TitleDiff == nil &&
		schemaDiff.FormatDiff == nil &&
		schemaDiff.DescriptionDiff == nil &&
		!schemaDiff.EnumDiff
}

func diffSchema(schema1 *openapi3.SchemaRef, schema2 *openapi3.SchemaRef) SchemaDiff {

	value1, value2, status := getSchemaValues(schema1, schema2)

	if status != SchemaStatusOK {
		return getSchemaDiff(status)
	}

	result := SchemaDiff{}

	result.OneOfDiff = getDiffSchemas(value1.OneOf, value2.OneOf)
	result.AnyOfDiff = getDiffSchemas(value1.AnyOf, value2.AnyOf)
	result.AllOfDiff = getDiffSchemas(value1.AllOf, value2.AllOf)
	result.NotDiff = !diffSchema(value1.Not, value2.Not).empty()
	result.TypeDiff = getValueDiff(value1.Type, value2.Type)
	result.TitleDiff = getValueDiff(value1.Title, value2.Title)
	result.FormatDiff = getValueDiff(value1.Format, value2.Format)
	result.DescriptionDiff = getValueDiff(value1.Description, value2.Description)
	result.EnumDiff = getEnumDiff(value1.Enum, value2.Enum)

	return result
}

type SchemaStatus int

const (
	SchemaStatusOK SchemaStatus = iota
	SchemaStatusNoSchemas
	SchemaStatusSchemaAdded
	SchemaStatusSchemaDeleted
	SchemaStatusNoValues
	SchemaStatusValueAdded
	SchemaStatusValueDeleted
)

func getSchemaValues(schema1 *openapi3.SchemaRef, schema2 *openapi3.SchemaRef) (*openapi3.Schema, *openapi3.Schema, SchemaStatus) {

	if schema1 == nil && schema2 == nil {
		return nil, nil, SchemaStatusNoSchemas
	}

	if schema1 == nil && schema2 != nil {
		return nil, nil, SchemaStatusSchemaAdded
	}

	if schema1 != nil && schema2 == nil {
		return nil, nil, SchemaStatusSchemaDeleted
	}

	value1 := schema1.Value
	value2 := schema2.Value

	if value1 == nil && value2 == nil {
		return nil, nil, SchemaStatusNoValues
	}

	if value1 == nil && value2 != nil {
		return nil, nil, SchemaStatusValueAdded
	}

	if value1 != nil && value2 == nil {
		return nil, nil, SchemaStatusValueDeleted
	}

	return value1, value2, SchemaStatusOK
}

func getSchemaDiff(schemaStatus SchemaStatus) SchemaDiff {
	switch schemaStatus {
	case SchemaStatusSchemaAdded:
		return SchemaDiff{SchemaAdded: true}
	case SchemaStatusSchemaDeleted:
		return SchemaDiff{SchemaDeleted: true}
	case SchemaStatusValueAdded:
		return SchemaDiff{ValueAdded: true}
	case SchemaStatusValueDeleted:
		return SchemaDiff{ValueDeleted: true}
	}

	// all other cases -> empty diff
	return SchemaDiff{}
}

func getEnumDiff(enum1 []interface{}, enum2 []interface{}) bool {
	return !reflect.DeepEqual(enum1, enum2)
}

func getDiffSchemas(schemaRefs1 openapi3.SchemaRefs, schemaRefs2 openapi3.SchemaRefs) bool {

	for _, schemaRef1 := range schemaRefs1 {
		if schemaRef1 != nil && schemaRef1.Value != nil {
			if !findSchema(schemaRef1, schemaRefs2) {
				return true
			}
		}
	}
	return false
}

func findSchema(schemaRef1 *openapi3.SchemaRef, schemaRefs2 openapi3.SchemaRefs) bool {
	// TODO: optimize with a map
	for _, schemaRef2 := range schemaRefs2 {
		if schemaRef2 == nil || schemaRef2.Value == nil {
			continue
		}

		if diff := diffSchema(schemaRef1, schemaRef2); diff.empty() {
			return true
		}
	}

	return false
}
