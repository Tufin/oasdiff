package diff

import (
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/tufin/oasdiff/utils"
)

// RequiredPropertiesDiff describes the changes between a pair of lists of required properties
type RequiredPropertiesDiff struct {
	StringsDiff
}

// Empty indicates whether a change was found in this element
func (diff *RequiredPropertiesDiff) Empty() bool {
	if diff == nil {
		return true
	}

	return diff.StringsDiff.Empty()
}

func propDeleted(property string, schema1, schema2 *openapi3.Schema) bool {
	if schema1 == nil || schema2 == nil {
		return false
	}

	_, ok1 := schema1.Properties[property]
	_, ok2 := schema2.Properties[property]

	return ok1 && !ok2
}

func propSunsetAllowed(property string, schema1 *openapi3.Schema) bool {
	if schema1 == nil {
		return false
	}

	schemaRef, ok := schema1.Properties[property]
	if !ok || schemaRef == nil || schemaRef.Value == nil {
		return false
	}

	return SunsetAllowed(schemaRef.Value.Deprecated, schemaRef.Value.Extensions)
}

func getRequiredPropertiesDiff(config *Config, state *state, schema1, schema2 *openapi3.Schema) *RequiredPropertiesDiff {
	diff := getRequiredPropertiesDiffInternal(schema1.Required, schema2.Required)

	if diff.Empty() {
		return nil
	}

	return diff
}

func getRequiredPropertiesDiffInternal(strings1, strings2 utils.StringList) *RequiredPropertiesDiff {
	if stringsDiff := getStringsDiff(strings1, strings2); stringsDiff != nil {
		return &RequiredPropertiesDiff{
			StringsDiff: *stringsDiff,
		}
	}
	return nil
}
