package checker_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/checker"
	"github.com/tufin/oasdiff/diff"
)

// CL: changing required response property to optional
func TestResponsePropertyBecameOptionalCheck(t *testing.T) {
	s1, err := open("../data/checker/response_property_became_optional_base.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/response_property_became_optional_revision.yaml")
	require.NoError(t, err)
	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.ResponsePropertyBecameOptionalCheck), d, osm, checker.ERR)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{
		Id:          "response-property-became-optional",
		Text:        "the response property 'data/name' became optional for the status '200'",
		Comment:     "",
		Level:       checker.ERR,
		Operation:   "POST",
		Path:        "/api/v1.0/groups",
		Source:      "../data/checker/response_property_became_optional_revision.yaml",
		OperationId: "createOneGroup",
	}, errs[0])
}

// CL: changing write-only required response property to optional
func TestResponseWriteOnlyPropertyBecameOptionalCheck(t *testing.T) {
	s1, err := open("../data/checker/response_property_became_optional_base.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/response_property_became_optional_revision.yaml")
	require.NoError(t, err)
	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)

	s1.Spec.Components.Schemas["GroupView"].Value.Properties["data"].Value.Properties["name"].Value.WriteOnly = true

	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.ResponsePropertyBecameOptionalCheck), d, osm, checker.INFO)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{
		Id:          "response-write-only-property-became-optional",
		Text:        "the response write-only property 'data/name' became optional for the status '200'",
		Comment:     "",
		Level:       checker.INFO,
		Operation:   "POST",
		Path:        "/api/v1.0/groups",
		Source:      "../data/checker/response_property_became_optional_revision.yaml",
		OperationId: "createOneGroup",
	}, errs[0])
}
