package checker_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/checker"
	"github.com/tufin/oasdiff/diff"
	"github.com/tufin/oasdiff/load"
)

// CL: increasing minLength value of request parameter
func TestRequestParameterMinLengthIncreasedCheck(t *testing.T) {
	s1, err := open("../data/checker/request_parameter_min_length_increased_base.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/request_parameter_min_length_increased_revision.yaml")
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)

	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.RequestParameterMinLengthUpdatedCheck), d, osm, checker.INFO)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{
		Id:          checker.RequestParameterMinLengthIncreasedId,
		Args:        []any{"query", "name", uint64(3), uint64(5)},
		Level:       checker.ERR,
		Operation:   "POST",
		Path:        "/test",
		Source:      load.NewSource("../data/checker/request_parameter_min_length_increased_revision.yaml"),
		OperationId: "createTest",
	}, errs[0])
}

// CL: decreasing minLength value of request parameter
func TestRequestParameterMinLengthDecreasedCheck(t *testing.T) {
	s1, err := open("../data/checker/request_parameter_min_length_increased_revision.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/request_parameter_min_length_increased_base.yaml")
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)

	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.RequestParameterMinLengthUpdatedCheck), d, osm, checker.INFO)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{
		Id:          checker.RequestParameterMinLengthDecreasedId,
		Args:        []any{"query", "name", uint64(5), uint64(3)},
		Level:       checker.INFO,
		Operation:   "POST",
		Path:        "/test",
		Source:      load.NewSource("../data/checker/request_parameter_min_length_increased_base.yaml"),
		OperationId: "createTest",
	}, errs[0])
}
