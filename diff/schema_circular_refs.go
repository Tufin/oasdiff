package diff

import "github.com/getkin/kin-openapi/openapi3"

type circularRefStatus int

const (
	circularRefStatusNone circularRefStatus = iota
	circularRefStatusDiff
	circularRefStatusNoDiff
)

const circularRefExtension = "x-oasdiff-circular-ref"

func getCircularRefsDiff(schema1, schema2 *openapi3.SchemaRef) circularRefStatus {

	if schema1 == nil || schema2 == nil ||
		schema1.Value == nil || schema2.Value == nil {
		return circularRefStatusNone
	}

	_, circular1 := schema1.Value.Extensions[circularRefExtension]
	_, circular2 := schema2.Value.Extensions[circularRefExtension]

	// neither are circular
	if !circular1 && !circular2 {
		return circularRefStatusNone
	}

	// one ref is circular but the other isn't
	if circular1 != circular2 {
		return circularRefStatusDiff
	}

	// now we know that both refs are circular

	// if they reference the same schema name, we consider them to be different
	if schema1.Ref != schema2.Ref {
		return circularRefStatusDiff
	}

	// otherwise, consider them equal
	return circularRefStatusNoDiff
}

func incRefCount(schema *openapi3.Schema) {
	if schema == nil {
		return
	}

	if schema.Extensions == nil {
		schema.Extensions = map[string]interface{}{}
	}

	if i, ok := schema.Extensions[circularRefExtension]; ok {
		schema.Extensions[circularRefExtension] = i.(int) + 1
	} else {
		schema.Extensions[circularRefExtension] = 1
	}
}

func decRefCount(schema *openapi3.Schema) {
	if schema == nil {
		return
	}

	if i, ok := schema.Extensions[circularRefExtension]; ok {
		if i == 1 {
			delete(schema.Extensions, circularRefExtension)
		} else {
			schema.Extensions[circularRefExtension] = i.(int) - 1
		}
	}
}
