package checker_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/checker"
	"github.com/tufin/oasdiff/diff"
)

// CL: changing pattern of request parameters
func TestRequestParameterPatternChanged(t *testing.T) {
	s1, err := open("../data/checker/request_parameter_pattern_added_or_changed_base.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/request_parameter_pattern_added_or_changed_base.yaml")
	require.NoError(t, err)

	s2.Spec.Paths["/test"].Post.Parameters[0].Value.Schema.Value.Pattern = "^[\\w\\s]+$"
	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.RequestParameterPatternAddedOrChangedCheck), d, osm, checker.WARN)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{
		Id:          checker.RequestParameterPatternChangedId,
		Text:        "changed the pattern of the 'query' request parameter 'category' from '^\\w+$' to '^[\\w\\s]+$'",
		Args:        []any{"query", "category", "^\\w+$", "^[\\w\\s]+$"},
		Comment:     "This is a warning because it is difficult to automatically analyze if the new pattern is a superset of the previous pattern (e.g. changed from '[0-9]+' to '[0-9]*')",
		Level:       checker.WARN,
		Operation:   "POST",
		Path:        "/test",
		Source:      "../data/checker/request_parameter_pattern_added_or_changed_base.yaml",
		OperationId: "",
	}, errs[0])
}

// CL: adding pattern to request parameters
func TestRequestParameterPatternAdded(t *testing.T) {
	s1, err := open("../data/checker/request_parameter_pattern_added_or_changed_revision.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/request_parameter_pattern_added_or_changed_base.yaml")
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.RequestParameterPatternAddedOrChangedCheck), d, osm, checker.WARN)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{
		Id:          checker.RequestParameterPatternAddedId,
		Text:        "added the pattern '^\\w+$' to the 'query' request parameter 'category'",
		Args:        []any{"^\\w+$", "query", "category"},
		Comment:     "This is a warning because it is difficult to automatically analyze if the new pattern is a superset of the previous pattern (e.g. changed from '[0-9]+' to '[0-9]*')",
		Level:       checker.WARN,
		Operation:   "POST",
		Path:        "/test",
		Source:      "../data/checker/request_parameter_pattern_added_or_changed_base.yaml",
		OperationId: "",
	}, errs[0])
}

// CL: removing pattern from request parameters
func TestRequestParameterPatternRemoved(t *testing.T) {
	s1, err := open("../data/checker/request_parameter_pattern_added_or_changed_base.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/request_parameter_pattern_added_or_changed_revision.yaml")
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.RequestParameterPatternAddedOrChangedCheck), d, osm, checker.INFO)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{
		Id:          checker.RequestParameterPatternRemovedId,
		Text:        "removed the pattern '^\\w+$' from the 'query' request parameter 'category'",
		Args:        []any{"^\\w+$", "query", "category"},
		Level:       checker.INFO,
		Operation:   "POST",
		Path:        "/test",
		Source:      "../data/checker/request_parameter_pattern_added_or_changed_revision.yaml",
		OperationId: "",
	}, errs[0])
}
