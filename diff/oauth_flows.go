package diff

import "github.com/getkin/kin-openapi/openapi3"

// OAuthFlowsDiff is a diff between oauth flows objects: https://swagger.io/specification/#oauth-flows-object
type OAuthFlowsDiff struct {
	ExtensionsDiff        *ExtensionsDiff `json:"extensions,omitempty" yaml:"extensions,omitempty"`
	ImplicitDiff          *OAuthFlowDiff  `json:"implicit,omitempty" yaml:"implicit,omitempty"`
	PasswordDiff          *OAuthFlowDiff  `json:"password,omitempty" yaml:"password,omitempty"`
	ClientCredentialsDiff *OAuthFlowDiff  `json:"clientCredentials,omitempty" yaml:"clientCredentials,omitempty"`
	AuthorizationCodeDiff *OAuthFlowDiff  `json:"authorizationCode,omitempty" yaml:"authorizationCode,omitempty"`
}

// Empty return true if there is no diff
func (diff *OAuthFlowsDiff) Empty() bool {
	return diff == nil || *diff == OAuthFlowsDiff{}
}

func getOAuthFlowsDiff(config *Config, flows1, flows2 *openapi3.OAuthFlows) *OAuthFlowsDiff {
	diff := getOAuthFlowsDiffInternal(config, flows1, flows2)
	if diff.Empty() {
		return nil
	}
	return diff
}

func getOAuthFlowsDiffInternal(config *Config, flows1, flows2 *openapi3.OAuthFlows) *OAuthFlowsDiff {

	if flows1 == nil && flows2 == nil {
		return nil
	}

	result := OAuthFlowsDiff{}

	result.ExtensionsDiff = getExtensionsDiff(config, flows1.ExtensionProps, flows2.ExtensionProps)
	result.ImplicitDiff = getOAuthFlowDiff(config, flows1.Implicit, flows2.Implicit)
	result.PasswordDiff = getOAuthFlowDiff(config, flows1.Password, flows2.Password)
	result.ClientCredentialsDiff = getOAuthFlowDiff(config, flows1.ClientCredentials, flows2.ClientCredentials)
	result.AuthorizationCodeDiff = getOAuthFlowDiff(config, flows1.AuthorizationCode, flows2.AuthorizationCode)

	return &result
}
