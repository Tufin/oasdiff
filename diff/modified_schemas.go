package diff

import (
	"github.com/getkin/kin-openapi/openapi3"
)

// ModifiedSchemas is map of schema names to their respective diffs
type ModifiedSchemas map[string]*SchemaDiff

// ToStringList returns the modified schema names
func (modifiedSchemas ModifiedSchemas) ToStringList() StringList {
	keys := make(StringList, len(modifiedSchemas))
	i := 0
	for k := range modifiedSchemas {
		keys[i] = k
		i++
	}
	return keys
}

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
