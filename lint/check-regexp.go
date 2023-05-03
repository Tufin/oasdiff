package lint

import (
	"regexp"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/tufin/oasdiff/load"
)

type regexCtx struct {
	source string
	cache  map[string]error
}

func newRegexCtx(source string) *regexCtx {
	return &regexCtx{
		source: source,
		cache:  map[string]error{},
	}
}

func (context *regexCtx) test(pattern string) error {
	if result, ok := context.cache[pattern]; ok {
		return result
	}
	_, err := regexp.Compile(pattern)
	context.cache[pattern] = err
	return err
}

// *** THIS IS A TEMPORARY IMPLEMENTATION ***
// SHOULD USE ECMA 262, SEE: https://swagger.io/docs/specification/data-models/data-types/#pattern

func RegexCheck(source string, s *load.OpenAPISpecInfo) []*Error {

	result := make([]*Error, 0)

	if s == nil || s.Spec == nil {
		return result
	}

	context := newRegexCtx(source)

	for _, path := range s.Spec.Paths {
		result = append(result, checkParameters(path.Parameters, context)...)
		result = append(result, checkOperations(path.Operations(), context)...)
	}

	return result
}

func checkOperations(operations map[string]*openapi3.Operation, context *regexCtx) []*Error {
	result := make([]*Error, 0)
	for _, op := range operations {

		result = append(result, checkParameters(op.Parameters, context)...)

		if op.RequestBody != nil {
			for _, mediaType := range op.RequestBody.Value.Content {
				result = append(result, checkSchemaRef(mediaType.Schema, context)...)
			}
		}

		for _, response := range op.Responses {
			for _, mediaType := range response.Value.Content {
				result = append(result, checkSchemaRef(mediaType.Schema, context)...)
			}
			for _, header := range response.Value.Headers {
				result = append(result, checkSchemaRef(header.Value.Schema, context)...)
			}
		}

		for _, callback := range op.Callbacks {
			for _, pathItem := range *callback.Value {
				result = append(result, checkParameters(pathItem.Parameters, context)...)
				result = append(result, checkOperations(pathItem.Operations(), context)...)
			}
		}
	}
	return result
}

func checkParameters(parameters openapi3.Parameters, context *regexCtx) []*Error {
	result := make([]*Error, 0)
	for _, parameter := range parameters {
		if parameter.Value == nil {
			continue
		}
		if parameter.Value.Schema != nil {
			result = append(result, checkSchemaRef(parameter.Value.Schema, context)...)
		}
		for _, mediaType := range parameter.Value.Content {
			if mediaType.Schema != nil {
				result = append(result, checkSchemaRef(mediaType.Schema, context)...)
			}
		}
	}
	return result
}

func checkSchema(schema *openapi3.Schema, context *regexCtx) []*Error {
	result := make([]*Error, 0)

	if err := checkRegex(schema.Pattern, context); err != nil {
		result = append(result, err)
	}

	for _, subSchema := range schema.OneOf {
		result = append(result, checkSchemaRef(subSchema, context)...)
	}
	for _, subSchema := range schema.AnyOf {
		result = append(result, checkSchemaRef(subSchema, context)...)
	}
	for _, subSchema := range schema.AllOf {
		result = append(result, checkSchemaRef(subSchema, context)...)
	}
	if schema.Not != nil {
		result = append(result, checkSchemaRef(schema.Not, context)...)
	}
	if schema.Items != nil {
		result = append(result, checkSchemaRef(schema.Items, context)...)
	}
	for _, subSchema := range schema.Properties {
		result = append(result, checkSchemaRef(subSchema, context)...)
	}
	if schema.AdditionalProperties.Schema != nil {
		result = append(result, checkSchemaRef(schema.AdditionalProperties.Schema, context)...)
	}
	return result
}

func checkSchemaRef(schema *openapi3.SchemaRef, context *regexCtx) []*Error {
	return checkSchema(schema.Value, context)
}

func checkRegex(pattern string, context *regexCtx) *Error {
	if pattern == "" {
		return nil
	}

	if err := context.test(pattern); err != nil {
		return &Error{
			Id:     "invalid-regex-pattern",
			Level:  LEVEL_ERROR,
			Text:   err.Error(),
			Source: context.source,
		}
	}

	return nil
}
