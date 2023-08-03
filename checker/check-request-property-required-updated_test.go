package checker_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/checker"
	"github.com/tufin/oasdiff/diff"
)

// CL: changing request property required value to true
func TestRequestPropertyMarkedRequired(t *testing.T) {
	s1, err := open("../data/checker/request_property_became_required_base.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/request_property_became_required_base.yaml")
	require.NoError(t, err)

	s1.Spec.Paths["/products"].Post.RequestBody.Value.Content["application/json"].Schema.Value.Required = []string{""}
	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.RequestPropertyRequiredUpdatedCheck), d, osm, checker.INFO)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{
		Id:          "request-property-became-required",
		Text:        "the request property 'name' became required",
		Level:       checker.ERR,
		Operation:   "POST",
		Path:        "/products",
		Source:      "../data/checker/request_property_became_required_base.yaml",
		OperationId: "addProduct",
	}, errs[0])
}

// CL: changing request property required value to false
func TestRequestPropertyMarkedOptional(t *testing.T) {
	s1, err := open("../data/checker/request_property_became_required_base.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/request_property_became_required_base.yaml")
	require.NoError(t, err)

	s2.Spec.Paths["/products"].Post.RequestBody.Value.Content["application/json"].Schema.Value.Required = []string{""}
	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.RequestPropertyRequiredUpdatedCheck), d, osm, checker.INFO)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{
		Id:          "request-property-became-optional",
		Text:        "the request property 'name' became optional",
		Level:       checker.INFO,
		Operation:   "POST",
		Path:        "/products",
		Source:      "../data/checker/request_property_became_required_base.yaml",
		OperationId: "addProduct",
	}, errs[0])
}
