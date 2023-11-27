package checker_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/checker"
	"github.com/tufin/oasdiff/diff"
)

// CL: changing request property pattern
func TestRequestPropertyPatternChanged(t *testing.T) {
	s1, err := open("../data/checker/request_property_pattern_added_or_changed_base.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/request_property_pattern_added_or_changed_revision.yaml")
	require.NoError(t, err)

	s2.Spec.Paths["/test"].Post.RequestBody.Value.Content["application/json"].Schema.Value.Properties["name"].Value.Pattern = "^[\\w\\s]+$"

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.RequestPropertyPatternUpdatedCheck), d, osm, checker.INFO)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{
		Id:          checker.RequestPropertyPatternChangedId,
		Args:        []any{"name", "^\\w+$", "^[\\w\\s]+$"},
		Level:       checker.WARN,
		Operation:   "POST",
		Path:        "/test",
		Source:      "../data/checker/request_property_pattern_added_or_changed_revision.yaml",
		OperationId: "",
		Comment:     "This is a warning because it is difficult to automatically analyze if the new pattern is a superset of the previous pattern (e.g. changed from '[0-9]+' to '[0-9]*')",
	}, errs[0])
}

// CL: adding request property pattern
func TestRequestPropertyPatternAdded(t *testing.T) {
	s1, err := open("../data/checker/request_property_pattern_added_or_changed_revision.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/request_property_pattern_added_or_changed_base.yaml")
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.RequestPropertyPatternUpdatedCheck), d, osm, checker.INFO)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{
		Id:          checker.RequestPropertyPatternAddedId,
		Args:        []any{"^\\w+$", "name"},
		Level:       checker.WARN,
		Operation:   "POST",
		Path:        "/test",
		Source:      "../data/checker/request_property_pattern_added_or_changed_base.yaml",
		OperationId: "",
		Comment:     "This is a warning because it is difficult to automatically analyze if the new pattern is a superset of the previous pattern (e.g. changed from '[0-9]+' to '[0-9]*')",
	}, errs[0])
}

// CL: removing request property pattern
func TestRequestPropertyPatternRemoved(t *testing.T) {
	s1, err := open("../data/checker/request_property_pattern_added_or_changed_base.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/request_property_pattern_added_or_changed_revision.yaml")
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.RequestPropertyPatternUpdatedCheck), d, osm, checker.INFO)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{
		Id:          checker.RequestPropertyPatternRemovedId,
		Args:        []any{"^\\w+$", "name"},
		Level:       checker.INFO,
		Operation:   "POST",
		Path:        "/test",
		Source:      "../data/checker/request_property_pattern_added_or_changed_revision.yaml",
		OperationId: "",
	}, errs[0])
}
