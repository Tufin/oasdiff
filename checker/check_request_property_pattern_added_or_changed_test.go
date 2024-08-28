package checker_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/checker"
	"github.com/tufin/oasdiff/diff"
	"github.com/tufin/oasdiff/load"
)

// CL: changing request property pattern
func TestRequestPropertyPatternChanged(t *testing.T) {
	s1, err := open("../data/checker/request_property_pattern_added_or_changed_base.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/request_property_pattern_added_or_changed_revision.yaml")
	require.NoError(t, err)

	s2.Spec.Paths.Value("/test").Post.RequestBody.Value.Content["application/json"].Schema.Value.Properties["name"].Value.Pattern = "^[\\w\\s]+$"

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.RequestPropertyPatternUpdatedCheck), d, osm, checker.INFO)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{
		Id:        checker.RequestPropertyPatternChangedId,
		Args:      []any{"name", "application/json", "^\\w+$", "^[\\w\\s]+$"},
		Level:     checker.WARN,
		Operation: "POST",
		Path:      "/test",
		Source:    load.NewSource("../data/checker/request_property_pattern_added_or_changed_revision.yaml"),
		Comment:   checker.PatternChangedCommentId,
	}, errs[0])
	require.Equal(t, "changed pattern of the request property 'name' of media-type 'application/json' from '^\\w+$' to '^[\\w\\s]+$'", errs[0].GetUncolorizedText(checker.NewDefaultLocalizer()))
	require.Equal(t, "This is a warning because it is difficult to automatically analyze if the new pattern is a superset of the previous pattern (e.g. changed from '[0-9]+' to '[0-9]*')", errs[0].GetComment(checker.NewDefaultLocalizer()))
}

// CL: generalizing request property pattern
func TestRequestPropertyPatternGeneralized(t *testing.T) {
	s1, err := open("../data/checker/request_property_pattern_added_or_changed_base.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/request_property_pattern_added_or_changed_revision.yaml")
	require.NoError(t, err)

	s2.Spec.Paths.Value("/test").Post.RequestBody.Value.Content["application/json"].Schema.Value.Properties["name"].Value.Pattern = ".*"

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.RequestPropertyPatternUpdatedCheck), d, osm, checker.INFO)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{
		Id:        checker.RequestPropertyPatternGeneralizedId,
		Args:      []any{"name", "application/json", "^\\w+$", ".*"},
		Level:     checker.INFO,
		Operation: "POST",
		Path:      "/test",
		Source:    load.NewSource("../data/checker/request_property_pattern_added_or_changed_revision.yaml"),
	}, errs[0])
	require.Equal(t, "changed pattern of the request property 'name' of media-type 'application/json' from '^\\w+$' to a more general pattern '.*'", errs[0].GetUncolorizedText(checker.NewDefaultLocalizer()))
}

// CL: adding request property pattern
func TestRequestPropertyPatternAdded(t *testing.T) {
	s1, err := open("../data/checker/request_property_pattern_added_or_changed_revision.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/request_property_pattern_added_or_changed_base.yaml")
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.RequestPropertyPatternUpdatedCheck), d, osm, checker.INFO)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{
		Id:        checker.RequestPropertyPatternAddedId,
		Args:      []any{"^\\w+$", "name", "application/json"},
		Level:     checker.WARN,
		Operation: "POST",
		Path:      "/test",
		Source:    load.NewSource("../data/checker/request_property_pattern_added_or_changed_base.yaml"),
		Comment:   checker.PatternChangedCommentId,
	}, errs[0])
	require.Equal(t, "This is a warning because it is difficult to automatically analyze if the new pattern is a superset of the previous pattern (e.g. changed from '[0-9]+' to '[0-9]*')", errs[0].GetComment(checker.NewDefaultLocalizer()))
	require.Equal(t, "added pattern '^\\w+$' to the request property 'name' of media-type 'application/json'", errs[0].GetUncolorizedText(checker.NewDefaultLocalizer()))
}

// CL: removing request property pattern
func TestRequestPropertyPatternRemoved(t *testing.T) {
	s1, err := open("../data/checker/request_property_pattern_added_or_changed_base.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/request_property_pattern_added_or_changed_revision.yaml")
	require.NoError(t, err)

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.RequestPropertyPatternUpdatedCheck), d, osm, checker.INFO)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{
		Id:        checker.RequestPropertyPatternRemovedId,
		Args:      []any{"^\\w+$", "name", "application/json"},
		Level:     checker.INFO,
		Operation: "POST",
		Path:      "/test",
		Source:    load.NewSource("../data/checker/request_property_pattern_added_or_changed_revision.yaml"),
	}, errs[0])
	require.Equal(t, "removed pattern '^\\w+$' from the request property 'name' of media-type 'application/json'", errs[0].GetUncolorizedText(checker.NewDefaultLocalizer()))
}
