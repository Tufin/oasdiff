package checker_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/checker"
	"github.com/tufin/oasdiff/diff"
	"github.com/tufin/oasdiff/load"
)

// CL: changing response body property default value
func TestResponsePropertyDefaultValueUpdatedCheck(t *testing.T) {
	s1, err := open("../data/checker/response_property_default_value_changed_base.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/response_property_default_value_changed_revision.yaml")
	require.NoError(t, err)
	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.ResponsePropertyDefaultValueChangedCheck), d, osm, checker.INFO)
	require.Len(t, errs, 2)

	require.Equal(t, checker.ApiChange{
		Id:          checker.ResponsePropertyDefaultValueChangedId,
		Args:        []any{"2020-01-01T00:00:00Z", "2020-02-01T00:00:00Z", "created", "200"},
		Level:       checker.INFO,
		Operation:   "POST",
		Path:        "/api/v1.0/groups",
		Source:      load.NewSource("../data/checker/response_property_default_value_changed_revision.yaml"),
		OperationId: "createOneGroup",
	}, errs[0])
	require.Equal(t, "changed default value from '2020-01-01T00:00:00Z' to '2020-02-01T00:00:00Z' for property 'created' for response status '200'", errs[0].GetUncolorizedText(checker.NewDefaultLocalizer()))

	require.Equal(t, checker.ApiChange{
		Id:          checker.ResponsePropertyDefaultValueChangedId,
		Args:        []any{false, true, "enabled", "200"},
		Level:       checker.INFO,
		Operation:   "POST",
		Path:        "/api/v1.0/groups",
		Source:      load.NewSource("../data/checker/response_property_default_value_changed_revision.yaml"),
		OperationId: "createOneGroup",
	}, errs[1])
	require.Equal(t, "changed default value from 'false' to 'true' for property 'enabled' for response status '200'", errs[1].GetUncolorizedText(checker.NewDefaultLocalizer()))
}

// CL: changing response body default value
func TestResponseSchemaDefaultValueUpdatedCheck(t *testing.T) {
	s1, err := open("../data/checker/response_property_default_value_changed_base.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/response_property_default_value_changed_base.yaml")
	require.NoError(t, err)

	s2.Spec.Paths.Value("/api/v1.0/groups").Post.Responses.Value("404").Value.Content["text/plain"].Schema.Value.Default = "new default value"
	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.ResponsePropertyDefaultValueChangedCheck), d, osm, checker.INFO)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{
		Id:          checker.ResponseBodyDefaultValueChangedId,
		Args:        []any{"text/plain", "Error", "new default value", "404"},
		Level:       checker.INFO,
		Operation:   "POST",
		Path:        "/api/v1.0/groups",
		Source:      load.NewSource("../data/checker/response_property_default_value_changed_base.yaml"),
		OperationId: "createOneGroup",
	}, errs[0])
	require.Equal(t, "response body 'text/plain' default value changed from 'Error' to 'new default value' for status '404'", errs[0].GetUncolorizedText(checker.NewDefaultLocalizer()))
}

// CL: adding response body default value or response body property default value
func TestResponsePropertyDefaultValueAddedCheck(t *testing.T) {
	s1, err := open("../data/checker/response_property_default_value_changed_base.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/response_property_default_value_changed_base.yaml")
	require.NoError(t, err)

	s1.Spec.Paths.Value("/api/v1.0/groups").Post.Responses.Value("404").Value.Content["text/plain"].Schema.Value.Default = nil
	s1.Spec.Components.Schemas["GroupView"].Value.Properties["data"].Value.Properties["created"].Value.Default = nil

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.ResponsePropertyDefaultValueChangedCheck), d, osm, checker.INFO)
	require.Len(t, errs, 2)

	require.Equal(t, checker.ApiChange{
		Id:          checker.ResponseBodyDefaultValueAddedId,
		Args:        []any{"text/plain", "Error", "404"},
		Level:       checker.INFO,
		Operation:   "POST",
		Path:        "/api/v1.0/groups",
		Source:      load.NewSource("../data/checker/response_property_default_value_changed_base.yaml"),
		OperationId: "createOneGroup",
	}, errs[0])
	require.Equal(t, "response body 'text/plain' default value 'Error' was added for status '404'", errs[0].GetUncolorizedText(checker.NewDefaultLocalizer()))

	require.Equal(t, checker.ApiChange{
		Id:          checker.ResponsePropertyDefaultValueAddedId,
		Args:        []any{"2020-01-01T00:00:00Z", "created", "200"},
		Level:       checker.INFO,
		Operation:   "POST",
		Path:        "/api/v1.0/groups",
		Source:      load.NewSource("../data/checker/response_property_default_value_changed_base.yaml"),
		OperationId: "createOneGroup",
	}, errs[1])
	require.Equal(t, "added default value '2020-01-01T00:00:00Z' to response property 'created' for response status '200'", errs[1].GetUncolorizedText(checker.NewDefaultLocalizer()))
}

// CL: removing response body default value or response body property default value
func TestResponsePropertyDefaultValueRemovedCheck(t *testing.T) {
	s1, err := open("../data/checker/response_property_default_value_changed_base.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/response_property_default_value_changed_base.yaml")
	require.NoError(t, err)

	s2.Spec.Paths.Value("/api/v1.0/groups").Post.Responses.Value("404").Value.Content["text/plain"].Schema.Value.Default = nil
	s2.Spec.Components.Schemas["GroupView"].Value.Properties["data"].Value.Properties["created"].Value.Default = nil

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.ResponsePropertyDefaultValueChangedCheck), d, osm, checker.INFO)
	require.Len(t, errs, 2)

	require.Equal(t, checker.ApiChange{
		Id:          checker.ResponseBodyDefaultValueRemovedId,
		Args:        []any{"text/plain", "Error", "404"},
		Level:       checker.INFO,
		Operation:   "POST",
		Path:        "/api/v1.0/groups",
		Source:      load.NewSource("../data/checker/response_property_default_value_changed_base.yaml"),
		OperationId: "createOneGroup",
	}, errs[0])
	require.Equal(t, "response body 'text/plain' default value 'Error' was removed for status '404'", errs[0].GetUncolorizedText(checker.NewDefaultLocalizer()))

	require.Equal(t, checker.ApiChange{
		Id:          checker.ResponsePropertyDefaultValueRemovedId,
		Args:        []any{"2020-01-01T00:00:00Z", "created", "200"},
		Level:       checker.INFO,
		Operation:   "POST",
		Path:        "/api/v1.0/groups",
		Source:      load.NewSource("../data/checker/response_property_default_value_changed_base.yaml"),
		OperationId: "createOneGroup",
	}, errs[1])
	require.Equal(t, "removed default value '2020-01-01T00:00:00Z' of property 'created' for response status '200'", errs[1].GetUncolorizedText(checker.NewDefaultLocalizer()))
}
