package diff

import (
	"github.com/getkin/kin-openapi/openapi3"
)

// OAuthFlowsDiff describes the changes between a pair of oauth flows objects: https://swagger.io/specification/#oauth-flows-object
type OAuthFlowsDiff struct {
	Added                 bool            `json:"added,omitempty" yaml:"added,omitempty"`
	Deleted               bool            `json:"deleted,omitempty" yaml:"deleted,omitempty"`
	ExtensionsDiff        *ExtensionsDiff `json:"extensions,omitempty" yaml:"extensions,omitempty"`
	ImplicitDiff          *OAuthFlowDiff  `json:"implicit,omitempty" yaml:"implicit,omitempty"`
	PasswordDiff          *OAuthFlowDiff  `json:"password,omitempty" yaml:"password,omitempty"`
	ClientCredentialsDiff *OAuthFlowDiff  `json:"clientCredentials,omitempty" yaml:"clientCredentials,omitempty"`
	AuthorizationCodeDiff *OAuthFlowDiff  `json:"authorizationCode,omitempty" yaml:"authorizationCode,omitempty"`
}

// Empty indicates whether a change was found in this element
func (diff *OAuthFlowsDiff) Empty() bool {
	return diff == nil || *diff == OAuthFlowsDiff{}
}

func getOAuthFlowsDiff(config *Config, state *state, flows1, flows2 *openapi3.OAuthFlows) *OAuthFlowsDiff {
	diff := getOAuthFlowsDiffInternal(config, state, flows1, flows2)

	if diff.Empty() {
		return nil
	}

	return diff
}

func getOAuthFlowsDiffInternal(config *Config, state *state, flows1, flows2 *openapi3.OAuthFlows) *OAuthFlowsDiff {

	if flows1 == nil && flows2 == nil {
		return nil
	}

	if flows1 == nil && flows2 != nil {
		return &OAuthFlowsDiff{
			Added: true,
		}
	}

	if flows1 != nil && flows2 == nil {
		return &OAuthFlowsDiff{
			Deleted: true,
		}
	}

	result := OAuthFlowsDiff{}

	result.ExtensionsDiff = getExtensionsDiff(config, state, flows1.Extensions, flows2.Extensions)
	result.ImplicitDiff = getOAuthFlowDiff(config, state, flows1.Implicit, flows2.Implicit)
	result.PasswordDiff = getOAuthFlowDiff(config, state, flows1.Password, flows2.Password)
	result.ClientCredentialsDiff = getOAuthFlowDiff(config, state, flows1.ClientCredentials, flows2.ClientCredentials)
	result.AuthorizationCodeDiff = getOAuthFlowDiff(config, state, flows1.AuthorizationCode, flows2.AuthorizationCode)

	return &result
}
