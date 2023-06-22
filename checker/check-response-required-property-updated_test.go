package checker_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/checker"
	"github.com/tufin/oasdiff/diff"
)

// CL: Adding a required property to response body is detected
func TestResponseRequiredPropertyAdded(t *testing.T) {
	s1, _ := open("../data/checker/response_required_property_added_base.yaml")
	s2, err := open("../data/checker/response_required_property_added_revision.yaml")
	require.Empty(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.ResponseRequiredPropertyUpdatedCheck), d, osm, checker.INFO)
	require.NotEmpty(t, errs)
	require.Equal(t, checker.BackwardCompatibilityErrors{
		{
			Id:          "response-required-property-added",
			Text:        "added the required property 'data/new' to the response with the '200' status",
			Comment:     "",
			Level:       checker.INFO,
			Operation:   "POST",
			Path:        "/api/v1.0/groups",
			Source:      "../data/checker/response_required_property_added_revision.yaml",
			OperationId: "createOneGroup",
		}}, errs)
}

// CL: Removing an existent property that was required in response body is detected
func TestResponseRequiredPropertyRemoved(t *testing.T) {
	s1, _ := open("../data/checker/response_required_property_added_revision.yaml")
	s2, err := open("../data/checker/response_required_property_added_base.yaml")
	require.Empty(t, err)

	s2.Spec.Components.Schemas["GroupView"].Value.Properties["data"].Value.Required = []string{"name", "id"}
	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.ResponseRequiredPropertyUpdatedCheck), d, osm, checker.INFO)
	require.NotEmpty(t, errs)
	require.Equal(t, checker.BackwardCompatibilityErrors{
		{
			Id:          "response-required-property-removed",
			Text:        "removed the required property 'data/new' from the response with the '200' status",
			Comment:     "",
			Level:       checker.ERR,
			Operation:   "POST",
			Path:        "/api/v1.0/groups",
			Source:      "../data/checker/response_required_property_added_base.yaml",
			OperationId: "createOneGroup",
		}}, errs)
}
