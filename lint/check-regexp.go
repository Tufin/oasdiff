package lint

import (
	"regexp"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/tufin/oasdiff/load"
)

// *** THIS IS A TEMPORARY IMPLEMENTATION ***
// SHOULD USE ECMA 262, SEE: https://swagger.io/docs/specification/data-models/data-types/#pattern

func RegexCheck(source string, s *load.OpenAPISpecInfo) []*Error {

	result := make([]*Error, 0)

	if s == nil || s.Spec == nil {
		return result
	}

	for _, path := range s.Spec.Paths {
		result = append(result, checkParameters(path.Parameters, source)...)
		result = append(result, checkOperations(path.Operations(), source)...)
	}

	return result
}

func checkOperations(operations map[string]*openapi3.Operation, source string) []*Error {
	result := make([]*Error, 0)
	for _, op := range operations {

		result = append(result, checkParameters(op.Parameters, source)...)

		if op.RequestBody != nil {
			for _, mediaType := range op.RequestBody.Value.Content {
				result = append(result, checkSchema(mediaType.Schema, source)...)
			}
		}

		for _, response := range op.Responses {
			for _, mediaType := range response.Value.Content {
				result = append(result, checkSchema(mediaType.Schema, source)...)
			}
			for _, header := range response.Value.Headers {
				result = append(result, checkSchema(header.Value.Schema, source)...)
			}
		}

		for _, callback := range op.Callbacks {
			for _, pathItem := range *callback.Value {
				result = append(result, checkParameters(pathItem.Parameters, source)...)
				result = append(result, checkOperations(pathItem.Operations(), source)...)
			}
		}
	}
	return result
}

func checkParameters(parameters openapi3.Parameters, source string) []*Error {
	result := make([]*Error, 0)
	for _, parameter := range parameters {
		if parameter.Value == nil {
			continue
		}
		if parameter.Value.Schema != nil {
			result = append(result, checkSchema(parameter.Value.Schema, source)...)
		}
		for _, mediaType := range parameter.Value.Content {
			result = append(result, checkSchema(mediaType.Schema, source)...)
		}
	}
	return result
}

func checkSchema(schema *openapi3.SchemaRef, source string) []*Error {
	result := make([]*Error, 0)
	if err := checkRegex(schema.Value.Pattern, source); err != nil {
		result = append(result, err)
	}
	for _, schema := range schema.Value.Properties {
		result = append(result, checkSchema(schema, source)...)
	}
	return result
}

func checkRegex(pattern string, source string) *Error {
	if pattern == "" {
		return nil
	}

	if _, err := regexp.Compile(pattern); err != nil {
		return &Error{
			Id:     "invalid-regex-pattern",
			Level:  LEVEL_ERROR,
			Text:   err.Error(),
			Source: source,
		}
	}

	return nil
}
