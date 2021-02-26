package diff

import (
	"github.com/getkin/kin-openapi/openapi3"
)

// ServerDiff is a diff between server objects: https://swagger.io/specification/#server-object
type ServerDiff struct {

	// TODO: diff ExtensionProps
	URLDiff         *ValueDiff `json:"urlType,omitempty"`
	DescriptionDiff *ValueDiff `json:"description,omitempty"`
	// TODO: diff Variables
}

func (diff ServerDiff) empty() bool {
	return diff == ServerDiff{}
}

func getServerDiff(value1, value2 *openapi3.Server) ServerDiff {
	result := ServerDiff{}

	result.URLDiff = getValueDiff(value1.URL, value2.URL)
	result.DescriptionDiff = getValueDiff(value1.Description, value2.Description)

	return result
}
