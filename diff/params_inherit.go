package diff

import "github.com/getkin/kin-openapi/openapi3"

func paramsInherit(pathParams, opParams openapi3.Parameters) openapi3.Parameters {
	for _, pathParam := range pathParams {
		opParams = paramInherit(opParams, pathParam.Value)
	}
	return opParams
}

func paramInherit(opParams openapi3.Parameters, pathParam *openapi3.Parameter) openapi3.Parameters {
	if opParams.GetByInAndName(pathParam.In, pathParam.Name) == nil {
		opParams = append(opParams, &openapi3.ParameterRef{
			Value: pathParam,
		})
	}
	return opParams
}
