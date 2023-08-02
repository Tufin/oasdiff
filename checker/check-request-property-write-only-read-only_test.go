package checker_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/checker"
	"github.com/tufin/oasdiff/diff"
)

// CL: changing optional request property to write-only
func TestRequestOptionalPropertyBecameWriteOnly(t *testing.T) {
	s1, err := open("../data/checker/request_optional_property_write_only_read_only_base.yaml")
	require.NoError(t, err)

	s2, err := open("../data/checker/request_optional_property_write_only_read_only_base.yaml")
	require.NoError(t, err)

	s2.Spec.Paths["/api/v1.0/groups"].Post.RequestBody.Value.Content["application/json"].Schema.Value.Properties["name"].Value.WriteOnly = true

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.RequestPropertyWriteOnlyReadOnlyCheck), d, osm, checker.INFO)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{
		Id:          "request-optional-property-became-write-only",
		Text:        "the request optional property 'name' became write-only",
		Operation:   "POST",
		Path:        "/api/v1.0/groups",
		Level:       checker.INFO,
		Source:      "../data/checker/request_optional_property_write_only_read_only_base.yaml",
		OperationId: "createOneGroup",
	}, errs[0])
}

// CL: changing optional request property to not write-only
func TestRequestOptionalPropertyBecameNotWriteOnly(t *testing.T) {
	s1, err := open("../data/checker/request_optional_property_write_only_read_only_base.yaml")
	require.NoError(t, err)

	s2, err := open("../data/checker/request_optional_property_write_only_read_only_base.yaml")
	require.NoError(t, err)

	s1.Spec.Paths["/api/v1.0/groups"].Post.RequestBody.Value.Content["application/json"].Schema.Value.Properties["name"].Value.WriteOnly = true

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.RequestPropertyWriteOnlyReadOnlyCheck), d, osm, checker.INFO)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{
		Id:          "request-optional-property-became-not-write-only",
		Text:        "the request optional property 'name' became not write-only",
		Level:       checker.INFO,
		Operation:   "POST",
		Path:        "/api/v1.0/groups",
		Source:      "../data/checker/request_optional_property_write_only_read_only_base.yaml",
		OperationId: "createOneGroup",
	}, errs[0])
}

// CL: changing optional request property to read-only
func TestRequestOptionalPropertyBecameReadOnly(t *testing.T) {
	s1, err := open("../data/checker/request_optional_property_write_only_read_only_base.yaml")
	require.NoError(t, err)

	s2, err := open("../data/checker/request_optional_property_write_only_read_only_base.yaml")
	require.NoError(t, err)

	s2.Spec.Paths["/api/v1.0/groups"].Post.RequestBody.Value.Content["application/json"].Schema.Value.Properties["name"].Value.ReadOnly = true

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.RequestPropertyWriteOnlyReadOnlyCheck), d, osm, checker.INFO)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{
		Id:          "request-optional-property-became-read-only",
		Text:        "the request optional property 'name' became read-only",
		Level:       checker.INFO,
		Operation:   "POST",
		Path:        "/api/v1.0/groups",
		Source:      "../data/checker/request_optional_property_write_only_read_only_base.yaml",
		OperationId: "createOneGroup",
	}, errs[0])
}

// CL: changing optional request property to not read-only
func TestRequestOptionalPropertyBecameNonReadOnly(t *testing.T) {
	s1, err := open("../data/checker/request_optional_property_write_only_read_only_base.yaml")
	require.NoError(t, err)

	s2, err := open("../data/checker/request_optional_property_write_only_read_only_base.yaml")
	require.NoError(t, err)

	s1.Spec.Paths["/api/v1.0/groups"].Post.RequestBody.Value.Content["application/json"].Schema.Value.Properties["name"].Value.ReadOnly = true

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.RequestPropertyWriteOnlyReadOnlyCheck), d, osm, checker.INFO)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{
		Id:          "request-optional-property-became-not-read-only",
		Text:        "the request optional property 'name' became not read-only",
		Level:       checker.INFO,
		Operation:   "POST",
		Path:        "/api/v1.0/groups",
		Source:      "../data/checker/request_optional_property_write_only_read_only_base.yaml",
		OperationId: "createOneGroup",
	}, errs[0])
}

// CL: changing required request property to write-only
func TestRequestRequiredPropertyBecameWriteOnly(t *testing.T) {
	s1, err := open("../data/checker/request_optional_property_write_only_read_only_base.yaml")
	require.NoError(t, err)

	s2, err := open("../data/checker/request_optional_property_write_only_read_only_base.yaml")
	require.NoError(t, err)

	s2.Spec.Paths["/api/v1.0/groups"].Post.RequestBody.Value.Content["application/json"].Schema.Value.Properties["id"].Value.WriteOnly = true

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.RequestPropertyWriteOnlyReadOnlyCheck), d, osm, checker.INFO)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{
		Id:          "request-required-property-became-write-only",
		Text:        "the request required property 'id' became write-only",
		Operation:   "POST",
		Path:        "/api/v1.0/groups",
		Level:       checker.INFO,
		Source:      "../data/checker/request_optional_property_write_only_read_only_base.yaml",
		OperationId: "createOneGroup",
	}, errs[0])
}

// CL: changing required request property to not write-only
func TestRequestRequiredPropertyBecameNotWriteOnly(t *testing.T) {
	s1, err := open("../data/checker/request_optional_property_write_only_read_only_base.yaml")
	require.NoError(t, err)

	s2, err := open("../data/checker/request_optional_property_write_only_read_only_base.yaml")
	require.NoError(t, err)

	s1.Spec.Paths["/api/v1.0/groups"].Post.RequestBody.Value.Content["application/json"].Schema.Value.Properties["id"].Value.WriteOnly = true

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.RequestPropertyWriteOnlyReadOnlyCheck), d, osm, checker.INFO)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{
		Id:          "request-required-property-became-not-write-only",
		Text:        "the request required property 'id' became not write-only",
		Level:       checker.INFO,
		Operation:   "POST",
		Path:        "/api/v1.0/groups",
		Source:      "../data/checker/request_optional_property_write_only_read_only_base.yaml",
		OperationId: "createOneGroup",
	}, errs[0])
}

// CL: changing required request property to read-only
func TestRequestRequiredPropertyBecameReadOnly(t *testing.T) {
	s1, err := open("../data/checker/request_optional_property_write_only_read_only_base.yaml")
	require.NoError(t, err)

	s2, err := open("../data/checker/request_optional_property_write_only_read_only_base.yaml")
	require.NoError(t, err)

	s2.Spec.Paths["/api/v1.0/groups"].Post.RequestBody.Value.Content["application/json"].Schema.Value.Properties["id"].Value.ReadOnly = true

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.RequestPropertyWriteOnlyReadOnlyCheck), d, osm, checker.INFO)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{
		Id:          "request-required-property-became-read-only",
		Text:        "the request required property 'id' became read-only",
		Level:       checker.INFO,
		Operation:   "POST",
		Path:        "/api/v1.0/groups",
		Source:      "../data/checker/request_optional_property_write_only_read_only_base.yaml",
		OperationId: "createOneGroup",
	}, errs[0])
}

// CL: changing required request property to not read-only
func TestRequestRequiredPropertyBecameNonReadOnly(t *testing.T) {
	s1, err := open("../data/checker/request_optional_property_write_only_read_only_base.yaml")
	require.NoError(t, err)

	s2, err := open("../data/checker/request_optional_property_write_only_read_only_base.yaml")
	require.NoError(t, err)

	s1.Spec.Paths["/api/v1.0/groups"].Post.RequestBody.Value.Content["application/json"].Schema.Value.Properties["id"].Value.ReadOnly = true

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.RequestPropertyWriteOnlyReadOnlyCheck), d, osm, checker.INFO)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{
		Id:          "request-required-property-became-not-read-only",
		Text:        "the request required property 'id' became not read-only",
		Level:       checker.INFO,
		Operation:   "POST",
		Path:        "/api/v1.0/groups",
		Source:      "../data/checker/request_optional_property_write_only_read_only_base.yaml",
		OperationId: "createOneGroup",
	}, errs[0])
}
