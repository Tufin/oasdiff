package diff

import (
	"github.com/getkin/kin-openapi/openapi3"
)

// PathDiff is a diff between path item objects: https://swagger.io/specification/#path-item-object
type PathDiff struct {
	SummaryDiff     *ValueDiff      `json:"summary,omitempty"`
	DescriptionDiff *ValueDiff      `json:"description,omitempty"`
	OperationsDiff  *OperationsDiff `json:"operations,omitempty"`
	ServersDiff     *ServersDiff    `json:"servers,omitempty"`
	ParametersDiff  *ParametersDiff `json:"parameters,omitempty"`
}

func newPathDiff() *PathDiff {
	return &PathDiff{}
}

func (pathDiff *PathDiff) empty() bool {
	return pathDiff == nil || *pathDiff == *newPathDiff()
}

func getPathDiff(config *Config, pathItem1, pathItem2 *openapi3.PathItem) *PathDiff {
	result := newPathDiff()

	result.SummaryDiff = getValueDiff(pathItem1.Summary, pathItem2.Summary)
	result.DescriptionDiff = getValueDiff(pathItem1.Description, pathItem2.Description)
	result.OperationsDiff = getOperationsDiff(config, pathItem1, pathItem2)

	if diff := getServersDiff(config, &pathItem1.Servers, &pathItem2.Servers); !diff.empty() {
		result.ServersDiff = diff
	}

	if diff := getParametersDiff(config, pathItem1.Parameters, pathItem2.Parameters); !diff.empty() {
		result.ParametersDiff = diff
	}

	return result
}
