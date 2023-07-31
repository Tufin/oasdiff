package checker_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/checker"
	"github.com/tufin/oasdiff/diff"
)

// CL: changing request body type
func TestRequestBodyTypeChangedCheck(t *testing.T) {
	s1, err := open("../data/checker/request_property_type_changed_base.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/request_property_type_changed_base.yaml")
	require.NoError(t, err)

	s2.Spec.Paths["/pets"].Post.RequestBody.Value.Content["application/json"].Schema.Value.Type = "array"

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)

	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.RequestPropertyTypeChangedCheck), d, osm, checker.ERR)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{
		Id:          "request-body-type-changed",
		Level:       checker.ERR,
		Text:        "the request's body type/format changed from 'object'/'none' to 'array'/'none'",
		Operation:   "POST",
		Path:        "/pets",
		Source:      "../data/checker/request_property_type_changed_base.yaml",
		OperationId: "addPet",
	}, errs[0])
}

// CL: changing request body type
func TestRequestBodyFormatChangedCheck(t *testing.T) {
	s1, err := open("../data/checker/request_property_type_changed_base.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/request_property_type_changed_base.yaml")
	require.NoError(t, err)

	s2.Spec.Paths["/pets"].Post.RequestBody.Value.Content["application/json"].Schema.Value.Format = "uuid"

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)

	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.RequestPropertyTypeChangedCheck), d, osm, checker.ERR)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{
		Id:          "request-body-type-changed",
		Level:       checker.ERR,
		Text:        "the request's body type/format changed from 'object'/'none' to 'object'/'uuid'",
		Operation:   "POST",
		Path:        "/pets",
		Source:      "../data/checker/request_property_type_changed_base.yaml",
		OperationId: "addPet",
	}, errs[0])
}

// CL: changing request property type
func TestRequestPropertyTypeChangedCheck(t *testing.T) {
	s1, err := open("../data/checker/request_property_type_changed_base.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/request_property_type_changed_revision.yaml")
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)

	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.RequestPropertyTypeChangedCheck), d, osm, checker.ERR)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{
		Id:          "request-property-type-changed",
		Level:       checker.ERR,
		Text:        "the 'age' request property type/format changed from 'integer'/'int32' to 'string'/'string'",
		Operation:   "POST",
		Path:        "/pets",
		Source:      "../data/checker/request_property_type_changed_revision.yaml",
		OperationId: "addPet",
	}, errs[0])
}

// CL: changing request property format
func TestRequestPropertyFormatChangedCheck(t *testing.T) {
	s1, err := open("../data/checker/request_property_type_changed_base.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/request_property_type_changed_base.yaml")
	require.NoError(t, err)

	s2.Spec.Paths["/pets"].Post.RequestBody.Value.Content["application/json"].Schema.Value.Properties["age"].Value.Format = "uuid"

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)

	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.RequestPropertyTypeChangedCheck), d, osm, checker.ERR)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{
		Id:          "request-property-type-changed",
		Level:       checker.ERR,
		Text:        "the 'age' request property type/format changed from 'integer'/'int32' to 'integer'/'uuid'",
		Operation:   "POST",
		Path:        "/pets",
		Source:      "../data/checker/request_property_type_changed_base.yaml",
		OperationId: "addPet",
	}, errs[0])
}

func TestRequestPropertyFormatChangedCheckNonBreaking(t *testing.T) {
	s1, err := open("../data/checker/request_property_type_changed_base.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/request_property_type_changed_base.yaml")
	require.NoError(t, err)

	s1.Spec.Paths["/pets"].Post.RequestBody.Value.Content["application/json"].Schema.Value.Properties["age"].Value.Type = "integer"
	s2.Spec.Paths["/pets"].Post.RequestBody.Value.Content["application/json"].Schema.Value.Properties["age"].Value.Type = "number"

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)

	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.RequestPropertyTypeChangedCheck), d, osm, checker.INFO)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{
		Id:          "request-property-type-changed",
		Level:       checker.INFO,
		Text:        "the 'age' request property type/format changed from 'integer'/'int32' to 'number'/'int32'",
		Operation:   "POST",
		Path:        "/pets",
		Source:      "../data/checker/request_property_type_changed_base.yaml",
		OperationId: "addPet",
	}, errs[0])
}
