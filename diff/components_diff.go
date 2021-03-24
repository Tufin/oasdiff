package diff

import "github.com/getkin/kin-openapi/openapi3"

// ComponentsDiff describes the changes between a pair of component objects: https://swagger.io/specification/#components-object
type ComponentsDiff struct {
	SchemasDiff         *SchemasDiff         `json:"schemas,omitempty" yaml:"schemas,omitempty"`
	ParametersDiff      *ParametersDiff      `json:"parameters,omitempty" yaml:"parameters,omitempty"`
	HeadersDiff         *HeadersDiff         `json:"headers,omitempty" yaml:"headers,omitempty"`
	RequestBodiesDiff   *RequestBodiesDiff   `json:"requestBodies,omitempty" yaml:"requestBodies,omitempty"`
	ResponsesDiff       *ResponsesDiff       `json:"responses,omitempty" yaml:"responses,omitempty"`
	SecuritySchemesDiff *SecuritySchemesDiff `json:"securitySchemes,omitempty" yaml:"securitySchemes,omitempty"`
	ExamplesDiff        *ExamplesDiff        `json:"examples,omitempty" yaml:"examples,omitempty"`
	LinksDiff           *LinksDiff           `json:"links,omitempty" yaml:"links,omitempty"`
	CallbacksDiff       *CallbacksDiff       `json:"callbacks,omitempty" yaml:"callbacks,omitempty"`
}

func getComponentsDiff(config *Config, s1, s2 openapi3.Components) (ComponentsDiff, error) {

	result := ComponentsDiff{}
	var err error

	result.SchemasDiff, err = getSchemasDiff(config, s1.Schemas, s2.Schemas)
	if err != nil {
		return result, err
	}

	result.ParametersDiff, err = getParametersDiff(config, toParameters(s1.Parameters), toParameters(s2.Parameters))
	if err != nil {
		return result, err
	}

	result.HeadersDiff, err = getHeadersDiff(config, s1.Headers, s2.Headers)
	if err != nil {
		return result, err
	}

	result.RequestBodiesDiff, err = getRequestBodiesDiff(config, s1.RequestBodies, s2.RequestBodies)
	if err != nil {
		return result, err
	}

	result.ResponsesDiff, err = getResponsesDiff(config, s1.Responses, s2.Responses)
	if err != nil {
		return result, err
	}

	result.SecuritySchemesDiff, err = getSecuritySchemesDiff(config, s1.SecuritySchemes, s2.SecuritySchemes)
	if err != nil {
		return result, err
	}

	result.ExamplesDiff, err = getExamplesDiff(config, s1.Examples, s2.Examples)
	if err != nil {
		return result, err
	}

	result.LinksDiff, err = getLinksDiff(config, s1.Links, s2.Links)
	if err != nil {
		return result, err
	}

	result.CallbacksDiff, err = getCallbacksDiff(config, s1.Callbacks, s2.Callbacks)
	if err != nil {
		return result, err
	}

	return result, nil
}
