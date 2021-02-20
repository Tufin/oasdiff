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
	diff.filterByRegex(filter)

	return Result{
		Diff:    diff,
		Summary: diff.getSummary(),
	}
}
