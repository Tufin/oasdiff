package checker_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/checker"
	"github.com/tufin/oasdiff/diff"
)

// CL: changing request path parameter type
func TestRequestPathParamTypeChanged(t *testing.T) {
	s1, err := open("../data/checker/request_parameter_type_changed_base.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/request_parameter_type_changed_base.yaml")
	require.NoError(t, err)

	s2.Spec.Paths.Value("/api/v1.0/groups").Post.Parameters[0].Value.Schema.Value.Type = "int"

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.RequestParameterTypeChangedCheck), d, osm, checker.INFO)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{
		Id:          checker.RequestParameterTypeChangedId,
		Args:        []any{"path", "groupId", "string", "", "int", ""},
		Level:       checker.ERR,
		Operation:   "POST",
		Path:        "/api/v1.0/groups",
		Source:      "../data/checker/request_parameter_type_changed_base.yaml",
		OperationId: "createOneGroup",
	}, errs[0])
}

// CL: changing request query parameter type
func TestRequestQueryParamTypeChanged(t *testing.T) {
	s1, err := open("../data/checker/request_parameter_type_changed_base.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/request_parameter_type_changed_base.yaml")
	require.NoError(t, err)

	s2.Spec.Paths.Value("/api/v1.0/groups").Post.Parameters[1].Value.Schema.Value.Type = "int"

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.RequestParameterTypeChangedCheck), d, osm, checker.INFO)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{
		Id:          checker.RequestParameterTypeChangedId,
		Args:        []any{"query", "token", "string", "uuid", "int", "uuid"},
		Level:       checker.ERR,
		Operation:   "POST",
		Path:        "/api/v1.0/groups",
		Source:      "../data/checker/request_parameter_type_changed_base.yaml",
		OperationId: "createOneGroup",
	}, errs[0])
}

// CL: changing request header parameter type
func TestRequestQueryHeaderTypeChanged(t *testing.T) {
	s1, err := open("../data/checker/request_parameter_type_changed_base.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/request_parameter_type_changed_base.yaml")
	require.NoError(t, err)

	s2.Spec.Paths.Value("/api/v1.0/groups").Post.Parameters[2].Value.Schema.Value.Type = "int"

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.RequestParameterTypeChangedCheck), d, osm, checker.INFO)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{
		Id:          checker.RequestParameterTypeChangedId,
		Args:        []any{"header", "X-Request-ID", "string", "uuid", "int", "uuid"},
		Level:       checker.ERR,
		Operation:   "POST",
		Path:        "/api/v1.0/groups",
		Source:      "../data/checker/request_parameter_type_changed_base.yaml",
		OperationId: "createOneGroup",
	}, errs[0])
}

// CL: changing request path parameter format
func TestRequestPathParamFormatChanged(t *testing.T) {
	s1, err := open("../data/checker/request_parameter_type_changed_base.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/request_parameter_type_changed_base.yaml")
	require.NoError(t, err)

	s2.Spec.Paths.Value("/api/v1.0/groups").Post.Parameters[0].Value.Schema.Value.Format = "uuid"

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.RequestParameterTypeChangedCheck), d, osm, checker.INFO)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{
		Id:          checker.RequestParameterTypeChangedId,
		Args:        []any{"path", "groupId", "string", "", "string", "uuid"},
		Level:       checker.ERR,
		Operation:   "POST",
		Path:        "/api/v1.0/groups",
		Source:      "../data/checker/request_parameter_type_changed_base.yaml",
		OperationId: "createOneGroup",
	}, errs[0])
}

// CL: changing request query parameter format
func TestRequestQueryParamFormatChanged(t *testing.T) {
	s1, err := open("../data/checker/request_parameter_type_changed_base.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/request_parameter_type_changed_base.yaml")
	require.NoError(t, err)

	s2.Spec.Paths.Value("/api/v1.0/groups").Post.Parameters[1].Value.Schema.Value.Format = "uri"

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.RequestParameterTypeChangedCheck), d, osm, checker.INFO)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{
		Id:          checker.RequestParameterTypeChangedId,
		Args:        []any{"query", "token", "string", "uuid", "string", "uri"},
		Level:       checker.ERR,
		Operation:   "POST",
		Path:        "/api/v1.0/groups",
		Source:      "../data/checker/request_parameter_type_changed_base.yaml",
		OperationId: "createOneGroup",
	}, errs[0])
}

// CL: changing request header parameter format
func TestRequestQueryHeaderFormatChanged(t *testing.T) {
	s1, err := open("../data/checker/request_parameter_type_changed_base.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/request_parameter_type_changed_base.yaml")
	require.NoError(t, err)

	s2.Spec.Paths.Value("/api/v1.0/groups").Post.Parameters[2].Value.Schema.Value.Format = "uri"

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.RequestParameterTypeChangedCheck), d, osm, checker.INFO)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{
		Id:          checker.RequestParameterTypeChangedId,
		Args:        []any{"header", "X-Request-ID", "string", "uuid", "string", "uri"},
		Level:       checker.ERR,
		Operation:   "POST",
		Path:        "/api/v1.0/groups",
		Source:      "../data/checker/request_parameter_type_changed_base.yaml",
		OperationId: "createOneGroup",
	}, errs[0])
}
