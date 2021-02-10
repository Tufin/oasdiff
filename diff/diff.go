package diff

import (
	"strings"

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

type PathItemPair struct {
	PathItem1 *openapi3.PathItem
	PathItem2 *openapi3.PathItem
}

type PathItemPairs map[string]*PathItemPair

func diffEndpoints(paths1 openapi3.Paths, paths2 openapi3.Paths, prefix string) (openapi3.Paths, openapi3.Paths, PathItemPairs) {

	added := openapi3.Paths{}
	deleted := openapi3.Paths{}
	other := PathItemPairs{}

	for endpoint1, pathItem1 := range paths1 {
		if pathItem2, ok := findEndpoint(endpoint1, "", paths2); ok {
			other[endpoint1] = &PathItemPair{
				PathItem1: pathItem1,
				PathItem2: pathItem2,
			}
		} else {
			deleted[endpoint1] = pathItem1
		}
	}

	for endpoint2, pathItem2 := range paths2 {
		if _, ok := findEndpoint(endpoint2, "", paths1); !ok {
			added[endpoint2] = pathItem2
		}
	}

	return added, deleted, other
}

func findEndpoint(entrypoint string, prefix string, paths openapi3.Paths) (*openapi3.PathItem, bool) {
	noSlash, withSlash := combine(entrypoint, prefix)

	if pathItem, ok := paths[noSlash]; ok {
		return pathItem, true
	}

	if pathItem, ok := paths[withSlash]; ok {
		return pathItem, true
	}

	return nil, false
}

func combine(s string, prefix string) (string, string) {
	s = strings.TrimSuffix(s, "/")
	s = strings.TrimPrefix(s, prefix)

	return s, s + "/"
}

func equalParams(param1 *openapi3.Parameter, param2 *openapi3.Parameter) bool {
	return param1.Name == param2.Name && param1.In == param2.In
}
