package diff

import "github.com/getkin/kin-openapi/openapi3"

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

// Breaking indicates whether this element includes a breaking change
func (diff *SchemaListDiff) Breaking() bool {
	return !diff.Empty()
}

func getSchemaListsDiff(config *Config, schemaRefs1, schemaRefs2 openapi3.SchemaRefs) (*SchemaListDiff, error) {
	diff, err := getSchemaListsDiffInternal(config, schemaRefs1, schemaRefs2)
	if err != nil {
		return nil, err
	}

	if diff.Empty() {
		return nil, nil
	}

	if config.BreakingOnly && !diff.Breaking() {
		return nil, nil
	}

	return diff, nil
}

func getSchemaListsDiffInternal(config *Config, schemaRefs1, schemaRefs2 openapi3.SchemaRefs) (*SchemaListDiff, error) {

	added, err := schemaRefsContained(config, schemaRefs1, schemaRefs2)
	if err != nil {
		return nil, err
	}

	deleted, err := schemaRefsContained(config, schemaRefs2, schemaRefs1)
	if err != nil {
		return nil, err
	}

	if len(added) == 1 && len(deleted) == 1 {
		d, err := getSchemaDiff(config, added[0], deleted[0])
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

func schemaRefsContained(config *Config, schemaRefs1, schemaRefs2 openapi3.SchemaRefs) ([]*openapi3.SchemaRef, error) {

	result := []*openapi3.SchemaRef{}

	for _, schemaRef1 := range schemaRefs1 {
		found, err := findSchema(config, schemaRef1, schemaRefs2)
		if err != nil {
			return nil, err
		}
		if !found {
			result = append(result, schemaRef1)
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
