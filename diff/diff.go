package diff

import (
	"github.com/getkin/kin-openapi/openapi3"
)

// Diff finds changes between two OAS specs
func Diff(s1 *openapi3.Swagger, s2 *openapi3.Swagger, prefix string) *DiffResult {

	result := newDiffResult()

	addedEndpoints, deletedEndpoints, otherEndpoints := diffEndpoints(s1.Paths, s2.Paths, prefix)

	for endpoint := range addedEndpoints {
		result.addAddedEndpoint(endpoint)
	}

	for endpoint := range deletedEndpoints {
		result.addDeletedEndpoint(endpoint)
	}

	for endpoint, pathItemPair := range otherEndpoints {
		result.addModifiedEndpoint(endpoint, pathItemPair.PathItem1, pathItemPair.PathItem2)
	}

	return result
}
