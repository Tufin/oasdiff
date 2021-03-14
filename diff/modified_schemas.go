package diff

import (
	"github.com/getkin/kin-openapi/openapi3"
)

// ModifiedSchemas is map of schema names to their respective diffs
type ModifiedSchemas map[string]*SchemaDiff

func (modifiedSchemas ModifiedSchemas) addSchemaDiff(config *Config, schema1 string, schemaRef1, schemaRef2 *openapi3.SchemaRef) error {

	diff, err := getSchemaDiff(config, schemaRef1, schemaRef2)
	if err != nil {
		return err
	}
	if !diff.Empty() {
		modifiedSchemas[schema1] = diff
	}

	return nil
}
