package diff

import (
	"github.com/getkin/kin-openapi/openapi3"
)

// ServerDiff is a diff between server objects: https://swagger.io/specification/#server-object
type ServerDiff struct {
	ExtensionsDiff  *ExtensionsDiff `json:"extensions,omitempty" yaml:"extensions,omitempty"`
	URLDiff         *ValueDiff      `json:"urlType,omitempty" yaml:"urlType,omitempty"`
	DescriptionDiff *ValueDiff      `json:"description,omitempty" yaml:"description,omitempty"`
	// Variables
}

// Empty return true if there is no diff
func (diff *ServerDiff) Empty() bool {
	return diff == nil || *diff == ServerDiff{}
}

func getServerDiff(config *Config, value1, value2 *openapi3.Server) *ServerDiff {
	diff := getServerDiffInternal(config, value1, value2)
	if diff.Empty() {
		return nil
	}
	return diff
}

func getServerDiffInternal(config *Config, value1, value2 *openapi3.Server) *ServerDiff {
	result := ServerDiff{}

	result.ExtensionsDiff = getExtensionsDiff(config, value1.ExtensionProps, value2.ExtensionProps)
	result.URLDiff = getValueDiff(value1.URL, value2.URL)
	result.DescriptionDiff = getValueDiff(value1.Description, value2.Description)

	return &result
}
