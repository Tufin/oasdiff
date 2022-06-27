package diff

import (
	"github.com/getkin/kin-openapi/openapi3"
)

// ModifiedSchemas is map of schema names to their respective diffs
type ModifiedSchemas map[string]*SchemaDiff

func (modifiedSchemas ModifiedSchemas) addSchemaDiff(config *Config, state *state, schemaName string, schemaRef1, schemaRef2 *openapi3.SchemaRef) error {

	diff, err := getSchemaDiff(config, state, schemaRef1, schemaRef2)
	if err != nil {
		return err
	}
	if !diff.Empty() {
		modifiedSchemas[schemaName] = diff
	}

	return nil
}

func (modifiedSchemas ModifiedSchemas) combine(other ModifiedSchemas) ModifiedSchemas {
	result := ModifiedSchemas{}

	for ref, d := range modifiedSchemas {
		result[ref] = d
	}
	for ref, d := range other {
		result[ref] = d
	}
	return result
}
