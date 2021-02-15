package diff

import (
	"github.com/getkin/kin-openapi/openapi3"
)

type ModifiedSchemas map[string]*SchemaDiff

func (modifiedSchemas ModifiedSchemas) addSchemaDiff(schema1 string, schemaRef1 *openapi3.SchemaRef, schemaRef2 *openapi3.SchemaRef) {

	diff := diffSchema(schemaRef1, schemaRef2)
	if !diff.empty() {
		modifiedSchemas[schema1] = &diff
	}
}
