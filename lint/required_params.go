package lint

import (
	"fmt"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/tufin/oasdiff/load"
)

func RequiredParamsCheck(source string, s *load.SpecInfo) []*Error {
	result := make([]*Error, 0)

	if s == nil || s.Spec == nil {
		return result
	}

	for path, pathItem := range s.Spec.Paths.Map() {
		for _, parameter := range pathItem.Parameters {
			if !parameter.Value.Required {
				continue
			}

			if parameter.Value.Schema != nil && parameter.Value.Schema.Value.Default != nil {
				result = append(result, &Error{
					Id:     "required-param-with-default",
					Level:  LEVEL_ERROR,
					Text:   fmt.Sprintf("required path parameter %q shouldn't have a default value: %s", parameter.Value.Name, path),
					Source: source,
				})
			}
		}
		for method, op := range pathItem.Operations() {
			result = append(result, checkOperationRequiredParams(path, method, op, source)...)
		}
	}

	return result
}

func checkOperationRequiredParams(path, method string, op *openapi3.Operation, source string) []*Error {
	result := make([]*Error, 0)

	for _, parameter := range op.Parameters {
		if !parameter.Value.Required {
			continue
		}

		if parameter.Value.Schema != nil && parameter.Value.Schema.Value.Default != nil {
			result = append(result, &Error{
				Id:     "required-param-with-default",
				Level:  LEVEL_ERROR,
				Text:   fmt.Sprintf("required path parameter %q shouldn't have a default value: %s %s", parameter.Value.Name, method, path),
				Source: source,
			})
		}
	}

	return result
}
