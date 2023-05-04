package diff

import (
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/tufin/oasdiff/utils"
)

type pathItemPair struct {
	PathItem1     *openapi3.PathItem
	PathItem2     *openapi3.PathItem
	PathParamsMap PathParamsMap
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
		if pathItem2, pathParamsMap, ok := findEndpoint(endpoint1, paths2); ok {
			other[endpoint1] = &pathItemPair{
				PathItem1:     pathItem1,
				PathItem2:     pathItem2,
				PathParamsMap: pathParamsMap,
			}
		} else {
			deleted[endpoint1] = pathItem1
		}
	}

	for endpoint2, pathItem2 := range paths2 {
		if _, _, ok := findEndpoint(endpoint2, paths1); !ok {
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

func findEndpoint(endpoint string, paths openapi3.Paths) (*openapi3.PathItem, PathParamsMap, bool) {
	if pathItem, ok := paths[endpoint]; ok {
		return pathItem, PathParamsMap{}, true
	}

	return findNormalizedEndpoint(endpoint, paths)
}

/*
findNormalizedEndpoint finds a corresponding path ignoring differences in template variable names

This implementation is based on Paths.Find in openapi3
*/
func findNormalizedEndpoint(key string, paths openapi3.Paths) (*openapi3.PathItem, PathParamsMap, bool) {
	normalizedPath, expected, pathParams1 := utils.NormalizeTemplatedPath(key)
	for path, pathItem := range paths {
		pathNormalized, got, pathParams2 := utils.NormalizeTemplatedPath(path)
		if got == expected && pathNormalized == normalizedPath {
			if pathParamsMap, ok := NewPathParamsMap(pathParams1, pathParams2); ok {
				return pathItem, pathParamsMap, true
			}
		}
	}
	return nil, nil, false
}
