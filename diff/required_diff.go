package diff

import (
	"github.com/getkin/kin-openapi/openapi3"
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

func (diff *RequiredPropertiesDiff) removeNonBreaking(state *state) {
	if diff.Empty() {
		return
	}

	if diff.StringsDiff.Empty() {
		return
	}

	switch state.direction {
	case directionRequest:
		// if this is part of the request, then required properties can be deleted without breaking the client
		diff.Deleted = nil
	case directionResponse:
		// if this is part of the response, then required properties can be added without breaking the client
		diff.Added = nil
	}
}

func (diff *RequiredPropertiesDiff) removeReadOnly(state *state, schema1, schema2 *openapi3.Schema) {
	if diff.Empty() || diff.StringsDiff.Empty() || state.direction == directionResponse {
		// readonly properties are only valid for responses
		return
	}
	added := make(StringList, 0)
	for _, v := range diff.Added {
		if p, ok := schema2.Properties[v]; ok && !p.Value.ReadOnly {
			added = append(added, v)
		}
	}
	diff.Added = added
	deleted := make(StringList, 0)
	for _, v := range diff.Deleted {
		if p, ok := schema1.Properties[v]; ok && !p.Value.ReadOnly {
			deleted = append(deleted, v)
		}
	}
	diff.Deleted = deleted
}

func (diff *RequiredPropertiesDiff) removeWriteOnly(state *state, schema1, schema2 *openapi3.Schema) {
	if diff.Empty() || diff.StringsDiff.Empty() || state.direction == directionRequest {
		// writeOnly properties are only valid for requests
		return
	}
	added := make(StringList, 0)
	for _, v := range diff.Added {
		if p, ok := schema2.Properties[v]; ok && !p.Value.WriteOnly {
			added = append(added, v)
		}
	}
	diff.Added = added
	deleted := make(StringList, 0)
	for _, v := range diff.Deleted {
		if p, ok := schema1.Properties[v]; ok && !p.Value.WriteOnly {
			deleted = append(deleted, v)
		}
	}
	diff.Deleted = deleted
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

func (diff *RequiredPropertiesDiff) removeSunsetProperties(state *state, schema1, schema2 *openapi3.Schema) {
	if diff.Empty() {
		return
	}

	if state.direction != directionResponse {
		return
	}

	deleted := make(StringList, 0)
	for _, property := range diff.Deleted {
		// if property was sunset then making it optional is not breaking
		if propDeleted(property, schema1, schema2) && propSunsetAllowed(property, schema1) {
			continue
		} else {
			deleted = append(deleted, property)
		}
	}
	diff.Deleted = deleted
}

func getRequiredPropertiesDiff(config *Config, state *state, schema1, schema2 *openapi3.Schema) *RequiredPropertiesDiff {
	diff := getRequiredPropertiesDiffInternal(schema1.Required, schema2.Required)

	if config.BreakingOnly {
		diff.removeReadOnly(state, schema1, schema2)
		diff.removeWriteOnly(state, schema1, schema2)
		diff.removeSunsetProperties(state, schema1, schema2)
		diff.removeNonBreaking(state)
	}

	if diff.Empty() {
		return nil
	}

	return diff
}

func getRequiredPropertiesDiffInternal(strings1, strings2 StringList) *RequiredPropertiesDiff {
	if stringsDiff := getStringsDiff(strings1, strings2); stringsDiff != nil {
		return &RequiredPropertiesDiff{
			StringsDiff: *stringsDiff,
		}
	}
	return nil
}
