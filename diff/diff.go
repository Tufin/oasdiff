package diff

import "github.com/getkin/kin-openapi/openapi3"

type Diff struct {
	PathDiff   *PathDiff             `json:"endpoints,omitempty"`
	SchemaDiff *SchemaCollectionDiff `json:"schemas,omitempty"`
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

	diff.setPathDiff(diffPaths(s1.Paths, s2.Paths, prefix))
	diff.setSchemaDiff(diffSchemaCollection(s1.Components.Schemas, s2.Components.Schemas))

	return diff
}

func (diff *Diff) setPathDiff(pathDiff *PathDiff) {
	diff.PathDiff = nil

	if !pathDiff.empty() {
		diff.PathDiff = pathDiff
	}
}

func (diff *Diff) setSchemaDiff(schemaCollectionDiff *SchemaCollectionDiff) {
	diff.SchemaDiff = nil

	if !schemaCollectionDiff.empty() {
		diff.SchemaDiff = schemaCollectionDiff
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

		diff.setPathDiff(diff.PathDiff.filterByRegex(filter))
	}
}
