package checker_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/checker"
	"github.com/tufin/oasdiff/diff"
)

// CL: increasing maxItems of request parameters
func TestRequestParameterMaxItemsIncreased(t *testing.T) {
	s1, err := open("../data/checker/request_parameter_max_items_updated_base.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/request_parameter_max_items_updated_revision.yaml")
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.RequestParameterMaxItemsUpdatedCheck), d, osm, checker.INFO)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{
		Id:          checker.RequestParameterMaxItemsIncreasedId,
		Text:        "for the 'query' request parameter 'category', the maxItems was increased from '10' to '20'",
		Comment:     "",
		Level:       checker.INFO,
		Operation:   "POST",
		Path:        "/api/v1.0/groups",
		Source:      "../data/checker/request_parameter_max_items_updated_revision.yaml",
		OperationId: "createOneGroup",
	}, errs[0])
}

// CL: decreasing maxItems of request parameters
func TestRequestParameterMaxItemsDecreased(t *testing.T) {
	s1, err := open("../data/checker/request_parameter_max_items_updated_revision.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/request_parameter_max_items_updated_base.yaml")
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.RequestParameterMaxItemsUpdatedCheck), d, osm, checker.ERR)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{
		Id:          checker.RequestParameterMaxItemsDecreasedId,
		Text:        "for the 'query' request parameter 'category', the maxItems was decreased from '20' to '10'",
		Comment:     "",
		Level:       checker.ERR,
		Operation:   "POST",
		Path:        "/api/v1.0/groups",
		Source:      "../data/checker/request_parameter_max_items_updated_base.yaml",
		OperationId: "createOneGroup",
	}, errs[0])
}
