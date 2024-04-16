package diff

import (
	"github.com/getkin/kin-openapi/openapi3"
)

type ModifiedSchema struct {
	Base     Schema      `json:"base,omitempty" yaml:"base,omitempty"`
	Revision Schema      `json:"revision,omitempty" yaml:"revision,omitempty"`
	Diff     *SchemaDiff `json:"diff,omitempty" yaml:"diff,omitempty"`
}

type ModifiedSchemas []*ModifiedSchema

// ModifiedSchemasOld is map of schema names to their respective diffs
type ModifiedSchemasOld map[string]*SchemaDiff

func (modifiedSchemas ModifiedSchemas) addSchemaDiff(config *Config, state *state, schemaName string, schemaRef1, schemaRef2 *openapi3.SchemaRef) (ModifiedSchemas, error) {

	diff, err := getSchemaDiff(config, state, schemaRef1, schemaRef2)
	if err != nil {
		return nil, err
	}
	if !diff.Empty() {
		modifiedSchemas = append(modifiedSchemas, &ModifiedSchema{
			Base: Schema{
				// Index: schemaRef1.Index,
				Title: schemaRef1.Value.Title,
			},
			Revision: Schema{
				// Index: schemaRef2.Index,
				Title: schemaRef2.Value.Title,
			},
			Diff: diff,
		})
	}

	return modifiedSchemas, nil
}

func (modifiedSchemas ModifiedSchemas) combine(other ModifiedSchemas) ModifiedSchemas {
	return append(modifiedSchemas, other...)
}
