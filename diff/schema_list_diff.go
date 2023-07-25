package diff

import (
	"fmt"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/tufin/oasdiff/utils"
)

/*
SchemaListDiff describes the changes between a pair of lists of schema objects: https://swagger.io/specification/#schema-object
The result is a combination of two diffs:
1. Diff of schemas with a $ref: number of added/deleted schemas; modified=diff of schemas with the same $ref
2. Diff of schemas without a $ref (inline schemas): number of added/deleted schemas; modified=only if exactly one schema was added and one deleted, the Modified field will show a diff between them
*/
type SchemaListDiff struct {
	Added    utils.StringList `json:"added,omitempty" yaml:"added,omitempty"`
	Deleted  utils.StringList `json:"deleted,omitempty" yaml:"deleted,omitempty"`
	Modified ModifiedSchemas  `json:"modified,omitempty" yaml:"modified,omitempty"`
}

// Empty indicates whether a change was found in this element
func (diff *SchemaListDiff) Empty() bool {
	if diff == nil {
		return true
	}

	return len(diff.Added) == 0 &&
		len(diff.Deleted) == 0 &&
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

func (diff SchemaListDiff) combine(other SchemaListDiff) (*SchemaListDiff, error) {

	return &SchemaListDiff{
		Added:    append(diff.Added, other.Added...),
		Deleted:  append(diff.Deleted, other.Deleted...),
		Modified: diff.Modified.combine(other.Modified),
	}, nil
}

func getSchemaListsDiffInternal(config *Config, state *state, schemaRefs1, schemaRefs2 openapi3.SchemaRefs) (*SchemaListDiff, error) {

	if len(schemaRefs1) == 0 && len(schemaRefs2) == 0 {
		return nil, nil
	}

	diffRefs, err := getSchemaListsRefsDiff(config, state, schemaRefs1, schemaRefs2, isSchemaRef)
	if err != nil {
		return nil, err
	}

	diffInline, err := getSchemaListsInlineDiff(config, state, schemaRefs1, schemaRefs2, isSchemaInline)
	if err != nil {
		return nil, err
	}

	return diffRefs.combine(diffInline)
}

type schemaRefsFilter func(schemaRef *openapi3.SchemaRef) bool

// getSchemaListsRefsDiff compares schemas by $ref name
func getSchemaListsRefsDiff(config *Config, state *state, schemaRefs1, schemaRefs2 openapi3.SchemaRefs, filter schemaRefsFilter) (SchemaListDiff, error) {
	added := utils.StringList{}
	deleted := utils.StringList{}
	modified := ModifiedSchemas{}

	schemaMap2 := toSchemaRefMap(schemaRefs2, filter)
	for _, schema1 := range schemaRefs1 {
		if !filter(schema1) {
			continue
		}
		ref := schema1.Ref
		if schema2, found := schemaMap2[ref]; found {
			schemaMap2.delete(ref)
			if err := modified.addSchemaDiff(config, state, ref, schema1, schema2.schemaRef); err != nil {
				return SchemaListDiff{}, err
			}
		} else {
			schemaName := ref[strings.LastIndex(ref, "/")+1:]
			deleted = append(deleted, schemaName)
		}
	}

	schemaMap1 := toSchemaRefMap(schemaRefs1, filter)
	for _, schema2 := range schemaRefs2 {
		if !filter(schema2) {
			continue
		}
		ref := schema2.Ref
		if _, found := schemaMap1[ref]; found {
			schemaMap1.delete(ref)
		} else {
			schemaName := ref[strings.LastIndex(ref, "/")+1:]
			added = append(added, schemaName)
		}
	}
	return SchemaListDiff{
		Added:    added,
		Deleted:  deleted,
		Modified: modified,
	}, nil
}

// getSchemaListsRefsDiff compares schemas by their syntax
func getSchemaListsInlineDiff(config *Config, state *state, schemaRefs1, schemaRefs2 openapi3.SchemaRefs, filter schemaRefsFilter) (SchemaListDiff, error) {

	addedIdx, addedSchemas, err := getGroupDifference(config, state, schemaRefs2, schemaRefs1, filter, "RevisionSchema")
	if err != nil {
		return SchemaListDiff{}, err
	}

	deletedIdx, deletedSchemas, err := getGroupDifference(config, state, schemaRefs1, schemaRefs2, filter, "BaseSchema")
	if err != nil {
		return SchemaListDiff{}, err
	}

	if len(addedIdx) == 1 && len(deletedIdx) == 1 {
		d, err := getSchemaDiff(config, state, schemaRefs1[deletedIdx[0]], schemaRefs2[addedIdx[0]])
		if err != nil {
			return SchemaListDiff{}, err
		}

		if d.Empty() {
			return SchemaListDiff{}, err
		}

		return SchemaListDiff{
			Modified: ModifiedSchemas{fmt.Sprintf("#%d", 1+deletedIdx[0]): d},
		}, nil
	}

	return SchemaListDiff{
		Added:   addedSchemas,
		Deleted: deletedSchemas,
	}, nil
}

func getGroupDifference(config *Config, state *state, schemaRefs1, schemaRefs2 openapi3.SchemaRefs, filter schemaRefsFilter, inlineSchemaPrefix string) ([]int, []string, error) {

	notContainedIdx := []int{}
	notContainedSchemas := []string{}
	matched := map[int]struct{}{}

	for index1, schemaRef1 := range schemaRefs1 {
		if !filter(schemaRef1) {
			continue
		}

		if found, index2, err := findIndenticalSchema(config, state, schemaRef1, schemaRefs2, matched, filter); err != nil {
			return nil, nil, err
		} else if !found {
			notContainedIdx = append(notContainedIdx, index1)
			schemaName := schemaRef1.Ref
			if schemaName == "" {
				schemaName = fmt.Sprintf("%s[%d]", inlineSchemaPrefix, index1)
			}
			notContainedSchemas = append(notContainedSchemas, schemaName)
		} else {
			matched[index2] = struct{}{}
		}
	}
	return notContainedIdx, notContainedSchemas, nil
}

func findIndenticalSchema(config *Config, state *state, schemaRef1 *openapi3.SchemaRef, schemasRefs2 openapi3.SchemaRefs, matched map[int]struct{}, filter schemaRefsFilter) (bool, int, error) {
	for index2, schemaRef2 := range schemasRefs2 {
		if !filter(schemaRef1) {
			continue
		}
		if alreadyMatched(index2, matched) {
			continue
		}

		if schemaDiff, err := getSchemaDiff(config, state, schemaRef1, schemaRef2); err != nil {
			return false, 0, err
		} else if schemaDiff.Empty() {
			return true, index2, nil
		}
	}

	return false, 0, nil
}

func alreadyMatched(index int, matched map[int]struct{}) bool {
	_, found := matched[index]
	return found
}

func isSchemaInline(schemaRef *openapi3.SchemaRef) bool {
	if schemaRef == nil {
		return false
	}
	return schemaRef.Ref == ""
}

func isSchemaRef(schemaRef *openapi3.SchemaRef) bool {
	if schemaRef == nil {
		return false
	}
	return schemaRef.Ref != ""
}
