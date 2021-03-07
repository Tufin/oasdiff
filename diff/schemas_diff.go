package diff

import (
	"github.com/getkin/kin-openapi/openapi3"
)

// SchemasDiff is a diff between two sets of schema objects: https://swagger.io/specification/#schema-object
type SchemasDiff struct {
	Added    StringList      `json:"added,omitempty"`
	Deleted  StringList      `json:"deleted,omitempty"`
	Modified ModifiedSchemas `json:"modified,omitempty"`
}

func (schemasDiff *SchemasDiff) empty() bool {
	if schemasDiff == nil {
		return true
	}

	return len(schemasDiff.Added) == 0 &&
		len(schemasDiff.Deleted) == 0 &&
		len(schemasDiff.Modified) == 0
}

func newSchemasDiff() *SchemasDiff {
	return &SchemasDiff{
		Added:    StringList{},
		Deleted:  StringList{},
		Modified: ModifiedSchemas{},
	}
}

type schemaRefPair struct {
	SchemaRef1 *openapi3.SchemaRef
	SchemaRef2 *openapi3.SchemaRef
}

type schemaRefPairs map[string]*schemaRefPair

func getSchemasDiff(config *Config, schemas1, schemas2 openapi3.Schemas) *SchemasDiff {

	result := newSchemasDiff()

	addedSchemas, deletedSchemas, otherSchemas := diffSchemas(schemas1, schemas2)

	for schema := range addedSchemas {
		result.addAddedSchema(schema)
	}

	for schema := range deletedSchemas {
		result.addDeletedSchema(schema)
	}

	for schema, schemaRefPair := range otherSchemas {
		result.addModifiedSchema(config, schema, schemaRefPair.SchemaRef1, schemaRefPair.SchemaRef2)
	}

	if result.empty() {
		return nil
	}

	return result
}

func diffSchemas(schemas1, schemas2 openapi3.Schemas) (openapi3.Schemas, openapi3.Schemas, schemaRefPairs) {

	added := openapi3.Schemas{}
	deleted := openapi3.Schemas{}
	other := schemaRefPairs{}

	for schemaName1, schemaRef1 := range schemas1 {
		schemaRef2, ok := schemas2[schemaName1]
		if !ok {
			deleted[schemaName1] = schemaRef1
			continue
		}

		other[schemaName1] = &schemaRefPair{
			SchemaRef1: schemaRef1,
			SchemaRef2: schemaRef2,
		}
	}

	for schemaName2, schemaRef2 := range schemas2 {
		_, ok := schemas1[schemaName2]
		if !ok {
			added[schemaName2] = schemaRef2
		}
	}

	return added, deleted, other
}

func (schemasDiff *SchemasDiff) addAddedSchema(schema string) {
	schemasDiff.Added = append(schemasDiff.Added, schema)
}

func (schemasDiff *SchemasDiff) addDeletedSchema(schema string) {
	schemasDiff.Deleted = append(schemasDiff.Deleted, schema)
}

func (schemasDiff *SchemasDiff) addModifiedSchema(config *Config, schema1 string, schemaRef1, schemaRef2 *openapi3.SchemaRef) {
	schemasDiff.Modified.addSchemaDiff(config, schema1, schemaRef1, schemaRef2)
}

func (schemasDiff *SchemasDiff) getSummary() *SummaryDetails {
	return &SummaryDetails{
		Added:    len(schemasDiff.Added),
		Deleted:  len(schemasDiff.Deleted),
		Modified: len(schemasDiff.Modified),
	}
}
