package checker_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/checker"
	"github.com/tufin/oasdiff/diff"
)

// CL: changing required response property to write-only
func TestResponseRequiredPropertyBecameWriteOnly(t *testing.T) {
	s1, err := open("../data/checker/response_required_property_write_only_read_only_base.yaml")
	require.NoError(t, err)

	s2, err := open("../data/checker/response_required_property_write_only_read_only_base.yaml")
	require.NoError(t, err)

	s2.Spec.Components.Schemas["GroupView"].Value.Properties["data"].Value.Properties["name"].Value.WriteOnly = true
	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.ResponseRequiredPropertyWriteOnlyReadOnlyCheck), d, osm, checker.INFO)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{

		Id:          "response-required-property-became-write-only",
		Text:        "the response required property 'data/name' became write-only for the status '200'",
		Comment:     "",
		Level:       checker.INFO,
		Operation:   "POST",
		Path:        "/api/v1.0/groups",
		Source:      "../data/checker/response_required_property_write_only_read_only_base.yaml",
		OperationId: "createOneGroup",
	}, errs[0])
}

// CL: changing required response property to not write-only
func TestResponseRequiredPropertyBecameNotWriteOnly(t *testing.T) {
	s1, err := open("../data/checker/response_required_property_write_only_read_only_base.yaml")
	require.NoError(t, err)

	s2, err := open("../data/checker/response_required_property_write_only_read_only_base.yaml")
	require.NoError(t, err)

	s2.Spec.Components.Schemas["GroupView"].Value.Properties["data"].Value.Properties["writeOnlyName"].Value.WriteOnly = false
	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.ResponseRequiredPropertyWriteOnlyReadOnlyCheck), d, osm, checker.INFO)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{

		Id:          checker.ResponseRequiredPropertyBecameNonWriteOnlyId,
		Text:        "the response required property 'data/writeOnlyName' became not write-only for the status '200'",
		Comment:     "It is valid only if the property was always returned before the specification has been changed",
		Level:       checker.WARN,
		Operation:   "POST",
		Path:        "/api/v1.0/groups",
		Source:      "../data/checker/response_required_property_write_only_read_only_base.yaml",
		OperationId: "createOneGroup",
	}, errs[0])
}

// CL: changing required response property to read-only
func TestResponseRequiredPropertyBecameReadOnly(t *testing.T) {
	s1, err := open("../data/checker/response_required_property_write_only_read_only_base.yaml")
	require.NoError(t, err)

	s2, err := open("../data/checker/response_required_property_write_only_read_only_base.yaml")
	require.NoError(t, err)

	s1.Spec.Components.Schemas["GroupView"].Value.Properties["data"].Value.Properties["id"].Value.ReadOnly = false
	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.ResponseRequiredPropertyWriteOnlyReadOnlyCheck), d, osm, checker.INFO)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{

		Id:          checker.ResponseRequiredPropertyBecameReadOnlyId,
		Text:        "the response required property 'data/id' became read-only for the status '200'",
		Comment:     "",
		Level:       checker.INFO,
		Operation:   "POST",
		Path:        "/api/v1.0/groups",
		Source:      "../data/checker/response_required_property_write_only_read_only_base.yaml",
		OperationId: "createOneGroup",
	}, errs[0])
}

// CL: changing required response property to not read-only
func TestResponseRequiredPropertyBecameNonReadOnly(t *testing.T) {
	s1, err := open("../data/checker/response_required_property_write_only_read_only_base.yaml")
	require.NoError(t, err)

	s2, err := open("../data/checker/response_required_property_write_only_read_only_base.yaml")
	require.NoError(t, err)

	s2.Spec.Components.Schemas["GroupView"].Value.Properties["data"].Value.Properties["id"].Value.ReadOnly = false

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.ResponseRequiredPropertyWriteOnlyReadOnlyCheck), d, osm, checker.INFO)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{

		Id:          checker.ResponseRequiredPropertyBecameNonReadOnlyId,
		Text:        "the response required property 'data/id' became not read-only for the status '200'",
		Comment:     "",
		Level:       checker.INFO,
		Operation:   "POST",
		Path:        "/api/v1.0/groups",
		Source:      "../data/checker/response_required_property_write_only_read_only_base.yaml",
		OperationId: "createOneGroup",
	}, errs[0])
}
