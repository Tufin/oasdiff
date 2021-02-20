package diff

import (
	"github.com/getkin/kin-openapi/openapi3"
)

type ModifiedEndpoints map[string]*EndpointDiff // key is endpoint (path)

func (modifiedEndpoints ModifiedEndpoints) addEndpointDiff(entrypoint1 string, pathItem1 *openapi3.PathItem, pathItem2 *openapi3.PathItem) {

	if diff := diffEndpoint(pathItem1, pathItem2); !diff.empty() {
		modifiedEndpoints[entrypoint1] = diff
	}
}
