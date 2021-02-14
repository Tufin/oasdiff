package diff

import (
	"github.com/getkin/kin-openapi/openapi3"
)

// Diff returns the diff between two OAS (swagger) specs
func Diff(s1 *openapi3.Swagger, s2 *openapi3.Swagger, prefix string) *DiffResult {

	result := newDiffResult()

	if pathDiff := diffPaths(s1.Paths, s2.Paths, prefix); !pathDiff.empty() {
		result.PathDiff = pathDiff
	}

	if schemaDiff := diffSchemaCollection(s1.Components.Schemas, s2.Components.Schemas); !schemaDiff.empty() {
		result.SchemaDiff = schemaDiff
	}

	return result
}
