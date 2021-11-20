package diff

import "github.com/getkin/kin-openapi/openapi3"

// LinkDiff describes the changes between a pair of link objects: https://swagger.io/specification/#link-object
type LinkDiff struct {
	ExtensionsDiff   *ExtensionsDiff   `json:"extensions,omitempty" yaml:"extensions,omitempty"`
	OperationIDDiff  *ValueDiff        `json:"operationId,omitempty" yaml:"operationId,omitempty"`
	OperationRefDiff *ValueDiff        `json:"operationRef,omitempty" yaml:"operationRef,omitempty"`
	DescriptionDiff  *ValueDiff        `json:"description,omitempty" yaml:"description,omitempty"`
	ParametersDiff   *InterfaceMapDiff `json:"parameters,omitempty" yaml:"parameters,omitempty"`
	ServerDiff       *ServerDiff       `json:"server,omitempty" yaml:"server,omitempty"`
	RequestBodyDiff  *ValueDiff        `json:"requestBody,omitempty" yaml:"requestBody,omitempty"`
}

// Empty indicates whether a change was found in this element
func (diff *LinkDiff) Empty() bool {
	return diff == nil || *diff == LinkDiff{}
}

// Breaking indicates whether this element includes a breaking change
func (diff *LinkDiff) Breaking() bool {
	return false
}

func getLinkDiff(config *Config, link1, link2 *openapi3.Link) (*LinkDiff, error) {
	diff, err := getLinkDiffInternal(config, link1, link2)
	if err != nil {
		return nil, err
	}

	if diff.Empty() {
		return nil, nil
	}

	if config.BreakingOnly && !diff.Breaking() {
		return nil, nil
	}

	return diff, nil
}

func getLinkDiffInternal(config *Config, link1, link2 *openapi3.Link) (*LinkDiff, error) {
	result := LinkDiff{}

	result.ExtensionsDiff = getExtensionsDiff(config, link1.ExtensionProps, link2.ExtensionProps)
	result.OperationIDDiff = getValueDiff(config, false, link1.OperationID, link2.OperationID)
	result.OperationRefDiff = getValueDiff(config, false, link1.OperationRef, link2.OperationRef)
	result.DescriptionDiff = getValueDiffConditional(config, false, config.ExcludeDescription, link1.Description, link2.Description)
	result.ParametersDiff = getInterfaceMapDiff(config, true, link1.Parameters, link2.Parameters, StringSet{})
	result.ServerDiff = getServerDiff(config, link1.Server, link2.Server)
	result.RequestBodyDiff = getValueDiff(config, false, link1.RequestBody, link2.RequestBody)

	return &result, nil
}
