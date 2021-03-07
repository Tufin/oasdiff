package diff

import "github.com/getkin/kin-openapi/openapi3"

// SecuritySchemeDiff is a diff between security scheme objects: https://swagger.io/specification/#security-scheme-object
type SecuritySchemeDiff struct {
	ExtensionProps   *ExtensionsDiff `json:"extensions,omitempty"`
	TypeDiff         *ValueDiff      `json:"type,omitempty"`
	DescriptionDiff  *ValueDiff      `json:"description,omitempty"`
	NameDiff         *ValueDiff      `json:"name,omitempty"`
	InDiff           *ValueDiff      `json:"in,omitempty"`
	SchemeDiff       *ValueDiff      `json:"scheme,omitempty"`
	BearerFormatDiff *ValueDiff      `json:"bearerFormat,omitempty"`
	// Flows
	OpenIDConnectURLDiff *ValueDiff `json:"openIDConnectURL,omitempty"`
}

func (securitySchemeDiff SecuritySchemeDiff) empty() bool {
	return securitySchemeDiff == SecuritySchemeDiff{}
}

func diffSecuritySchemes(config *Config, scheme1, scheme2 *openapi3.SecurityScheme) SecuritySchemeDiff {
	result := SecuritySchemeDiff{}

	if diff := getExtensionsDiff(config, scheme1.ExtensionProps, scheme2.ExtensionProps); !diff.empty() {
		result.ExtensionProps = diff
	}

	result.TypeDiff = getValueDiff(scheme1.Type, scheme2.Type)
	result.DescriptionDiff = getValueDiff(scheme1.Description, scheme2.Description)
	result.NameDiff = getValueDiff(scheme1.Name, scheme2.Name)
	result.InDiff = getValueDiff(scheme1.In, scheme2.In)
	result.SchemeDiff = getValueDiff(scheme1.Scheme, scheme2.Scheme)
	result.BearerFormatDiff = getValueDiff(scheme1.BearerFormat, scheme2.BearerFormat)
	// Flows            *OAuthFlows `json:"flows,omitempty" yaml:"flows,omitempty"`
	result.OpenIDConnectURLDiff = getValueDiff(scheme1.OpenIdConnectUrl, scheme2.OpenIdConnectUrl)

	return result
}
