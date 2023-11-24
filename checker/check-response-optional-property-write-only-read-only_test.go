package checker_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/checker"
	"github.com/tufin/oasdiff/diff"
)

// CL: changing optional response property to write-only
func TestResponseOptionalPropertyBecameWriteOnly(t *testing.T) {
	s1, err := open("../data/checker/response_optional_property_write_only_read_only_base.yaml")
	require.NoError(t, err)

	s2, err := open("../data/checker/response_optional_property_write_only_read_only_base.yaml")
	require.NoError(t, err)

	s2.Spec.Components.Schemas["GroupView"].Value.Properties["data"].Value.Properties["name"].Value.WriteOnly = true
	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.ResponseOptionalPropertyWriteOnlyReadOnlyCheck), d, osm, checker.INFO)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{

		Id:          checker.ResponseOptionalPropertyBecameWriteOnlyId,
		Text:        "the response optional property 'data/name' became write-only for the status '200'",
		Args:        []any{"data/name", "200"},
		Comment:     "",
		Level:       checker.INFO,
		Operation:   "POST",
		Path:        "/api/v1.0/groups",
		Source:      "../data/checker/response_optional_property_write_only_read_only_base.yaml",
		OperationId: "createOneGroup",
	}, errs[0])
}

// CL: changing optional response property to not write-only
func TestResponseOptionalPropertyBecameNotWriteOnly(t *testing.T) {
	s1, err := open("../data/checker/response_optional_property_write_only_read_only_base.yaml")
	require.NoError(t, err)

	s2, err := open("../data/checker/response_optional_property_write_only_read_only_base.yaml")
	require.NoError(t, err)

	s2.Spec.Components.Schemas["GroupView"].Value.Properties["data"].Value.Properties["writeOnlyName"].Value.WriteOnly = false
	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.ResponseOptionalPropertyWriteOnlyReadOnlyCheck), d, osm, checker.INFO)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{

		Id:          checker.ResponseOptionalPropertyBecameNonWriteOnlyId,
		Text:        "the response optional property 'data/writeOnlyName' became not write-only for the status '200'",
		Args:        []any{"data/writeOnlyName", "200"},
		Level:       checker.INFO,
		Operation:   "POST",
		Path:        "/api/v1.0/groups",
		Source:      "../data/checker/response_optional_property_write_only_read_only_base.yaml",
		OperationId: "createOneGroup",
	}, errs[0])
}

// CL: changing optional response property to read-only
func TestResponseOptionalPropertyBecameReadOnly(t *testing.T) {
	s1, err := open("../data/checker/response_optional_property_write_only_read_only_base.yaml")
	require.NoError(t, err)

	s2, err := open("../data/checker/response_optional_property_write_only_read_only_base.yaml")
	require.NoError(t, err)

	s1.Spec.Components.Schemas["GroupView"].Value.Properties["data"].Value.Properties["id"].Value.ReadOnly = false
	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.ResponseOptionalPropertyWriteOnlyReadOnlyCheck), d, osm, checker.INFO)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{

		Id:          checker.ResponseOptionalPropertyBecameReadOnlyId,
		Text:        "the response optional property 'data/id' became read-only for the status '200'",
		Args:        []any{"data/id", "200"},
		Comment:     "",
		Level:       checker.INFO,
		Operation:   "POST",
		Path:        "/api/v1.0/groups",
		Source:      "../data/checker/response_optional_property_write_only_read_only_base.yaml",
		OperationId: "createOneGroup",
	}, errs[0])
}

// CL: changing optional response property to not read-only
func TestResponseOptionalPropertyBecameNonReadOnly(t *testing.T) {
	s1, err := open("../data/checker/response_optional_property_write_only_read_only_base.yaml")
	require.NoError(t, err)

	s2, err := open("../data/checker/response_optional_property_write_only_read_only_base.yaml")
	require.NoError(t, err)

	s2.Spec.Components.Schemas["GroupView"].Value.Properties["data"].Value.Properties["id"].Value.ReadOnly = false

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.ResponseOptionalPropertyWriteOnlyReadOnlyCheck), d, osm, checker.INFO)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{

		Id:          checker.ResponseOptionalPropertyBecameNonReadOnlyId,
		Text:        "the response optional property 'data/id' became not read-only for the status '200'",
		Args:        []any{"data/id", "200"},
		Comment:     "",
		Level:       checker.INFO,
		Operation:   "POST",
		Path:        "/api/v1.0/groups",
		Source:      "../data/checker/response_optional_property_write_only_read_only_base.yaml",
		OperationId: "createOneGroup",
	}, errs[0])
}
