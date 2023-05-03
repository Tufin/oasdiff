package lint

import (
	"github.com/tufin/oasdiff/load"
)

func PathParamsCheck(source string, s *load.OpenAPISpecInfo) []Error {
	result := make([]Error, 0)

	if s == nil || s.Spec == nil {
		return result
	}

	// for path, pathItem := range s.Spec.Paths {
	// _, _, pathParams := utils.NormalizeTemplatedPath(path)

	// actualPathParams := []string{}
	// for _, parameter := range pathItem.Parameters {
	// 	if parameter.Value.In != openapi3.ParameterInPath {
	// 		continue
	// 	}
	// 	actualPathParams = append(actualPathParams, parameter.Value.Name)
	// }

	// if len(utils.StringList(pathParams).Minus(actualPathParams)) > 0 {

	// }
	// }

	return result
}
