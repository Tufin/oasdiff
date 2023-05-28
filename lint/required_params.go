package lint

import (
	"fmt"

	"github.com/tufin/oasdiff/load"
)

func RequiredParamsCheck(source string, s *load.OpenAPISpecInfo) []*Error {
	result := make([]*Error, 0)

	if s == nil || s.Spec == nil {
		return result
	}

	for path, pathItem := range s.Spec.Paths {
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
	}

	return result
}
