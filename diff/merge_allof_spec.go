package diff

import (
	"github.com/getkin/kin-openapi/openapi3"
)

func MergeSpec(spec openapi3.T) (openapi3.T, error) {
	schemas, err := mergeSchemas(spec.Components.Schemas)
	if err != nil {
		return spec, err
	}
	spec.Components.Schemas = schemas
	parameters, err := mergeParametersMap(spec.Components.Parameters)
	if err != nil {
		return spec, err
	}
	spec.Components.Parameters = parameters
	headers, err := mergeHeaders(spec.Components.Headers)
	if err != nil {
		return spec, err
	}
	spec.Components.Headers = headers
	requestBodies, err := mergeRequestBodies(spec.Components.RequestBodies)
	if err != nil {
		return spec, err
	}
	spec.Components.RequestBodies = requestBodies
	responses, err := mergeResponses(spec.Components.Responses)
	if err != nil {
		return spec, err
	}
	spec.Components.Responses = responses
	callbacks, err := mergeCallbacks(spec.Components.Callbacks)
	if err != nil {
		return spec, err
	}
	spec.Components.Callbacks = callbacks

	for _, v := range spec.Paths {
		if v == nil {
			continue
		}
		pathItem, err := mergePathItem(*v)
		if err != nil {
			return spec, err
		}
		*v = pathItem
	}
	return spec, nil
}

func mergeOperation(operation openapi3.Operation) (openapi3.Operation, error) {
	parameteres, err := mergeParameters(operation.Parameters)
	if err != nil {
		return operation, err
	}
	operation.Parameters = parameteres
	if operation.RequestBody != nil && operation.RequestBody.Value != nil {
		content, err := mergeContent(operation.RequestBody.Value.Content)
		if err != nil {
			return operation, err
		}
		operation.RequestBody.Value.Content = content
	}
	responses, err := mergeResponses(operation.Responses)
	if err != nil {
		return operation, err
	}
	operation.Responses = responses
	callbacks, err := mergeCallbacks(operation.Callbacks)
	if err != nil {
		return operation, err
	}
	operation.Callbacks = callbacks
	return operation, nil
}

func mergePathItem(pathItem openapi3.PathItem) (openapi3.PathItem, error) {
	operations := []*openapi3.Operation{
		pathItem.Connect, pathItem.Delete, pathItem.Get, pathItem.Head,
		pathItem.Options, pathItem.Patch, pathItem.Post, pathItem.Put, pathItem.Trace,
	}

	for _, op := range operations {
		if op != nil {
			mergedOp, err := mergeOperation(*op)
			if err != nil {
				return pathItem, err
			}
			*op = mergedOp
		}
	}

	parameters, err := mergeParameters(pathItem.Parameters)
	if err != nil {
		return pathItem, err
	}
	pathItem.Parameters = parameters
	return pathItem, nil
}

func mergeCallbacks(callbacks openapi3.Callbacks) (openapi3.Callbacks, error) {
	for _, v := range callbacks {
		for _, pathItem := range *v.Value {
			m, err := mergePathItem(*pathItem)
			if err != nil {
				return callbacks, err
			}
			*pathItem = m
		}
	}
	return callbacks, nil
}

func mergeSchemas(schemas openapi3.Schemas) (openapi3.Schemas, error) {
	for _, s := range schemas {
		if s == nil || s.Value == nil {
			continue
		}
		m, err := Merge(*s.Value)
		if err != nil {
			return schemas, err
		}
		s.Value = m
	}
	return schemas, nil
}

func mergeResponses(responses openapi3.Responses) (openapi3.Responses, error) {
	for _, v := range responses {
		if v == nil || v.Value == nil {
			continue
		}
		content, err := mergeContent(v.Value.Content)
		if err != nil {
			return responses, err
		}
		v.Value.Content = content
		headers, err := mergeHeaders(v.Value.Headers)
		if err != nil {
			return responses, err
		}
		v.Value.Headers = headers
	}
	return responses, nil
}

func mergeRequestBodies(rb openapi3.RequestBodies) (openapi3.RequestBodies, error) {
	for _, v := range rb {
		if v == nil || v.Value == nil {
			continue
		}
		content, err := mergeContent(v.Value.Content)
		if err != nil {
			return rb, err
		}
		v.Value.Content = content
	}
	return rb, nil
}

func mergeParameters(parameters openapi3.Parameters) (openapi3.Parameters, error) {
	for _, v := range parameters {
		if v == nil || v.Value == nil {
			continue
		}
		parameter, err := mergeParameter(*v.Value)
		if err != nil {
			return parameters, err
		}
		*v.Value = parameter
	}
	return parameters, nil
}

func mergeParametersMap(parameters openapi3.ParametersMap) (openapi3.ParametersMap, error) {
	for _, v := range parameters {
		if v == nil || v.Value == nil {
			continue
		}
		parameter, err := mergeParameter(*v.Value)
		if err != nil {
			return parameters, err
		}
		*v.Value = parameter
	}
	return parameters, nil
}

func mergeParameter(p openapi3.Parameter) (openapi3.Parameter, error) {
	if p.Schema == nil || p.Schema.Value == nil {
		return p, nil
	}
	m, err := Merge(*p.Schema.Value)
	if err != nil {
		return p, err
	}
	p.Schema.Value = m
	content, err := mergeContent(p.Content)
	if err != nil {
		return p, err
	}
	p.Content = content
	return p, nil
}

func mergeContent(content openapi3.Content) (openapi3.Content, error) {
	for _, mediaType := range content {
		if mediaType == nil || mediaType.Schema == nil || mediaType.Schema.Value == nil {
			continue
		}
		m, err := Merge(*mediaType.Schema.Value)
		if err != nil {
			return content, err
		}
		mediaType.Schema.Value = m
		for _, encoding := range mediaType.Encoding {
			if encoding == nil {
				continue
			}
			headers, err := mergeHeaders(encoding.Headers)
			if err != nil {
				return content, err
			}
			encoding.Headers = headers
		}
	}
	return content, nil
}

func mergeHeaders(headers openapi3.Headers) (openapi3.Headers, error) {
	for _, v := range headers {
		if v == nil || v.Value == nil {
			continue
		}
		parameter, err := mergeParameter(v.Value.Parameter)
		if err != nil {
			return headers, err
		}
		v.Value.Parameter = parameter
	}
	return headers, nil
}
