package diff

import (
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
	PropertiesDiff  bool       `json:"propertiesDiff,omitempty"`
}

func (schemaDiff SchemaDiff) empty() bool {
	return !schemaDiff.SchemaAdded && !schemaDiff.SchemaDeleted && !schemaDiff.ValueAdded && !schemaDiff.ValueDeleted &&
		!schemaDiff.OneOfDiff && !schemaDiff.AnyOfDiff && !schemaDiff.AllOfDiff && !schemaDiff.NotDiff &&
		schemaDiff.TypeDiff == nil &&
		schemaDiff.TitleDiff == nil &&
		schemaDiff.FormatDiff == nil &&
		schemaDiff.DescriptionDiff == nil &&
		!schemaDiff.EnumDiff &&
		!schemaDiff.PropertiesDiff
}

func diffSchema(schema1 *openapi3.SchemaRef, schema2 *openapi3.SchemaRef) SchemaDiff {

	value1, value2, status := getSchemaValues(schema1, schema2)

	if status != schemaStatusOK {
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
	result.PropertiesDiff = getDiffSchemaMap(value1.Properties, value2.Properties)

	return result
}

type schemaStatus int

const (
	schemaStatusOK schemaStatus = iota
	schemaStatusNoSchemas
	schemaStatusSchemaAdded
	schemaStatusSchemaDeleted
	schemaStatusNoValues
	schemaStatusValueAdded
	schemaStatusValueDeleted
)

func getSchemaValues(schema1 *openapi3.SchemaRef, schema2 *openapi3.SchemaRef) (*openapi3.Schema, *openapi3.Schema, schemaStatus) {

	if schema1 == nil && schema2 == nil {
		return nil, nil, schemaStatusNoSchemas
	}

	if schema1 == nil && schema2 != nil {
		return nil, nil, schemaStatusSchemaAdded
	}

	if schema1 != nil && schema2 == nil {
		return nil, nil, schemaStatusSchemaDeleted
	}

	value1 := schema1.Value
	value2 := schema2.Value

	if value1 == nil && value2 == nil {
		return nil, nil, schemaStatusNoValues
	}

	if value1 == nil && value2 != nil {
		return nil, nil, schemaStatusValueAdded
	}

	if value1 != nil && value2 == nil {
		return nil, nil, schemaStatusValueDeleted
	}

	return value1, value2, schemaStatusOK
}

func getSchemaDiff(status schemaStatus) SchemaDiff {
	switch status {
	case schemaStatusSchemaAdded:
		return SchemaDiff{SchemaAdded: true}
	case schemaStatusSchemaDeleted:
		return SchemaDiff{SchemaDeleted: true}
	case schemaStatusValueAdded:
		return SchemaDiff{ValueAdded: true}
	case schemaStatusValueDeleted:
		return SchemaDiff{ValueDeleted: true}
	}

	// all other cases -> empty diff
	return SchemaDiff{}
}

func getDiffSchemaMap(schemas1 openapi3.Schemas, schemas2 openapi3.Schemas) bool {

	for schemaName1, schemaRef1 := range schemas1 {
		schemaRef2, ok := schemas2[schemaName1]
		if !ok {
			return true
		}

		if diff := diffSchema(schemaRef1, schemaRef2); !diff.empty() {
			return true
		}
	}

	// TODO: handle added schemas
	return false
}

func getDiffSchemas(schemaRefs1 openapi3.SchemaRefs, schemaRefs2 openapi3.SchemaRefs) bool {

	for _, schemaRef1 := range schemaRefs1 {
		if schemaRef1 != nil && schemaRef1.Value != nil {
			if !findSchema(schemaRef1, schemaRefs2) {
				return true
			}
		}
	}

	// TODO: handle added schemas
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
