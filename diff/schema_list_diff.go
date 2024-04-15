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
1. Diff of schemas with a $ref: added/deleted schema names; modified=diff of schemas with the same $ref
2. Diff of schemas without a $ref (inline schemas): added/deleted schemas (base/revision + index in the list of schemas); modified=diff of schemas with the same title or a single modified inline schema with no title
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

func getSchemaListsInlineDiff(config *Config, state *state, schemaRefs1, schemaRefs2 openapi3.SchemaRefs, filter schemaRefsFilter) (SchemaListDiff, error) {

	// find schemas in revision that have no matching schema in the base
	addedIdx, err := getNonContainedInlineSchemas(config, state, schemaRefs2, schemaRefs1, filter, "RevisionSchema")
	if err != nil {
		return SchemaListDiff{}, err
	}

	// find schemas in base that have no matching schema in the revision
	deletedIdx, err := getNonContainedInlineSchemas(config, state, schemaRefs1, schemaRefs2, filter, "BaseSchema")
	if err != nil {
		return SchemaListDiff{}, err
	}

	// match schemas by title
	addedIdx, deletedIdx, modifiedSchemas, err := compareByTitle(config, state, addedIdx, deletedIdx, schemaRefs1, schemaRefs2)
	if err != nil {
		return SchemaListDiff{}, err
	}

	// special case: single modified schema with no title
	if isSingleModifiedCase(schemaRefs1, schemaRefs2, addedIdx, deletedIdx) {
		d, err := getSchemaDiff(config, state, schemaRefs1[deletedIdx[0]], schemaRefs2[addedIdx[0]])
		if err != nil {
			return SchemaListDiff{}, err
		}

		modifiedSchemas[fmt.Sprintf("#%d", 1+deletedIdx[0])] = d
		addedIdx = []int{}
		deletedIdx = []int{}
	}

	return SchemaListDiff{
		Added:    getSchemaNames("RevisionSchema", addedIdx, schemaRefs2),
		Deleted:  getSchemaNames("BaseSchema", deletedIdx, schemaRefs1),
		Modified: modifiedSchemas,
	}, nil
}

func isSingleModifiedCase(schemaRefs1, schemaRefs2 openapi3.SchemaRefs, addedIdx, deletedIdx []int) bool {
	return len(addedIdx) == 1 &&
		len(deletedIdx) == 1 &&
		schemaRefs1[deletedIdx[0]].Value.Title == "" &&
		schemaRefs2[addedIdx[0]].Value.Title == ""
}

func compareByTitle(config *Config, state *state, addedIdx, deletedIdx []int, schemaRefs1, schemaRefs2 openapi3.SchemaRefs) ([]int, []int, ModifiedSchemas, error) {

	addedMatched, deletedMatched := matchByTitle(config, state, addedIdx, deletedIdx, schemaRefs1, schemaRefs2)

	modifiedSchemas := ModifiedSchemas{}
	for _, addedId := range addedIdx {
		deletedId, found := addedMatched[addedId]
		if !found {
			continue
		}

		d, err := getSchemaDiff(config, state, schemaRefs1[deletedId], schemaRefs2[addedId])
		if err != nil {
			return nil, nil, nil, err
		}
		modifiedSchemas[schemaRefs1[deletedId].Value.Title] = d
	}

	return deleteMatched(addedIdx, addedMatched), deleteMatched(deletedIdx, deletedMatched), modifiedSchemas, nil
}

func deleteMatched(idx []int, addedMatched map[int]int) []int {
	addedIdxRemaining := []int{}
	for _, addedId := range idx {
		if _, found := addedMatched[addedId]; !found {
			addedIdxRemaining = append(addedIdxRemaining, addedId)
		}
	}
	return addedIdxRemaining
}

func matchByTitle(config *Config, state *state, addedIdx, deletedIdx []int, schemaRefs1, schemaRefs2 openapi3.SchemaRefs) (map[int]int, map[int]int) {

	addedMatched := map[int]int{}
	deletedMatched := map[int]int{}

	matchedTitles := utils.StringSet{}
	matchedTitles.Add("") // empty title is not allowed

	for _, addedId := range addedIdx {
		title := schemaRefs2[addedId].Value.Title
		if matchedTitles.Contains(title) {
			// title already matched, skip
			continue
		}
		for _, deletedId := range deletedIdx {
			if title == schemaRefs1[deletedId].Value.Title {
				addedMatched[addedId] = deletedId
				deletedMatched[deletedId] = addedId
				matchedTitles.Add(title)
				break
			}
		}
	}

	return addedMatched, deletedMatched
}

func getNonContainedInlineSchemas(config *Config, state *state, schemaRefs1, schemaRefs2 openapi3.SchemaRefs, filter schemaRefsFilter, inlineSchemaPrefix string) ([]int, error) {

	notContainedIdx := []int{}
	matched := map[int]struct{}{}

	for index1, schemaRef1 := range schemaRefs1 {
		if !filter(schemaRef1) {
			continue
		}

		if found, index2, err := findIndenticalSchema(config, state, schemaRef1, schemaRefs2, matched, filter); err != nil {
			return nil, err
		} else if !found {
			notContainedIdx = append(notContainedIdx, index1)
		} else {
			matched[index2] = struct{}{}
		}
	}
	return notContainedIdx, nil
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

func getSchemaNames(inlineSchemaPrefix string, indexes []int, schemaRefs openapi3.SchemaRefs) utils.StringList {
	result := utils.StringList{}
	for _, index := range indexes {
		result = append(result, getSchemaName(inlineSchemaPrefix, schemaRefs[index], index))
	}
	return result
}

func getSchemaName(inlineSchemaPrefix string, schemaRef *openapi3.SchemaRef, index int) string {
	schemaName := fmt.Sprintf("%s[%d]", inlineSchemaPrefix, index)
	if schemaRef != nil && schemaRef.Value != nil && schemaRef.Value.Title != "" {
		return fmt.Sprintf("%s:%s", schemaName, schemaRef.Value.Title)
	}
	return schemaName
}
