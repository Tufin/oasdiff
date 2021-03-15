package diff

import (
	"fmt"

	"github.com/getkin/kin-openapi/openapi3"
)

/*
SchemaDiff is a diff between schema objects: https://swagger.io/specification/#schema-object.
Boolean fields specify whether the property in question was changed between the two versions.
Pointer fields specify not only the presence of a change but also the old and new values of the property.
*/
type SchemaDiff struct {
	ExtensionsDiff                  *ExtensionsDiff `json:"extensions,omitempty" yaml:"extensions,omitempty"`
	SchemaAdded                     bool            `json:"schemaAdded,omitempty" yaml:"schemaAdded,omitempty"`
	SchemaDeleted                   bool            `json:"schemaDeleted,omitempty" yaml:"schemaDeleted,omitempty"`
	ValueAdded                      bool            `json:"valueAdded,omitempty" yaml:"valueAdded,omitempty"`
	ValueDeleted                    bool            `json:"valueDeleted,omitempty" yaml:"valueDeleted,omitempty"`
	OneOfDiff                       bool            `json:"oneOf,omitempty" yaml:"oneOf,omitempty"`
	AnyOfDiff                       bool            `json:"anyOf,omitempty" yaml:"anyOf,omitempty"`
	AllOfDiff                       bool            `json:"allOf,omitempty" yaml:"allOf,omitempty"`
	NotDiff                         *SchemaDiff     `json:"not,omitempty" yaml:"not,omitempty"`
	TypeDiff                        *ValueDiff      `json:"type,omitempty" yaml:"type,omitempty"`
	TitleDiff                       *ValueDiff      `json:"title,omitempty" yaml:"title,omitempty"`
	FormatDiff                      *ValueDiff      `json:"format,omitempty" yaml:"format,omitempty"`
	DescriptionDiff                 *ValueDiff      `json:"description,omitempty" yaml:"description,omitempty"`
	EnumDiff                        *EnumDiff       `json:"enum,omitempty" yaml:"enum,omitempty"`
	DefaultDiff                     *ValueDiff      `json:"default,omitempty" yaml:"default,omitempty"`
	ExampleDiff                     *ValueDiff      `json:"example,omitempty" yaml:"example,omitempty"`
	AdditionalPropertiesAllowedDiff *ValueDiff      `json:"additionalPropertiesAllowed,omitempty" yaml:"additionalPropertiesAllowed,omitempty"`
	UniqueItemsDiff                 *ValueDiff      `json:"uniqueItems,omitempty" yaml:"uniqueItems,omitempty"`
	ExclusiveMinDiff                *ValueDiff      `json:"exclusiveMin,omitempty" yaml:"exclusiveMin,omitempty"`
	ExclusiveMaxDiff                *ValueDiff      `json:"exclusiveMax,omitempty" yaml:"exclusiveMax,omitempty"`
	NullableDiff                    *ValueDiff      `json:"nullable,omitempty" yaml:"nullable,omitempty"`
	ReadOnlyDiff                    *ValueDiff      `json:"readOnly,omitempty" yaml:"readOnly,omitempty"`
	WriteOnlyDiff                   *ValueDiff      `json:"writeOnly,omitempty" yaml:"writeOnly,omitempty"`
	AllowEmptyValueDiff             *ValueDiff      `json:"allowEmptyValue,omitempty" yaml:"allowEmptyValue,omitempty"`
	XMLDiff                         *ValueDiff      `json:"XML,omitempty" yaml:"XML,omitempty"`
	DeprecatedDiff                  *ValueDiff      `json:"deprecated,omitempty" yaml:"deprecated,omitempty"`
	MinDiff                         *ValueDiff      `json:"min,omitempty" yaml:"min,omitempty"`
	MaxDiff                         *ValueDiff      `json:"max,omitempty" yaml:"max,omitempty"`
	MultipleOf                      *ValueDiff      `json:"multipleOf,omitempty" yaml:"multipleOf,omitempty"`
	MinLength                       *ValueDiff      `json:"minLength,omitempty" yaml:"minLength,omitempty"`
	MaxLength                       *ValueDiff      `json:"maxLength,omitempty" yaml:"maxLength,omitempty"`
	Pattern                         *ValueDiff      `json:"pattern,omitempty" yaml:"pattern,omitempty"`
	MinItems                        *ValueDiff      `json:"minItems,omitempty" yaml:"minItems,omitempty"`
	MaxItems                        *ValueDiff      `json:"maxItems,omitempty" yaml:"maxItems,omitempty"`
	Items                           *SchemaDiff     `json:"items,omitempty" yaml:"items,omitempty"`
	Required                        *StringsDiff    `json:"required,omitempty" yaml:"required,omitempty"`
	PropertiesDiff                  *SchemasDiff    `json:"properties,omitempty" yaml:"properties,omitempty"`
	MinProps                        *ValueDiff      `json:"minProps,omitempty" yaml:"minProps,omitempty"`
	MaxProps                        *ValueDiff      `json:"maxProps,omitempty" yaml:"maxProps,omitempty"`
	AdditionalProperties            *SchemaDiff     `json:"additionalProperties,omitempty" yaml:"additionalProperties,omitempty"`
}

// Empty indicates whether a change was found in this element
func (diff *SchemaDiff) Empty() bool {
	return diff == nil || *diff == SchemaDiff{}
}

func getSchemaDiff(config *Config, schema1, schema2 *openapi3.SchemaRef) (*SchemaDiff, error) {
	diff, err := getSchemaDiffInternal(config, schema1, schema2)
	if err != nil {
		return nil, err
	}
	if diff.Empty() {
		return nil, nil
	}
	return diff, nil
}

func getSchemaDiffInternal(config *Config, schema1, schema2 *openapi3.SchemaRef) (*SchemaDiff, error) {

	status := getSchemaStatus(schema1, schema2)
	if status != schemaStatusOK {
		return toSchemaDiff(status), nil
	}

	value1, err := derefSchema(schema1)
	if err != nil {
		return nil, err
	}
	value2, err := derefSchema(schema2)
	if err != nil {
		return nil, err
	}

	result := SchemaDiff{}

	result.ExtensionsDiff = getExtensionsDiff(config, value1.ExtensionProps, value2.ExtensionProps)
	result.OneOfDiff, err = getDiffSchemas(config, value1.OneOf, value2.OneOf)
	if err != nil {
		return nil, err
	}
	result.AnyOfDiff, err = getDiffSchemas(config, value1.AnyOf, value2.AnyOf)
	if err != nil {
		return nil, err
	}
	result.AllOfDiff, err = getDiffSchemas(config, value1.AllOf, value2.AllOf)
	if err != nil {
		return nil, err
	}
	result.NotDiff, err = getSchemaDiff(config, value1.Not, value2.Not)
	if err != nil {
		return nil, err
	}
	result.TypeDiff = getValueDiff(value1.Type, value2.Type)
	result.TitleDiff = getValueDiff(value1.Title, value2.Title)
	result.FormatDiff = getValueDiff(value1.Format, value2.Format)
	result.DescriptionDiff = getValueDiff(value1.Description, value2.Description)
	result.EnumDiff = getEnumDiff(value1.Enum, value2.Enum)
	result.DefaultDiff = getValueDiff(value1.Default, value2.Default)

	if config.IncludeExamples {
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
	result.Items, err = getSchemaDiff(config, value1.Items, value2.Items)
	if err != nil {
		return nil, err
	}

	result.Required = getStringsDiff(value1.Required, value2.Required)
	result.PropertiesDiff, err = getSchemasDiff(config, value1.Properties, value2.Properties)
	if err != nil {
		return nil, err
	}

	result.MinProps = getValueDiff(value1.MinProps, value2.MinProps)
	result.MaxProps = getValueDiff(value1.MaxProps, value2.MaxProps)
	result.AdditionalProperties, err = getSchemaDiff(config, value1.AdditionalProperties, value2.AdditionalProperties)
	if err != nil {
		return nil, err
	}
	// Discriminator

	return &result, nil
}

type schemaStatus int

const (
	schemaStatusOK schemaStatus = iota
	schemaStatusNoSchemas
	schemaStatusSchemaAdded
	schemaStatusSchemaDeleted
)

func getSchemaStatus(schema1, schema2 *openapi3.SchemaRef) schemaStatus {

	if schema1 == nil && schema2 == nil {
		return schemaStatusNoSchemas
	}

	if schema1 == nil && schema2 != nil {
		return schemaStatusSchemaAdded
	}

	if schema1 != nil && schema2 == nil {
		return schemaStatusSchemaDeleted
	}

	return schemaStatusOK
}

func derefSchema(ref *openapi3.SchemaRef) (*openapi3.Schema, error) {

	if ref == nil || ref.Value == nil {
		return nil, fmt.Errorf("schema reference is nil")
	}

	return ref.Value, nil
}

func toSchemaDiff(status schemaStatus) *SchemaDiff {
	switch status {
	case schemaStatusSchemaAdded:
		return &SchemaDiff{SchemaAdded: true}
	case schemaStatusSchemaDeleted:
		return &SchemaDiff{SchemaDeleted: true}
	}

	// all other cases -> empty diff
	return nil
}

func getDiffSchemas(config *Config, schemaRefs1, schemaRefs2 openapi3.SchemaRefs) (bool, error) {

	contained1, err := schemaRefsContained(config, schemaRefs1, schemaRefs2)
	if err != nil {
		return false, err
	}
	contained2, err := schemaRefsContained(config, schemaRefs2, schemaRefs1)
	if err != nil {
		return false, err
	}

	return !contained1 || !contained2, nil
}

func schemaRefsContained(config *Config, schemaRefs1, schemaRefs2 openapi3.SchemaRefs) (bool, error) {
	for _, schemaRef1 := range schemaRefs1 {
		found, err := findSchema(config, schemaRef1, schemaRefs2)
		if err != nil {
			return false, err
		}
		if !found {
			return false, nil
		}
	}
	return true, nil
}

func findSchema(config *Config, schemaRef1 *openapi3.SchemaRef, schemaRefs2 openapi3.SchemaRefs) (bool, error) {
	// TODO: optimize with a map
	for _, schemaRef2 := range schemaRefs2 {
		diff, err := getSchemaDiff(config, schemaRef1, schemaRef2)
		if err != nil {
			return false, err
		}
		if diff.Empty() {
			return true, nil
		}
	}

	return false, nil
}
