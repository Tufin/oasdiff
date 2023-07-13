package checker_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/checker"
	"github.com/tufin/oasdiff/diff"
)

// CL: Adding a new global security to the API
func TestAPIGlobalSecurityyAdded(t *testing.T) {
	s1, err := open("../data/checker/api_security_global_added_base.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/api_security_global_added_revision.yaml")
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.APISecurityUpdatedCheck), d, osm, checker.INFO)
	require.Len(t, errs, 1)
	require.Equal(t, checker.SecurityChange{
		Id:      "api-global-security-added",
		Text:    "the security scheme 'petstore_auth' was added to the API",
		Comment: "",
		Level:   checker.INFO,
		Source:  "",
	}, errs[0])
}

// CL: Removing a global security from the API
func TestAPIGlobalSecurityyDeleted(t *testing.T) {
	s1, err := open("../data/checker/api_security_global_added_revision.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/api_security_global_added_base.yaml")
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.APISecurityUpdatedCheck), d, osm, checker.INFO)
	require.Len(t, errs, 1)
	require.Equal(t, checker.SecurityChange{
		Id:      "api-global-security-removed",
		Text:    "the security scheme 'petstore_auth' was removed from the API",
		Comment: "",
		Level:   checker.INFO,
		Source:  "",
	}, errs[0])
}

// CL: Removing a security scope from an API global security
func TestAPIGlobalSecurityScopeRemoved(t *testing.T) {
	s1, err := open("../data/checker/api_security_global_added_revision.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/api_security_global_added_revision.yaml")
	require.NoError(t, err)

	s2.Spec.Security[0]["petstore_auth"] = s2.Spec.Security[0]["petstore_auth"][:1]
	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.APISecurityUpdatedCheck), d, osm, checker.INFO)
	require.Len(t, errs, 1)
	require.Equal(t, checker.SecurityChange{
		Id:      "api-global-security-scope-removed",
		Text:    "the security scope 'read:pets' was removed from the global security scheme 'petstore_auth'",
		Comment: "",
		Level:   checker.INFO,
		Source:  "",
	}, errs[0])
}

// CL: Adding a security scope from an API global security
func TestAPIGlobalSecurityScopeAdded(t *testing.T) {
	s1, err := open("../data/checker/api_security_global_added_revision.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/api_security_global_added_revision.yaml")
	require.NoError(t, err)

	s1.Spec.Security[0]["petstore_auth"] = s2.Spec.Security[0]["petstore_auth"][:1]
	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.APISecurityUpdatedCheck), d, osm, checker.INFO)
	require.Len(t, errs, 1)
	require.Equal(t, checker.SecurityChange{
		Id:      "api-global-security-scope-added",
		Text:    "the security scope 'read:pets' was added to the global security scheme 'petstore_auth'",
		Comment: "",
		Level:   checker.INFO,
		Source:  "",
	}, errs[0])
}

// CL: Adding a new security to the API endpoint
func TestAPISecurityAdded(t *testing.T) {
	s1, err := open("../data/checker/api_security_added_base.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/api_security_added_revision.yaml")
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.APISecurityUpdatedCheck), d, osm, checker.INFO)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{
		Id:          "api-security-added",
		Text:        "the endpoint scheme security 'petstore_auth' was added to the API",
		Comment:     "",
		Level:       checker.INFO,
		Operation:   "POST",
		Path:        "/subscribe",
		Source:      "../data/checker/api_security_added_revision.yaml",
		OperationId: "",
	}, errs[0])
}

// CL: Removing a new security to the API endpoint
func TestAPISecurityDeleted(t *testing.T) {
	s1, err := open("../data/checker/api_security_added_revision.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/api_security_added_base.yaml")
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.APISecurityUpdatedCheck), d, osm, checker.INFO)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{
		Id:          "api-security-removed",
		Text:        "the endpoint scheme security 'petstore_auth' was removed from the API",
		Comment:     "",
		Level:       checker.INFO,
		Operation:   "POST",
		Path:        "/subscribe",
		Source:      "../data/checker/api_security_added_base.yaml",
		OperationId: "",
	}, errs[0])
}

// CL: Removing a security scope from an API endpoint security
func TestAPISecurityScopeRemoved(t *testing.T) {
	s1, err := open("../data/checker/api_security_updated_base.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/api_security_updated_revision.yaml")
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.APISecurityUpdatedCheck), d, osm, checker.INFO)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{
		Id:          "api-security-scope-removed",
		Text:        "the security scope 'read:pets' was removed from the endpoint's security scheme 'petstore_auth'",
		Comment:     "",
		Level:       checker.INFO,
		Operation:   "POST",
		Path:        "/subscribe",
		Source:      "../data/checker/api_security_updated_revision.yaml",
		OperationId: "",
	}, errs[0])
}

// CL: Adding a security scope to an API endpoint security
func TestAPISecurityScopeAdded(t *testing.T) {
	s1, err := open("../data/checker/api_security_updated_revision.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/api_security_updated_base.yaml")
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.APISecurityUpdatedCheck), d, osm, checker.INFO)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{
		Id:          "api-security-scope-added",
		Text:        "the security scope 'read:pets' was added to the endpoint's security scheme 'petstore_auth'",
		Comment:     "",
		Level:       checker.INFO,
		Operation:   "POST",
		Path:        "/subscribe",
		Source:      "../data/checker/api_security_updated_base.yaml",
		OperationId: "",
	}, errs[0])
}
