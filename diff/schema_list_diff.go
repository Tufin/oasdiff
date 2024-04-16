package diff

import (
	"fmt"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/tufin/oasdiff/utils"
)

/*
SchemaListDiff describes the changes between a pair of subschemas under AllOf, AnyOf or OneOf
[oneOf, anyOf, allOf]: https://swagger.io/docs/specification/data-models/oneof-anyof-allof-not/
[Schema Objects]: https://swagger.io/specification/#schema-object
The SchemaListDiff is a combination of two diffs:

 1. Diff of schemas with a $ref:
    - schemas with the same $ref across base and revision are compared to each other and based on the result are considered as modified or unmodified
    - other schemas are considered added/deleted

 2. Diff of schemas without a $ref ("inline schemas"):
    Unlike schemas with $ref, inline schemas are not identified by a unique name and are therefor compared by their content.
    - syntactically identical schemas across base and revision are considered unmodified
    - schemas with the same title across base and revision are compared to each other and based on the result are considered as modified or unmodified
    - other schemas are considered added/deleted
    - special case: if there remains exactly one added schema and one deleted schema without titles, they will be be compared to eachother and considered as modified or unmodified

The SchemaListDiff format:
- Under Deleted: 1-based index in the deletred schema + schema title (if exists)
- Under Added: 1-based index in the added schema + schema title (if exists)
- Under Modified: schema title + diff of the schema objects
*/
type SchemaListDiff struct {
	Added    Schemas         `json:"added,omitempty" yaml:"added,omitempty"`
	Deleted  Schemas         `json:"deleted,omitempty" yaml:"deleted,omitempty"`
	Modified ModifiedSchemas `json:"modified,omitempty" yaml:"modified,omitempty"`
}

func NewSchemaListDiff() *SchemaListDiff {
	return &SchemaListDiff{
		Added:    Schemas{},
		Deleted:  Schemas{},
		Modified: ModifiedSchemas{},
	}
}

func (diff *SchemaListDiff) appendAddedSchema(schemaName string) {
	diff.Added = append(diff.Added, Schema{Title: schemaName})
}

func (diff *SchemaListDiff) appendDeletedSchema(schemaName string) {
	diff.Deleted = append(diff.Deleted, Schema{Title: schemaName})
}

type Schema struct {
	Index int    `json:"index" yaml:"index"`
	Title string `json:"title,omitempty" yaml:"title,omitempty"`
}

type Schemas []Schema

func getSchemas(indexes []int, schemaRefs openapi3.SchemaRefs) Schemas {
	result := Schemas{}
	for _, index := range indexes {
		result = append(result, Schema{
			Index: index,
			Title: schemaRefs[index].Value.Title,
		})
	}
	return result
}

func (schemas Schemas) String() string {
	result := ""
	for _, schema := range schemas {
		result += fmt.Sprintf("%d: %s\n", schema.Index, schema.Title)
	}
	return result
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
func getSchemaListsRefsDiff(config *Config, state *state, schemaRefs1, schemaRefs2 openapi3.SchemaRefs, filter schemaRefsFilter) (*SchemaListDiff, error) {
	result := NewSchemaListDiff()

	schemaMap2 := toSchemaRefMap(schemaRefs2, filter)
	for _, schema1 := range schemaRefs1 {
		if !filter(schema1) {
			continue
		}
		ref := schema1.Ref
		if schema2, found := schemaMap2[ref]; found {
			schemaMap2.delete(ref)
			var err error
			result.Modified, err = result.Modified.addSchemaDiff(config, state, ref, schema1, schema2.schemaRef)
			if err != nil {
				return result, err
			}
		} else {
			result.appendDeletedSchema(ref[strings.LastIndex(ref, "/")+1:])
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
			result.appendAddedSchema(ref[strings.LastIndex(ref, "/")+1:])
		}
	}
	return result, nil
}

func getSchemaListsInlineDiff(config *Config, state *state, schemaRefs1, schemaRefs2 openapi3.SchemaRefs, filter schemaRefsFilter) (SchemaListDiff, error) {

	// find schemas in revision that have no matching schema in the base
	addedIdx, err := getNonContainedInlineSchemas(config, state, schemaRefs2, schemaRefs1, filter)
	if err != nil {
		return SchemaListDiff{}, err
	}

	// find schemas in base that have no matching schema in the revision
	deletedIdx, err := getNonContainedInlineSchemas(config, state, schemaRefs1, schemaRefs2, filter)
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
		var err error
		modifiedSchemas, err = modifiedSchemas.addSchemaDiff(config, state, fmt.Sprintf("#%d", 1+deletedIdx[0]), schemaRefs1[deletedIdx[0]], schemaRefs2[addedIdx[0]])
		if err != nil {
			return SchemaListDiff{}, err
		}
		addedIdx = []int{}
		deletedIdx = []int{}
	}

	return SchemaListDiff{
		Added:    getSchemas(addedIdx, schemaRefs2),
		Deleted:  getSchemas(deletedIdx, schemaRefs1),
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

		var err error
		modifiedSchemas, err = modifiedSchemas.addSchemaDiff(config, state, schemaRefs1[deletedId].Value.Title, schemaRefs1[deletedId], schemaRefs2[addedId])
		if err != nil {
			return nil, nil, nil, err
		}
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

func getNonContainedInlineSchemas(config *Config, state *state, schemaRefs1, schemaRefs2 openapi3.SchemaRefs, filter schemaRefsFilter) ([]int, error) {

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
