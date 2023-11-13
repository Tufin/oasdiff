package checker_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/checker"
	"github.com/tufin/oasdiff/diff"
)

// CL: adding a new required request property
func TestRequiredRequestPropertyAdded(t *testing.T) {
	s1, err := open("../data/checker/request_property_added_base.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/request_property_added_revision.yaml")
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.RequestPropertyUpdatedCheck), d, osm, checker.INFO)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{
		Id:          checker.NewRequiredRequestPropertyId,
		Text:        "added the new required request property 'description'",
		Level:       checker.ERR,
		Operation:   "POST",
		Path:        "/products",
		Source:      "../data/checker/request_property_added_revision.yaml",
		OperationId: "addProduct",
	}, errs[0])
}

// CL: adding two new request properties, one required, one optional
func TestRequiredRequestPropertiesAdded(t *testing.T) {
	s1, err := open("../data/checker/request_property_added_base.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/request_property_added_revision2.yaml")
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.RequestPropertyUpdatedCheck), d, osm, checker.INFO)
	require.ElementsMatch(t, []checker.ApiChange{
		{
			Id:          checker.NewRequiredRequestPropertyId,
			Text:        "added the new required request property 'description'",
			Level:       checker.ERR,
			Operation:   "POST",
			Path:        "/products",
			Source:      "../data/checker/request_property_added_revision2.yaml",
			OperationId: "addProduct",
		},
		{
			Id:          checker.NewOptionalRequestPropertyId,
			Text:        "added the new optional request property 'info'",
			Level:       checker.INFO,
			Operation:   "POST",
			Path:        "/products",
			Source:      "../data/checker/request_property_added_revision2.yaml",
			OperationId: "addProduct",
		}}, errs)
}

// CL: adding a new optional request property
func TestRequiredOptionalPropertyAdded(t *testing.T) {
	s1, err := open("../data/checker/request_property_added_base.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/request_property_added_revision.yaml")
	require.NoError(t, err)

	s2.Spec.Paths["/products"].Post.RequestBody.Value.Content["application/json"].Schema.Value.Required = []string{"name"}
	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.RequestPropertyUpdatedCheck), d, osm, checker.INFO)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{
		Id:          checker.NewOptionalRequestPropertyId,
		Text:        "added the new optional request property 'description'",
		Level:       checker.INFO,
		Operation:   "POST",
		Path:        "/products",
		Source:      "../data/checker/request_property_added_revision.yaml",
		OperationId: "addProduct",
	}, errs[0])
}

// CL: removing a required request property
func TestRequiredRequestPropertyRemoved(t *testing.T) {
	s1, err := open("../data/checker/request_property_added_revision.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/request_property_added_base.yaml")
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.RequestPropertyUpdatedCheck), d, osm, checker.INFO)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{
		Id:          checker.RequestPropertyRemovedId,
		Text:        "removed the request property 'description'",
		Level:       checker.WARN,
		Operation:   "POST",
		Path:        "/products",
		Source:      "../data/checker/request_property_added_base.yaml",
		OperationId: "addProduct",
	}, errs[0])
}
