package diff

import (
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

type pathItemPair struct {
	PathItem1 *openapi3.PathItem
	PathItem2 *openapi3.PathItem
}

type pathItemPairs map[string]*pathItemPair

func getPathItemsDiff(paths1, paths2 openapi3.Paths, config *Config) (openapi3.Paths, openapi3.Paths, pathItemPairs) {
	return getPathItemsDiffInternal(addPrefix(paths1, config.PathPrefixBase), addPrefix(paths2, config.PathPrefixRevision))
}

func getPathItemsDiffInternal(paths1, paths2 openapi3.Paths) (openapi3.Paths, openapi3.Paths, pathItemPairs) {

	added := openapi3.Paths{}
	deleted := openapi3.Paths{}
	other := pathItemPairs{}

	for endpoint1, pathItem1 := range paths1 {
		if pathItem2, ok := findEndpoint(endpoint1, paths2); ok {
			other[endpoint1] = &pathItemPair{
				PathItem1: pathItem1,
				PathItem2: pathItem2,
			}
		} else {
			deleted[endpoint1] = pathItem1
		}
	}

	for endpoint2, pathItem2 := range paths2 {
		if _, ok := findEndpoint(endpoint2, paths1); !ok {
			added[endpoint2] = pathItem2
		}
	}

	return added, deleted, other
}

func addPrefix(paths openapi3.Paths, prefix string) openapi3.Paths {
	result := make(openapi3.Paths, len(paths))
	for path, pathItem := range paths {
		result[prefix+path] = pathItem
	}
	return result
}

func findEndpoint(entrypoint string, paths openapi3.Paths) (*openapi3.PathItem, bool) {
	noSlash, withSlash := combine(entrypoint)

	if pathItem, ok := paths[noSlash]; ok {
		return pathItem, true
	}

	if pathItem, ok := paths[withSlash]; ok {
		return pathItem, true
	}

	return nil, false
}

func combine(s string) (string, string) {
	s = strings.TrimSuffix(s, "/")

	return s, s + "/"
}
