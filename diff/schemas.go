package diff

import (
	"github.com/getkin/kin-openapi/openapi3"
)

type SchemaCollectionDiff struct {
	AddedSchemas    []string `json:"addedSchemas,omitempty"`
	DeletedSchemas  []string `json:"deletedSchemas,omitempty"`
	ModifiedSchemas []string `json:"modifiedSchemas,omitempty"`
}

func (diff *SchemaCollectionDiff) empty() bool {
	return len(diff.AddedSchemas) == 0 &&
		len(diff.DeletedSchemas) == 0 &&
		len(diff.ModifiedSchemas) == 0
}

func newSchemaDiff() *SchemaCollectionDiff {
	return &SchemaCollectionDiff{
		AddedSchemas:    []string{},
		DeletedSchemas:  []string{},
		ModifiedSchemas: []string{},
	}
}

func diffSchemaCollection(schemas1 openapi3.Schemas, schemas2 openapi3.Schemas) *SchemaCollectionDiff {

	result := newSchemaDiff()

	addedSchemas, deletedSchemas, modifiedSchemas := diffSchemas(schemas1, schemas2)

	for schema := range addedSchemas {
		result.addAddedSchema(schema)
	}

	for schema := range deletedSchemas {
		result.addDeletedSchema(schema)
	}

	for schema := range modifiedSchemas {
		result.addModifiedSchema(schema)
	}

	return result
}

func diffSchemas(schemas1 openapi3.Schemas, schemas2 openapi3.Schemas) (openapi3.Schemas, openapi3.Schemas, openapi3.Schemas) {

	added := openapi3.Schemas{}
	deleted := openapi3.Schemas{}
	modified := openapi3.Schemas{}

	for schemaName1, schemaRef1 := range schemas1 {
		schemaRef2, ok := schemas2[schemaName1]
		if !ok {
			deleted[schemaName1] = schemaRef1
			continue
		}

		if diff := diffSchema(schemaRef1, schemaRef2); !diff.empty() {
			modified[schemaName1] = schemaRef1
		}
	}

	for schemaName2, schemaRef2 := range schemas2 {
		_, ok := schemas1[schemaName2]
		if !ok {
			added[schemaName2] = schemaRef2
		}
	}

	return added, deleted, modified
}

func (diff *SchemaCollectionDiff) addAddedSchema(schema string) {
	diff.AddedSchemas = append(diff.AddedSchemas, schema)
}

func (diff *SchemaCollectionDiff) addDeletedSchema(schema string) {
	diff.DeletedSchemas = append(diff.DeletedSchemas, schema)
}

func (diff *SchemaCollectionDiff) addModifiedSchema(schema string) {
	diff.ModifiedSchemas = append(diff.ModifiedSchemas, schema)
}

func (diff *SchemaCollectionDiff) getSummary() SchemaDiffSummary {
	return SchemaDiffSummary{
		AddedSchemas:    len(diff.AddedSchemas),
		DeletedSchemas:  len(diff.DeletedSchemas),
		ModifiedSchemas: len(diff.ModifiedSchemas),
	}
}
