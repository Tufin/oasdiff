package lint

import (
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/tufin/oasdiff/load"
	"github.com/tufin/oasdiff/utils"
)

type state struct {
	source      string
	cache       map[string]error
	visitedRefs utils.VisitedRefs
}

func newState(source string) *state {
	return &state{
		source:      source,
		cache:       map[string]error{},
		visitedRefs: utils.VisitedRefs{},
	}
}

func SchemaCheck(source string, spec *load.SpecInfo) []*Error {
	result := make([]*Error, 0)

	if spec == nil || spec.Spec == nil {
		return result
	}

	s := newState(source)

	for _, path := range spec.Spec.Paths.Map() {
		result = append(result, checkParameters(path.Parameters, s)...)
		result = append(result, checkOperations(path.Operations(), s)...)
	}

	return result
}

func checkOperations(operations map[string]*openapi3.Operation, s *state) []*Error {
	result := make([]*Error, 0)
	for _, op := range operations {

		result = append(result, checkParameters(op.Parameters, s)...)

		if op.RequestBody != nil {
			for _, mediaType := range op.RequestBody.Value.Content {
				result = append(result, checkSchemaRef(mediaType.Schema, s)...)
			}
		}

		for _, response := range op.Responses.Map() {
			for _, mediaType := range response.Value.Content {
				result = append(result, checkSchemaRef(mediaType.Schema, s)...)
			}
			for _, header := range response.Value.Headers {
				result = append(result, checkSchemaRef(header.Value.Schema, s)...)
			}
		}

		for _, callback := range op.Callbacks {
			for _, pathItem := range callback.Value.Map() {
				result = append(result, checkParameters(pathItem.Parameters, s)...)
				result = append(result, checkOperations(pathItem.Operations(), s)...)
			}
		}
	}
	return result
}

func checkParameters(parameters openapi3.Parameters, s *state) []*Error {
	result := make([]*Error, 0)
	for _, parameter := range parameters {
		if parameter.Value == nil {
			continue
		}
		if parameter.Value.Schema != nil {
			result = append(result, checkSchemaRef(parameter.Value.Schema, s)...)
		}
		for _, mediaType := range parameter.Value.Content {
			if mediaType.Schema != nil {
				result = append(result, checkSchemaRef(mediaType.Schema, s)...)
			}
		}
	}
	return result
}

func checkSchema(schema *openapi3.Schema, s *state) []*Error {
	result := make([]*Error, 0)

	result = append(result, runCheckers(schema, s)...)

	for _, subSchema := range schema.OneOf {
		result = append(result, checkSchemaRef(subSchema, s)...)
	}
	for _, subSchema := range schema.AnyOf {
		result = append(result, checkSchemaRef(subSchema, s)...)
	}
	for _, subSchema := range schema.AllOf {
		result = append(result, checkSchemaRef(subSchema, s)...)
	}
	if schema.Not != nil {
		result = append(result, checkSchemaRef(schema.Not, s)...)
	}
	if schema.Items != nil {
		result = append(result, checkSchemaRef(schema.Items, s)...)
	}
	for _, subSchema := range schema.Properties {
		result = append(result, checkSchemaRef(subSchema, s)...)
	}
	if schema.AdditionalProperties.Schema != nil {
		result = append(result, checkSchemaRef(schema.AdditionalProperties.Schema, s)...)
	}
	return result
}

func checkSchemaRef(schema *openapi3.SchemaRef, s *state) []*Error {
	if s.visitedRefs.IsVisited(schema.Ref) {
		return nil
	}
	// mark visited schema references to avoid infinite loops
	if schema.Ref != "" {
		s.visitedRefs.Add(schema.Ref)
		defer s.visitedRefs.Remove(schema.Ref)
	}

	return checkSchema(schema.Value, s)
}

func runCheckers(schema *openapi3.Schema, s *state) []*Error {
	result := make([]*Error, 0)

	if err := checkRegex(schema.Pattern, s); err != nil {
		result = append(result, err)
	}

	if err := checkRequireProperties(schema, s); err != nil {
		result = append(result, err)
	}

	return result
}
