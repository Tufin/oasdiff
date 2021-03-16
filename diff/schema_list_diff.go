package diff

import "github.com/getkin/kin-openapi/openapi3"

/*
SchemaListDiff is a diff between two lists of schema objects: https://swagger.io/specification/#schema-object
Unlike other diff structs, this one doesn't indicate the exact change but only tells us how many schema objects where added and/or deleted
*/
type SchemaListDiff struct {
	Added   int `json:"added,omitempty" yaml:"added,omitempty"`
	Deleted int `json:"deleted,omitempty" yaml:"deleted,omitempty"`
}

func newSchemaListDiff() *SchemaListDiff {
	return &SchemaListDiff{}
}

// Empty indicates whether a change was found in this element
func (diff *SchemaListDiff) Empty() bool {
	return diff == nil || *diff == SchemaListDiff{}
}

func getSchemaListsDiff(config *Config, schemaRefs1, schemaRefs2 openapi3.SchemaRefs) (*SchemaListDiff, error) {
	diff, err := getSchemaListsDiffInternal(config, schemaRefs1, schemaRefs2)
	if err != nil {
		return nil, err
	}
	if diff.Empty() {
		return nil, nil
	}
	return diff, nil
}

func getSchemaListsDiffInternal(config *Config, schemaRefs1, schemaRefs2 openapi3.SchemaRefs) (*SchemaListDiff, error) {

	result := newSchemaListDiff()
	var err error

	result.Added, err = schemaRefsContained(config, schemaRefs1, schemaRefs2)
	if err != nil {
		return nil, err
	}

	result.Deleted, err = schemaRefsContained(config, schemaRefs2, schemaRefs1)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func schemaRefsContained(config *Config, schemaRefs1, schemaRefs2 openapi3.SchemaRefs) (int, error) {

	result := 0

	for _, schemaRef1 := range schemaRefs1 {
		found, err := findSchema(config, schemaRef1, schemaRefs2)
		if err != nil {
			return 0, err
		}
		if !found {
			result++
		}
	}
	return result, nil
}

func findSchema(config *Config, schemaRef1 *openapi3.SchemaRef, schemaRefs2 openapi3.SchemaRefs) (bool, error) {
	// TODO: optimize with a map
	for _, schemaRef2 := range schemaRefs2 {
		diff, err := getSchemaDiff(config, schemaRef1, schemaRef2)
		if err != nil {
			return false, err
		}
		if diff.Empty() {
			return true, nil
		}
	}

	return false, nil
}
