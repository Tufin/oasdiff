package diff

import "github.com/getkin/kin-openapi/openapi3"

// Diff describes the changes between a pair of OpenAPI specifications including a summary of the changes
type Diff struct {
	SpecDiff *SpecDiff `json:"spec,omitempty" yaml:"spec,omitempty"`
	Summary  *Summary  `json:"summary,omitempty" yaml:"summary,omitempty"`
}

/*
Get calculates the diff between a pair of OpenAPI specifications.

Note that Get expects OpenAPI References (https://swagger.io/docs/specification/using-ref/) to be resolved.  
References are normally resolved automatically when you load the spec.  
In other cases you can resolve refs using https://pkg.go.dev/github.com/getkin/kin-openapi/openapi3#SwaggerLoader.ResolveRefsIn.
*/
func Get(config *Config, s1, s2 *openapi3.Swagger) (Diff, error) {
	diff, err := getDiff(config, s1, s2)
	if err != nil {
		return Diff{}, err
	}

	return Diff{
		SpecDiff: diff,
		Summary:  diff.getSummary(),
	}, nil
}
