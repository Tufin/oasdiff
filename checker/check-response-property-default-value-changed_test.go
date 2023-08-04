package checker_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/checker"
	"github.com/tufin/oasdiff/diff"
)

// CL: changing response body property default value
func TestResponsePropertyDefaultValueUpdatedCheck(t *testing.T) {
	s1, err := open("../data/checker/response_property_default_value_changed_base.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/response_property_default_value_changed_revision.yaml")
	require.NoError(t, err)
	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.ResponsePropertyDefaultValueChangedCheck), d, osm, checker.INFO)
	require.Len(t, errs, 2)
	require.ElementsMatch(t, []checker.ApiChange{{
		Id:          "response-property-default-value-changed",
		Text:        "the 'created' response's property default value changed from '2020-01-01T00:00:00Z' to '2020-02-01T00:00:00Z' for the status '200'",
		Comment:     "",
		Level:       checker.INFO,
		Operation:   "POST",
		Path:        "/api/v1.0/groups",
		Source:      "../data/checker/response_property_default_value_changed_revision.yaml",
		OperationId: "createOneGroup",
	}, {
		Id:          "response-property-default-value-changed",
		Text:        "the 'enabled' response's property default value changed from 'false' to 'true' for the status '200'",
		Comment:     "",
		Level:       checker.INFO,
		Operation:   "POST",
		Path:        "/api/v1.0/groups",
		Source:      "../data/checker/response_property_default_value_changed_revision.yaml",
		OperationId: "createOneGroup",
	}}, errs)
}

// CL: changing response body default value
func TestResponseSchemaDefaultValueUpdatedCheck(t *testing.T) {
	s1, err := open("../data/checker/response_property_default_value_changed_base.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/response_property_default_value_changed_base.yaml")
	require.NoError(t, err)

	s2.Spec.Paths["/api/v1.0/groups"].Post.Responses["404"].Value.Content["text/plain"].Schema.Value.Default = "new default value"
	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.ResponsePropertyDefaultValueChangedCheck), d, osm, checker.INFO)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{
		Id:          "response-body-default-value-changed",
		Text:        "the response body 'text/plain' default value changed from 'Error' to 'new default value' for the status '404'",
		Comment:     "",
		Level:       checker.INFO,
		Operation:   "POST",
		Path:        "/api/v1.0/groups",
		Source:      "../data/checker/response_property_default_value_changed_base.yaml",
		OperationId: "createOneGroup",
	}, errs[0])
}

// CL: adding response body default value or response body property default value
func TestResponsePropertyDefaultValueAddedCheck(t *testing.T) {
	s1, err := open("../data/checker/response_property_default_value_changed_base.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/response_property_default_value_changed_base.yaml")
	require.NoError(t, err)

	s1.Spec.Paths["/api/v1.0/groups"].Post.Responses["404"].Value.Content["text/plain"].Schema.Value.Default = nil
	s1.Spec.Components.Schemas["GroupView"].Value.Properties["data"].Value.Properties["created"].Value.Default = nil

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.ResponsePropertyDefaultValueChangedCheck), d, osm, checker.INFO)
	require.Len(t, errs, 2)
	require.ElementsMatch(t, []checker.ApiChange{{
		Id:          "response-body-default-value-added",
		Text:        "the response body 'text/plain' default value 'Error' was added for the status '404'",
		Comment:     "",
		Level:       checker.INFO,
		Operation:   "POST",
		Path:        "/api/v1.0/groups",
		Source:      "../data/checker/response_property_default_value_changed_base.yaml",
		OperationId: "createOneGroup",
	}, {
		Id:          "response-property-default-value-added",
		Text:        "the 'created' response's property default value '2020-01-01T00:00:00Z' was added for the status '200'",
		Comment:     "",
		Level:       checker.INFO,
		Operation:   "POST",
		Path:        "/api/v1.0/groups",
		Source:      "../data/checker/response_property_default_value_changed_base.yaml",
		OperationId: "createOneGroup",
	}}, errs)
}

// CL: removing response body default value or response body property default value
func TestResponsePropertyDefaultValueRemovedCheck(t *testing.T) {
	s1, err := open("../data/checker/response_property_default_value_changed_base.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/response_property_default_value_changed_base.yaml")
	require.NoError(t, err)

	s2.Spec.Paths["/api/v1.0/groups"].Post.Responses["404"].Value.Content["text/plain"].Schema.Value.Default = nil
	s2.Spec.Components.Schemas["GroupView"].Value.Properties["data"].Value.Properties["created"].Value.Default = nil

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.ResponsePropertyDefaultValueChangedCheck), d, osm, checker.INFO)
	require.Len(t, errs, 2)
	require.ElementsMatch(t, []checker.ApiChange{{
		Id:          "response-body-default-value-removed",
		Text:        "the response body 'text/plain' default value 'Error' was removed for the status '404'",
		Comment:     "",
		Level:       checker.INFO,
		Operation:   "POST",
		Path:        "/api/v1.0/groups",
		Source:      "../data/checker/response_property_default_value_changed_base.yaml",
		OperationId: "createOneGroup",
	}, {
		Id:          "response-property-default-value-removed",
		Text:        "the 'created' response's property default value '2020-01-01T00:00:00Z' was removed for the status '200'",
		Comment:     "",
		Level:       checker.INFO,
		Operation:   "POST",
		Path:        "/api/v1.0/groups",
		Source:      "../data/checker/response_property_default_value_changed_base.yaml",
		OperationId: "createOneGroup",
	}}, errs)
}
