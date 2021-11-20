package diff

import "github.com/getkin/kin-openapi/openapi3"

// SecuritySchemeDiff describes the changes between a pair of security scheme objects: https://swagger.io/specification/#security-scheme-object
type SecuritySchemeDiff struct {
	ExtensionsDiff       *ExtensionsDiff `json:"extensions,omitempty" yaml:"extensions,omitempty"`
	TypeDiff             *ValueDiff      `json:"type,omitempty" yaml:"type,omitempty"`
	DescriptionDiff      *ValueDiff      `json:"description,omitempty" yaml:"description,omitempty"`
	NameDiff             *ValueDiff      `json:"name,omitempty" yaml:"name,omitempty"`
	InDiff               *ValueDiff      `json:"in,omitempty" yaml:"in,omitempty"`
	SchemeDiff           *ValueDiff      `json:"scheme,omitempty" yaml:"scheme,omitempty"`
	BearerFormatDiff     *ValueDiff      `json:"bearerFormat,omitempty" yaml:"bearerFormat,omitempty"`
	OAuthFlowsDiff       *OAuthFlowsDiff `json:"OAuthFlows,omitempty" yaml:"OAuthFlows,omitempty"`
	OpenIDConnectURLDiff *ValueDiff      `json:"openIDConnectURL,omitempty" yaml:"openIDConnectURL,omitempty"`
}

// Empty indicates whether a change was found in this element
func (diff *SecuritySchemeDiff) Empty() bool {
	return diff == nil || *diff == SecuritySchemeDiff{}
}

// Breaking indicates whether this element includes a breaking change
func (diff *SecuritySchemeDiff) Breaking() bool {
	if diff.Empty() {
		return false
	}

	return diff.NameDiff.Breaking() ||
		diff.InDiff.Breaking() ||
		diff.SchemeDiff.Breaking() ||
		diff.BearerFormatDiff.Breaking() ||
		diff.OAuthFlowsDiff.Breaking() ||
		diff.OpenIDConnectURLDiff.Breaking()
}

func getSecuritySchemeDiff(config *Config, scheme1, scheme2 *openapi3.SecurityScheme) *SecuritySchemeDiff {
	diff := getSecuritySchemeDiffInternal(config, scheme1, scheme2)

	if diff.Empty() {
		return nil
	}

	if config.BreakingOnly && !diff.Breaking() {
		return nil
	}

	return diff
}

func getSecuritySchemeDiffInternal(config *Config, scheme1, scheme2 *openapi3.SecurityScheme) *SecuritySchemeDiff {
	result := SecuritySchemeDiff{}

	result.ExtensionsDiff = getExtensionsDiff(config, scheme1.ExtensionProps, scheme2.ExtensionProps)
	result.TypeDiff = getValueDiff(config, false, scheme1.Type, scheme2.Type)
	result.DescriptionDiff = getValueDiffConditional(config, false, config.ExcludeDescription, scheme1.Description, scheme2.Description)
	result.NameDiff = getValueDiff(config, false, scheme1.Name, scheme2.Name)
	result.InDiff = getValueDiff(config, false, scheme1.In, scheme2.In)
	result.SchemeDiff = getValueDiff(config, false, scheme1.Scheme, scheme2.Scheme)
	result.BearerFormatDiff = getValueDiff(config, false, scheme1.BearerFormat, scheme2.BearerFormat)
	result.OAuthFlowsDiff = getOAuthFlowsDiff(config, scheme1.Flows, scheme2.Flows)
	result.OpenIDConnectURLDiff = getValueDiff(config, false, scheme1.OpenIdConnectUrl, scheme2.OpenIdConnectUrl)

	return &result
}
