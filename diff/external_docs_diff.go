package diff

import "github.com/getkin/kin-openapi/openapi3"

// ExternalDocsDiff describes the changes between a pair of external documentation objects: https://swagger.io/specification/#external-documentation-object
type ExternalDocsDiff struct {
	Added           bool            `json:"added,omitempty" yaml:"added,omitempty"`
	Deleted         bool            `json:"deleted,omitempty" yaml:"deleted,omitempty"`
	ExtensionsDiff  *ExtensionsDiff `json:"extensions,omitempty" yaml:"extensions,omitempty"`
	DescriptionDiff *ValueDiff      `json:"description,omitempty" yaml:"description,omitempty"`
	URLDiff         *ValueDiff      `json:"url,omitempty" yaml:"url,omitempty"`
}

func newExternalDocsDiff() *ExternalDocsDiff {
	return &ExternalDocsDiff{}
}

// Empty indicates whether a change was found in this element
func (diff *ExternalDocsDiff) Empty() bool {
	return diff == nil || *diff == ExternalDocsDiff{}
}

func getExternalDocsDiff(config *Config, state *state, docs1, docs2 *openapi3.ExternalDocs) *ExternalDocsDiff {
	diff := getExternalDocsDiffInternal(config, state, docs1, docs2)

	if diff.Empty() {
		return nil
	}

	return diff
}

func getExternalDocsDiffInternal(config *Config, state *state, docs1, docs2 *openapi3.ExternalDocs) *ExternalDocsDiff {
	result := newExternalDocsDiff()

	if docs1 == nil && docs2 == nil {
		return result
	}

	if docs1 == nil && docs2 != nil {
		result.Added = true
		return result
	}

	if docs1 != nil && docs2 == nil {
		result.Deleted = true
		return result
	}

	result.ExtensionsDiff = getExtensionsDiff(config, state, docs1.Extensions, docs2.Extensions)
	result.DescriptionDiff = getValueDiffConditional(config.ExcludeDescription, docs1.Description, docs2.Description)
	result.URLDiff = getValueDiff(docs1.URL, docs2.URL)

	return result
}
