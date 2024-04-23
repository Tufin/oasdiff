package diff

import (
	"fmt"

	"github.com/getkin/kin-openapi/openapi3"
)

// ModifiedSubschemas is list of modified subschemas with their diffs
type ModifiedSubschemas []*ModifiedSubschema

// ModifiedSubschema is a single modified subschema with its diff
type ModifiedSubschema struct {
	Base     Subschema   `json:"base" yaml:"base"`
	Revision Subschema   `json:"revision" yaml:"revision"`
	Diff     *SchemaDiff `json:"diff" yaml:"diff"`
}

// Subschema identifies a subschema by its index and title
type Subschema struct {
	Index     int    `json:"index" yaml:"index"`
	Component string `json:"component,omitempty" yaml:"component,omitempty"`
	Title     string `json:"title,omitempty" yaml:"title,omitempty"`
}

// Subschemas is a list of subschemas
type Subschemas []Subschema

func getSubschemas(indexes []int, schemaRefs openapi3.SchemaRefs) Subschemas {
	result := Subschemas{}
	for _, index := range indexes {
		result = append(result, Subschema{
			Index: index,
			Title: schemaRefs[index].Value.Title,
		})
	}
	return result
}

func (schemas Subschemas) String() string {
	result := ""
	for _, schema := range schemas {
		result += fmt.Sprintf("%d: %s\n", schema.Index, schema.Title)
	}
	return result
}

func (modifiedSchemas ModifiedSubschemas) addSchemaDiff(config *Config, state *state, schemaRef1, schemaRef2 *openapi3.SchemaRef, index1, index2 int) (ModifiedSubschemas, error) {

	diff, err := getSchemaDiff(config, state, schemaRef1, schemaRef2)
	if err != nil {
		return nil, err
	}
	if !diff.Empty() {
		modifiedSchemas = append(modifiedSchemas, &ModifiedSubschema{
			Base: Subschema{
				Index:     index1,
				Component: getComponentName(schemaRef1),
				Title:     schemaRef1.Value.Title,
			},
			Revision: Subschema{
				Index:     index2,
				Component: getComponentName(schemaRef2),
				Title:     schemaRef2.Value.Title,
			},
			Diff: diff,
		})
	}

	return modifiedSchemas, nil
}

func (modifiedSchemas ModifiedSubschemas) combine(other ModifiedSubschemas) ModifiedSubschemas {
	return append(modifiedSchemas, other...)
}
