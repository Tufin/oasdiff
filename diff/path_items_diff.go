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
	return getPathItemsDiffInternal(
		rewritePrefix(paths1, config.PathStripPrefixBase, config.PathPrefixBase),
		rewritePrefix(paths2, config.PathStripPrefixRevision, config.PathPrefixRevision))
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

func rewritePrefix(paths openapi3.Paths, strip, prepend string) openapi3.Paths {
	result := make(openapi3.Paths, len(paths))
	for path, pathItem := range paths {
		result[prepend+strings.TrimPrefix(path, strip)] = pathItem
	}
	return result
}

func findEndpoint(entrypoint string, paths openapi3.Paths) (*openapi3.PathItem, bool) {
	pathItem, ok := paths[entrypoint]
	return pathItem, ok
}
