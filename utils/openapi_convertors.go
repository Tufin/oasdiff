package utils

import "github.com/getkin/kin-openapi/openapi3"

func ResponseBodiesToResponses(responseBodies openapi3.ResponseBodies) *openapi3.Responses {
	result := openapi3.NewResponsesWithCapacity(len(responseBodies))
	for k, v := range responseBodies {
		result.Set(k, v)
	}
	return result
}
