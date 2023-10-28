package checker_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/checker"
	"github.com/tufin/oasdiff/diff"
)

// CL: changing request parameter default value
func TestRequestParameterDefaultValueChanged(t *testing.T) {
	s1, err := open("../data/checker/request_parameter_default_value_changed_base.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/request_parameter_default_value_changed_revision.yaml")
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.RequestParameterDefaultValueChangedCheck), d, osm, checker.ERR)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{
		Id:          "request-parameter-default-value-changed",
		Text:        "for the 'query' request parameter 'category', default value was changed from 'default_category' to 'updated_category'",
		Comment:     "",
		Level:       checker.ERR,
		Operation:   "POST",
		Path:        "/api/v1.0/groups",
		Source:      "../data/checker/request_parameter_default_value_changed_revision.yaml",
		OperationId: "createOneGroup",
	}, errs[0])
}

// CL: adding request parameter default value
func TestRequestParameterDefaultValueAdded(t *testing.T) {
	s1, err := open("../data/checker/request_parameter_default_value_changed_base.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/request_parameter_default_value_changed_base.yaml")
	require.NoError(t, err)

	s1.Spec.Paths["/api/v1.0/groups"].Post.Parameters[1].Value.Schema.Value.Default = nil

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.RequestParameterDefaultValueChangedCheck), d, osm, checker.ERR)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{
		Id:          "request-parameter-default-value-added",
		Text:        "for the 'query' request parameter 'category', default value 'default_category' was added",
		Comment:     "",
		Level:       checker.ERR,
		Operation:   "POST",
		Path:        "/api/v1.0/groups",
		Source:      "../data/checker/request_parameter_default_value_changed_base.yaml",
		OperationId: "createOneGroup",
	}, errs[0])
}

// CL: removing request parameter default value
func TestRequestParameterDefaultValueRemoved(t *testing.T) {
	s1, err := open("../data/checker/request_parameter_default_value_changed_base.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/request_parameter_default_value_changed_base.yaml")
	require.NoError(t, err)

	s2.Spec.Paths["/api/v1.0/groups"].Post.Parameters[1].Value.Schema.Value.Default = nil

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.RequestParameterDefaultValueChangedCheck), d, osm, checker.ERR)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{
		Id:          "request-parameter-default-value-removed",
		Text:        "for the 'query' request parameter 'category', default value 'default_category' was removed",
		Comment:     "",
		Level:       checker.ERR,
		Operation:   "POST",
		Path:        "/api/v1.0/groups",
		Source:      "../data/checker/request_parameter_default_value_changed_base.yaml",
		OperationId: "createOneGroup",
	}, errs[0])
}
