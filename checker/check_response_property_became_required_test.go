package checker_test

import (
	"testing"

	"github.com/oasdiff/oasdiff/checker"
	"github.com/oasdiff/oasdiff/diff"
	"github.com/oasdiff/oasdiff/load"
	"github.com/stretchr/testify/require"
)

// CL: changing optional response property to required
func TestResponsePropertyBecameRequiredlCheck(t *testing.T) {
	s1, err := open("../data/checker/response_property_became_optional_revision.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/response_property_became_optional_base.yaml")
	require.NoError(t, err)
	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.ResponsePropertyBecameRequiredCheck), d, osm, checker.INFO)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{
		Id:          checker.ResponsePropertyBecameRequiredId,
		Args:        []any{"data/name", "200"},
		Level:       checker.INFO,
		Operation:   "POST",
		Path:        "/api/v1.0/groups",
		Source:      load.NewSource("../data/checker/response_property_became_optional_base.yaml"),
		OperationId: "createOneGroup",
	}, errs[0])
}

// CL: changing optional response write-only property to required
func TestResponseWriteOnlyPropertyBecameRequiredCheck(t *testing.T) {
	s1, err := open("../data/checker/response_property_became_optional_revision.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/response_property_became_optional_base.yaml")
	require.NoError(t, err)
	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)

	s1.Spec.Components.Schemas["GroupView"].Value.Properties["data"].Value.Properties["name"].Value.WriteOnly = true

	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.ResponsePropertyBecameRequiredCheck), d, osm, checker.INFO)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{
		Id:          checker.ResponseWriteOnlyPropertyBecameRequiredId,
		Args:        []any{"data/name", "200"},
		Level:       checker.INFO,
		Operation:   "POST",
		Path:        "/api/v1.0/groups",
		Source:      load.NewSource("../data/checker/response_property_became_optional_base.yaml"),
		OperationId: "createOneGroup",
	}, errs[0])
}
