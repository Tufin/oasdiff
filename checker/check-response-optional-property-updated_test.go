package checker_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/checker"
	"github.com/tufin/oasdiff/diff"
)

// CL: Removing an optional write-only property from a response
func TestResponseOptionalPropertyUpdatedCheck(t *testing.T) {
	s1, err := open("../data/checker/response_optional_property_removed_base.yaml")
	require.Empty(t, err)
	s2, err := open("../data/checker/response_optional_property_removed_revision.yaml")
	require.Empty(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.ResponseOptionalPropertyUpdatedCheck), d, osm, checker.INFO)
	require.NotEmpty(t, errs)
	require.Equal(t, checker.BackwardCompatibilityErrors{
		{
			Id:          "response-optional-property-removed",
			Text:        "removed the optional property 'data/id' from the response with the '200' status",
			Comment:     "",
			Level:       checker.WARN,
			Operation:   "POST",
			Path:        "/api/v1.0/groups",
			Source:      "../data/checker/response_optional_property_removed_revision.yaml",
			OperationId: "createOneGroup",
		},
	}, errs)
}

// CL: Adding an optional write-only property to a response
func TestResponseOptionalPropertyAddedCheck(t *testing.T) {
	s1, err := open("../data/checker/response_optional_property_removed_revision.yaml")
	require.Empty(t, err)
	s2, err := open("../data/checker/response_optional_property_removed_base.yaml")
	require.Empty(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.ResponseOptionalPropertyUpdatedCheck), d, osm, checker.INFO)

	require.NotEmpty(t, errs)
	require.Equal(t, checker.BackwardCompatibilityErrors{
		{
			Id:          "response-optional-property-added",
			Text:        "added the optional property 'data/id' to the response with the '200' status",
			Comment:     "",
			Level:       checker.INFO,
			Operation:   "POST",
			Path:        "/api/v1.0/groups",
			Source:      "../data/checker/response_optional_property_removed_base.yaml",
			OperationId: "createOneGroup",
		},
	}, errs)
}

// CL: Removing an optional write-only property from a response
func TestResponseOptionalWriteOnlyPropertyRemovedCheck(t *testing.T) {
	s1, err := open("../data/checker/response_optional_property_removed_base.yaml")
	require.Empty(t, err)
	s2, err := open("../data/checker/response_optional_property_removed_revision.yaml")
	require.Empty(t, err)

	s1.Spec.Paths["/api/v1.0/groups"].Post.Responses["200"].Value.Content["application/json"].Schema.Value.Properties["data"].Value.Properties["id"].Value.WriteOnly = true
	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.ResponseOptionalPropertyUpdatedCheck), d, osm, checker.INFO)
	require.NotEmpty(t, errs)
	require.Equal(t, checker.BackwardCompatibilityErrors{
		{
			Id:          "response-optional-write-only-property-removed",
			Text:        "removed the optional write-only property 'data/id' from the response with the '200' status",
			Level:       checker.INFO,
			Operation:   "POST",
			Path:        "/api/v1.0/groups",
			Source:      "../data/checker/response_optional_property_removed_revision.yaml",
			OperationId: "createOneGroup",
		},
	}, errs)
}

// CL: Adding an optional write-only property to a response
func TestResponseOptionalWriteOnlyPropertyAddedCheck(t *testing.T) {
	s1, err := open("../data/checker/response_optional_property_removed_revision.yaml")
	require.Empty(t, err)
	s2, err := open("../data/checker/response_optional_property_removed_base.yaml")
	require.Empty(t, err)

	s2.Spec.Paths["/api/v1.0/groups"].Post.Responses["200"].Value.Content["application/json"].Schema.Value.Properties["data"].Value.Properties["id"].Value.WriteOnly = true
	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.ResponseOptionalPropertyUpdatedCheck), d, osm, checker.INFO)

	require.NotEmpty(t, errs)
	require.Equal(t, checker.BackwardCompatibilityErrors{
		{
			Id:          "response-optional-write-only-property-added",
			Text:        "added the optional write-only property 'data/id' to the response with the '200' status",
			Level:       checker.INFO,
			Operation:   "POST",
			Path:        "/api/v1.0/groups",
			Source:      "../data/checker/response_optional_property_removed_base.yaml",
			OperationId: "createOneGroup",
		},
	}, errs)
}
