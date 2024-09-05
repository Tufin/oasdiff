package checker_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/checker"
	"github.com/tufin/oasdiff/diff"
	"github.com/tufin/oasdiff/load"
)

// CL: changing request property required value to true
func TestRequestPropertyMarkedRequired(t *testing.T) {
	s1, err := open("../data/checker/request_property_became_required_base.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/request_property_became_required_base.yaml")
	require.NoError(t, err)

	s1.Spec.Paths.Value("/products").Post.RequestBody.Value.Content["application/json"].Schema.Value.Required = []string{""}
	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.RequestPropertyRequiredUpdatedCheck), d, osm, checker.INFO)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{
		Id:          checker.RequestPropertyBecameRequiredId,
		Args:        []any{"name", "application/json"},
		Level:       checker.ERR,
		Operation:   "POST",
		Path:        "/products",
		Source:      load.NewSource("../data/checker/request_property_became_required_base.yaml"),
		OperationId: "addProduct",
	}, errs[0])
	require.Equal(t, "request property 'name' of media-type 'application/json' became required", errs[0].GetUncolorizedText(checker.NewDefaultLocalizer()))
}

// CL: changing request property required value to false
func TestRequestPropertyMarkedOptional(t *testing.T) {
	s1, err := open("../data/checker/request_property_became_required_base.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/request_property_became_required_base.yaml")
	require.NoError(t, err)

	s2.Spec.Paths.Value("/products").Post.RequestBody.Value.Content["application/json"].Schema.Value.Required = []string{""}
	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.RequestPropertyRequiredUpdatedCheck), d, osm, checker.INFO)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{
		Id:          checker.RequestPropertyBecameOptionalId,
		Args:        []any{"name", "application/json"},
		Level:       checker.INFO,
		Operation:   "POST",
		Path:        "/products",
		Source:      load.NewSource("../data/checker/request_property_became_required_base.yaml"),
		OperationId: "addProduct",
	}, errs[0])
	require.Equal(t, "request property 'name' of media-type 'application/json' became optional", errs[0].GetUncolorizedText(checker.NewDefaultLocalizer()))
}

// CL: making request property required, while also giving it a default value
func TestRequestPropertyWithDefaultMarkedRequired(t *testing.T) {
	s1, err := open("../data/checker/request_property_became_required_base.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/request_property_became_required_base.yaml")
	require.NoError(t, err)

	s1.Spec.Paths.Value("/products").Post.RequestBody.Value.Content["application/json"].Schema.Value.Required = []string{""}
	s2.Spec.Paths.Value("/products").Post.RequestBody.Value.Content["application/json"].Schema.Value.Properties["name"].Value.Default = "default"
	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.RequestPropertyRequiredUpdatedCheck), d, osm, checker.INFO)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{
		Id:          checker.RequestPropertyBecameRequiredWithDefaultId,
		Args:        []any{"name", "application/json"},
		Level:       checker.INFO,
		Operation:   "POST",
		Path:        "/products",
		Source:      load.NewSource("../data/checker/request_property_became_required_base.yaml"),
		OperationId: "addProduct",
	}, errs[0])
	require.Equal(t, "request property 'name' with a default value became required in media-type 'application/json'", errs[0].GetUncolorizedText(checker.NewDefaultLocalizer()))
}
