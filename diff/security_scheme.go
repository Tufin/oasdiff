package diff

import "github.com/getkin/kin-openapi/openapi3"

// SecuritySchemeDiff is a diff between security scheme objects: https://swagger.io/specification/#security-scheme-object
type SecuritySchemeDiff struct {
	ExtensionsDiff       *ExtensionsDiff `json:"extensions,omitempty"`
	TypeDiff             *ValueDiff      `json:"type,omitempty"`
	DescriptionDiff      *ValueDiff      `json:"description,omitempty"`
	NameDiff             *ValueDiff      `json:"name,omitempty"`
	InDiff               *ValueDiff      `json:"in,omitempty"`
	SchemeDiff           *ValueDiff      `json:"scheme,omitempty"`
	BearerFormatDiff     *ValueDiff      `json:"bearerFormat,omitempty"`
	OAuthFlowsDiff       *OAuthFlowsDiff `json:"OAuthFlows,omitempty"`
	OpenIDConnectURLDiff *ValueDiff      `json:"openIDConnectURL,omitempty"`
}

func (diff SecuritySchemeDiff) empty() bool {
	return diff == SecuritySchemeDiff{}
}

func getSecuritySchemeDiff(config *Config, scheme1, scheme2 *openapi3.SecurityScheme) SecuritySchemeDiff {
	result := SecuritySchemeDiff{}

	result.ExtensionsDiff = getExtensionsDiff(config, scheme1.ExtensionProps, scheme2.ExtensionProps)
	result.TypeDiff = getValueDiff(scheme1.Type, scheme2.Type)
	result.DescriptionDiff = getValueDiff(scheme1.Description, scheme2.Description)
	result.NameDiff = getValueDiff(scheme1.Name, scheme2.Name)
	result.InDiff = getValueDiff(scheme1.In, scheme2.In)
	result.SchemeDiff = getValueDiff(scheme1.Scheme, scheme2.Scheme)
	result.BearerFormatDiff = getValueDiff(scheme1.BearerFormat, scheme2.BearerFormat)
	result.OAuthFlowsDiff = getOAuthFlowsDiff(config, scheme1.Flows, scheme2.Flows)
	result.OpenIDConnectURLDiff = getValueDiff(scheme1.OpenIdConnectUrl, scheme2.OpenIdConnectUrl)

	return result
}
