package checker_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/checker"
	"github.com/tufin/oasdiff/diff"
)

// CL: increasing maximum value of request parameter
func TestRequestParameterMaxIncreased(t *testing.T) {
	s1, err := open("../data/checker/request_parameter_max_updated_base.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/request_parameter_max_updated_revision.yaml")
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.RequestParameterMaxUpdatedCheck), d, osm, checker.INFO)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{
		Id:          "request-parameter-max-increased",
		Text:        "for the 'query' request parameter 'category', the max was increased from '5.00' to '10.00'",
		Level:       checker.INFO,
		Operation:   "POST",
		Path:        "/api/v1.0/groups",
		Source:      "../data/checker/request_parameter_max_updated_revision.yaml",
		OperationId: "createOneGroup",
	}, errs[0])
}

// CL: decreasing maximum value of request parameter
func TestRequestParameterMaxDecreased(t *testing.T) {
	s1, err := open("../data/checker/request_parameter_max_updated_revision.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/request_parameter_max_updated_base.yaml")
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.RequestParameterMaxUpdatedCheck), d, osm, checker.INFO)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{
		Id:          "request-parameter-max-decreased",
		Text:        "for the 'query' request parameter 'category', the max was decreased from '10.00' to '5.00'",
		Level:       checker.ERR,
		Operation:   "POST",
		Path:        "/api/v1.0/groups",
		Source:      "../data/checker/request_parameter_max_updated_base.yaml",
		OperationId: "createOneGroup",
	}, errs[0])
}
