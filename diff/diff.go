package diff

import "github.com/getkin/kin-openapi/openapi3"

// Result contains a diff and a summary
type Result struct {
	Diff    *Diff    `json:"diff,omitempty"`
	Summary *Summary `json:"summary,omitempty"`
}

// Run returns the diff between two OAS specs including a summary
func Run(s1 *openapi3.Swagger, s2 *openapi3.Swagger, prefix string, filter string) Result {
	diff := getDiff(s1, s2, prefix)
	diff.FilterByRegex(filter)

	return Result{
		Diff:    diff,
		Summary: diff.getSummary(),
	}
}

func getDiff(s1 *openapi3.Swagger, s2 *openapi3.Swagger, prefix string) *Diff {

	result := newDiff()

	if pathDiff := diffPaths(s1.Paths, s2.Paths, prefix); !pathDiff.empty() {
		result.PathDiff = pathDiff
	}

	if schemaDiff := diffSchemaCollection(s1.Components.Schemas, s2.Components.Schemas); !schemaDiff.empty() {
		result.SchemaDiff = schemaDiff
	}

	return result
}
