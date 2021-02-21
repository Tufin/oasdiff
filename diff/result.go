package diff

import "github.com/getkin/kin-openapi/openapi3"

type Result struct {
	Diff    *Diff    `json:"diff,omitempty"`
	Summary *Summary `json:"summary,omitempty"`
}

/*
Run calculates the diff between two OAS specs including a summary.
Prefix is an optional path prefix that exists in s1 endpoints but not in s2.
If filter isn't empty, the diff will only include endpoints that match this regex.
*/
func Run(s1, s2 *openapi3.Swagger, prefix string, filter string) Result {
	diff := getDiff(s1, s2, prefix)
	diff.filterByRegex(filter)

	return Result{
		Diff:    diff,
		Summary: diff.getSummary(),
	}
}
