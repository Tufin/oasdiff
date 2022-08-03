package diff

import (
	"errors"
	"github.com/getkin/kin-openapi/openapi3"
)

// SchemaDiff describes the changes between a pair of schema objects: https://swagger.io/specification/#schema-object
type SchemaDiff struct {
	SchemaAdded                     bool                    `json:"schemaAdded,omitempty" yaml:"schemaAdded,omitempty"`
	SchemaDeleted                   bool                    `json:"schemaDeleted,omitempty" yaml:"schemaDeleted,omitempty"`
	CircularRefDiff                 bool                    `json:"circularRef,omitempty" yaml:"circularRef,omitempty"`
	ExtensionsDiff                  *ExtensionsDiff         `json:"extensions,omitempty" yaml:"extensions,omitempty"`
	OneOfDiff                       *SchemaListDiff         `json:"oneOf,omitempty" yaml:"oneOf,omitempty"`
	AnyOfDiff                       *SchemaListDiff         `json:"anyOf,omitempty" yaml:"anyOf,omitempty"`
	AllOfDiff                       *SchemaListDiff         `json:"allOf,omitempty" yaml:"allOf,omitempty"`
	NotDiff                         *SchemaDiff             `json:"not,omitempty" yaml:"not,omitempty"`
	TypeDiff                        *ValueDiff              `json:"type,omitempty" yaml:"type,omitempty"`
	TitleDiff                       *ValueDiff              `json:"title,omitempty" yaml:"title,omitempty"`
	FormatDiff                      *ValueDiff              `json:"format,omitempty" yaml:"format,omitempty"`
	DescriptionDiff                 *ValueDiff              `json:"description,omitempty" yaml:"description,omitempty"`
	EnumDiff                        *EnumDiff               `json:"enum,omitempty" yaml:"enum,omitempty"`
	DefaultDiff                     *ValueDiff              `json:"default,omitempty" yaml:"default,omitempty"`
	ExampleDiff                     *ValueDiff              `json:"example,omitempty" yaml:"example,omitempty"`
	ExternalDocsDiff                *ExternalDocsDiff       `json:"externalDocs,omitempty" yaml:"externalDocs,omitempty"`
	AdditionalPropertiesAllowedDiff *ValueDiff              `json:"additionalPropertiesAllowed,omitempty" yaml:"additionalPropertiesAllowed,omitempty"`
	UniqueItemsDiff                 *ValueDiff              `json:"uniqueItems,omitempty" yaml:"uniqueItems,omitempty"`
	ExclusiveMinDiff                *ValueDiff              `json:"exclusiveMin,omitempty" yaml:"exclusiveMin,omitempty"`
	ExclusiveMaxDiff                *ValueDiff              `json:"exclusiveMax,omitempty" yaml:"exclusiveMax,omitempty"`
	NullableDiff                    *ValueDiff              `json:"nullable,omitempty" yaml:"nullable,omitempty"`
	ReadOnlyDiff                    *ValueDiff              `json:"readOnly,omitempty" yaml:"readOnly,omitempty"`
	WriteOnlyDiff                   *ValueDiff              `json:"writeOnly,omitempty" yaml:"writeOnly,omitempty"`
	AllowEmptyValueDiff             *ValueDiff              `json:"allowEmptyValue,omitempty" yaml:"allowEmptyValue,omitempty"`
	XMLDiff                         *ValueDiff              `json:"XML,omitempty" yaml:"XML,omitempty"`
	DeprecatedDiff                  *ValueDiff              `json:"deprecated,omitempty" yaml:"deprecated,omitempty"`
	MinDiff                         *ValueDiff              `json:"min,omitempty" yaml:"min,omitempty"`
	MaxDiff                         *ValueDiff              `json:"max,omitempty" yaml:"max,omitempty"`
	MultipleOfDiff                  *ValueDiff              `json:"multipleOf,omitempty" yaml:"multipleOf,omitempty"`
	MinLengthDiff                   *ValueDiff              `json:"minLength,omitempty" yaml:"minLength,omitempty"`
	MaxLengthDiff                   *ValueDiff              `json:"maxLength,omitempty" yaml:"maxLength,omitempty"`
	PatternDiff                     *ValueDiff              `json:"pattern,omitempty" yaml:"pattern,omitempty"`
	MinItemsDiff                    *ValueDiff              `json:"minItems,omitempty" yaml:"minItems,omitempty"`
	MaxItemsDiff                    *ValueDiff              `json:"maxItems,omitempty" yaml:"maxItems,omitempty"`
	ItemsDiff                       *SchemaDiff             `json:"items,omitempty" yaml:"items,omitempty"`
	RequiredDiff                    *RequiredPropertiesDiff `json:"required,omitempty" yaml:"required,omitempty"`
	PropertiesDiff                  *SchemasDiff            `json:"properties,omitempty" yaml:"properties,omitempty"`
	MinPropsDiff                    *ValueDiff              `json:"minProps,omitempty" yaml:"minProps,omitempty"`
	MaxPropsDiff                    *ValueDiff              `json:"maxProps,omitempty" yaml:"maxProps,omitempty"`
	AdditionalPropertiesDiff        *SchemaDiff             `json:"additionalProperties,omitempty" yaml:"additionalProperties,omitempty"`
	DiscriminatorDiff               *DiscriminatorDiff      `json:"discriminatorDiff,omitempty" yaml:"discriminatorDiff,omitempty"`
}

// Empty indicates whether a change was found in this element
func (diff *SchemaDiff) Empty() bool {
	return diff == nil || *diff == SchemaDiff{}
}

func (diff *SchemaDiff) removeNonBreaking(state *state, schema2 *openapi3.SchemaRef) {

	if diff.Empty() {
		return
	}

	diff.SchemaDeleted = false
	diff.TitleDiff = nil
	diff.ExtensionsDiff = nil
	diff.DescriptionDiff = nil
	diff.ExampleDiff = nil
	diff.ExternalDocsDiff = nil

	if !diff.UniqueItemsDiff.CompareWithDefault(false, true, false) { // TODO: check default value
		diff.UniqueItemsDiff = nil
	}

	if !diff.NullableDiff.CompareWithDefault(true, false, false) { // TODO: check default value
		diff.NullableDiff = nil
	}

	if !diff.ReadOnlyDiff.CompareWithDefault(false, true, false) { // TODO: Relevant only for Schema "properties" definitions
		diff.ReadOnlyDiff = nil
	}

	if !diff.WriteOnlyDiff.CompareWithDefault(false, true, false) { // TODO: Relevant only for Schema "properties" definitions
		diff.WriteOnlyDiff = nil
	}

	if !diff.AllowEmptyValueDiff.CompareWithDefault(true, false, false) {
		diff.AllowEmptyValueDiff = nil
	}

	diff.DeprecatedDiff = nil

	// Number
	if !diff.MinDiff.minBreakingFloat64() { // *float64
		diff.MinDiff = nil
	}

	if !diff.MaxDiff.maxBreakingFloat64() { // *float64
		diff.MaxDiff = nil
	}

	// TODO: check MultipleOf

	if !diff.ExclusiveMinDiff.CompareWithDefault(false, true, false) { // TODO: check default value
		diff.ExclusiveMinDiff = nil
	}

	if !diff.ExclusiveMaxDiff.CompareWithDefault(false, true, false) { // TODO: check default value
		diff.ExclusiveMaxDiff = nil
	}

	// String
	if !diff.MinLengthDiff.minBreakingUInt64() { // uint64
		diff.MinLengthDiff = nil
	}

	if !diff.MaxLengthDiff.maxBreakingUInt64() { // *uint64
		diff.MaxLengthDiff = nil
	}

	// Array
	if !diff.MinItemsDiff.minBreakingUInt64() { // uint64
		diff.MinItemsDiff = nil
	}

	if !diff.MaxItemsDiff.maxBreakingUInt64() { // *uint64
		diff.MaxItemsDiff = nil
	}

	// Object
	diff.removeChangedButNonRequiredProperties(state, schema2)

	if !diff.AdditionalPropertiesAllowedDiff.CompareWithDefault(true, false, true) {
		diff.AdditionalPropertiesAllowedDiff = nil
	}

	if !diff.MinPropsDiff.minBreakingUInt64() { // uint64
		diff.MinPropsDiff = nil
	}

	if !diff.MaxPropsDiff.maxBreakingUInt64() { // *uint64
		diff.MaxPropsDiff = nil
	}
}

// removeChangedButNonRequiredProperties deletes non-required property changes that don't break client
// In request: remove added but non-required properties
// In response: remove deleted but non-required properties
func (diff *SchemaDiff) removeChangedButNonRequiredProperties(state *state, schema2 *openapi3.SchemaRef) {
	if diff.Empty() || diff.PropertiesDiff.Empty() {
		return
	}

	if schema2 == nil || schema2.Value == nil {
		return
	}

	requiredMap := StringList(schema2.Value.Required).toStringSet()
	changedSet := diff.PropertiesDiff.getBreakingSetByDirection(state.direction)

	newList := StringList{}
	for _, property := range *changedSet {
		if _, ok := requiredMap[property]; ok {
			newList = append(newList, property)
		}
	}
	*changedSet = newList

	if diff.PropertiesDiff.Empty() {
		diff.PropertiesDiff = nil
	}
}

func getSchemaDiff(config *Config, state *state, schema1, schema2 *openapi3.SchemaRef) (*SchemaDiff, error) {

	if diff, ok := state.cache.get(state.direction, schema1, schema2); ok {
		return diff, nil
	}

	diff, err := getSchemaDiffInternal(config, state, schema1, schema2)
	if err != nil {
		return nil, err
	}

	if config.BreakingOnly {
		diff.removeNonBreaking(state, schema2)
	}

	if diff.Empty() {
		diff = nil
	}

	state.cache.add(state.direction, schema1, schema2, diff)
	return diff, nil
}

func getSchemaDiffInternal(config *Config, state *state, schema1, schema2 *openapi3.SchemaRef) (*SchemaDiff, error) {
	if status := getSchemaStatus(schema1, schema2); status != schemaStatusOK {
		switch status {
		case schemaStatusNoSchemas:
			return nil, nil
		case schemaStatusSchemaAdded:
			return &SchemaDiff{SchemaAdded: true}, nil
		case schemaStatusSchemaDeleted:
			return &SchemaDiff{SchemaDeleted: true}, nil
		}
	}

	if status := getCircularRefsDiff(state.visitedSchemasBase, state.visitedSchemasRevision, schema1, schema2); status != circularRefStatusNone {
		switch status {
		case circularRefStatusDiff:
			return &SchemaDiff{CircularRefDiff: true}, nil
		case circularRefStatusNoDiff:
			return nil, nil
		}
	}

	value1, err := derefSchema(schema1)
	if err != nil {
		return nil, err
	}

	value2, err := derefSchema(schema2)
	if err != nil {
		return nil, err
	}

	// mark visited schema references to avoid infinite loops
	if schema1.Ref != "" {
		state.visitedSchemasBase.add(schema1.Ref)
		defer state.visitedSchemasBase.remove(schema1.Ref)
	}

	if schema2.Ref != "" {
		state.visitedSchemasRevision.add(schema2.Ref)
		defer state.visitedSchemasRevision.remove(schema2.Ref)
	}

	result := SchemaDiff{}

	result.ExtensionsDiff = getExtensionsDiff(config, state, value1.ExtensionProps, value2.ExtensionProps)
	result.OneOfDiff, err = getSchemaListsDiff(config, state, value1.OneOf, value2.OneOf)
	if err != nil {
		return nil, err
	}
	result.AnyOfDiff, err = getSchemaListsDiff(config, state, value1.AnyOf, value2.AnyOf)
	if err != nil {
		return nil, err
	}
	result.AllOfDiff, err = getSchemaListsDiff(config, state, value1.AllOf, value2.AllOf)
	if err != nil {
		return nil, err
	}
	result.NotDiff, err = getSchemaDiff(config, state, value1.Not, value2.Not)
	if err != nil {
		return nil, err
	}
	result.TypeDiff = getValueDiff(value1.Type, value2.Type)
	result.TitleDiff = getValueDiff(value1.Title, value2.Title)
	result.FormatDiff = getValueDiff(value1.Format, value2.Format)
	result.DescriptionDiff = getValueDiffConditional(config.ExcludeDescription, value1.Description, value2.Description)
	result.EnumDiff = getEnumDiff(config, state, value1.Enum, value2.Enum)
	result.DefaultDiff = getValueDiff(value1.Default, value2.Default)
	result.ExampleDiff = getValueDiffConditional(config.ExcludeExamples, value1.Example, value2.Example)
	result.ExternalDocsDiff = getExternalDocsDiff(config, state, value1.ExternalDocs, value2.ExternalDocs)
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
	result.MultipleOfDiff = getFloat64RefDiff(value1.MultipleOf, value2.MultipleOf)
	result.MinLengthDiff = getValueDiff(value1.MinLength, value2.MinLength)
	result.MaxLengthDiff = getUInt64RefDiff(value1.MaxLength, value2.MaxLength)
	result.PatternDiff = getValueDiff(value1.Pattern, value2.Pattern)
	// compiledPattern is derived from pattern -> no need to diff
	result.MinItemsDiff = getValueDiff(value1.MinItems, value2.MinItems)
	result.MaxItemsDiff = getUInt64RefDiff(value1.MaxItems, value2.MaxItems)
	result.ItemsDiff, err = getSchemaDiff(config, state, value1.Items, value2.Items)
	if err != nil {
		return nil, err
	}

	// Object
	result.RequiredDiff = getRequiredPropertiesDiff(config, state, value1, value2)
	result.PropertiesDiff, err = getSchemasDiff(config, state, value1.Properties, value2.Properties)
	if err != nil {
		return nil, err
	}

	result.MinPropsDiff = getValueDiff(value1.MinProps, value2.MinProps)
	result.MaxPropsDiff = getUInt64RefDiff(value1.MaxProps, value2.MaxProps)
	result.AdditionalPropertiesDiff, err = getSchemaDiff(config, state, value1.AdditionalProperties, value2.AdditionalProperties)
	if err != nil {
		return nil, err
	}

	result.DiscriminatorDiff = getDiscriminatorDiff(config, state, value1.Discriminator, value2.Discriminator)

	return &result, nil
}

type schemaStatus int

const (
	schemaStatusOK schemaStatus = iota
	schemaStatusNoSchemas
	schemaStatusSchemaAdded
	schemaStatusSchemaDeleted
	schemaStatusCircularRefDiff
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
		return nil, errors.New("schema reference is nil")
	}

	return ref.Value, nil
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
