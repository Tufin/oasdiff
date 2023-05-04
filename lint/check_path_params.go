package lint

import (
	"fmt"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/tufin/oasdiff/load"
	"github.com/tufin/oasdiff/utils"
)

type pathParamsCtx struct{}

func PathParamsCheck(source string, s *load.OpenAPISpecInfo) []*Error {
	result := make([]*Error, 0)

	if s == nil || s.Spec == nil {
		return result
	}

	for path, pathItem := range s.Spec.Paths {
		pathParamsFromURL := getPathParamsFromURL(path)

		pathParams := utils.StringSet{}
		for _, parameter := range pathItem.Parameters {
			if parameter.Value.In == openapi3.ParameterInPath {
				pathParams.Add(parameter.Value.Name)
			}
		}

		context := pathParamsCtx{}
		for method, op := range pathItem.Operations() {
			result = append(result, context.checkOperation(pathParamsFromURL, pathParams, path, method, op, source)...)
		}
	}

	return result
}

func getPathParamsFromURL(path string) utils.StringSet {
	_, _, pathParams := utils.NormalizeTemplatedPath(path)
	return utils.StringList(pathParams).ToStringSet()
}

func (context *pathParamsCtx) checkOperation(pathParamsFromURL, pathParams utils.StringSet, path string, method string, op *openapi3.Operation, source string) []*Error {
	result := make([]*Error, 0)

	opParams := pathParams.Copy()
	for _, parameter := range op.Parameters {
		if parameter.Value.In == openapi3.ParameterInPath {
			opParams.Add(parameter.Value.Name)
		}
	}

	for param := range opParams.Minus(pathParamsFromURL) {
		result = append(result, &Error{
			Id:     "path-param-extra",
			Level:  LEVEL_ERROR,
			Text:   fmt.Sprintf("path parameter %q appears in the parameters section of the operation or path but is missing in the URL: %s %s", param, method, path),
			Source: source,
		})
	}

	for param := range pathParamsFromURL.Minus(opParams) {
		result = append(result, &Error{
			Id:     "path-param-missing",
			Level:  LEVEL_ERROR,
			Text:   fmt.Sprintf("path parameter %q appears in the URL path but is missing from the parameters section of the operation and path: %s %s", param, method, path),
			Source: source,
		})
	}

	return result
}
