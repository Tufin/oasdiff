package checker_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/checker"
	"github.com/tufin/oasdiff/diff"
)

// CL: Changing request path parameter type
func TestRequestPathParamTypeChanged(t *testing.T) {
	s1, err := open("../data/checker/request_parameter_type_changed_base.yaml")
	require.Empty(t, err)
	s2, err := open("../data/checker/request_parameter_type_changed_base.yaml")
	require.Empty(t, err)

	s2.Spec.Paths["/api/v1.0/groups"].Post.Parameters[0].Value.Schema.Value.Type = "int"

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.RequestParameterTypeChangedCheck), d, osm, checker.INFO)
	require.Len(t, errs, 1)
	require.Equal(t, checker.BackwardCompatibilityError{
		Id:          "request-parameter-type-changed",
		Text:        "for the 'path' request parameter 'groupId', the type/format was changed from 'string'/'none' to 'int'/'none'",
		Comment:     "",
		Level:       checker.ERR,
		Operation:   "POST",
		Path:        "/api/v1.0/groups",
		Source:      "../data/checker/request_parameter_type_changed_base.yaml",
		OperationId: "createOneGroup",
	}, errs[0])
}

// CL: Changing request query parameter type
func TestRequestQueryParamTypeChanged(t *testing.T) {
	s1, err := open("../data/checker/request_parameter_type_changed_base.yaml")
	require.Empty(t, err)
	s2, err := open("../data/checker/request_parameter_type_changed_base.yaml")
	require.Empty(t, err)

	s2.Spec.Paths["/api/v1.0/groups"].Post.Parameters[1].Value.Schema.Value.Type = "int"

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.RequestParameterTypeChangedCheck), d, osm, checker.INFO)
	require.Len(t, errs, 1)
	require.Equal(t, checker.BackwardCompatibilityError{
		Id:          "request-parameter-type-changed",
		Text:        "for the 'query' request parameter 'token', the type/format was changed from 'string'/'uuid' to 'int'/'uuid'",
		Comment:     "",
		Level:       checker.ERR,
		Operation:   "POST",
		Path:        "/api/v1.0/groups",
		Source:      "../data/checker/request_parameter_type_changed_base.yaml",
		OperationId: "createOneGroup",
	}, errs[0])
}

// CL: Changing request header parameter type
func TestRequestQueryHeaderTypeChanged(t *testing.T) {
	s1, err := open("../data/checker/request_parameter_type_changed_base.yaml")
	require.Empty(t, err)
	s2, err := open("../data/checker/request_parameter_type_changed_base.yaml")
	require.Empty(t, err)

	s2.Spec.Paths["/api/v1.0/groups"].Post.Parameters[2].Value.Schema.Value.Type = "int"

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.RequestParameterTypeChangedCheck), d, osm, checker.INFO)
	require.Len(t, errs, 1)
	require.Equal(t, checker.BackwardCompatibilityError{
		Id:          "request-parameter-type-changed",
		Text:        "for the 'header' request parameter 'X-Request-ID', the type/format was changed from 'string'/'uuid' to 'int'/'uuid'",
		Comment:     "",
		Level:       checker.ERR,
		Operation:   "POST",
		Path:        "/api/v1.0/groups",
		Source:      "../data/checker/request_parameter_type_changed_base.yaml",
		OperationId: "createOneGroup",
	}, errs[0])
}
