package diff

import (
	"github.com/getkin/kin-openapi/openapi3"
)

type ModifiedEndpoints map[string]*EndpointDiff // key is endpoint (path)

func (endpointDiffResult ModifiedEndpoints) addEndpointDiff(entrypoint1 string, pathItem1 *openapi3.PathItem, pathItem2 *openapi3.PathItem) {

	diff := diffEndpoint(pathItem1, pathItem2)
	if !diff.empty() {
		endpointDiffResult[entrypoint1] = diff
	}
}
