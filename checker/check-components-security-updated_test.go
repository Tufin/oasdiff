package checker_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/checker"
	"github.com/tufin/oasdiff/diff"
)

// CL: Changing security component oauth's url
func TestComponentSecurityOauthURLUpdated(t *testing.T) {
	s1, _ := open("../data/checker/component_security_updated_base.yaml")
	s2, err := open("../data/checker/component_security_updated_base.yaml")
	require.Empty(t, err)

	s2.Spec.Components.SecuritySchemes["petstore_auth"].Value.Flows.Implicit.AuthorizationURL = "http://example.new.org/api/oauth/dialog"

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.APIComponentsSecurityUpdatedCheck), d, osm, checker.INFO)
	require.NotEmpty(t, errs)
	require.Equal(t, checker.BackwardCompatibilityErrors{
		{
			Id:          "api-security-component-oauth-url-changed",
			Text:        "the component security schema 'petstore_auth' oauth url changed from 'http://example.org/api/oauth/dialog' to 'http://example.new.org/api/oauth/dialog'",
			Comment:     "",
			Level:       checker.INFO,
			Operation:   "N/A",
			Path:        "N/A",
			Source:      "N/A",
			OperationId: "N/A",
		}}, errs)
}

// CL: Changing security component type
func TestComponentSecurityTypeUpdated(t *testing.T) {
	s1, _ := open("../data/checker/component_security_updated_base.yaml")
	s2, err := open("../data/checker/component_security_updated_base.yaml")
	require.Empty(t, err)

	s2.Spec.Components.SecuritySchemes["petstore_auth"].Value.Type = "http"

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.APIComponentsSecurityUpdatedCheck), d, osm, checker.INFO)
	require.NotEmpty(t, errs)
	require.Equal(t, checker.BackwardCompatibilityErrors{
		{
			Id:          "api-security-component-type-changed",
			Text:        "the component security schema 'petstore_auth' type changed from 'oauth2' to 'http'",
			Comment:     "",
			Level:       checker.INFO,
			Operation:   "N/A",
			Path:        "N/A",
			Source:      "N/A",
			OperationId: "N/A",
		}}, errs)
}

// CL: Adding a new security component
func TestComponentSecurityAdded(t *testing.T) {
	s1, _ := open("../data/checker/component_security_updated_base.yaml")
	s2, err := open("../data/checker/component_security_updated_revision.yaml")
	require.Empty(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.APIComponentsSecurityUpdatedCheck), d, osm, checker.INFO)
	require.NotEmpty(t, errs)
	require.Equal(t, checker.BackwardCompatibilityErrors{
		{
			Id:          "api-security-component-added",
			Text:        "the component security schema 'BasicAuth' was added",
			Comment:     "",
			Level:       checker.INFO,
			Operation:   "N/A",
			Path:        "N/A",
			Source:      "N/A",
			OperationId: "N/A",
		}}, errs)
}

// CL: Removing a new security component
func TestComponentSecurityRemoved(t *testing.T) {
	s1, _ := open("../data/checker/component_security_updated_revision.yaml")
	s2, err := open("../data/checker/component_security_updated_base.yaml")
	require.Empty(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.APIComponentsSecurityUpdatedCheck), d, osm, checker.INFO)
	require.NotEmpty(t, errs)
	require.Equal(t, checker.BackwardCompatibilityErrors{
		{
			Id:          "api-security-component-removed",
			Text:        "the component security schema 'BasicAuth' was removed",
			Comment:     "",
			Level:       checker.INFO,
			Operation:   "N/A",
			Path:        "N/A",
			Source:      "N/A",
			OperationId: "N/A",
		}}, errs)
}

// CL: Adding a new oauth security scope
func TestComponentSecurityOauthScopeAdded(t *testing.T) {
	s1, _ := open("../data/checker/component_security_updated_base.yaml")
	s2, err := open("../data/checker/component_security_updated_base.yaml")
	require.Empty(t, err)

	s2.Spec.Components.SecuritySchemes["petstore_auth"].Value.Flows.Implicit.Scopes["admin:pets"] = "grants access to admin operations"

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.APIComponentsSecurityUpdatedCheck), d, osm, checker.INFO)
	require.NotEmpty(t, errs)
	require.Equal(t, checker.BackwardCompatibilityErrors{
		{
			Id:          "api-security-component-oauth-scope-added",
			Text:        "the component security schema 'petstore_auth' oauth scope 'admin:pets' was added",
			Comment:     "",
			Level:       checker.INFO,
			Operation:   "N/A",
			Path:        "N/A",
			Source:      "N/A",
			OperationId: "N/A",
		}}, errs)
}

// CL: Removing a new oauth security scope
func TestComponentSecurityOauthScopeRemoved(t *testing.T) {
	s1, _ := open("../data/checker/component_security_updated_base.yaml")
	s2, err := open("../data/checker/component_security_updated_base.yaml")
	require.Empty(t, err)

	// Add to s1 so that it's deletion is identified
	s1.Spec.Components.SecuritySchemes["petstore_auth"].Value.Flows.Implicit.Scopes["admin:pets"] = "grants access to admin operations"

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.APIComponentsSecurityUpdatedCheck), d, osm, checker.INFO)
	require.NotEmpty(t, errs)
	require.Equal(t, checker.BackwardCompatibilityErrors{
		{
			Id:          "api-security-component-oauth-scope-removed",
			Text:        "the component security schema 'petstore_auth' oauth scope 'admin:pets' was removed",
			Comment:     "",
			Level:       checker.INFO,
			Operation:   "N/A",
			Path:        "N/A",
			Source:      "N/A",
			OperationId: "N/A",
		}}, errs)
}

// CL: Removing a new oauth security scope
func TestComponentSecurityOauthScopeUpdated(t *testing.T) {
	s1, _ := open("../data/checker/component_security_updated_base.yaml")
	s2, err := open("../data/checker/component_security_updated_base.yaml")
	require.Empty(t, err)

	s2.Spec.Components.SecuritySchemes["petstore_auth"].Value.Flows.Implicit.Scopes["read:pets"] = "grants access to pets (deprecated)"

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.APIComponentsSecurityUpdatedCheck), d, osm, checker.INFO)
	require.NotEmpty(t, errs)
	require.Equal(t, checker.BackwardCompatibilityErrors{
		{
			Id:          "api-security-component-oauth-scope-changed",
			Text:        "the component security schema 'petstore_auth' oauth scope 'read:pets' was updated from 'read your pets' to 'grants access to pets (deprecated)'",
			Comment:     "",
			Level:       checker.INFO,
			Operation:   "N/A",
			Path:        "N/A",
			Source:      "N/A",
			OperationId: "N/A",
		}}, errs)
}
