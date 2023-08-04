package checker_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/checker"
	"github.com/tufin/oasdiff/diff"
)

// CL: removing an enum value from request parameter
func TestRequestParameterEnumValueRemovedCheck(t *testing.T) {
	s1, err := open("../data/checker/request_parameter_enum_value_updated_base.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/request_parameter_enum_value_updated_revision.yaml")
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)

	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.RequestParameterEnumValueUpdatedCheck), d, osm, checker.ERR)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{
		Id:          "request-parameter-enum-value-removed",
		Text:        "removed the enum value 'available' from the 'query' request parameter 'status'",
		Level:       checker.ERR,
		Operation:   "GET",
		Path:        "/test",
		Source:      "../data/checker/request_parameter_enum_value_updated_revision.yaml",
		OperationId: "getTest",
	}, errs[0])
}

// CL: adding an enum value to request parameter
func TestRequestParameterEnumValueAddedCheck(t *testing.T) {
	s1, err := open("../data/checker/request_parameter_enum_value_updated_revision.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/request_parameter_enum_value_updated_base.yaml")
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)

	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.RequestParameterEnumValueUpdatedCheck), d, osm, checker.INFO)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{
		Id:          "request-parameter-enum-value-added",
		Text:        "added the new enum value 'available' to the 'query' request parameter 'status'",
		Level:       checker.INFO,
		Operation:   "GET",
		Path:        "/test",
		Source:      "../data/checker/request_parameter_enum_value_updated_base.yaml",
		OperationId: "getTest",
	}, errs[0])
}
