package diff

import "github.com/getkin/kin-openapi/openapi3"

// Diff describes the changes between two OAS specs
type Diff struct {
	PathDiff   *PathsDiff   `json:"paths,omitempty"`
	SchemaDiff *SchemasDiff `json:"schemas,omitempty"`
}

func newDiff() *Diff {
	return &Diff{}
}

func (diff *Diff) empty() bool {
	return diff.PathDiff == nil &&
		diff.SchemaDiff == nil
}

func getDiff(s1, s2 *openapi3.Swagger, prefix string) *Diff {

	diff := newDiff()

	diff.setPath(getPathsDiff(s1.Paths, s2.Paths, prefix))
	diff.setSchemas(getSchemasDiff(s1.Components.Schemas, s2.Components.Schemas))

	return diff
}

func (diff *Diff) setPath(paths *PathsDiff) {
	diff.PathDiff = nil

	if !paths.empty() {
		diff.PathDiff = paths
	}
}

func (diff *Diff) setSchemas(schemas *SchemasDiff) {
	diff.SchemaDiff = nil

	if !schemas.empty() {
		diff.SchemaDiff = schemas
	}
}

func (diff *Diff) getSummary() *Summary {

	result := Summary{
		Diff: !diff.empty(),
	}

	if diff.PathDiff != nil {
		result.PathSummary = diff.PathDiff.getSummary()
	}

	if diff.SchemaDiff != nil {
		result.SchemaSummary = diff.SchemaDiff.getSummary()
	}

	return &result
}

func (diff *Diff) filterByRegex(filter string) {
	if diff.PathDiff != nil {

		diff.setPath(diff.PathDiff.filterByRegex(filter))
	}
}
