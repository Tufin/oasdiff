package diff

import (
	"fmt"
	"reflect"

	"github.com/getkin/kin-openapi/openapi3"
)

/*
SchemaListDiff describes the changes between a pair of lists of schema objects: https://swagger.io/specification/#schema-object
The result is a combination of two diffs:
1. Diff of schemas with a $ref: number of added/deleted schemas; modified=diff of schemas with the same $ref
2. Diff of schemas without a $ref (inline schemas): number of added/deleted schemas; modified=only if exactly one schema was added and one deleted, the Modified field will show a diff between them
*/
type SchemaListDiff struct {
	Added    int             `json:"added,omitempty" yaml:"added,omitempty"`
	Deleted  int             `json:"deleted,omitempty" yaml:"deleted,omitempty"`
	Modified ModifiedSchemas `json:"modified,omitempty" yaml:"modified,omitempty"`
}

// Empty indicates whether a change was found in this element
func (diff *SchemaListDiff) Empty() bool {
	if diff == nil {
		return true
	}

	return diff.Added == 0 &&
		diff.Deleted == 0 &&
		len(diff.Modified) == 0
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

type SchemaRefMap map[string]*openapi3.SchemaRef

func toSchemaRefsMap(schemaRefs openapi3.SchemaRefs) SchemaRefMap {
	result := SchemaRefMap{}
	for _, schemaRef := range schemaRefs {
		if !isSchemaInline(schemaRef) {
			result[schemaRef.Ref] = schemaRef
		}
	}
	return result
}

func (diff SchemaListDiff) combine(other SchemaListDiff) (*SchemaListDiff, error) {

	return &SchemaListDiff{
		Added:    diff.Added + other.Added,
		Deleted:  diff.Deleted + other.Deleted,
		Modified: diff.Modified.combine(other.Modified),
	}, nil
}

func getSchemaListsDiffInternal(config *Config, state *state, schemaRefs1, schemaRefs2 openapi3.SchemaRefs) (*SchemaListDiff, error) {

	diffRefs, err := getSchemaListsRefsDiff(config, state, toSchemaRefsMap(schemaRefs1), toSchemaRefsMap(schemaRefs2))
	if err != nil {
		return nil, err
	}

	diffInline, err := getSchemaListsInlineDiff(config, state, schemaRefs1, schemaRefs2)
	if err != nil {
		return nil, err
	}

	return diffRefs.combine(diffInline)
}

// getSchemaListsRefsDiff compares schemas that have a reference
func getSchemaListsRefsDiff(config *Config, state *state, schemaMap1, schemaMap2 SchemaRefMap) (SchemaListDiff, error) {

	deleted := 0
	modified := ModifiedSchemas{}
	for ref, schema1 := range schemaMap1 {
		if schema2, found := schemaMap2[ref]; found {
			if err := modified.addSchemaDiff(config, state, ref, schema1, schema2); err != nil {
				return SchemaListDiff{}, err
			}
		} else {
			deleted++
		}
	}

	added := 0
	for ref := range schemaMap2 {
		if _, found := schemaMap1[ref]; !found {
			added++
		}
	}
	return SchemaListDiff{
		Added:    added,
		Deleted:  deleted,
		Modified: modified,
	}, nil
}

// getSchemaListsRefsDiff compares schemas that don't have a reference (inline schemas)
func getSchemaListsInlineDiff(config *Config, state *state, schemaRefs1, schemaRefs2 openapi3.SchemaRefs) (SchemaListDiff, error) {

	added, err := getGroupDifference(schemaRefs2, schemaRefs1)
	if err != nil {
		return SchemaListDiff{}, err
	}

	deleted, err := getGroupDifference(schemaRefs1, schemaRefs2)
	if err != nil {
		return SchemaListDiff{}, err
	}

	if len(added) == 1 && len(deleted) == 1 {
		d, err := getSchemaDiff(config, state, schemaRefs1[added[0]], schemaRefs2[deleted[0]])
		if err != nil {
			return SchemaListDiff{}, err
		}

		if d.Empty() {
			return SchemaListDiff{}, err
		}

		return SchemaListDiff{
			Modified: ModifiedSchemas{fmt.Sprintf("#%d", 1+added[0]): d},
		}, nil
	}

	return SchemaListDiff{
		Added:   len(added),
		Deleted: len(deleted),
	}, nil
}

func getGroupDifference(schemaRefs1, schemaRefs2 openapi3.SchemaRefs) ([]int, error) {

	notContained := []int{}

	for i, schemaRef1 := range schemaRefs1 {
		if isSchemaInline(schemaRef1) {

			found, err := findIndenticalSchema(schemaRef1, schemaRefs2)
			if err != nil {
				return nil, err
			}
			if !found {
				notContained = append(notContained, i)
			}
		}
	}
	return notContained, nil
}

func findIndenticalSchema(schemaRef1 *openapi3.SchemaRef, schemasRefs2 openapi3.SchemaRefs) (bool, error) {
	for _, schemaRef2 := range schemasRefs2 {
		if isSchemaInline(schemaRef2) {
			// compare with DeepEqual rather than SchemaDiff to ensure an exact syntactical match
			if reflect.DeepEqual(schemaRef1, schemaRef2) {
				return true, nil
			}
		}
	}

	return false, nil
}

func isSchemaInline(schemaRef *openapi3.SchemaRef) bool {
	if schemaRef == nil {
		return false
	}
	return schemaRef.Ref == ""
}
