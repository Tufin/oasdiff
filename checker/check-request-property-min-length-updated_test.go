package checker_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/checker"
	"github.com/tufin/oasdiff/diff"
)

// CL: decreasing minLength of request property
func TestRequestPropertyMinLengthDecreased(t *testing.T) {
	s1, err := open("../data/checker/request_property_min_length_decreased_base.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/request_property_min_length_decreased_base.yaml")
	require.NoError(t, err)

	s2.Spec.Paths["/products"].Post.RequestBody.Value.Content["application/json"].Schema.Value.Properties["name"].Value.MinLength = uint64(2)

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.RequestPropertyMinLengthUpdatedCheck), d, osm, checker.INFO)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{
		Id:          "request-property-min-length-decreased",
		Text:        "the 'name' request property's minLength was decreased from '3' to '2'",
		Comment:     "",
		Level:       checker.INFO,
		Operation:   "POST",
		Path:        "/products",
		Source:      "../data/checker/request_property_min_length_decreased_base.yaml",
		OperationId: "addProduct",
	}, errs[0])
}

// CL: increasing minLength of request property
func TestRequestPropertyMinLengthIncreased(t *testing.T) {
	s1, err := open("../data/checker/request_property_min_length_decreased_base.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/request_property_min_length_decreased_base.yaml")
	require.NoError(t, err)

	s2.Spec.Paths["/products"].Post.RequestBody.Value.Content["application/json"].Schema.Value.Properties["name"].Value.MinLength = uint64(5)

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.RequestPropertyMinLengthUpdatedCheck), d, osm, checker.INFO)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{
		Id:          "request-property-min-length-increased",
		Text:        "the 'name' request property's maxLength was increased from '3' to '5'",
		Comment:     "",
		Level:       checker.ERR,
		Operation:   "POST",
		Path:        "/products",
		Source:      "../data/checker/request_property_min_length_decreased_base.yaml",
		OperationId: "addProduct",
	}, errs[0])
}

// CL: increasing minLength of request body
func TestRequestBodyMinLengthIncreased(t *testing.T) {
	s1, err := open("../data/checker/request_property_min_length_decreased_base.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/request_property_min_length_decreased_base.yaml")
	require.NoError(t, err)

	s2.Spec.Paths["/products"].Post.RequestBody.Value.Content["application/json"].Schema.Value.MinLength = uint64(100)

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.RequestPropertyMinLengthUpdatedCheck), d, osm, checker.INFO)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{
		Id:          "request-body-min-length-increased",
		Text:        "the request's body minLength was increased from '10' to '100'",
		Comment:     "",
		Level:       checker.ERR,
		Operation:   "POST",
		Path:        "/products",
		Source:      "../data/checker/request_property_min_length_decreased_base.yaml",
		OperationId: "addProduct",
	}, errs[0])
}

// CL: decreasing minLength of request body
func TestRequestBodyMinLengthDecreased(t *testing.T) {
	s1, err := open("../data/checker/request_property_min_length_decreased_base.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/request_property_min_length_decreased_base.yaml")
	require.NoError(t, err)

	s2.Spec.Paths["/products"].Post.RequestBody.Value.Content["application/json"].Schema.Value.MinLength = uint64(1)

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.RequestPropertyMinLengthUpdatedCheck), d, osm, checker.INFO)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{
		Id:          "request-body-min-length-decreased",
		Text:        "the request's body minLength was decreased from '10' to '1'",
		Comment:     "",
		Level:       checker.INFO,
		Operation:   "POST",
		Path:        "/products",
		Source:      "../data/checker/request_property_min_length_decreased_base.yaml",
		OperationId: "addProduct",
	}, errs[0])
}
