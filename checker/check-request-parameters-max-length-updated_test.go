package checker_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/checker"
	"github.com/tufin/oasdiff/diff"
)

func TestRequestParameterMaxLengthIncreasedCheck(t *testing.T) {
	s1, err := open("../data/checker/request_parameter_max_length_updated_base.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/request_parameter_max_length_updated_revision.yaml")
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.RequestParameterMaxLengthUpdatedCheck), d, osm, checker.INFO)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{
		Id:          "request-parameter-max-length-increased",
		Text:        "for the 'query' request parameter 'category', the maxLength was increased from '10' to '15'",
		Level:       checker.INFO,
		Operation:   "POST",
		Path:        "/test",
		Source:      "../data/checker/request_parameter_max_length_updated_revision.yaml",
		OperationId: "",
	}, errs[0])
}

func TestRequestParameterMaxLengthDecreasedCheck(t *testing.T) {
	s1, err := open("../data/checker/request_parameter_max_length_updated_revision.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/request_parameter_max_length_updated_base.yaml")
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.RequestParameterMaxLengthUpdatedCheck), d, osm, checker.INFO)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{
		Id:          "request-parameter-max-length-decreased",
		Text:        "for the 'query' request parameter 'category', the maxLength was decreased from '15' to '10'",
		Level:       checker.ERR,
		Operation:   "POST",
		Path:        "/test",
		Source:      "../data/checker/request_parameter_max_length_updated_base.yaml",
		OperationId: "",
	}, errs[0])
}
