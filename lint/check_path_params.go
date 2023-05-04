package lint

import (
	"fmt"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/tufin/oasdiff/load"
	"github.com/tufin/oasdiff/utils"
)

func PathParamsCheck(source string, s *load.OpenAPISpecInfo) []*Error {
	result := make([]*Error, 0)

	if s == nil || s.Spec == nil {
		return result
	}

	for path, pathItem := range s.Spec.Paths {
		pathParamsFromURL := getPathParamsFromURL(path)

		pathParams := utils.StringSet{}
		for _, parameter := range pathItem.Parameters {
			if parameter.Value.In != openapi3.ParameterInPath {
				continue
			}
			pathParams.Add(parameter.Value.Name)
		}

		for method, op := range pathItem.Operations() {

			opParams := utils.StringSet{}
			for _, parameter := range op.Parameters {
				if parameter.Value.In != openapi3.ParameterInPath {
					continue
				}
				opParams.Add(parameter.Value.Name)
			}

			for param := range opParams.Minus(pathParamsFromURL) {
				result = append(result, &Error{
					Id:     "path-param-mismatch",
					Level:  LEVEL_ERROR,
					Text:   fmt.Sprintf("path param %s doesn't appear in URL: %s %s", param, method, path),
					Source: source,
				})
			}

			for param := range pathParamsFromURL.Minus(opParams) {
				result = append(result, &Error{
					Id:     "path-param-mismatch",
					Level:  LEVEL_ERROR,
					Text:   fmt.Sprintf("path param %s appears in URL but is unused: %s %s", param, method, path),
					Source: source,
				})
			}
		}
	}

	return result
}

func getPathParamsFromURL(path string) utils.StringSet {
	_, _, pathParams := utils.NormalizeTemplatedPath(path)
	return utils.StringList(pathParams).ToStringSet()
}
