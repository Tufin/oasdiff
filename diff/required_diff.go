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

func (diff *RequiredPropertiesDiff) removeReadOnly(state *state, value1, value2 *openapi3.Schema) {
	if diff.Empty() || diff.StringsDiff.Empty() || state.direction == directionResponse {
		// readonly properties are only valid for responses
		return
	}
	added := make(StringList, 0)
	for _, v := range diff.Added {
		if p, ok := value2.Properties[v]; ok && !p.Value.ReadOnly {
			added = append(added, v)
		}
	}
	diff.Added = added
	deleted := make(StringList, 0)
	for _, v := range diff.Deleted {
		if p, ok := value1.Properties[v]; ok && !p.Value.ReadOnly {
			deleted = append(deleted, v)
		}
	}
	diff.Deleted = deleted
}

func (diff *RequiredPropertiesDiff) removeWriteOnly(state *state, value1, value2 *openapi3.Schema) {
	if diff.Empty() || diff.StringsDiff.Empty() || state.direction == directionRequest {
		// writeOnly properties are only valid for requests
		return
	}
	added := make(StringList, 0)
	for _, v := range diff.Added {
		if p, ok := value2.Properties[v]; ok && !p.Value.WriteOnly {
			added = append(added, v)
		}
	}
	diff.Added = added
	deleted := make(StringList, 0)
	for _, v := range diff.Deleted {
		if p, ok := value1.Properties[v]; ok && !p.Value.WriteOnly {
			deleted = append(deleted, v)
		}
	}
	diff.Deleted = deleted
}

func getRequiredPropertiesDiff(config *Config, state *state, value1, value2 *openapi3.Schema) *RequiredPropertiesDiff {
	diff := getRequiredPropertiesDiffInternal(value1.Required, value2.Required)

	diff.removeReadOnly(state, value1, value2)
	diff.removeWriteOnly(state, value1, value2)

	if config.BreakingOnly {
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
