package diff

import "github.com/getkin/kin-openapi/openapi3"

func toParameters(parametersMap openapi3.ParametersMap) openapi3.Parameters {

	result := make(openapi3.Parameters, len(parametersMap))

	i := 0
	for _, v := range parametersMap {
		result[i] = v
		i++
	}

	return result
}
