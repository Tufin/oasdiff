package checker_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/checker"
	"github.com/tufin/oasdiff/diff"
)

// CL: adding a new tag
func TestTagAdded(t *testing.T) {
	s1, err := open("../data/checker/tag_added_base.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/tag_added_base.yaml")
	require.NoError(t, err)

	s2.Spec.Paths["/api/v1.0/groups"].Post.Tags = []string{"newTag"}

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.APITagUpdatedCheck), d, osm, checker.INFO)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{
		Id:          "api-tag-added",
		Text:        "api tag 'newTag' added",
		Comment:     "",
		Level:       checker.INFO,
		Operation:   "POST",
		Path:        "/api/v1.0/groups",
		Source:      "../data/checker/tag_added_base.yaml",
		OperationId: "createOneGroup",
	}, errs[0])
}

// CL: removing an existing tag
func TestTagRemoved(t *testing.T) {
	s1, err := open("../data/checker/tag_removed_base.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/tag_removed_base.yaml")
	require.NoError(t, err)

	s2.Spec.Paths["/api/v1.0/groups"].Post.Tags = []string{}

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.APITagUpdatedCheck), d, osm, checker.INFO)
	require.NotEmpty(t, errs)
	require.Len(t, errs, 1)
	require.Equal(t, checker.ApiChange{
		Id:          "api-tag-removed",
		Text:        "api tag 'Test' removed",
		Comment:     "",
		Level:       checker.INFO,
		Operation:   "POST",
		Path:        "/api/v1.0/groups",
		Source:      "../data/checker/tag_removed_base.yaml",
		OperationId: "createOneGroup",
	}, errs[0])
}

// CL: updating an existing tag
func TestTagUpdated(t *testing.T) {
	s1, err := open("../data/checker/tag_removed_base.yaml")
	require.NoError(t, err)
	s2, err := open("../data/checker/tag_removed_base.yaml")
	require.NoError(t, err)

	s2.Spec.Paths["/api/v1.0/groups"].Post.Tags = []string{"newTag"}

	d, osm, err := diff.GetWithOperationsSourcesMap(getConfig(), s1, s2)
	require.NoError(t, err)
	errs := checker.CheckBackwardCompatibilityUntilLevel(singleCheckConfig(checker.APITagUpdatedCheck), d, osm, checker.INFO)
	require.NotEmpty(t, errs)
	require.Len(t, errs, 2)
	for cl := range errs {
		require.Equal(t, checker.INFO, errs[cl].GetLevel())
		if errs[cl].GetId() == "api-tag-removed" {
			require.Equal(t, checker.ApiChange{
				Id:          "api-tag-removed",
				Text:        "api tag 'Test' removed",
				Comment:     "",
				Level:       checker.INFO,
				Operation:   "POST",
				Path:        "/api/v1.0/groups",
				Source:      "../data/checker/tag_removed_base.yaml",
				OperationId: "createOneGroup",
			}, errs[cl])
		}

		if errs[cl].GetId() == "api-tag-added" {
			require.Equal(t, checker.ApiChange{
				Id:          "api-tag-added",
				Text:        "api tag 'newTag' added",
				Comment:     "",
				Level:       checker.INFO,
				Operation:   "POST",
				Path:        "/api/v1.0/groups",
				Source:      "../data/checker/tag_removed_base.yaml",
				OperationId: "createOneGroup",
			}, errs[cl])
		}
	}
}
