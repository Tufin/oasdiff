package diff

import (
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

// Diff finds s1-s2 endpoints (appear in s1 but not in s2)
func Diff(s1 *openapi3.Swagger, s2 *openapi3.Swagger, prefix string) *DiffResult {

	result := newDiffResult()

	for entrypoint1, pathItem1 := range s1.Paths {

		pathItem2, ok := findEndpoint(entrypoint1, prefix, s2.Paths)

		if !ok {
			result.addDeletedEndpoint(entrypoint1)
			continue
		}

		result.addModifiedEndpoint(entrypoint1, pathItem1, pathItem2)
	}
	return result
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
