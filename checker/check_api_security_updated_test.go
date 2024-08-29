package checker_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/checker"
	"github.com/tufin/oasdiff/diff"
	"github.com/tufin/oasdiff/load"
)

// CL: adding a new global security to the API
func TestAPIGlobalSecurityyAdded(t *testing.T) {
	s1, err := open("../data/checker/api_security_global_added_base.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/api_security_global_added_revision.yaml")
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.APISecurityUpdatedCheck), d, osm, checker.INFO)
	require.Len(t, errs, 1)
	require.Equal(t, checker.SecurityChange{
		Id:    checker.APIGlobalSecurityAddedCheckId,
		Args:  []any{"petstore_auth"},
		Level: checker.INFO,
	}, errs[0])
	require.Equal(t, "security scheme 'petstore_auth' was added", errs[0].GetUncolorizedText(checker.NewDefaultLocalizer()))
}

// CL: removing a global security from the API
func TestAPIGlobalSecurityyDeleted(t *testing.T) {
	s1, err := open("../data/checker/api_security_global_added_revision.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/api_security_global_added_base.yaml")
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.APISecurityUpdatedCheck), d, osm, checker.INFO)
	require.Len(t, errs, 1)
	require.Equal(t, checker.SecurityChange{
		Id:    checker.APIGlobalSecurityRemovedCheckId,
		Args:  []any{"petstore_auth"},
		Level: checker.INFO,
	}, errs[0])
	require.Equal(t, "removed security scheme 'petstore_auth'", errs[0].GetUncolorizedText(checker.NewDefaultLocalizer()))
}

// CL: removing a security scope from an API global security
func TestAPIGlobalSecurityScopeRemoved(t *testing.T) {
	s1, err := open("../data/checker/api_security_global_added_revision.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/api_security_global_added_revision.yaml")
	require.NoError(t, err)

	s2.Spec.Security[0]["petstore_auth"] = s2.Spec.Security[0]["petstore_auth"][:1]
	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.APISecurityUpdatedCheck), d, osm, checker.INFO)
	require.Len(t, errs, 1)
	require.Equal(t, checker.SecurityChange{
		Id:    checker.APIGlobalSecurityScopeRemovedId,
		Args:  []any{"read:pets", "petstore_auth"},
		Level: checker.INFO,
	}, errs[0])
	require.Equal(t, "removed security scope 'read:pets' from global security scheme 'petstore_auth'", errs[0].GetUncolorizedText(checker.NewDefaultLocalizer()))
}

// CL: adding a security scope from an API global security
func TestAPIGlobalSecurityScopeAdded(t *testing.T) {
	s1, err := open("../data/checker/api_security_global_added_revision.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/api_security_global_added_revision.yaml")
	require.NoError(t, err)

	s1.Spec.Security[0]["petstore_auth"] = s2.Spec.Security[0]["petstore_auth"][:1]
	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.APISecurityUpdatedCheck), d, osm, checker.INFO)
	require.Len(t, errs, 1)
	require.Equal(t, checker.SecurityChange{
		Id:    checker.APIGlobalSecurityScopeAddedId,
		Args:  []any{"read:pets", "petstore_auth"},
		Level: checker.INFO,
	}, errs[0])
	require.Equal(t, "security scope 'read:pets' was added to global security scheme 'petstore_auth'", errs[0].GetUncolorizedText(checker.NewDefaultLocalizer()))
}

// CL: adding a new security to the API endpoint
func TestAPISecurityAdded(t *testing.T) {
	s1, err := open("../data/checker/api_security_added_base.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/api_security_added_revision.yaml")
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.APISecurityUpdatedCheck), d, osm, checker.INFO)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{
		Id:        checker.APISecurityAddedCheckId,
		Args:      []any{"petstore_auth"},
		Level:     checker.INFO,
		Operation: "POST",
		Path:      "/subscribe",
		Source:    load.NewSource("../data/checker/api_security_added_revision.yaml"),
	}, errs[0])
	require.Equal(t, "endpoint scheme security 'petstore_auth' was added", errs[0].GetUncolorizedText(checker.NewDefaultLocalizer()))
}

// CL: removing a new security to the API endpoint
func TestAPISecurityDeleted(t *testing.T) {
	s1, err := open("../data/checker/api_security_added_revision.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/api_security_added_base.yaml")
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.APISecurityUpdatedCheck), d, osm, checker.INFO)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{
		Id:        checker.APISecurityRemovedCheckId,
		Args:      []any{"petstore_auth"},
		Level:     checker.INFO,
		Operation: "POST",
		Path:      "/subscribe",
		Source:    load.NewSource("../data/checker/api_security_added_base.yaml"),
	}, errs[0])
	require.Equal(t, "removed endpoint scheme security 'petstore_auth'", errs[0].GetUncolorizedText(checker.NewDefaultLocalizer()))
}

// CL: removing a security scope from an API endpoint security
func TestAPISecurityScopeRemoved(t *testing.T) {
	s1, err := open("../data/checker/api_security_updated_base.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/api_security_updated_revision.yaml")
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.APISecurityUpdatedCheck), d, osm, checker.INFO)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{
		Id:        checker.APISecurityScopeRemovedId,
		Args:      []any{"read:pets", "petstore_auth"},
		Level:     checker.INFO,
		Operation: "POST",
		Path:      "/subscribe",
		Source:    load.NewSource("../data/checker/api_security_updated_revision.yaml"),
	}, errs[0])
	require.Equal(t, "removed security scope 'read:pets' from endpoint security scheme 'petstore_auth'", errs[0].GetUncolorizedText(checker.NewDefaultLocalizer()))
}

// CL: adding a security scope to an API endpoint security
func TestAPISecurityScopeAdded(t *testing.T) {
	s1, err := open("../data/checker/api_security_updated_revision.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/api_security_updated_base.yaml")
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.APISecurityUpdatedCheck), d, osm, checker.INFO)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{
		Id:        checker.APISecurityScopeAddedId,
		Args:      []any{"read:pets", "petstore_auth"},
		Level:     checker.INFO,
		Operation: "POST",
		Path:      "/subscribe",
		Source:    load.NewSource("../data/checker/api_security_updated_base.yaml"),
	}, errs[0])
	require.Equal(t, "security scope 'read:pets' was added to endpoint security scheme 'petstore_auth'", errs[0].GetUncolorizedText(checker.NewDefaultLocalizer()))
}
