package headers

import (
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

// Lowercase replaces header names to lowercase
func Lowercase(spec *openapi3.T) {
	lowerHeaderNames(spec)
}

func lowerHeaderNames(spec *openapi3.T) {
	for _, path := range spec.Paths.Map() {

		for _, paramRef := range path.Parameters {
			lowerHeaderName(paramRef.Value)
		}

		for _, op := range path.Operations() {
			for _, paramRef := range op.Parameters {
				lowerHeaderName(paramRef.Value)
			}

			for _, responseRef := range op.Responses.Map() {
				for _, headerRef := range responseRef.Value.Headers {
					lowerHeaderName(&headerRef.Value.Parameter)
				}
			}
		}
	}
}

func lowerHeaderName(param *openapi3.Parameter) {
	if param.In != openapi3.ParameterInHeader {
		return
	}

	param.Name = strings.ToLower(param.Name)
}
