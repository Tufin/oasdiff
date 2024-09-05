package checker_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/checker"
	"github.com/tufin/oasdiff/diff"
	"github.com/tufin/oasdiff/load"
)

// CL: decreasing minLength of request property
func TestRequestPropertyMinLengthDecreased(t *testing.T) {
	s1, err := open("../data/checker/request_property_min_length_decreased_base.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/request_property_min_length_decreased_base.yaml")
	require.NoError(t, err)

	s2.Spec.Paths.Value("/products").Post.RequestBody.Value.Content["application/json"].Schema.Value.Properties["name"].Value.MinLength = uint64(2)

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.RequestPropertyMinLengthUpdatedCheck), d, osm, checker.INFO)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{
		Id:          checker.RequestPropertyMinLengthDecreasedId,
		Args:        []any{"name", "application/json", uint64(3), uint64(2)},
		Level:       checker.INFO,
		Operation:   "POST",
		Path:        "/products",
		Source:      load.NewSource("../data/checker/request_property_min_length_decreased_base.yaml"),
		OperationId: "addProduct",
	}, errs[0])
	require.Equal(t, "minLength value of request property 'name' of media-type 'application/json' was decreased from '3' to '2'", errs[0].GetUncolorizedText(checker.NewDefaultLocalizer()))
}

// CL: increasing minLength of request property
func TestRequestPropertyMinLengthIncreased(t *testing.T) {
	s1, err := open("../data/checker/request_property_min_length_decreased_base.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/request_property_min_length_decreased_base.yaml")
	require.NoError(t, err)

	s2.Spec.Paths.Value("/products").Post.RequestBody.Value.Content["application/json"].Schema.Value.Properties["name"].Value.MinLength = uint64(5)

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.RequestPropertyMinLengthUpdatedCheck), d, osm, checker.INFO)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{
		Id:          checker.RequestPropertyMinLengthIncreasedId,
		Args:        []any{"name", "application/json", uint64(3), uint64(5)},
		Level:       checker.ERR,
		Operation:   "POST",
		Path:        "/products",
		Source:      load.NewSource("../data/checker/request_property_min_length_decreased_base.yaml"),
		OperationId: "addProduct",
	}, errs[0])
	require.Equal(t, "minLength value of request property 'name' of media-type 'application/json' was increased from '3' to '5'", errs[0].GetUncolorizedText(checker.NewDefaultLocalizer()))
}

// CL: increasing minLength of request body
func TestRequestBodyMinLengthIncreased(t *testing.T) {
	s1, err := open("../data/checker/request_property_min_length_decreased_base.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/request_property_min_length_decreased_base.yaml")
	require.NoError(t, err)

	s2.Spec.Paths.Value("/products").Post.RequestBody.Value.Content["application/json"].Schema.Value.MinLength = uint64(100)

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.RequestPropertyMinLengthUpdatedCheck), d, osm, checker.INFO)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{
		Id:          checker.RequestBodyMinLengthIncreasedId,
		Args:        []any{"application/json", uint64(10), uint64(100)},
		Level:       checker.ERR,
		Operation:   "POST",
		Path:        "/products",
		Source:      load.NewSource("../data/checker/request_property_min_length_decreased_base.yaml"),
		OperationId: "addProduct",
	}, errs[0])
	require.Equal(t, "minLength value of media-type 'application/json' of request body was increased from '10' to '100'", errs[0].GetUncolorizedText(checker.NewDefaultLocalizer()))
}

// CL: decreasing minLength of request body
func TestRequestBodyMinLengthDecreased(t *testing.T) {
	s1, err := open("../data/checker/request_property_min_length_decreased_base.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/request_property_min_length_decreased_base.yaml")
	require.NoError(t, err)

	s2.Spec.Paths.Value("/products").Post.RequestBody.Value.Content["application/json"].Schema.Value.MinLength = uint64(1)

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.RequestPropertyMinLengthUpdatedCheck), d, osm, checker.INFO)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{
		Id:          checker.RequestBodyMinLengthDecreasedId,
		Args:        []any{"application/json", uint64(10), uint64(1)},
		Level:       checker.INFO,
		Operation:   "POST",
		Path:        "/products",
		Source:      load.NewSource("../data/checker/request_property_min_length_decreased_base.yaml"),
		OperationId: "addProduct",
	}, errs[0])
	require.Equal(t, "minLength value of media-type 'application/json' of request body was decreased from '10' to '1'", errs[0].GetUncolorizedText(checker.NewDefaultLocalizer()))
}
