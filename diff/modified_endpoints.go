package diff

import (
	"github.com/getkin/kin-openapi/openapi3"
)

// ModifiedEndpoints maps endpoints (paths) to thier diff
type ModifiedEndpoints map[string]*EndpointDiff

func (modifiedEndpoints ModifiedEndpoints) addEndpointDiff(entrypoint1 string, pathItem1, pathItem2 *openapi3.PathItem) {

	if diff := diffEndpoint(pathItem1, pathItem2); !diff.empty() {
		modifiedEndpoints[entrypoint1] = diff
	}
}
