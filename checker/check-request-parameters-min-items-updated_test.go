package checker_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/checker"
	"github.com/tufin/oasdiff/diff"
)

// CL: increasing minItems value of request parameter
func TestRequestParameterMinItemsIncreased(t *testing.T) {
	s1, err := open("../data/checker/request_parameter_min_items_increased_base.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/request_parameter_min_items_increased_revision.yaml")
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.RequestParameterMinItemsUpdatedCheck), d, osm, checker.ERR)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{
		Id:          checker.RequestParameterMinItemsIncreasedId,
		Text:        "for the 'query' request parameter 'category', the minItems was increased from '2' to '3'",
		Args:        []any{"query", "category", uint64(2), uint64(3)},
		Level:       checker.ERR,
		Operation:   "POST",
		Path:        "/api/v1.0/groups",
		Source:      "../data/checker/request_parameter_min_items_increased_revision.yaml",
		OperationId: "createOneGroup",
	}, errs[0])
}

// CL: decreasing minItems value of request parameter
func TestRequestParameterMinItemsDecreased(t *testing.T) {
	s1, err := open("../data/checker/request_parameter_min_items_increased_revision.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/request_parameter_min_items_increased_base.yaml")
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.RequestParameterMinItemsUpdatedCheck), d, osm, checker.INFO)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{
		Id:          checker.RequestParameterMinItemsDecreasedId,
		Text:        "for the 'query' request parameter 'category', the minItems was decreased from '3' to '2'",
		Args:        []any{"query", "category", uint64(3), uint64(2)},
		Level:       checker.INFO,
		Operation:   "POST",
		Path:        "/api/v1.0/groups",
		Source:      "../data/checker/request_parameter_min_items_increased_base.yaml",
		OperationId: "createOneGroup",
	}, errs[0])
}
