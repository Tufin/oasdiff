package diff

import (
	"fmt"

	"github.com/getkin/kin-openapi/openapi3"
)

// SchemaDiff describes the changes between a pair of schema objects: https://swagger.io/specification/#schema-object
type SchemaDiff struct {
	SchemaAdded                     bool               `json:"schemaAdded,omitempty" yaml:"schemaAdded,omitempty"`
	SchemaDeleted                   bool               `json:"schemaDeleted,omitempty" yaml:"schemaDeleted,omitempty"`
	ExtensionsDiff                  *ExtensionsDiff    `json:"extensions,omitempty" yaml:"extensions,omitempty"`
	OneOfDiff                       *SchemaListDiff    `json:"oneOf,omitempty" yaml:"oneOf,omitempty"`
	AnyOfDiff                       *SchemaListDiff    `json:"anyOf,omitempty" yaml:"anyOf,omitempty"`
	AllOfDiff                       *SchemaListDiff    `json:"allOf,omitempty" yaml:"allOf,omitempty"`
	NotDiff                         *SchemaDiff        `json:"not,omitempty" yaml:"not,omitempty"`
	TypeDiff                        *ValueDiff         `json:"type,omitempty" yaml:"type,omitempty"`
	TitleDiff                       *ValueDiff         `json:"title,omitempty" yaml:"title,omitempty"`
	FormatDiff                      *ValueDiff         `json:"format,omitempty" yaml:"format,omitempty"`
	DescriptionDiff                 *ValueDiff         `json:"description,omitempty" yaml:"description,omitempty"`
	EnumDiff                        *EnumDiff          `json:"enum,omitempty" yaml:"enum,omitempty"`
	DefaultDiff                     *ValueDiff         `json:"default,omitempty" yaml:"default,omitempty"`
	ExampleDiff                     *ValueDiff         `json:"example,omitempty" yaml:"example,omitempty"`
	ExternalDocsDiff                *ExternalDocsDiff  `json:"externalDocs,omitempty" yaml:"externalDocs,omitempty"`
	AdditionalPropertiesAllowedDiff *ValueDiff         `json:"additionalPropertiesAllowed,omitempty" yaml:"additionalPropertiesAllowed,omitempty"`
	UniqueItemsDiff                 *ValueDiff         `json:"uniqueItems,omitempty" yaml:"uniqueItems,omitempty"`
	ExclusiveMinDiff                *ValueDiff         `json:"exclusiveMin,omitempty" yaml:"exclusiveMin,omitempty"`
	ExclusiveMaxDiff                *ValueDiff         `json:"exclusiveMax,omitempty" yaml:"exclusiveMax,omitempty"`
	NullableDiff                    *ValueDiff         `json:"nullable,omitempty" yaml:"nullable,omitempty"`
	ReadOnlyDiff                    *ValueDiff         `json:"readOnly,omitempty" yaml:"readOnly,omitempty"`
	WriteOnlyDiff                   *ValueDiff         `json:"writeOnly,omitempty" yaml:"writeOnly,omitempty"`
	AllowEmptyValueDiff             *ValueDiff         `json:"allowEmptyValue,omitempty" yaml:"allowEmptyValue,omitempty"`
	XMLDiff                         *ValueDiff         `json:"XML,omitempty" yaml:"XML,omitempty"`
	DeprecatedDiff                  *ValueDiff         `json:"deprecated,omitempty" yaml:"deprecated,omitempty"`
	MinDiff                         *ValueDiff         `json:"min,omitempty" yaml:"min,omitempty"`
	MaxDiff                         *ValueDiff         `json:"max,omitempty" yaml:"max,omitempty"`
	MultipleOfDiff                  *ValueDiff         `json:"multipleOf,omitempty" yaml:"multipleOf,omitempty"`
	MinLengthDiff                   *ValueDiff         `json:"minLength,omitempty" yaml:"minLength,omitempty"`
	MaxLengthDiff                   *ValueDiff         `json:"maxLength,omitempty" yaml:"maxLength,omitempty"`
	PatternDiff                     *ValueDiff         `json:"pattern,omitempty" yaml:"pattern,omitempty"`
	MinItemsDiff                    *ValueDiff         `json:"minItems,omitempty" yaml:"minItems,omitempty"`
	MaxItemsDiff                    *ValueDiff         `json:"maxItems,omitempty" yaml:"maxItems,omitempty"`
	ItemsDiff                       *SchemaDiff        `json:"items,omitempty" yaml:"items,omitempty"`
	RequiredDiff                    *StringsDiff       `json:"required,omitempty" yaml:"required,omitempty"`
	PropertiesDiff                  *SchemasDiff       `json:"properties,omitempty" yaml:"properties,omitempty"`
	MinPropsDiff                    *ValueDiff         `json:"minProps,omitempty" yaml:"minProps,omitempty"`
	MaxPropsDiff                    *ValueDiff         `json:"maxProps,omitempty" yaml:"maxProps,omitempty"`
	AdditionalPropertiesDiff        *SchemaDiff        `json:"additionalProperties,omitempty" yaml:"additionalProperties,omitempty"`
	DiscriminatorDiff               *DiscriminatorDiff `json:"discriminatorDiff,omitempty" yaml:"discriminatorDiff,omitempty"`
}

// Empty indicates whether a change was found in this element
func (diff *SchemaDiff) Empty() bool {
	return diff == nil || *diff == SchemaDiff{}
}

// Breaking indicates whether this element includes a breaking change
func (diff *SchemaDiff) Breaking() bool {
	if diff.Empty() {
		return false
	}

	return diff.SchemaAdded ||
		diff.OneOfDiff.Breaking() ||
		diff.AnyOfDiff.Breaking() ||
		diff.AllOfDiff.Breaking() ||
		diff.NotDiff.Breaking() ||
		diff.TypeDiff.Breaking() ||
		diff.FormatDiff.Breaking() ||
		diff.EnumDiff.Breaking() ||
		diff.DefaultDiff.Breaking() ||
		diff.AdditionalPropertiesAllowedDiff.CompareWithDefault(true, false, true) ||
		diff.UniqueItemsDiff.CompareWithDefault(false, true, false) || // TODO: check default value
		diff.ExclusiveMinDiff.CompareWithDefault(false, true, false) || // TODO: check default value
		diff.ExclusiveMaxDiff.CompareWithDefault(false, true, false) || // TODO: check default value
		diff.NullableDiff.CompareWithDefault(true, false, false) || // TODO: check default value
		diff.ReadOnlyDiff.CompareWithDefault(false, true, false) || // TODO: Relevant only for Schema "properties" definitions
		diff.WriteOnlyDiff.CompareWithDefault(false, true, false) || // TODO: Relevant only for Schema "properties" definitions
		diff.AllowEmptyValueDiff.CompareWithDefault(true, false, false) ||
		diff.XMLDiff.Breaking() ||
		diff.DeprecatedDiff.CompareWithDefault(false, true, false) ||
		diff.MinDiff.GreaterUInt64Ref() ||
		diff.MaxDiff.SmallerUInt64() ||
		diff.MultipleOfDiff.Breaking() ||
		diff.MinLengthDiff.GreaterUInt64Ref() ||
		diff.MaxLengthDiff.SmallerUInt64() ||
		diff.PatternDiff.Breaking() ||
		diff.MinItemsDiff.GreaterUInt64Ref() ||
		diff.MaxItemsDiff.SmallerUInt64() ||
		diff.ItemsDiff.Breaking() ||
		diff.RequiredDiff.Breaking() ||
		diff.PropertiesDiff.Breaking() ||
		diff.MinPropsDiff.GreaterUInt64Ref() ||
		diff.MaxPropsDiff.SmallerUInt64() ||
		diff.AdditionalPropertiesDiff.Breaking() ||
		diff.DiscriminatorDiff.Breaking()
}

func getSchemaDiff(config *Config, schema1, schema2 *openapi3.SchemaRef) (*SchemaDiff, error) {
	diff, err := getSchemaDiffInternal(config, schema1, schema2)
	if err != nil {
		return nil, err
	}

	if diff.Empty() {
		return nil, nil
	}

	if config.BreakingOnly && !diff.Breaking() {
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
	result.OneOfDiff, err = getSchemaListsDiff(config, value1.OneOf, value2.OneOf)
	if err != nil {
		return nil, err
	}
	result.AnyOfDiff, err = getSchemaListsDiff(config, value1.AnyOf, value2.AnyOf)
	if err != nil {
		return nil, err
	}
	result.AllOfDiff, err = getSchemaListsDiff(config, value1.AllOf, value2.AllOf)
	if err != nil {
		return nil, err
	}
	result.NotDiff, err = getSchemaDiff(config, value1.Not, value2.Not)
	if err != nil {
		return nil, err
	}
	result.TypeDiff = getValueDiff(config, true, value1.Type, value2.Type)
	result.TitleDiff = getValueDiff(config, false, value1.Title, value2.Title)
	result.FormatDiff = getValueDiff(config, true, value1.Format, value2.Format)
	result.DescriptionDiff = getValueDiffConditional(config, false, config.ExcludeDescription, value1.Description, value2.Description)
	result.EnumDiff = getEnumDiff(config, value1.Enum, value2.Enum)
	result.DefaultDiff = getValueDiff(config, true, value1.Default, value2.Default)
	result.ExampleDiff = getValueDiffConditional(config, false, config.ExcludeExamples, value1.Example, value2.Example)
	result.ExternalDocsDiff = getExternalDocsDiff(config, value1.ExternalDocs, value2.ExternalDocs)
	result.AdditionalPropertiesAllowedDiff = getBoolRefDiff(config, true, value1.AdditionalPropertiesAllowed, value2.AdditionalPropertiesAllowed)
	result.UniqueItemsDiff = getValueDiff(config, true, value1.UniqueItems, value2.UniqueItems)
	result.ExclusiveMinDiff = getValueDiff(config, true, value1.ExclusiveMin, value2.ExclusiveMin)
	result.ExclusiveMaxDiff = getValueDiff(config, true, value1.ExclusiveMax, value2.ExclusiveMax)
	result.NullableDiff = getValueDiff(config, true, value1.Nullable, value2.Nullable)
	result.ReadOnlyDiff = getValueDiff(config, true, value1.ReadOnly, value2.ReadOnly)
	result.WriteOnlyDiff = getValueDiff(config, true, value1.WriteOnly, value2.WriteOnly)
	result.AllowEmptyValueDiff = getValueDiff(config, true, value1.AllowEmptyValue, value2.AllowEmptyValue)
	result.XMLDiff = getValueDiff(config, true, value1.XML, value2.XML)
	result.DeprecatedDiff = getValueDiff(config, true, value1.Deprecated, value2.Deprecated)
	result.MinDiff = getFloat64RefDiff(config, true, value1.Min, value2.Min)
	result.MaxDiff = getFloat64RefDiff(config, true, value1.Max, value2.Max)
	result.MultipleOfDiff = getFloat64RefDiff(config, true, value1.MultipleOf, value2.MultipleOf)
	result.MinLengthDiff = getValueDiff(config, true, value1.MinLength, value2.MinLength)
	result.MaxLengthDiff = getUInt64RefDiff(config, true, value1.MaxLength, value2.MaxLength)
	result.PatternDiff = getValueDiff(config, true, value1.Pattern, value2.Pattern)
	// compiledPattern is derived from pattern -> no need to diff
	result.MinItemsDiff = getValueDiff(config, true, value1.MinItems, value2.MinItems)
	result.MaxItemsDiff = getUInt64RefDiff(config, true, value1.MaxItems, value2.MaxItems)
	result.ItemsDiff, err = getSchemaDiff(config, value1.Items, value2.Items)
	if err != nil {
		return nil, err
	}

	result.RequiredDiff = getStringsDiff(config, true, value1.Required, value2.Required)
	result.PropertiesDiff, err = getSchemasDiff(config, value1.Properties, value2.Properties)
	if err != nil {
		return nil, err
	}

	result.MinPropsDiff = getValueDiff(config, true, value1.MinProps, value2.MinProps)
	result.MaxPropsDiff = getUInt64RefDiff(config, true, value1.MaxProps, value2.MaxProps)
	result.AdditionalPropertiesDiff, err = getSchemaDiff(config, value1.AdditionalProperties, value2.AdditionalProperties)
	if err != nil {
		return nil, err
	}

	result.DiscriminatorDiff = getDiscriminatorDiff(config, value1.Discriminator, value2.Discriminator)

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

// Patch applies the patch to a schema
func (diff *SchemaDiff) Patch(schema *openapi3.Schema) error {

	if diff.Empty() {
		return nil
	}

	if err := diff.TypeDiff.patchString(&schema.Type); err != nil {
		return err
	}

	if err := diff.TitleDiff.patchString(&schema.Title); err != nil {
		return err
	}

	if err := diff.FormatDiff.patchString(&schema.Format); err != nil {
		return err
	}

	if err := diff.DescriptionDiff.patchString(&schema.Description); err != nil {
		return err
	}

	diff.EnumDiff.Patch(&schema.Enum)

	if err := diff.MaxLengthDiff.patchUInt64Ref(&schema.MaxLength); err != nil {
		return err
	}

	if err := patchPattern(diff.PatternDiff, schema); err != nil {
		return err
	}

	return nil
}

// patchPattern uses "Schema.WithPattern" to ensure that schema.compiledPattern is updated too
func patchPattern(valueDiff *ValueDiff, schema *openapi3.Schema) error {
	return valueDiff.patchStringCB(func(s string) { schema.WithPattern(valueDiff.To.(string)) })
}
