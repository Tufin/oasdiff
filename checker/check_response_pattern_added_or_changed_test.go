package checker_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/checker"
	"github.com/tufin/oasdiff/diff"
	"github.com/tufin/oasdiff/load"
)

// CL: changing response property pattern
func TestResponsePropertyPatternChanged(t *testing.T) {
	s1, err := open("../data/checker/response_pattern_added_or_changed_base.yaml")
	require.NoError(t, err)

	s2, err := open("../data/checker/response_pattern_added_or_changed_revision.yaml")
	require.NoError(t, err)

	s2.Spec.Components.Schemas["GroupView"].Value.Properties["data"].Value.Properties["name"].Value.WriteOnly = true
	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.ResponsePatternAddedOrChangedCheck), d, osm, checker.INFO)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{
		Id:          checker.ResponsePropertyPatternChangedId,
		Args:        []any{"^[a-z]+$", "^(?:([a-z]+-)*([a-z]+)?)$", "data/created", "200"},
		Level:       checker.INFO,
		Operation:   "POST",
		Path:        "/api/v1.0/groups",
		Source:      load.NewSource("../data/checker/response_pattern_added_or_changed_revision.yaml"),
		OperationId: "createOneGroup",
	}, errs[0])
	require.Equal(t, "changed pattern from '^[a-z]+$' to '^(?:([a-z]+-)*([a-z]+)?)$' in property 'data/created' for response status '200'", errs[0].GetUncolorizedText(checker.NewDefaultLocalizer()))
}

// CL: adding response property pattern
func TestResponsePropertyPatternAdded(t *testing.T) {
	s1, err := open("../data/checker/response_pattern_added_or_changed_base.yaml")
	require.NoError(t, err)

	s2, err := open("../data/checker/response_pattern_added_or_changed_base.yaml")
	require.NoError(t, err)

	s1.Spec.Components.Schemas["GroupView"].Value.Properties["data"].Value.Properties["created"].Value.Pattern = ""
	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.ResponsePatternAddedOrChangedCheck), d, osm, checker.INFO)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{
		Id:          checker.ResponsePropertyPatternAddedId,
		Args:        []any{"^[a-z]+$", "data/created", "200"},
		Level:       checker.INFO,
		Operation:   "POST",
		Path:        "/api/v1.0/groups",
		Source:      load.NewSource("../data/checker/response_pattern_added_or_changed_base.yaml"),
		OperationId: "createOneGroup",
	}, errs[0])
	require.Equal(t, "added pattern '^[a-z]+$' to property 'data/created' for response status '200'", errs[0].GetUncolorizedText(checker.NewDefaultLocalizer()))
}

// CL: removing response property pattern
func TestResponsePropertyPatternRemoved(t *testing.T) {
	s1, err := open("../data/checker/response_pattern_added_or_changed_base.yaml")
	require.NoError(t, err)

	s2, err := open("../data/checker/response_pattern_added_or_changed_base.yaml")
	require.NoError(t, err)

	s2.Spec.Components.Schemas["GroupView"].Value.Properties["data"].Value.Properties["created"].Value.Pattern = ""
	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.ResponsePatternAddedOrChangedCheck), d, osm, checker.INFO)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{
		Id:          checker.ResponsePropertyPatternRemovedId,
		Args:        []any{"^[a-z]+$", "data/created", "200"},
		Level:       checker.INFO,
		Operation:   "POST",
		Path:        "/api/v1.0/groups",
		Source:      load.NewSource("../data/checker/response_pattern_added_or_changed_base.yaml"),
		OperationId: "createOneGroup",
	}, errs[0])
	require.Equal(t, "removed pattern '^[a-z]+$' from property 'data/created' for response status '200'", errs[0].GetUncolorizedText(checker.NewDefaultLocalizer()))
}
