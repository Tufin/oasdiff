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
	Base     BaseSubschema     `json:"base" yaml:"base"`
	Revision RevisionSubschema `json:"revision" yaml:"revision"`
	Diff     *SchemaDiff       `json:"diff" yaml:"diff"`
}

func (modifiedSchema *ModifiedSubschema) String() string {
	baseName := modifiedSchema.Base.String()
	revisonName := modifiedSchema.Revision.String()

	if baseName == revisonName {
		return baseName
	}

	return fmt.Sprintf("%s -> %s", baseName, revisonName)
}

// Subschema identifies a subschema by its index and title
type Subschema struct {
	Index     int    `json:"index" yaml:"index"`
	Component string `json:"component,omitempty" yaml:"component,omitempty"`
	Title     string `json:"title,omitempty" yaml:"title,omitempty"`
}

type BaseSubschema struct {
	Subschema
}

func (subschema BaseSubschema) String() string {
	return subschema.Subschema.String("BaseSchema")
}

type RevisionSubschema struct {
	Subschema
}

func (subschema RevisionSubschema) String() string {
	return subschema.Subschema.String("RevisionSchema")
}

func (subschema Subschema) String(prefix string) string {
	if subschema.Title != "" {
		return fmt.Sprintf("%s[%d]:%s", prefix, subschema.Index, subschema.Title)
	}

	if subschema.Component != "" {
		return fmt.Sprintf("#/components/schemas/%s", subschema.Component)
	}

	return fmt.Sprintf("schema #%d", subschema.Index)
}

// BaseSubschemas is a list of Base Subschemas
type BaseSubschemas []BaseSubschema

// RevisionSubschemas is a list of Revision Subschemas
type RevisionSubschemas []RevisionSubschema

func getBaseSubschemas(indexes []int, schemaRefs openapi3.SchemaRefs) BaseSubschemas {
	result := BaseSubschemas{}
	for _, index := range indexes {
		result = append(result, BaseSubschema{
			Subschema: Subschema{
				Index: index,
				Title: schemaRefs[index].Value.Title,
			},
		})
	}
	return result
}

func getRevisionSubschemas(indexes []int, schemaRefs openapi3.SchemaRefs) RevisionSubschemas {
	result := RevisionSubschemas{}
	for _, index := range indexes {
		result = append(result, RevisionSubschema{
			Subschema: Subschema{
				Index: index,
				Title: schemaRefs[index].Value.Title,
			},
		})
	}
	return result
}

func (schemas BaseSubschemas) String() string {
	names := make([]string, len(schemas))
	for i, schema := range schemas {
		names[i] = schema.String()
	}
	list := utils.StringList(names)
	return list.String()
}

func (schemas RevisionSubschemas) String() string {
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
			Base: BaseSubschema{
				Subschema: Subschema{
					Index:     index1,
					Component: getComponentName(schemaRef1),
					Title:     schemaRef1.Value.Title,
				},
			},
			Revision: RevisionSubschema{
				Subschema: Subschema{
					Index:     index2,
					Component: getComponentName(schemaRef2),
					Title:     schemaRef2.Value.Title,
				},
			},
			Diff: diff,
		})
	}

	return modifiedSchemas, nil
}

func (modifiedSchemas ModifiedSubschemas) combine(other ModifiedSubschemas) ModifiedSubschemas {
	return append(modifiedSchemas, other...)
}
