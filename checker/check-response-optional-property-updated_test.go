package checker_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/checker"
	"github.com/tufin/oasdiff/diff"
)

// CL: removing an optional write-only property from a response
func TestResponseOptionalPropertyUpdatedCheck(t *testing.T) {
	s1, err := open("../data/checker/response_optional_property_removed_base.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/response_optional_property_removed_revision.yaml")
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.ResponseOptionalPropertyUpdatedCheck), d, osm, checker.INFO)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{
		Id:          checker.ResponseOptionalPropertyRemovedId,
		Text:        "removed the optional property 'data/id' from the response with the '200' status",
		Comment:     "",
		Level:       checker.WARN,
		Operation:   "POST",
		Path:        "/api/v1.0/groups",
		Source:      "../data/checker/response_optional_property_removed_revision.yaml",
		OperationId: "createOneGroup",
	}, errs[0])
}

// CL: adding an optional write-only property to a response
func TestResponseOptionalPropertyAddedCheck(t *testing.T) {
	s1, err := open("../data/checker/response_optional_property_removed_revision.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/response_optional_property_removed_base.yaml")
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.ResponseOptionalPropertyUpdatedCheck), d, osm, checker.INFO)

	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{
		Id:          checker.ResponseOptionalPropertyAddedId,
		Text:        "added the optional property 'data/id' to the response with the '200' status",
		Comment:     "",
		Level:       checker.INFO,
		Operation:   "POST",
		Path:        "/api/v1.0/groups",
		Source:      "../data/checker/response_optional_property_removed_base.yaml",
		OperationId: "createOneGroup",
	}, errs[0])
}

// CL: removing an optional write-only property from a response
func TestResponseOptionalWriteOnlyPropertyRemovedCheck(t *testing.T) {
	s1, err := open("../data/checker/response_optional_property_removed_base.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/response_optional_property_removed_revision.yaml")
	require.NoError(t, err)

	s1.Spec.Paths["/api/v1.0/groups"].Post.Responses["200"].Value.Content["application/json"].Schema.Value.Properties["data"].Value.Properties["id"].Value.WriteOnly = true
	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.ResponseOptionalPropertyUpdatedCheck), d, osm, checker.INFO)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{
		Id:          checker.ResponseOptionalWriteOnlyPropertyRemovedId,
		Text:        "removed the optional write-only property 'data/id' from the response with the '200' status",
		Level:       checker.INFO,
		Operation:   "POST",
		Path:        "/api/v1.0/groups",
		Source:      "../data/checker/response_optional_property_removed_revision.yaml",
		OperationId: "createOneGroup",
	}, errs[0])
}

// CL: adding an optional write-only property to a response
func TestResponseOptionalWriteOnlyPropertyAddedCheck(t *testing.T) {
	s1, err := open("../data/checker/response_optional_property_removed_revision.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/response_optional_property_removed_base.yaml")
	require.NoError(t, err)

	s2.Spec.Paths["/api/v1.0/groups"].Post.Responses["200"].Value.Content["application/json"].Schema.Value.Properties["data"].Value.Properties["id"].Value.WriteOnly = true
	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.ResponseOptionalPropertyUpdatedCheck), d, osm, checker.INFO)

	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{
		Id:          checker.ResponseOptionalWriteOnlyPropertyAddedId,
		Text:        "added the optional write-only property 'data/id' to the response with the '200' status",
		Level:       checker.INFO,
		Operation:   "POST",
		Path:        "/api/v1.0/groups",
		Source:      "../data/checker/response_optional_property_removed_base.yaml",
		OperationId: "createOneGroup",
	}, errs[0])
}
