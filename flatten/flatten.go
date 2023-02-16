package flatten

import (
	"github.com/getkin/kin-openapi/openapi3"
)

// Handle replaces objects under AllOf with a flattened equivalent
func Handle(schema openapi3.Schema) *openapi3.Schema {
	if !isListOfObjects(&schema) {
		return &schema
	}

	result := schema
	result.AllOf = openapi3.SchemaRefs{schema.AllOf[0]}
	for _, schema := range schema.AllOf[1:] {
		add(&result, schema.Value)
	}

	return &result
}

func isListOfObjects(schema *openapi3.Schema) bool {
	if schema == nil || schema.AllOf == nil {
		return false
	}

	for _, subSchema := range schema.AllOf {
		if subSchema.Value.Type != "object" {
			return false
		}
	}

	return true
}

func add(result, schema *openapi3.Schema) {
	if schema == nil || schema.Properties == nil {
		return
	}

	for name, schema := range schema.Properties {
		result.AllOf[0].Value.Properties[name] = schema
	}
}
