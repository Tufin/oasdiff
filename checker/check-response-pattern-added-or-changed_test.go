package checker_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/checker"
	"github.com/tufin/oasdiff/diff"
)

// CL: changing response property pattern
func TestResponsePropertyPatternChanged(t *testing.T) {
	s1, err := open("../data/checker/response_pattern_added_or_changed_base.yaml")
	require.NoError(t, err)

	s2, err := open("../data/checker/response_pattern_added_or_changed_revision.yaml")
	require.NoError(t, err)

	s2.Spec.Components.Schemas["GroupView"].Value.Properties["data"].Value.Properties["name"].Value.WriteOnly = true
	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.ResponsePatternAddedOrChangedCheck), d, osm, checker.INFO)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{

		Id:          "response-property-pattern-changed",
		Text:        "the 'data/created' response's property pattern was changed from '^[a-z]+$' to '^(?:([a-z]+-)*([a-z]+)?)$' for the status '200'",
		Comment:     "",
		Level:       checker.INFO,
		Operation:   "POST",
		Path:        "/api/v1.0/groups",
		Source:      "../data/checker/response_pattern_added_or_changed_revision.yaml",
		OperationId: "createOneGroup",
	}, errs[0])
}

// CL: adding response property pattern
func TestResponsePropertyPatternAdded(t *testing.T) {
	s1, err := open("../data/checker/response_pattern_added_or_changed_base.yaml")
	require.NoError(t, err)

	s2, err := open("../data/checker/response_pattern_added_or_changed_base.yaml")
	require.NoError(t, err)

	s1.Spec.Components.Schemas["GroupView"].Value.Properties["data"].Value.Properties["created"].Value.Pattern = ""
	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.ResponsePatternAddedOrChangedCheck), d, osm, checker.INFO)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{

		Id:          "response-property-pattern-added",
		Text:        "the 'data/created' response's property pattern '^[a-z]+$' was added for the status '200'",
		Comment:     "",
		Level:       checker.INFO,
		Operation:   "POST",
		Path:        "/api/v1.0/groups",
		Source:      "../data/checker/response_pattern_added_or_changed_base.yaml",
		OperationId: "createOneGroup",
	}, errs[0])
}

// CL: removing response property pattern
func TestResponsePropertyPatternRemoved(t *testing.T) {
	s1, err := open("../data/checker/response_pattern_added_or_changed_base.yaml")
	require.NoError(t, err)

	s2, err := open("../data/checker/response_pattern_added_or_changed_base.yaml")
	require.NoError(t, err)

	s2.Spec.Components.Schemas["GroupView"].Value.Properties["data"].Value.Properties["created"].Value.Pattern = ""
	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.ResponsePatternAddedOrChangedCheck), d, osm, checker.INFO)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{

		Id:          "response-property-pattern-removed",
		Text:        "the 'data/created' response's property pattern '^[a-z]+$' was removed for the status '200'",
		Comment:     "",
		Level:       checker.INFO,
		Operation:   "POST",
		Path:        "/api/v1.0/groups",
		Source:      "../data/checker/response_pattern_added_or_changed_base.yaml",
		OperationId: "createOneGroup",
	}, errs[0])
}
