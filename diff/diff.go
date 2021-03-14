package diff

import "github.com/getkin/kin-openapi/openapi3"

// Diff describes changes between two OpenAPI specifications including a summary of the changes
type Diff struct {
	SpecDiff *SpecDiff `json:"spec,omitempty" yaml:"spec,omitempty"`
	Summary  *Summary  `json:"summary,omitempty" yaml:"summary,omitempty"`
}

/*
Get calculates the diff between two OpenAPI specifications.
References should be resolved before calling this function.
This is normally done automatically by loaders.
See https://swagger.io/docs/specification/using-ref/ and https://pkg.go.dev/github.com/getkin/kin-openapi/openapi3#SwaggerLoader.ResolveRefsIn.
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
