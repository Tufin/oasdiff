package diff

import (
	"github.com/getkin/kin-openapi/openapi3"
)

/*
SchemaDiff is a diff between schema objects: https://swagger.io/specification/#schema-object.
Boolean fields specify whether the property in question was changed between the two versions.
Pointer fields specify not only the presence of a change but also the old and new values of the property.
*/
type SchemaDiff struct {
	ExtensionProps                  *ExtensionsDiff `json:"extensions,omitempty"`
	SchemaAdded                     bool            `json:"schemaAdded,omitempty"`
	SchemaDeleted                   bool            `json:"schemaDeleted,omitempty"`
	ValueAdded                      bool            `json:"valueAdded,omitempty"`
	ValueDeleted                    bool            `json:"valueDeleted,omitempty"`
	OneOfDiff                       bool            `json:"oneOf,omitempty"`
	AnyOfDiff                       bool            `json:"anyOf,omitempty"`
	AllOfDiff                       bool            `json:"allOf,omitempty"`
	NotDiff                         bool            `json:"not,omitempty"`
	TypeDiff                        *ValueDiff      `json:"type,omitempty"`
	TitleDiff                       *ValueDiff      `json:"title,omitempty"`
	FormatDiff                      *ValueDiff      `json:"format,omitempty"`
	DescriptionDiff                 *ValueDiff      `json:"description,omitempty"`
	EnumDiff                        *EnumDiff       `json:"enum,omitempty"`
	DefaultDiff                     *ValueDiff      `json:"default,omitempty"`
	ExampleDiff                     *ValueDiff      `json:"example,omitempty"`
	AdditionalPropertiesAllowedDiff *ValueDiff      `json:"additionalPropertiesAllowed,omitempty"`
	UniqueItemsDiff                 *ValueDiff      `json:"uniqueItems,omitempty"`
	ExclusiveMinDiff                *ValueDiff      `json:"exclusiveMin,omitempty"`
	ExclusiveMaxDiff                *ValueDiff      `json:"exclusiveMax,omitempty"`
	NullableDiff                    *ValueDiff      `json:"nullable,omitempty"`
	ReadOnlyDiff                    *ValueDiff      `json:"readOnly,omitempty"`
	WriteOnlyDiff                   *ValueDiff      `json:"writeOnly,omitempty"`
	AllowEmptyValueDiff             *ValueDiff      `json:"allowEmptyValue,omitempty"`
	XMLDiff                         *ValueDiff      `json:"XML,omitempty"`
	DeprecatedDiff                  *ValueDiff      `json:"deprecated,omitempty"`
	MinDiff                         *ValueDiff      `json:"min,omitempty"`
	MaxDiff                         *ValueDiff      `json:"max,omitempty"`
	MultipleOf                      *ValueDiff      `json:"multipleOf,omitempty"`
	MinLength                       *ValueDiff      `json:"minLength,omitempty"`
	MaxLength                       *ValueDiff      `json:"maxLength,omitempty"`
	Pattern                         *ValueDiff      `json:"pattern,omitempty"`
	MinItems                        *ValueDiff      `json:"minItems,omitempty"`
	MaxItems                        *ValueDiff      `json:"maxItems,omitempty"`
	Items                           bool            `json:"items,omitempty"`
	Required                        *StringsDiff    `json:"required,omitempty"`
	PropertiesDiff                  *SchemasDiff    `json:"properties,omitempty"`
	MinProps                        *ValueDiff      `json:"minProps,omitempty"`
	MaxProps                        *ValueDiff      `json:"maxProps,omitempty"`
	AdditionalProperties            bool            `json:"additionalProperties,omitempty"`
}

func (schemaDiff SchemaDiff) empty() bool {
	return schemaDiff == SchemaDiff{}
}

func getSchemaDiff(config *Config, schema1, schema2 *openapi3.SchemaRef) SchemaDiff {

	value1, value2, status := getSchemaValues(schema1, schema2)

	if status != schemaStatusOK {
		return toSchemaDiff(status)
	}

	result := SchemaDiff{}

	if diff := getExtensionsDiff(value1.ExtensionProps, value2.ExtensionProps); !diff.empty() {
		result.ExtensionProps = diff
	}

	result.OneOfDiff = getDiffSchemas(config, value1.OneOf, value2.OneOf)
	result.AnyOfDiff = getDiffSchemas(config, value1.AnyOf, value2.AnyOf)
	result.AllOfDiff = getDiffSchemas(config, value1.AllOf, value2.AllOf)
	result.NotDiff = !getSchemaDiff(config, value1.Not, value2.Not).empty()
	result.TypeDiff = getValueDiff(value1.Type, value2.Type)
	result.TitleDiff = getValueDiff(value1.Title, value2.Title)
	result.FormatDiff = getValueDiff(value1.Format, value2.Format)
	result.DescriptionDiff = getValueDiff(value1.Description, value2.Description)
	result.EnumDiff = getEnumDiff(value1.Enum, value2.Enum)
	result.DefaultDiff = getValueDiff(value1.Default, value2.Default)

	if config.Examples {
		result.ExampleDiff = getValueDiff(value1.Example, value2.Example)
	}

	// ExternalDocs
	result.AdditionalPropertiesAllowedDiff = getBoolRefDiff(value1.AdditionalPropertiesAllowed, value2.AdditionalPropertiesAllowed)
	result.UniqueItemsDiff = getValueDiff(value1.UniqueItems, value2.UniqueItems)
	result.ExclusiveMinDiff = getValueDiff(value1.ExclusiveMin, value2.ExclusiveMin)
	result.ExclusiveMaxDiff = getValueDiff(value1.ExclusiveMax, value2.ExclusiveMax)
	result.NullableDiff = getValueDiff(value1.Nullable, value2.Nullable)
	result.ReadOnlyDiff = getValueDiff(value1.ReadOnly, value2.ReadOnly)
	result.WriteOnlyDiff = getValueDiff(value1.WriteOnly, value2.WriteOnly)
	result.AllowEmptyValueDiff = getValueDiff(value1.AllowEmptyValue, value2.AllowEmptyValue)
	result.XMLDiff = getValueDiff(value1.XML, value2.XML)
	result.DeprecatedDiff = getValueDiff(value1.Deprecated, value2.Deprecated)
	result.MinDiff = getFloat64RefDiff(value1.Min, value2.Min)
	result.MaxDiff = getFloat64RefDiff(value1.Max, value2.Max)
	result.MultipleOf = getFloat64RefDiff(value1.MultipleOf, value2.MultipleOf)
	result.MinLength = getValueDiff(value1.MinLength, value2.MinLength)
	result.MaxLength = getValueDiff(value1.MaxLength, value2.MaxLength)
	result.Pattern = getValueDiff(value1.Pattern, value2.Pattern)
	// compiledPattern is derived from pattern -> no need to diff
	result.MinItems = getValueDiff(value1.MinItems, value2.MinItems)
	result.MaxItems = getValueDiff(value1.MaxItems, value2.MaxItems)
	result.Items = !getSchemaDiff(config, value1.Items, value2.Items).empty()
	result.Required = getStringsDiff(value1.Required, value2.Required)
	if diff := getSchemasDiff(config, value1.Properties, value2.Properties); !diff.empty() {
		result.PropertiesDiff = diff
	}
	result.MinProps = getValueDiff(value1.MinProps, value2.MinProps)
	result.MaxProps = getValueDiff(value1.MaxProps, value2.MaxProps)
	result.AdditionalProperties = !getSchemaDiff(config, value1.AdditionalProperties, value2.AdditionalProperties).empty()
	// Discriminator

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

func getSchemaValues(schema1, schema2 *openapi3.SchemaRef) (*openapi3.Schema, *openapi3.Schema, schemaStatus) {

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

func toSchemaDiff(status schemaStatus) SchemaDiff {
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

func getDiffSchemas(config *Config, schemaRefs1, schemaRefs2 openapi3.SchemaRefs) bool {

	return !schemaRefsContained(config, schemaRefs1, schemaRefs2) || !schemaRefsContained(config, schemaRefs2, schemaRefs1)
}

func schemaRefsContained(config *Config, schemaRefs1, schemaRefs2 openapi3.SchemaRefs) bool {
	for _, schemaRef1 := range schemaRefs1 {
		if schemaRef1 != nil && schemaRef1.Value != nil {
			if !findSchema(config, schemaRef1, schemaRefs2) {
				return false
			}
		}
	}
	return true
}

func findSchema(config *Config, schemaRef1 *openapi3.SchemaRef, schemaRefs2 openapi3.SchemaRefs) bool {
	// TODO: optimize with a map
	for _, schemaRef2 := range schemaRefs2 {
		if schemaRef2 == nil || schemaRef2.Value == nil {
			continue
		}

		if diff := getSchemaDiff(config, schemaRef1, schemaRef2); diff.empty() {
			return true
		}
	}

	return false
}
