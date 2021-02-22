package diff

import "github.com/getkin/kin-openapi/openapi3"

// Diff describes the changes between two OAS specs
type Diff struct {
	PathsDiff      *PathsDiff      `json:"paths,omitempty"`
	SchemasDiff    *SchemasDiff    `json:"schemas,omitempty"`
	ParametersDiff *ParametersDiff `json:"parameters,omitempty"`
}

func newDiff() *Diff {
	return &Diff{}
}

func (diff Diff) empty() bool {
	return diff == Diff{}
}

func getDiff(s1, s2 *openapi3.Swagger, prefix string) *Diff {

	diff := newDiff()

	diff.setPathsDiff(getPathsDiff(s1.Paths, s2.Paths, prefix))
	diff.setSchemasDiff(getSchemasDiff(s1.Components.Schemas, s2.Components.Schemas))
	diff.setParametersDiff(getParametersDiff(toParameters(s1.Components.Parameters), toParameters(s2.Components.Parameters)))

	return diff
}

func (diff *Diff) setPathsDiff(pathsDiff *PathsDiff) {
	diff.PathsDiff = nil

	if !pathsDiff.empty() {
		diff.PathsDiff = pathsDiff
	}
}

func (diff *Diff) setSchemasDiff(schemasDiff *SchemasDiff) {
	diff.SchemasDiff = nil

	if !schemasDiff.empty() {
		diff.SchemasDiff = schemasDiff
	}
}

func (diff *Diff) setParametersDiff(parametersDiff *ParametersDiff) {
	diff.ParametersDiff = nil

	if !parametersDiff.empty() {
		diff.ParametersDiff = parametersDiff
	}
}

func (diff *Diff) getSummary() *Summary {

	result := Summary{
		Diff: !diff.empty(),
	}

	if diff.PathsDiff != nil {
		result.PathSummary = diff.PathsDiff.getSummary()
	}

	if diff.SchemasDiff != nil {
		result.SchemaSummary = diff.SchemasDiff.getSummary()
	}

	if diff.ParametersDiff != nil {
		result.ParameterSummary = diff.ParametersDiff.getSummary()
	}

	return &result
}

func (diff *Diff) filterByRegex(filter string) {
	if diff.PathsDiff != nil {

		diff.setPathsDiff(diff.PathsDiff.filterByRegex(filter))
	}
}
