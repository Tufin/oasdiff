package checker_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/checker"
	"github.com/tufin/oasdiff/diff"
)

// CL: increasing minimum value of request parameter
func TestRequestParameterMinIncreased(t *testing.T) {
	s1, err := open("../data/checker/request_parameter_min_increased_base.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/request_parameter_min_increased_revision.yaml")
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.RequestParameterMinUpdatedCheck), d, osm, checker.ERR)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{
		Id:          checker.RequestParameterMinIncreasedId,
		Text:        "for the 'path' request parameter 'groupId', the min was increased from '1.00' to '10.00'",
		Level:       checker.ERR,
		Operation:   "POST",
		Path:        "/api/v1.0/groups",
		Source:      "../data/checker/request_parameter_min_increased_revision.yaml",
		OperationId: "createOneGroup",
	}, errs[0])
}

// CL: decreasing minimum value of request parameter
func TestRequestParameterMinDecreased(t *testing.T) {
	s1, err := open("../data/checker/request_parameter_min_increased_revision.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/request_parameter_min_increased_base.yaml")
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.RequestParameterMinUpdatedCheck), d, osm, checker.INFO)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{
		Id:          checker.RequestParameterMinDecreasedId,
		Text:        "for the 'path' request parameter 'groupId', the min was decreased from '10.00' to '1.00'",
		Level:       checker.INFO,
		Operation:   "POST",
		Path:        "/api/v1.0/groups",
		Source:      "../data/checker/request_parameter_min_increased_base.yaml",
		OperationId: "createOneGroup",
	}, errs[0])
}
