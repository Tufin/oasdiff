package diff

import (
	"reflect"

	"github.com/getkin/kin-openapi/openapi3"
)

/*
SchemaListDiff describes the changes between a pair of lists of schema objects: https://swagger.io/specification/#schema-object
Unlike other diff structs, this one doesn't indicate the exact change but only tells us how many schema objects where added and/or deleted.
As a special case, when exactly one schema was added and one deleted, the Modified field will show a diff between them.
*/
type SchemaListDiff struct {
	Added    int         `json:"added,omitempty" yaml:"added,omitempty"`
	Deleted  int         `json:"deleted,omitempty" yaml:"deleted,omitempty"`
	Modified *SchemaDiff `json:"modified,omitempty" yaml:"modified,omitempty"`
}

// Empty indicates whether a change was found in this element
func (diff *SchemaListDiff) Empty() bool {
	return diff == nil || *diff == SchemaListDiff{}
}

func getSchemaListsDiff(config *Config, state *state, schemaRefs1, schemaRefs2 openapi3.SchemaRefs) (*SchemaListDiff, error) {
	diff, err := getSchemaListsDiffInternal(config, state, schemaRefs1, schemaRefs2)
	if err != nil {
		return nil, err
	}

	if diff.Empty() {
		return nil, nil
	}

	return diff, nil
}

func getSchemaListsDiffInternal(config *Config, state *state, schemaRefs1, schemaRefs2 openapi3.SchemaRefs) (*SchemaListDiff, error) {

	added, err := getGroupDifference(schemaRefs1, schemaRefs2)
	if err != nil {
		return nil, err
	}

	deleted, err := getGroupDifference(schemaRefs2, schemaRefs1)
	if err != nil {
		return nil, err
	}

	if len(added) == 1 && len(deleted) == 1 {
		d, err := getSchemaDiff(config, state, schemaRefs1[added[0]], schemaRefs2[deleted[0]])
		if err != nil {
			return nil, err
		}
		return &SchemaListDiff{
			Modified: d,
		}, nil
	}

	return &SchemaListDiff{
		Added:   len(added),
		Deleted: len(deleted),
	}, nil
}

func getGroupDifference(schemaRefs1, schemaRefs2 openapi3.SchemaRefs) ([]int, error) {

	notContained := []int{}

	// TODO: optimize with a map
	for i, schemaRef1 := range schemaRefs1 {
		found, err := findSchema(schemaRef1, schemaRefs2)
		if err != nil {
			return nil, err
		}
		if !found {
			notContained = append(notContained, i)
		}
	}
	return notContained, nil
}

func findSchema(schema1 *openapi3.SchemaRef, schemas2 openapi3.SchemaRefs) (bool, error) {
	for _, schema2 := range schemas2 {
		// compare with DeepEqual rather than SchemaDiff to ensure an exact syntactical match
		if reflect.DeepEqual(schema1, schema2) {
			return true, nil
		}
	}

	return false, nil
}
