package diff

import "github.com/getkin/kin-openapi/openapi3"

// CallbackDiff is a diff between two OAS callbacks
type CallbackDiff struct {
}

func (callbackDiff CallbackDiff) empty() bool {
	return callbackDiff == CallbackDiff{}
}

func diffCallbackValues(callback1, callback2 *openapi3.Callback) CallbackDiff {
	result := CallbackDiff{}

	pathItems1 := openapi3.Paths(*callback1)
	pathItems2 := openapi3.Paths(*callback2)

	getEndpointsDiff(pathItems1, pathItems2, "")

	return result
}
