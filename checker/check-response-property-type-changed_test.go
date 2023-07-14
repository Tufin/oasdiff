package checker_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/checker"
	"github.com/tufin/oasdiff/diff"
)

// CL: Changing a response schema type
func TestResponseSchemaTypeChangedCheck(t *testing.T) {
	s1, err := open("../data/checker/response_schema_type_changed_base.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/response_schema_type_changed_revision.yaml")
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.ResponsePropertyTypeChangedCheck), d, osm, checker.ERR)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{
		Id:          "response-body-type-changed",
		Text:        "the response's body type/format changed from 'string'/'none' to 'object'/'none' for status '200'",
		Comment:     "",
		Level:       checker.ERR,
		Operation:   "POST",
		Path:        "/api/v1.0/groups",
		Source:      "../data/checker/response_schema_type_changed_revision.yaml",
		OperationId: "createOneGroup",
	}, errs[0])
}

// CL: Changing a response property schema type
func TestResponsePropertyTypeChangedCheck(t *testing.T) {
	s1, err := open("../data/checker/response_schema_type_changed_revision.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/response_schema_type_changed_revision.yaml")
	require.NoError(t, err)

	s2.Spec.Paths["/api/v1.0/groups"].Post.Responses["200"].Value.Content["application/json"].Schema.Value.Properties["data"].Value.Properties["name"].Value.Type = "integer"

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.ResponsePropertyTypeChangedCheck), d, osm, checker.ERR)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{
		Id:          "response-property-type-changed",
		Text:        "the response's property type/format changed from 'string'/'none' to 'integer'/'none' for status '200'",
		Comment:     "",
		Level:       checker.ERR,
		Operation:   "POST",
		Path:        "/api/v1.0/groups",
		Source:      "../data/checker/response_schema_type_changed_revision.yaml",
		OperationId: "createOneGroup",
	}, errs[0])
}

// CL: Changing a response property schema format
func TestResponsePropertyFormatChangedCheck(t *testing.T) {
	s1, err := open("../data/checker/response_schema_format_changed_base.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/response_schema_format_changed_base.yaml")
	require.NoError(t, err)

	s2.Spec.Paths["/api/v1.0/groups"].Post.Responses["200"].Value.Content["application/json"].Schema.Value.Properties["data"].Value.Properties["name"].Value.Format = "uuid"

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.ResponsePropertyTypeChangedCheck), d, osm, checker.ERR)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{
		Id:          "response-property-type-changed",
		Text:        "the response's property type/format changed from 'string'/'hostname' to 'string'/'uuid' for status '200'",
		Comment:     "",
		Level:       checker.ERR,
		Operation:   "POST",
		Path:        "/api/v1.0/groups",
		Source:      "../data/checker/response_schema_format_changed_base.yaml",
		OperationId: "createOneGroup",
	}, errs[0])
}
