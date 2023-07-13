package checker_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/checker"
	"github.com/tufin/oasdiff/diff"
)

// CL: Required response property became write-only
func TestResponseRequiredPropertyBecameWriteOnly(t *testing.T) {
	s1, _ := open("../data/checker/response_required_property_write_only_read_only_base.yaml")
	s2, err := open("../data/checker/response_required_property_write_only_read_only_base.yaml")
	require.Empty(t, err)

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

// CL: Required response property became not write-only
func TestResponseRequiredPropertyBecameNotWriteOnly(t *testing.T) {
	s1, _ := open("../data/checker/response_required_property_write_only_read_only_base.yaml")
	s2, err := open("../data/checker/response_required_property_write_only_read_only_base.yaml")
	require.Empty(t, err)

	s2.Spec.Components.Schemas["GroupView"].Value.Properties["data"].Value.Properties["writeOnlyName"].Value.WriteOnly = false
	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.ResponseRequiredPropertyWriteOnlyReadOnlyCheck), d, osm, checker.INFO)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{

		Id:          "response-required-property-became-not-write-only",
		Text:        "the response required property 'data/writeOnlyName' became not write-only for the status '200'",
		Comment:     "It is valid only if the property was always returned before the specification has been changed",
		Level:       checker.WARN,
		Operation:   "POST",
		Path:        "/api/v1.0/groups",
		Source:      "../data/checker/response_required_property_write_only_read_only_base.yaml",
		OperationId: "createOneGroup",
	}, errs[0])
}

// CL: Required response property became read-only
func TestResponseRequiredPropertyBecameReadOnly(t *testing.T) {
	s1, _ := open("../data/checker/response_required_property_write_only_read_only_base.yaml")
	s2, err := open("../data/checker/response_required_property_write_only_read_only_base.yaml")
	require.Empty(t, err)

	s1.Spec.Components.Schemas["GroupView"].Value.Properties["data"].Value.Properties["id"].Value.ReadOnly = false
	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.ResponseRequiredPropertyWriteOnlyReadOnlyCheck), d, osm, checker.INFO)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{

		Id:          "response-required-property-became-read-only",
		Text:        "the response required property 'data/id' became read-only for the status '200'",
		Comment:     "",
		Level:       checker.INFO,
		Operation:   "POST",
		Path:        "/api/v1.0/groups",
		Source:      "../data/checker/response_required_property_write_only_read_only_base.yaml",
		OperationId: "createOneGroup",
	}, errs[0])
}

// CL: Required response property became not read-only
func TestResponseRequiredPropertyBecameNonReadOnly(t *testing.T) {
	s1, _ := open("../data/checker/response_required_property_write_only_read_only_base.yaml")
	s2, err := open("../data/checker/response_required_property_write_only_read_only_base.yaml")
	require.Empty(t, err)

	s2.Spec.Components.Schemas["GroupView"].Value.Properties["data"].Value.Properties["id"].Value.ReadOnly = false

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.ResponseRequiredPropertyWriteOnlyReadOnlyCheck), d, osm, checker.INFO)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{

		Id:          "response-required-property-became-not-read-only",
		Text:        "the response required property 'data/id' became not read-only for the status '200'",
		Comment:     "",
		Level:       checker.INFO,
		Operation:   "POST",
		Path:        "/api/v1.0/groups",
		Source:      "../data/checker/response_required_property_write_only_read_only_base.yaml",
		OperationId: "createOneGroup",
	}, errs[0])
}
