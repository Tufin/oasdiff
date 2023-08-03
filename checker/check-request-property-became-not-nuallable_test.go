package checker_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/checker"
	"github.com/tufin/oasdiff/diff"
)

// CL: changing request property to not nullable
func TestRequestPropertyBecameNotNullable(t *testing.T) {
	s1, err := open("../data/checker/request_property_became_nullable_revision.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/request_property_became_nullable_base.yaml")
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.RequestPropertyBecameNotNullableCheck), d, osm, checker.INFO)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{
		Id:          "request-property-became-not-nullable",
		Text:        "the request property 'name' became not nullable",
		Comment:     "",
		Level:       checker.ERR,
		Operation:   "POST",
		Path:        "/products",
		Source:      "../data/checker/request_property_became_nullable_base.yaml",
		OperationId: "addProduct",
	}, errs[0])
}

// CL: changing request property to nullable
func TestRequestPropertyBecameNullable(t *testing.T) {
	s1, err := open("../data/checker/request_property_became_nullable_base.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/request_property_became_nullable_revision.yaml")
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.RequestPropertyBecameNotNullableCheck), d, osm, checker.INFO)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{
		Id:          "request-property-became-nullable",
		Text:        "the request property 'name' became nullable",
		Comment:     "",
		Level:       checker.INFO,
		Operation:   "POST",
		Path:        "/products",
		Source:      "../data/checker/request_property_became_nullable_revision.yaml",
		OperationId: "addProduct",
	}, errs[0])

}

// CL: changing request body to nullable
func TestRequestBodyBecameNullable(t *testing.T) {
	s1, err := open("../data/checker/request_property_became_nullable_base.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/request_property_became_nullable_base.yaml")
	require.NoError(t, err)

	s2.Spec.Paths["/products"].Post.RequestBody.Value.Content["application/json"].Schema.Value.Nullable = true

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)

	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.RequestPropertyBecameNotNullableCheck), d, osm, checker.INFO)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{
		Id:          "request-body-became-nullable",
		Text:        "the request's body became nullable",
		Comment:     "",
		Level:       checker.INFO,
		Operation:   "POST",
		Path:        "/products",
		Source:      "../data/checker/request_property_became_nullable_base.yaml",
		OperationId: "addProduct",
	}, errs[0])
}

// CL: changing request body to not nullable
func TestRequestBodyBecameNotNullable(t *testing.T) {
	s1, err := open("../data/checker/request_property_became_nullable_base.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/request_property_became_nullable_base.yaml")
	require.NoError(t, err)

	s1.Spec.Paths["/products"].Post.RequestBody.Value.Content["application/json"].Schema.Value.Nullable = true

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)

	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.RequestPropertyBecameNotNullableCheck), d, osm, checker.ERR)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{
		Id:          "request-body-became-not-nullable",
		Text:        "the request's body became not nullable",
		Comment:     "",
		Level:       checker.ERR,
		Operation:   "POST",
		Path:        "/products",
		Source:      "../data/checker/request_property_became_nullable_base.yaml",
		OperationId: "addProduct",
	}, errs[0])
}
