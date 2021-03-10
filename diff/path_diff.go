package diff

import (
	"github.com/getkin/kin-openapi/openapi3"
)

// PathDiff is a diff between path item objects: https://swagger.io/specification/#path-item-object
type PathDiff struct {
	SummaryDiff     *ValueDiff      `json:"summary,omitempty" yaml:"summary,omitempty"`
	DescriptionDiff *ValueDiff      `json:"description,omitempty" yaml:"description,omitempty"`
	OperationsDiff  *OperationsDiff `json:"operations,omitempty" yaml:"operations,omitempty"`
	ServersDiff     *ServersDiff    `json:"servers,omitempty" yaml:"servers,omitempty"`
	ParametersDiff  *ParametersDiff `json:"parameters,omitempty" yaml:"parameters,omitempty"`
}

func newPathDiff() *PathDiff {
	return &PathDiff{}
}

// Empty return true if there is no diff
func (pathDiff *PathDiff) Empty() bool {
	if pathDiff == nil {
		return true
	}

	return pathDiff == nil || *pathDiff == *newPathDiff()
}

func getPathDiff(config *Config, pathItem1, pathItem2 *openapi3.PathItem) *PathDiff {
	diff := getPathDiffInternal(config, pathItem1, pathItem2)
	if diff.Empty() {
		return nil
	}
	return diff
}

func getPathDiffInternal(config *Config, pathItem1, pathItem2 *openapi3.PathItem) *PathDiff {
	result := newPathDiff()

	result.SummaryDiff = getValueDiff(pathItem1.Summary, pathItem2.Summary)
	result.DescriptionDiff = getValueDiff(pathItem1.Description, pathItem2.Description)
	result.OperationsDiff = getOperationsDiff(config, pathItem1, pathItem2)
	result.ServersDiff = getServersDiff(config, &pathItem1.Servers, &pathItem2.Servers)
	result.ParametersDiff = getParametersDiff(config, pathItem1.Parameters, pathItem2.Parameters)

	return result
}
