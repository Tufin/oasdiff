package checker_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/checker"
	"github.com/tufin/oasdiff/diff"
	"github.com/tufin/oasdiff/load"
)

// CL: adding a new tag
func TestTagAdded(t *testing.T) {
	s1, err := open("../data/checker/tag_added_base.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/tag_added_base.yaml")
	require.NoError(t, err)

	s2.Spec.Paths.Value("/api/v1.0/groups").Post.Tags = []string{"newTag"}

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.APITagUpdatedCheck), d, osm, checker.INFO)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{
		Id:          checker.APITagAddedId,
		Args:        []any{"newTag"},
		Level:       checker.INFO,
		Operation:   "POST",
		Path:        "/api/v1.0/groups",
		Source:      load.NewSource("../data/checker/tag_added_base.yaml"),
		OperationId: "createOneGroup",
	}, errs[0])
	require.Equal(t, "added api tag 'newTag'", errs[0].GetUncolorizedText(checker.NewDefaultLocalizer()))
}

// CL: removing an existing tag
func TestTagRemoved(t *testing.T) {
	s1, err := open("../data/checker/tag_removed_base.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/tag_removed_base.yaml")
	require.NoError(t, err)

	s2.Spec.Paths.Value("/api/v1.0/groups").Post.Tags = []string{}

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.APITagUpdatedCheck), d, osm, checker.INFO)
	require.NotEmpty(t, errs)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{
		Id:          checker.APITagRemovedId,
		Args:        []any{"Test"},
		Level:       checker.INFO,
		Operation:   "POST",
		Path:        "/api/v1.0/groups",
		Source:      load.NewSource("../data/checker/tag_removed_base.yaml"),
		OperationId: "createOneGroup",
	}, errs[0])
	require.Equal(t, "removed api tag 'Test'", errs[0].GetUncolorizedText(checker.NewDefaultLocalizer()))

}

// CL: updating an existing tag
func TestTagUpdated(t *testing.T) {
	s1, err := open("../data/checker/tag_removed_base.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/tag_removed_base.yaml")
	require.NoError(t, err)

	s2.Spec.Paths.Value("/api/v1.0/groups").Post.Tags = []string{"newTag"}

	d, osm, err := diff.GetWithOperationsSourcesMap(diff.NewConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.APITagUpdatedCheck), d, osm, checker.INFO)
	require.NotEmpty(t, errs)
	require.Len(t, errs, 2)
	for cl := range errs {
		require.Equal(t, checker.INFO, errs[cl].GetLevel())
		if errs[cl].GetId() == checker.APITagRemovedId {
			require.Equal(t, checker.ApiChange{
				Id:          checker.APITagRemovedId,
				Args:        []any{"Test"},
				Level:       checker.INFO,
				Operation:   "POST",
				Path:        "/api/v1.0/groups",
				Source:      load.NewSource("../data/checker/tag_removed_base.yaml"),
				OperationId: "createOneGroup",
			}, errs[cl])
		}

		if errs[cl].GetId() == checker.APITagAddedId {
			require.Equal(t, checker.ApiChange{
				Id:          checker.APITagAddedId,
				Args:        []any{"newTag"},
				Level:       checker.INFO,
				Operation:   "POST",
				Path:        "/api/v1.0/groups",
				Source:      load.NewSource("../data/checker/tag_removed_base.yaml"),
				OperationId: "createOneGroup",
			}, errs[cl])
		}
	}
}
