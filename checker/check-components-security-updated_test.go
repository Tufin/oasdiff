package checker_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/checker"
	"github.com/tufin/oasdiff/diff"
)

// CL: changing security component oauth's url
func TestComponentSecurityOauthURLUpdated(t *testing.T) {
	s1, err := open("../data/checker/component_security_updated_base.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/component_security_updated_base.yaml")
	require.NoError(t, err)

	s2.Spec.Components.SecuritySchemes["petstore_auth"].Value.Flows.Implicit.AuthorizationURL = "http://example.new.org/api/oauth/dialog"

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.APIComponentsSecurityUpdatedCheck), d, osm, checker.INFO)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ComponentChange{
		Id:      checker.APIComponentsSecurityComponentOauthUrlUpdatedId,
		Text:    "the component security scheme 'petstore_auth' oauth url changed from 'http://example.org/api/oauth/dialog' to 'http://example.new.org/api/oauth/dialog'",
		Comment: "",
		Level:   checker.INFO,
		Source:  "",
	}, errs[0])
}

// CL: changing security component type
func TestComponentSecurityTypeUpdated(t *testing.T) {
	s1, err := open("../data/checker/component_security_updated_base.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/component_security_updated_base.yaml")
	require.NoError(t, err)

	s2.Spec.Components.SecuritySchemes["petstore_auth"].Value.Type = "http"

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.APIComponentsSecurityUpdatedCheck), d, osm, checker.INFO)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ComponentChange{
		Id:      checker.APIComponentsSecurityTypeUpdatedId,
		Text:    "the component security scheme 'petstore_auth' type changed from 'oauth2' to 'http'",
		Comment: "",
		Level:   checker.INFO,
		Source:  "",
	}, errs[0])
}

// CL: adding a new security component
func TestComponentSecurityAdded(t *testing.T) {
	s1, err := open("../data/checker/component_security_updated_base.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/component_security_updated_revision.yaml")
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.APIComponentsSecurityUpdatedCheck), d, osm, checker.INFO)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ComponentChange{
		Id:      checker.APIComponentsSecurityAddedId,
		Text:    "the component security scheme 'BasicAuth' was added",
		Comment: "",
		Level:   checker.INFO,
		Source:  "",
	}, errs[0])
}

// CL: removing a new security component
func TestComponentSecurityRemoved(t *testing.T) {
	s1, err := open("../data/checker/component_security_updated_revision.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/component_security_updated_base.yaml")
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.APIComponentsSecurityUpdatedCheck), d, osm, checker.INFO)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ComponentChange{
		Id:      checker.APIComponentsSecurityRemovedId,
		Text:    "the component security scheme 'BasicAuth' was removed",
		Comment: "",
		Level:   checker.INFO,
		Source:  "",
	}, errs[0])
}

// CL: adding a new oauth security scope
func TestComponentSecurityOauthScopeAdded(t *testing.T) {
	s1, err := open("../data/checker/component_security_updated_base.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/component_security_updated_base.yaml")
	require.NoError(t, err)

	s2.Spec.Components.SecuritySchemes["petstore_auth"].Value.Flows.Implicit.Scopes["admin:pets"] = "grants access to admin operations"

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.APIComponentsSecurityUpdatedCheck), d, osm, checker.INFO)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ComponentChange{
		Id:      checker.APIComponentSecurityOauthScopeAddedId,
		Text:    "the component security scheme 'petstore_auth' oauth scope 'admin:pets' was added",
		Comment: "",
		Level:   checker.INFO,
		Source:  "",
	}, errs[0])
}

// CL: removing a new oauth security scope
func TestComponentSecurityOauthScopeRemoved(t *testing.T) {
	s1, err := open("../data/checker/component_security_updated_base.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/component_security_updated_base.yaml")
	require.NoError(t, err)

	// Add to s1 so that it's deletion is identified
	s1.Spec.Components.SecuritySchemes["petstore_auth"].Value.Flows.Implicit.Scopes["admin:pets"] = "grants access to admin operations"

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.APIComponentsSecurityUpdatedCheck), d, osm, checker.INFO)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ComponentChange{
		Id:      "api-security-component-oauth-scope-removed",
		Text:    "the component security scheme 'petstore_auth' oauth scope 'admin:pets' was removed",
		Comment: "",
		Level:   checker.INFO,
		Source:  "",
	}, errs[0])
}

// CL: removing a new oauth security scope
func TestComponentSecurityOauthScopeUpdated(t *testing.T) {
	s1, err := open("../data/checker/component_security_updated_base.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/component_security_updated_base.yaml")
	require.NoError(t, err)

	s2.Spec.Components.SecuritySchemes["petstore_auth"].Value.Flows.Implicit.Scopes["read:pets"] = "grants access to pets (deprecated)"

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.APIComponentsSecurityUpdatedCheck), d, osm, checker.INFO)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ComponentChange{
		Id:      "api-security-component-oauth-scope-changed",
		Text:    "the component security scheme 'petstore_auth' oauth scope 'read:pets' was updated from 'read your pets' to 'grants access to pets (deprecated)'",
		Comment: "",
		Level:   checker.INFO,
		Source:  "",
	}, errs[0])
}
