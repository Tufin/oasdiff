package checker_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/checker"
	"github.com/tufin/oasdiff/diff"
)

// CL: decreasing request property maximum value
func TestRequestPropertyMaxDecreasedCheck(t *testing.T) {
	s1, err := open("../data/checker/request_property_max_decreased_base.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/request_property_max_decreased_base.yaml")
	require.NoError(t, err)

	max := float64(10)
	s2.Spec.Paths["/pets"].Post.RequestBody.Value.Content["application/json"].Schema.Value.Properties["name"].Value.Max = &max

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)

	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.RequestPropertyMaxDecreasedCheck), d, osm, checker.INFO)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{
		Id:          "request-property-max-decreased",
		Level:       checker.ERR,
		Text:        "the 'name' request property's max was decreased to '10.00'",
		Operation:   "POST",
		Path:        "/pets",
		Source:      "../data/checker/request_property_max_decreased_base.yaml",
		OperationId: "addPet",
	}, errs[0])
}

// CL: increasing request property maximum value
func TestRequestPropertyMaxIncreasingCheck(t *testing.T) {
	s1, err := open("../data/checker/request_property_max_decreased_base.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/request_property_max_decreased_base.yaml")
	require.NoError(t, err)

	max := float64(20)
	s2.Spec.Paths["/pets"].Post.RequestBody.Value.Content["application/json"].Schema.Value.Properties["name"].Value.Max = &max

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)

	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.RequestPropertyMaxDecreasedCheck), d, osm, checker.INFO)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{
		Id:          "request-property-max-increased",
		Level:       checker.INFO,
		Text:        "the 'name' request property's max was increased from '15.00' to '20.00'",
		Operation:   "POST",
		Path:        "/pets",
		Source:      "../data/checker/request_property_max_decreased_base.yaml",
		OperationId: "addPet",
	}, errs[0])
}

// CL: increasing request body maximum value
func TestRequestBodyMaxIncreasingCheck(t *testing.T) {
	s1, err := open("../data/checker/request_property_max_decreased_base.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/request_property_max_decreased_base.yaml")
	require.NoError(t, err)

	max := float64(20)
	newMax := float64(25)
	s1.Spec.Paths["/pets"].Post.RequestBody.Value.Content["application/json"].Schema.Value.Max = &max
	s2.Spec.Paths["/pets"].Post.RequestBody.Value.Content["application/json"].Schema.Value.Max = &newMax

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)

	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.RequestPropertyMaxDecreasedCheck), d, osm, checker.INFO)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{
		Id:          "request-body-max-increased",
		Level:       checker.INFO,
		Text:        "the request's body max was increased from '20.00' to '25.00'",
		Operation:   "POST",
		Path:        "/pets",
		Source:      "../data/checker/request_property_max_decreased_base.yaml",
		OperationId: "addPet",
	}, errs[0])
}

// CL: decreasing request body maximum value
func TestRequestBodyMaxDecreasedCheck(t *testing.T) {
	s1, err := open("../data/checker/request_property_max_decreased_base.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/request_property_max_decreased_base.yaml")
	require.NoError(t, err)

	max := float64(25)
	newMax := float64(20)
	s1.Spec.Paths["/pets"].Post.RequestBody.Value.Content["application/json"].Schema.Value.Max = &max
	s2.Spec.Paths["/pets"].Post.RequestBody.Value.Content["application/json"].Schema.Value.Max = &newMax

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)

	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.RequestPropertyMaxDecreasedCheck), d, osm, checker.INFO)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{
		Id:          "request-body-max-decreased",
		Level:       checker.ERR,
		Text:        "the request's body max was decreased to '20.00'",
		Operation:   "POST",
		Path:        "/pets",
		Source:      "../data/checker/request_property_max_decreased_base.yaml",
		OperationId: "addPet",
	}, errs[0])
}
