package checker_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/checker"
	"github.com/tufin/oasdiff/diff"
	"github.com/tufin/oasdiff/load"
)

// CL: adding a new required request property
func TestRequiredRequestPropertyAdded(t *testing.T) {
	s1, err := open("../data/checker/request_property_added_base.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/request_property_added_revision.yaml")
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.RequestPropertyUpdatedCheck), d, osm, checker.INFO)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{
		Id:          checker.NewRequiredRequestPropertyId,
		Args:        []any{"description", "application/json"},
		Level:       checker.ERR,
		Operation:   "POST",
		Path:        "/products",
		Source:      load.NewSource("../data/checker/request_property_added_revision.yaml"),
		OperationId: "addProduct",
	}, errs[0])
	require.Equal(t, "added required request property 'description' to media-type 'application/json'", errs[0].GetUncolorizedText(checker.NewDefaultLocalizer()))
}

// CL: adding two new request properties, one required, one optional
func TestRequiredRequestPropertiesAdded(t *testing.T) {
	s1, err := open("../data/checker/request_property_added_base.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/request_property_added_revision2.yaml")
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.RequestPropertyUpdatedCheck), d, osm, checker.INFO)
	require.Len(t, errs, 2)

	require.Equal(t, errs[0], checker.ApiChange{
		Id:          checker.NewRequiredRequestPropertyId,
		Args:        []any{"description", "application/json"},
		Level:       checker.ERR,
		Operation:   "POST",
		Path:        "/products",
		Source:      load.NewSource("../data/checker/request_property_added_revision2.yaml"),
		OperationId: "addProduct",
	})
	require.Equal(t, "added required request property 'description' to media-type 'application/json'", errs[0].GetUncolorizedText(checker.NewDefaultLocalizer()))

	require.Equal(t, errs[1], checker.ApiChange{
		Id:          checker.NewOptionalRequestPropertyId,
		Args:        []any{"info", "application/json"},
		Level:       checker.INFO,
		Operation:   "POST",
		Path:        "/products",
		Source:      load.NewSource("../data/checker/request_property_added_revision2.yaml"),
		OperationId: "addProduct",
	})
	require.Equal(t, "added optional request property 'info' to media-type 'application/json'", errs[1].GetUncolorizedText(checker.NewDefaultLocalizer()))
}

// CL: adding a new optional request property
func TestRequiredOptionalPropertyAdded(t *testing.T) {
	s1, err := open("../data/checker/request_property_added_base.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/request_property_added_revision.yaml")
	require.NoError(t, err)

	s2.Spec.Paths.Value("/products").Post.RequestBody.Value.Content["application/json"].Schema.Value.Required = []string{"name"}
	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.RequestPropertyUpdatedCheck), d, osm, checker.INFO)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{
		Id:          checker.NewOptionalRequestPropertyId,
		Args:        []any{"description", "application/json"},
		Level:       checker.INFO,
		Operation:   "POST",
		Path:        "/products",
		Source:      load.NewSource("../data/checker/request_property_added_revision.yaml"),
		OperationId: "addProduct",
	}, errs[0])
	require.Equal(t, "added optional request property 'description' to media-type 'application/json'", errs[0].GetUncolorizedText(checker.NewDefaultLocalizer()))
}

// CL: removing a required request property
func TestRequiredRequestPropertyRemoved(t *testing.T) {
	s1, err := open("../data/checker/request_property_added_revision.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/request_property_added_base.yaml")
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.RequestPropertyUpdatedCheck), d, osm, checker.INFO)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{
		Id:          checker.RequestPropertyRemovedId,
		Args:        []any{"description", "application/json"},
		Level:       checker.WARN,
		Operation:   "POST",
		Path:        "/products",
		Source:      load.NewSource("../data/checker/request_property_added_base.yaml"),
		OperationId: "addProduct",
	}, errs[0])
	require.Equal(t, "removed request property 'description' of media-type 'application/json'", errs[0].GetUncolorizedText(checker.NewDefaultLocalizer()))
}

// CL: adding a new required request property with a default value
func TestRequiredRequestPropertyAddedWithDefault(t *testing.T) {
	s1, err := open("../data/checker/request_property_added_base.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/request_property_added_with_default.yaml")
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.RequestPropertyUpdatedCheck), d, osm, checker.INFO)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{
		Id:          checker.NewRequiredRequestPropertyWithDefaultId,
		Args:        []any{"description", "application/json"},
		Level:       checker.INFO,
		Operation:   "POST",
		Path:        "/products",
		Source:      load.NewSource("../data/checker/request_property_added_with_default.yaml"),
		OperationId: "addProduct",
	}, errs[0])
	require.Equal(t, "added required request property 'description' with a default value to media-type 'application/json'", errs[0].GetUncolorizedText(checker.NewDefaultLocalizer()))
}
