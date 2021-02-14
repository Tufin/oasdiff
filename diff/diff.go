package diff

import (
	"github.com/getkin/kin-openapi/openapi3"
)

// Diff returns the diff between two OAS (swagger) specs
func Diff(s1 *openapi3.Swagger, s2 *openapi3.Swagger, prefix string) *DiffResult {

	result := newDiffResult()

	if pathsDiff := diffPaths(s1.Paths, s2.Paths, prefix); !pathsDiff.empty() {
		result.PathsDiff = pathsDiff
	}

	return result
}
