package checker_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/checker"
	"github.com/tufin/oasdiff/diff"
)

// CL: changing request property type/format
func TestRequestPropertyTypeChanged(t *testing.T) {
	s1, err := open("../data/checker/request_property_type_changed_base.yaml")
	require.NoError(t, err)

	s2, err := open("../data/checker/request_property_type_changed_base.yaml")
	require.NoError(t, err)

	s2.Spec.Paths["/example"].Post.RequestBody.Value.Content["application/json"].Schema.Value.Properties["amount"].Value.Type = "string"

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.RequestPropertyTypeChangedCheck), d, osm, checker.ERR)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{
		Id:          "request-property-type-changed",
		Text:        "the 'amount' request property type/format changed from 'integer'/'int32' to 'string'/'int32'",
		Level:       checker.ERR,
		Operation:   "POST",
		Path:        "/example",
		Source:      "../data/checker/request_property_type_changed_base.yaml",
		OperationId: "createExample",
	}, errs[0])
}

// CL: changing request body type/format
func TestRequestBodyTypeChanged(t *testing.T) {
	s1, err := open("../data/checker/request_property_type_changed_base.yaml")
	require.NoError(t, err)

	s2, err := open("../data/checker/request_property_type_changed_base.yaml")
	require.NoError(t, err)

	s2.Spec.Paths["/example"].Post.RequestBody.Value.Content["application/json"].Schema.Value.Type = "string"

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.RequestPropertyTypeChangedCheck), d, osm, checker.ERR)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{
		Id:          "request-body-type-changed",
		Text:        "the request's body type/format changed from 'object'/'none' to 'string'/'none'",
		Level:       checker.ERR,
		Operation:   "POST",
		Path:        "/example",
		Source:      "../data/checker/request_property_type_changed_base.yaml",
		OperationId: "createExample",
	}, errs[0])
}
