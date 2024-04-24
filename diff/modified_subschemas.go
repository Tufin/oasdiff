package diff

import (
	"fmt"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/tufin/oasdiff/utils"
)

// ModifiedSubschemas is list of modified subschemas with their diffs
type ModifiedSubschemas []*ModifiedSubschema

// ModifiedSubschema is a single modified subschema with its diff
type ModifiedSubschema struct {
	Base     Subschema   `json:"base" yaml:"base"`
	Revision Subschema   `json:"revision" yaml:"revision"`
	Diff     *SchemaDiff `json:"diff" yaml:"diff"`
}

// String returns a string representation of the modified subschema
func (modifiedSchema *ModifiedSubschema) String() string {
	baseName := modifiedSchema.Base.String()
	revisonName := modifiedSchema.Revision.String()

	if baseName == revisonName {
		return baseName
	}

	return fmt.Sprintf("%s -> %s", baseName, revisonName)
}

// Subschemas is a list of subschemas
type Subschemas []Subschema

// Subschema identifies a subschema by its index, component and title
type Subschema struct {
	Index     int    `json:"index" yaml:"index"`                             // zero-based index in the schema's subschemas
	Component string `json:"component,omitempty" yaml:"component,omitempty"` // component name if the subschema is a component
	Title     string `json:"title,omitempty" yaml:"title,omitempty"`         // title of the subschema
}

// String returns a string representation of the subschema
func (subschema Subschema) String() string {
	const prefix = "subschema"

	if subschema.Title != "" {
		return fmt.Sprintf("%s #%d: %s", prefix, subschema.Index+1, subschema.Title)
	}

	if subschema.Component != "" {
		return fmt.Sprintf("#/components/schemas/%s", subschema.Component)
	}

	return fmt.Sprintf("%s #%d", prefix, subschema.Index+1)
}

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

// String returns a string representation of the subschemas
func (schemas Subschemas) String() string {
	names := make([]string, len(schemas))
	for i, schema := range schemas {
		names[i] = schema.String()
	}
	list := utils.StringList(names)
	return list.String()
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
