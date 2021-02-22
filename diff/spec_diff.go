package diff

import "github.com/getkin/kin-openapi/openapi3"

// SpecDiff describes the changes between two OAS specs
type SpecDiff struct {
	PathsDiff      *PathsDiff      `json:"paths,omitempty"`      // deep diff of paths including their schemas, parameters, responses etc.
	SchemasDiff    *SchemasDiff    `json:"schemas,omitempty"`    // diff of top-level schemas (under components)
	ParametersDiff *ParametersDiff `json:"parameters,omitempty"` // diff of top-level parameters (under components)
	ResponsesDiff  *ResponsesDiff  `json:"responses,omitempty"`  // diff of top-level responses (under components)
}

func newSpecDiff() *SpecDiff {
	return &SpecDiff{}
}

func (specDiff SpecDiff) empty() bool {
	return specDiff == SpecDiff{}
}

func getDiff(s1, s2 *openapi3.Swagger, prefix string) *SpecDiff {

	diff := newSpecDiff()

	diff.setPathsDiff(getPathsDiff(s1.Paths, s2.Paths, prefix))
	diff.setSchemasDiff(getSchemasDiff(s1.Components.Schemas, s2.Components.Schemas))
	diff.setParametersDiff(getParametersDiff(toParameters(s1.Components.Parameters), toParameters(s2.Components.Parameters)))
	diff.setResponsesDiff(getResponsesDiff(s1.Components.Responses, s2.Components.Responses))

	return diff
}

func (specDiff *SpecDiff) setPathsDiff(pathsDiff *PathsDiff) {
	specDiff.PathsDiff = nil

	if !pathsDiff.empty() {
		specDiff.PathsDiff = pathsDiff
	}
}

func (specDiff *SpecDiff) setSchemasDiff(schemasDiff *SchemasDiff) {
	specDiff.SchemasDiff = nil

	if !schemasDiff.empty() {
		specDiff.SchemasDiff = schemasDiff
	}
}

func (specDiff *SpecDiff) setParametersDiff(parametersDiff *ParametersDiff) {
	specDiff.ParametersDiff = nil

	if !parametersDiff.empty() {
		specDiff.ParametersDiff = parametersDiff
	}
}

func (specDiff *SpecDiff) setResponsesDiff(responsesDiff *ResponsesDiff) {
	specDiff.ResponsesDiff = nil

	if !responsesDiff.empty() {
		specDiff.ResponsesDiff = responsesDiff
	}
}

func (specDiff *SpecDiff) getSummary() *Summary {

	result := Summary{
		Diff: !specDiff.empty(),
	}

	if specDiff.PathsDiff != nil {
		result.PathSummary = specDiff.PathsDiff.getSummary()
	}

	if specDiff.SchemasDiff != nil {
		result.SchemaSummary = specDiff.SchemasDiff.getSummary()
	}

	if specDiff.ParametersDiff != nil {
		result.ParameterSummary = specDiff.ParametersDiff.getSummary()
	}

	if specDiff.ResponsesDiff != nil {
		result.ResponsesSummary = specDiff.ResponsesDiff.getSummary()
	}

	return &result
}

func (specDiff *SpecDiff) filterByRegex(filter string) {
	if specDiff.PathsDiff != nil {

		specDiff.setPathsDiff(specDiff.PathsDiff.filterByRegex(filter))
	}
}
