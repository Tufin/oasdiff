package checker_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/checker"
	"github.com/tufin/oasdiff/diff"
	"github.com/tufin/oasdiff/load"
)

// CL: adding a required property to response body is detected
func TestResponseRequiredPropertyAdded(t *testing.T) {
	s1, err := open("../data/checker/response_required_property_added_base.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/response_required_property_added_revision.yaml")
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.ResponseRequiredPropertyUpdatedCheck), d, osm, checker.INFO)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{

		Id:          checker.ResponseRequiredPropertyAddedId,
		Args:        []any{"data/new", "200"},
		Level:       checker.INFO,
		Operation:   "POST",
		Path:        "/api/v1.0/groups",
		Source:      load.NewSource("../data/checker/response_required_property_added_revision.yaml"),
		OperationId: "createOneGroup",
	}, errs[0])
}

// CL: removing an existent property that was required in response body is detected
func TestResponseRequiredPropertyRemoved(t *testing.T) {
	s1, err := open("../data/checker/response_required_property_added_revision.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/response_required_property_added_base.yaml")
	require.NoError(t, err)

	s2.Spec.Components.Schemas["GroupView"].Value.Properties["data"].Value.Required = []string{"name", "id"}
	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.ResponseRequiredPropertyUpdatedCheck), d, osm, checker.INFO)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{
		Id:          checker.ResponseRequiredPropertyRemovedId,
		Args:        []any{"data/new", "200"},
		Level:       checker.ERR,
		Operation:   "POST",
		Path:        "/api/v1.0/groups",
		Source:      load.NewSource("../data/checker/response_required_property_added_base.yaml"),
		OperationId: "createOneGroup",
	}, errs[0])
}

// CL: adding a required write-only property to response body is detected
func TestResponseRequiredWriteOnlyPropertyAdded(t *testing.T) {
	s1, err := open("../data/checker/response_required_property_added_base.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/response_required_property_added_revision.yaml")
	require.NoError(t, err)

	s2.Spec.Components.Schemas["GroupView"].Value.Properties["data"].Value.Properties["new"].Value.WriteOnly = true

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.ResponseRequiredPropertyUpdatedCheck), d, osm, checker.INFO)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{

		Id:          checker.ResponseRequiredWriteOnlyPropertyAddedId,
		Args:        []any{"data/new", "200"},
		Level:       checker.INFO,
		Operation:   "POST",
		Path:        "/api/v1.0/groups",
		Source:      load.NewSource("../data/checker/response_required_property_added_revision.yaml"),
		OperationId: "createOneGroup",
	}, errs[0])
}

// CL: removing a required write-only property that was required in response body is detected
func TestResponseRequiredWriteOnlyPropertyRemoved(t *testing.T) {
	s1, err := open("../data/checker/response_required_property_added_revision.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/response_required_property_added_base.yaml")
	require.NoError(t, err)

	s1.Spec.Components.Schemas["GroupView"].Value.Properties["data"].Value.Properties["new"].Value.WriteOnly = true
	s2.Spec.Components.Schemas["GroupView"].Value.Properties["data"].Value.Required = []string{"name", "id"}
	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.ResponseRequiredPropertyUpdatedCheck), d, osm, checker.INFO)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{
		Id:          checker.ResponseRequiredWriteOnlyPropertyRemovedId,
		Args:        []any{"data/new", "200"},
		Level:       checker.INFO,
		Operation:   "POST",
		Path:        "/api/v1.0/groups",
		Source:      load.NewSource("../data/checker/response_required_property_added_base.yaml"),
		OperationId: "createOneGroup",
	}, errs[0])
}
