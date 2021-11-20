package diff

import (
	"github.com/getkin/kin-openapi/openapi3"
)

// ServerDiff describes the changes between a pair of server objects: https://swagger.io/specification/#server-object
type ServerDiff struct {
	Added           bool            `json:"added,omitempty" yaml:"added,omitempty"`
	Deleted         bool            `json:"deleted,omitempty" yaml:"deleted,omitempty"`
	ExtensionsDiff  *ExtensionsDiff `json:"extensions,omitempty" yaml:"extensions,omitempty"`
	URLDiff         *ValueDiff      `json:"urlType,omitempty" yaml:"urlType,omitempty"`
	DescriptionDiff *ValueDiff      `json:"description,omitempty" yaml:"description,omitempty"`
	VariablesDiff   *VariablesDiff  `json:"variables,omitempty" yaml:"variables,omitempty"`
}

// Empty indicates whether a change was found in this element
func (diff *ServerDiff) Empty() bool {
	return diff == nil || *diff == ServerDiff{}
}

// Breaking indicates whether this element includes a breaking change
func (diff *ServerDiff) Breaking() bool {
	return false
}

func getServerDiff(config *Config, value1, value2 *openapi3.Server) *ServerDiff {
	diff := getServerDiffInternal(config, value1, value2)

	if diff.Empty() {
		return nil
	}

	if config.BreakingOnly && !diff.Breaking() {
		return nil
	}

	return diff
}

func getServerDiffInternal(config *Config, value1, value2 *openapi3.Server) *ServerDiff {

	if value1 == nil && value2 == nil {
		return nil
	}

	result := ServerDiff{}

	if value1 == nil && value2 != nil {
		result.Added = true
		return &result
	}

	if value1 != nil && value2 == nil {
		result.Deleted = true
		return &result
	}

	result.ExtensionsDiff = getExtensionsDiff(config, value1.ExtensionProps, value2.ExtensionProps)
	result.URLDiff = getValueDiff(config, false, value1.URL, value2.URL)
	result.DescriptionDiff = getValueDiffConditional(config, false, config.ExcludeDescription, value1.Description, value2.Description)
	result.VariablesDiff = getVariablesDiff(config, value1.Variables, value2.Variables)

	return &result
}
