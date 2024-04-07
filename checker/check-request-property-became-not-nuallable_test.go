package checker_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/checker"
	"github.com/tufin/oasdiff/diff"
	"github.com/tufin/oasdiff/load"
)

// CL: changing request property to not nullable
func TestRequestPropertyBecameNotNullable(t *testing.T) {
	s1, err := open("../data/checker/request_property_became_nullable_revision.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/request_property_became_nullable_base.yaml")
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.RequestPropertyBecameNotNullableCheck), d, osm, checker.INFO)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{
		Id:          checker.RequestPropertyBecomeNotNullableId,
		Args:        []any{"name"},
		Level:       checker.ERR,
		Operation:   "POST",
		Path:        "/products",
		Source:      load.NewSource("../data/checker/request_property_became_nullable_base.yaml"),
		OperationId: "addProduct",
	}, errs[0])
}

// CL: changing request property to nullable
func TestRequestPropertyBecameNullable(t *testing.T) {
	s1, err := open("../data/checker/request_property_became_nullable_base.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/request_property_became_nullable_revision.yaml")
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.RequestPropertyBecameNotNullableCheck), d, osm, checker.INFO)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{
		Id:          checker.RequestPropertyBecomeNullableId,
		Args:        []any{"name"},
		Level:       checker.INFO,
		Operation:   "POST",
		Path:        "/products",
		Source:      load.NewSource("../data/checker/request_property_became_nullable_revision.yaml"),
		OperationId: "addProduct",
	}, errs[0])

}

// CL: changing request body to nullable
func TestRequestBodyBecameNullable(t *testing.T) {
	s1, err := open("../data/checker/request_property_became_nullable_base.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/request_property_became_nullable_base.yaml")
	require.NoError(t, err)

	s2.Spec.Paths.Value("/products").Post.RequestBody.Value.Content["application/json"].Schema.Value.Nullable = true

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)

	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.RequestPropertyBecameNotNullableCheck), d, osm, checker.INFO)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{
		Id:          checker.RequestBodyBecomeNullableId,
		Level:       checker.INFO,
		Operation:   "POST",
		Path:        "/products",
		Source:      load.NewSource("../data/checker/request_property_became_nullable_base.yaml"),
		OperationId: "addProduct",
	}, errs[0])
}

// CL: changing request body to not nullable
func TestRequestBodyBecameNotNullable(t *testing.T) {
	s1, err := open("../data/checker/request_property_became_nullable_base.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/request_property_became_nullable_base.yaml")
	require.NoError(t, err)

	s1.Spec.Paths.Value("/products").Post.RequestBody.Value.Content["application/json"].Schema.Value.Nullable = true

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)

	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.RequestPropertyBecameNotNullableCheck), d, osm, checker.ERR)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{
		Id:          checker.RequestBodyBecomeNotNullableId,
		Level:       checker.ERR,
		Operation:   "POST",
		Path:        "/products",
		Source:      load.NewSource("../data/checker/request_property_became_nullable_base.yaml"),
		OperationId: "addProduct",
	}, errs[0])
}
