package checker_test

import (
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/checker"
	"github.com/tufin/oasdiff/diff"
	"github.com/tufin/oasdiff/load"
	"github.com/tufin/oasdiff/utils"
)

// CL: changing request path parameter type
func TestRequestPathParamTypeChanged(t *testing.T) {
	s1, err := open("../data/checker/request_parameter_type_changed_base.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/request_parameter_type_changed_base.yaml")
	require.NoError(t, err)

	s2.Spec.Paths.Value("/api/v1.0/groups").Post.Parameters[0].Value.Schema.Value.Type = &openapi3.Types{"integer"}

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.RequestParameterTypeChangedCheck), d, osm, checker.INFO)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{
		Id:          checker.RequestParameterTypeChangedId,
		Args:        []any{"path", "groupId", utils.StringList{"string"}, "", utils.StringList{"integer"}, ""},
		Level:       checker.ERR,
		Operation:   "POST",
		Path:        "/api/v1.0/groups",
		Source:      load.NewSource("../data/checker/request_parameter_type_changed_base.yaml"),
		OperationId: "createOneGroup",
	}, errs[0])
}

// CL: changing request query parameter type
func TestRequestQueryParamTypeChanged(t *testing.T) {
	s1, err := open("../data/checker/request_parameter_type_changed_base.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/request_parameter_type_changed_base.yaml")
	require.NoError(t, err)

	s2.Spec.Paths.Value("/api/v1.0/groups").Post.Parameters[1].Value.Schema.Value.Type = &openapi3.Types{"integer"}

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.RequestParameterTypeChangedCheck), d, osm, checker.INFO)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{
		Id:          checker.RequestParameterTypeChangedId,
		Args:        []any{"query", "token", utils.StringList{"string"}, "uuid", utils.StringList{"integer"}, "uuid"},
		Level:       checker.ERR,
		Operation:   "POST",
		Path:        "/api/v1.0/groups",
		Source:      load.NewSource("../data/checker/request_parameter_type_changed_base.yaml"),
		OperationId: "createOneGroup",
	}, errs[0])
}

// CL: changing request header parameter type
func TestRequestQueryHeaderTypeChanged(t *testing.T) {
	s1, err := open("../data/checker/request_parameter_type_changed_base.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/request_parameter_type_changed_base.yaml")
	require.NoError(t, err)

	s2.Spec.Paths.Value("/api/v1.0/groups").Post.Parameters[2].Value.Schema.Value.Type = &openapi3.Types{"integer"}

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.RequestParameterTypeChangedCheck), d, osm, checker.INFO)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{
		Id:          checker.RequestParameterTypeChangedId,
		Args:        []any{"header", "X-Request-ID", utils.StringList{"string"}, "uuid", utils.StringList{"integer"}, "uuid"},
		Level:       checker.ERR,
		Operation:   "POST",
		Path:        "/api/v1.0/groups",
		Source:      load.NewSource("../data/checker/request_parameter_type_changed_base.yaml"),
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

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.RequestParameterTypeChangedCheck), d, osm, checker.INFO)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{
		Id:          checker.RequestParameterTypeChangedId,
		Args:        []any{"path", "groupId", utils.StringList{"string"}, "", utils.StringList{"string"}, "uuid"},
		Level:       checker.ERR,
		Operation:   "POST",
		Path:        "/api/v1.0/groups",
		Source:      load.NewSource("../data/checker/request_parameter_type_changed_base.yaml"),
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

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.RequestParameterTypeChangedCheck), d, osm, checker.INFO)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{
		Id:          checker.RequestParameterTypeChangedId,
		Args:        []any{"query", "token", utils.StringList{"string"}, "uuid", utils.StringList{"string"}, "uri"},
		Level:       checker.ERR,
		Operation:   "POST",
		Path:        "/api/v1.0/groups",
		Source:      load.NewSource("../data/checker/request_parameter_type_changed_base.yaml"),
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

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.RequestParameterTypeChangedCheck), d, osm, checker.INFO)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{
		Id:          checker.RequestParameterTypeChangedId,
		Args:        []any{"header", "X-Request-ID", utils.StringList{"string"}, "uuid", utils.StringList{"string"}, "uri"},
		Level:       checker.ERR,
		Operation:   "POST",
		Path:        "/api/v1.0/groups",
		Source:      load.NewSource("../data/checker/request_parameter_type_changed_base.yaml"),
		OperationId: "createOneGroup",
	}, errs[0])
}

// CL: changing request path parameter type by adding "string"
func TestRequestPathParamTypeAddString(t *testing.T) {
	s1, err := open("../data/checker/request_parameter_type_changed_base.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/request_parameter_type_changed_base.yaml")
	require.NoError(t, err)

	s1.Spec.Paths.Value("/api/v1.0/groups").Post.Parameters[0].Value.Schema.Value.Type = &openapi3.Types{"integer"}
	s2.Spec.Paths.Value("/api/v1.0/groups").Post.Parameters[0].Value.Schema.Value.Type = &openapi3.Types{"integer", "string"}

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.RequestParameterTypeChangedCheck), d, osm, checker.INFO)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{
		Id:          checker.RequestParameterTypeGeneralizedId,
		Args:        []any{"path", "groupId", utils.StringList{"integer"}, "", utils.StringList{"integer", "string"}, ""},
		Level:       checker.INFO,
		Operation:   "POST",
		Path:        "/api/v1.0/groups",
		Source:      load.NewSource("../data/checker/request_parameter_type_changed_base.yaml"),
		OperationId: "createOneGroup",
	}, errs[0])
}

// CL: changing request path parameter type by replacing "integer" with "number"
func TestRequestPathParamTypeIntegerToNumber(t *testing.T) {
	s1, err := open("../data/checker/request_parameter_type_changed_base.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/request_parameter_type_changed_base.yaml")
	require.NoError(t, err)

	s1.Spec.Paths.Value("/api/v1.0/groups").Post.Parameters[0].Value.Schema.Value.Type = &openapi3.Types{"integer", "string"}
	s2.Spec.Paths.Value("/api/v1.0/groups").Post.Parameters[0].Value.Schema.Value.Type = &openapi3.Types{"number", "string"}

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.RequestParameterTypeChangedCheck), d, osm, checker.INFO)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{
		Id:          checker.RequestParameterTypeGeneralizedId,
		Args:        []any{"path", "groupId", utils.StringList{"integer", "string"}, "", utils.StringList{"number", "string"}, ""},
		Level:       checker.INFO,
		Operation:   "POST",
		Path:        "/api/v1.0/groups",
		Source:      load.NewSource("../data/checker/request_parameter_type_changed_base.yaml"),
		OperationId: "createOneGroup",
	}, errs[0])
}

// BC: changing request's query param property type from number to string is breaking
func TestBreaking_ReqQueryParamTypeNumberToString(t *testing.T) {
	s1, err := open("../data/checker/request_parameter_property_type_changed_base_num.yaml")
	require.NoError(t, err)

	s2, err := open("../data/checker/request_parameter_property_type_changed_revision.yaml")
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.RequestParameterTypeChangedCheck), d, osm, checker.INFO)
	require.Len(t, errs, 1)
	require.Equal(t, checker.RequestParameterPropertyTypeChangedId, errs[0].GetId())
	require.Equal(t, "for the 'query' request parameter 'filters', the type/format of property 'groupId' was changed from 'number'/'' to 'string'/''", errs[0].GetUncolorizedText(checker.NewDefaultLocalizer()))
	require.Equal(t, checker.WARN, errs[0].GetLevel())
}

// BC: specializing request's query param property type from string to number is breaking
func TestBreaking_ReqQueryParamTypeStringToNumber(t *testing.T) {
	s1, err := open("../data/checker/request_parameter_property_type_changed_revision.yaml")
	require.NoError(t, err)

	s2, err := open("../data/checker/request_parameter_property_type_changed_base_num.yaml")
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.RequestParameterTypeChangedCheck), d, osm, checker.INFO)
	require.Len(t, errs, 1)
	require.Equal(t, checker.RequestParameterPropertyTypeSpecializedId, errs[0].GetId())
	require.Equal(t, "for the 'query' request parameter 'filters', the type/format of property 'groupId' was specialized from 'string'/'' to 'number'/''", errs[0].GetUncolorizedText(checker.NewDefaultLocalizer()))
	require.Equal(t, checker.ERR, errs[0].GetLevel())
}

// CL: generalizing request's query param property type from integer to number
func TestBreaking_ReqQueryParamTypeIntegerToNumber(t *testing.T) {
	s1, err := open("../data/checker/request_parameter_property_type_changed_base_int.yaml")
	require.NoError(t, err)

	s2, err := open("../data/checker/request_parameter_property_type_changed_base_num.yaml")
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.RequestParameterTypeChangedCheck), d, osm, checker.INFO)
	require.Len(t, errs, 1)
	require.Equal(t, checker.RequestParameterPropertyTypeGeneralizedId, errs[0].GetId())
	require.Equal(t, "for the 'query' request parameter 'filters', the type/format of property 'groupId' was generalized from 'integer'/'' to 'number'/''", errs[0].GetUncolorizedText(checker.NewDefaultLocalizer()))
	require.Equal(t, checker.INFO, errs[0].GetLevel())
}
